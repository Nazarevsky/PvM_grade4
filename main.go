package main

import (
	"fmt"
	"lab4/bot"
	"lab4/db"
)

func main() {
	controller, err := db.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = controller.MigrateUp()
	if err != nil {
		fmt.Println(err)
		return
	}

	tgBot, err := bot.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	tgBot.Start()
}
