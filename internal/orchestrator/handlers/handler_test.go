package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected []*Task
	}{
		{
			"1 + 2",
			[]*Task{
				{
					ID:           "1-1",
					Arg1:         1,
					Arg2:         2,
					Operation:    "+",
					ExpressionID: "1",
					Priority:     1,
				},
			},
		},
	}

	for _, test := range tests {
		result, err := parseExpression(test.input, "1")
		if err != nil {
			t.Errorf("parseExpression(%s) returned error: %v", test.input, err)
		}
		if len(result) != len(test.expected) {
			t.Errorf("parseExpression(%s) returned %d tasks, expected %d", test.input, len(result), len(test.expected))
		}
		for i := range result {
			if result[i].ID != test.expected[i].ID ||
				result[i].Arg1 != test.expected[i].Arg1 ||
				result[i].Arg2 != test.expected[i].Arg2 ||
				result[i].Operation != test.expected[i].Operation ||
				result[i].Priority != test.expected[i].Priority {
				t.Errorf("parseExpression(%s) = %v, expected %v", test.input, result[i], test.expected[i])
			}
		}
	}
}

func TestUpdateTaskArgs(t *testing.T) {
	tasks := map[string]*Task{
		"1-1": {
			ID:           "1-1",
			Arg1:         0,
			Arg2:         2,
			Operation:    "+",
			ExpressionID: "1",
			Priority:     1,
		},
		"1-2": {
			ID:           "1-2",
			Arg1:         0,
			Arg2:         3,
			Operation:    "*",
			ExpressionID: "1",
			Priority:     2,
		},
	}

	updateTaskArgs(tasks, "1-1", 3)

	if tasks["1-2"].Arg1 != 3 {
		t.Errorf("updateTaskArgs did not update Arg1 in task 1-2")
	}
}

func TestHandleCalculate(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(`{"expression": "1 + 2"}`))

	HandleCalculate(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("HandleCalculate returned status code %d, expected %d", w.Code, http.StatusCreated)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("HandleCalculate returned invalid JSON: %v", err)
	}

	if _, exists := response["id"]; !exists {
		t.Errorf("HandleCalculate did not return task ID")
	}
}

func TestHandleTaskGet(t *testing.T) {
	tasks["1-1"] = &Task{
		ID:           "1-1",
		Arg1:         1,
		Arg2:         2,
		Operation:    "+",
		ExpressionID: "1",
		Priority:     1,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/internal/task", nil)

	HandleTask(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("HandleTask returned status code %d, expected %d", w.Code, http.StatusOK)
	}

	var response struct {
		Task Task `json:"task"`
	}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("HandleTask returned invalid JSON: %v", err)
	}

	if response.Task.ID != "1-1" {
		t.Errorf("HandleTask returned task with ID %s, expected 1-1", response.Task.ID)
	}
}
