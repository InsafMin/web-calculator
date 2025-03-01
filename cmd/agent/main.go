package main

import (
	"github.com/InsafMin/web_calculator/internal/agent/worker"
	"log"
	"os"
	"strconv"
)

func main() {
	// Получаем количество горутин (вычислителей) из переменной окружения
	computingPowerStr := os.Getenv("COMPUTING_POWER")
	if computingPowerStr == "" {
		computingPowerStr = "4" // Значение по умолчанию
	}

	computingPower, err := strconv.Atoi(computingPowerStr)
	if err != nil {
		log.Fatalf("Invalid COMPUTING_POWER value: %v", err)
	}

	// Запускаем горутины для обработки задач
	for i := 0; i < computingPower; i++ {
		go worker.StartWorker()
	}

	// Бесконечный цикл для поддержания работы агента
	log.Println("Agent started with", computingPower, "workers")
	select {}
}
