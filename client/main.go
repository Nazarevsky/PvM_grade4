package main

import (
	"client/cli"
	"client/net/controller"
	"fmt"
)

func main() {
	messageController, err := controller.New("tcp", ":8000")
	if err != nil {
		fmt.Printf("Failed to create message controller: %s\n", err)
		return
	}

	cliController := cli.New(messageController)
	cliController.StartCli()
}
