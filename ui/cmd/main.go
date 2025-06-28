package main

import (
	"fmt"

	"chat_app/ui"
)

func main() {
	app := ui.NewApp()
	app.Run()
	fmt.Println("App running ....")
}
