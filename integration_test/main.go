package main

import (
	"fmt"
	"os"

	"github.com/Andrewmakmaer/SystemInfoFluxDaemon/internal/integration"
)

func main() {
	if err := integration.RunIntegrationTests(); err != nil {
		fmt.Printf("Ошибка при выполнении интеграционных тестов: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Все тесты прошли успешно!")
	os.Exit(0)
}
