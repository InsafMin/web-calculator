package main

import (
	"github.com/InsafMin/web_calculator/internal/orchestrator/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/calculate", handlers.HandleCalculate)
	http.HandleFunc("/api/v1/expressions", handlers.HandleGetExpressions)
	http.HandleFunc("/api/v1/expressions/", handlers.HandleGetExpression)
	http.HandleFunc("/internal/task", handlers.HandleTask)

	log.Println("Starting orchestrator on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
