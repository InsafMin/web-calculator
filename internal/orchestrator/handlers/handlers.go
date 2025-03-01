package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/InsafMin/web_calculator/pkg/calculator"
	"github.com/InsafMin/web_calculator/pkg/errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Expression struct {
	ID     string  `json:"id"`
	Expr   string  `json:"expression"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

type Task struct {
	ID            string        `json:"id"`
	Arg1          float64       `json:"arg1"`
	Arg2          float64       `json:"arg2"`
	Operation     string        `json:"operation"`
	OperationTime time.Duration `json:"operation_time"`
	ExpressionID  string        `json:"expression_id"`
	Priority      int           `json:"priority"`
	Done          chan bool     `json:"-"`
}

var (
	expressions = make(map[string]*Expression)
	tasks       = make(map[string]*Task)
	mutex       = &sync.Mutex{}
)

func HandleCalculate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
		return
	}

	id := fmt.Sprintf("%d", time.Now().UnixNano())

	expr := &Expression{
		ID:     id,
		Expr:   req.Expression,
		Status: "pending",
	}

	mutex.Lock()
	expressions[id] = expr
	mutex.Unlock()

	tasksList, err := parseExpression(req.Expression, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	mutex.Lock()
	for _, task := range tasksList {
		fmt.Printf("Added task: %+v\n", task)
		tasks[task.ID] = task
	}
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func getOperationTime(operation string) time.Duration {
	var timeMs int
	var err error

	switch operation {
	case "+":
		timeMs, err = strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	case "-":
		timeMs, err = strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	case "*":
		timeMs, err = strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	case "/":
		timeMs, err = strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
	default:
		return 0
	}

	if err != nil {
		switch operation {
		case "+":
			return 100 * time.Millisecond
		case "-":
			return 100 * time.Millisecond
		case "*":
			return 200 * time.Millisecond
		case "/":
			return 200 * time.Millisecond
		default:
			return 0
		}
	}

	return time.Duration(timeMs) * time.Millisecond
}

func HandleGetExpressions(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var exprs []Expression
	for _, expr := range expressions {
		exprs = append(exprs, *expr)
	}

	json.NewEncoder(w).Encode(map[string][]Expression{"expressions": exprs})
}

func HandleGetExpression(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/v1/expressions/"):]

	mutex.Lock()
	defer mutex.Unlock()

	expr, exists := expressions[id]
	if !exists {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]Expression{"expression": *expr})
}

func HandleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		mutex.Lock()
		defer mutex.Unlock()

		var nextTask *Task
		for _, task := range tasks {
			if nextTask == nil || task.Priority > nextTask.Priority {
				nextTask = task
			}
		}

		if nextTask != nil {
			log.Printf("Sending task to agent: %+v\n", nextTask)

			response := struct {
				Task Task `json:"task"`
			}{
				Task: *nextTask,
			}

			delete(tasks, nextTask.ID)

			taskCopy := *nextTask
			taskCopy.Done = nil

			response.Task = taskCopy

			json.NewEncoder(w).Encode(response)
			return
		}

		http.Error(w, "No tasks available", http.StatusNotFound)
	} else if r.Method == http.MethodPost {
		var req struct {
			ID     string  `json:"id"`
			Result float64 `json:"result"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		taskID := req.ID
		exprID := strings.Split(taskID, "-")[0] // ID выражения из ID задачи

		expr, exists := expressions[exprID]
		if !exists {
			http.Error(w, "Expression not found", http.StatusNotFound)
			return
		}

		updateTaskArgs(tasks, taskID, req.Result)

		allTasksDone := true
		for _, t := range tasks {
			if t.ExpressionID == exprID {
				allTasksDone = false
				break
			}
		}

		if allTasksDone {
			expr.Result = req.Result
			expr.Status = "done"
		}

		w.WriteHeader(http.StatusOK)
	}
}

func updateTaskArgs(tasks map[string]*Task, taskID string, result float64) {
	for _, task := range tasks {
		if task.Arg1 == 0 && strings.HasPrefix(taskID, task.ExpressionID) {
			task.Arg1 = result
		}
		if task.Arg2 == 0 && strings.HasPrefix(taskID, task.ExpressionID) {
			task.Arg2 = result
		}
	}
}

func parseExpression(expr string, exprID string) ([]*Task, error) {
	tokens, err := calculator.Tokenize(expr)
	if err != nil {
		return nil, err
	}

	rpn, err := calculator.ToRPN(tokens)
	if err != nil {
		return nil, err
	}

	bracketLevels := getBracketLevels(tokens, rpn)

	var tasks []*Task
	var stack []string
	taskMap := make(map[string]string)
	taskCounter := 1

	for _, token := range rpn {
		if _, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, token)
		} else if calculator.IsOperator(rune(token[0])) {
			if len(stack) < 2 {
				return nil, errors.ErrInvalidExpression
			}

			arg2 := stack[len(stack)-1]
			arg1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			taskID := fmt.Sprintf("%s-%d", exprID, taskCounter)

			if val, exists := taskMap[arg1]; exists {
				arg1 = val
			}
			if val, exists := taskMap[arg2]; exists {
				arg2 = val
			}

			priority := calculator.Priority(token) + bracketLevels[token]

			task := &Task{
				ID:            taskID,
				Arg1:          parseNumber(arg1),
				Arg2:          parseNumber(arg2),
				Operation:     token,
				ExpressionID:  exprID,
				Priority:      priority,
				OperationTime: getOperationTime(token),
				Done:          make(chan bool),
			}
			tasks = append(tasks, task)

			resultKey := fmt.Sprintf("task-%s", taskID)
			taskMap[resultKey] = taskID

			stack = append(stack, resultKey)

			taskCounter++
		}
	}

	return tasks, nil
}

func getBracketLevels(tokens []string, rpn []string) map[string]int {
	bracketLevels := make(map[string]int)

	currentLevel := 0
	tokenLevels := make(map[string]int)

	for _, token := range tokens {
		if token == "(" {
			currentLevel += 2
		} else if token == ")" {
			currentLevel -= 2
		} else if calculator.IsOperator(rune(token[0])) {
			tokenLevels[token] = currentLevel
		}
	}

	for _, token := range rpn {
		if calculator.IsOperator(rune(token[0])) {
			bracketLevels[token] = tokenLevels[token]
		}
	}

	return bracketLevels
}

func parseNumber(s string) float64 {
	if strings.HasPrefix(s, "task-") {
		return 0
	}
	num, _ := strconv.ParseFloat(s, 64)
	return num
}
