package main

import (
	"github.com/InsafMin/web_calculator/internal/agent/worker"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	computingPowerStr := os.Getenv("COMPUTING_POWER")
	if computingPowerStr == "" {
		computingPowerStr = "4" // Значение по умолчанию
	}

	computingPower, err := strconv.Atoi(computingPowerStr)
	if err != nil {
		log.Fatalf("Invalid COMPUTING_POWER value: %v", err)
	}

	log.Println("Waiting for orchestrator to start...")
	time.Sleep(5 * time.Second)

	for i := 0; i < computingPower; i++ {
		go worker.StartWorker()
	}

	log.Println("Agent started with", computingPower, "workers")
	select {}
}
