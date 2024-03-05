package cli

import (
	"client/net/controller"
	"client/net/messages"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

type Controller struct {
	Os  string
	Net controller.MessageController
}

func New(net controller.MessageController) Controller {
	return Controller{Os: runtime.GOOS, Net: net}
}

// Функція для очистки консолі
func (c *Controller) clearConsole() {
	if c.Os == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else if c.Os == "linux" || c.Os == "darwin" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		for i := 0; i < 3; i++ {
			println()
		}
	}
}

// Запуск CLI
func (c *Controller) StartCli() {
	errorMessage := ""
	for {
		c.clearConsole()
		fmt.Printf("Lab2. Sum of two numbers (Variant 2). Protocol: %s\n", c.Net.Protocol)
		println()

		if errorMessage != "" {
			fmt.Printf("An error occured: %s\n", errorMessage)
			errorMessage = ""
			println()

		}

		println("Enter the first number")
		var num1String string
		_, err := fmt.Scan(&num1String)
		if err != nil {
			errorMessage = fmt.Sprintf("failed to scan the first number: %s", err)
			continue
		}

		// Парсинг з рядка в число першого числа
		num1, err := strconv.ParseInt(num1String, 10, 64)
		if err != nil {
			errorMessage = fmt.Sprint("you need to enter a number for the first value")
			continue
		}

		println("Enter the second number")
		var num2String string
		_, err = fmt.Scan(&num2String)
		if err != nil {
			errorMessage = fmt.Sprintf("failed to scan the second number: %s", err)
			continue
		}

		// Парсинг з рядка в число другого числа
		num2, err := strconv.ParseInt(num2String, 10, 64)
		if err != nil {
			errorMessage = fmt.Sprint("you need to enter a number for the second value")
			continue
		}

		// Відправка запиту на сервер для отримання суми двох чисел
		sumResponseBytes, err := c.Net.SendSum(num1, num2)
		if err != nil {
			errorMessage = fmt.Sprintf("failed to send sum message: %s", err)
			continue
		}

		// Обробка відповіді суми
		var sumResponse messages.SumResponse
		failureMessage, err := controller.HandleResponse(sumResponseBytes, &sumResponse)
		if err != nil {
			errorMessage = fmt.Sprintf("failed to handle sum response: %s", err)
			continue
		}

		if failureMessage.Message != "" {
			errorMessage = fmt.Sprintf("server returned an error: %s", failureMessage.Message)
			continue
		}

		fmt.Printf("Result: %d + %d = %d\n", num1, num2, sumResponse.Result)
		println("To continue type in anything or Ctr+C to exit")

		var endMessage string
		fmt.Scan(&endMessage)
	}

}
