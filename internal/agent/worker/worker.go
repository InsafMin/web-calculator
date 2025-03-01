package worker

import (
	"encoding/json"
	"fmt"
	"github.com/InsafMin/web_calculator/pkg/calculator"
	"github.com/InsafMin/web_calculator/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

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
	globalMutex sync.Mutex
)

func StartWorker() {
	for {
		globalMutex.Lock()

		task, err := fetchTask()
		if err != nil {
			if errors.Is(err, errors.ErrNoTasksAvailable) {
				globalMutex.Unlock()
				time.Sleep(1 * time.Second)
				continue
			}
			fmt.Printf("Error fetching task: %v\n", err)
			globalMutex.Unlock()
			continue
		}

		result, err := executeTask(task)
		if err != nil {
			fmt.Printf("Error executing task %s: %v\n", task.ID, err)
			globalMutex.Unlock()
			continue
		}

		if err := sendResult(task.ID, result); err != nil {
			fmt.Printf("Error sending result for task %s: %v\n", task.ID, err)
			globalMutex.Unlock()
			continue
		}

		if task.Done != nil {
			close(task.Done)
		}

		log.Printf("Task finished successfully: %+v. Res = %f\n", task, result)

		globalMutex.Unlock()
	}
}

func fetchTask() (*Task, error) {
	orchestratorURL := os.Getenv("ORCHESTRATOR_URL")
	if orchestratorURL == "" {
		orchestratorURL = "http://localhost:8080"
	}

	resp, err := http.Get(orchestratorURL + "/internal/task")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.ErrNoTasksAvailable
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Task Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode task: %w", err)
	}

	return &response.Task, nil
}

func executeTask(task *Task) (float64, error) {
	time.Sleep(task.OperationTime)

	result, err := calculator.Resolve(task.Arg1, task.Arg2, task.Operation)
	if err != nil {
		return 0, fmt.Errorf("failed to resolve task: %w", err)
	}

	return result, nil
}

func sendResult(taskID string, result float64) error {
	payload := struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}{
		ID:     taskID,
		Result: result,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	orchestratorURL := os.Getenv("ORCHESTRATOR_URL")
	if orchestratorURL == "" {
		orchestratorURL = "http://localhost:8080"
	}

	resp, err := http.Post(orchestratorURL+"/internal/task", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to send result: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
