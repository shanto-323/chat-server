package main

import (
	"log"

	"chat_app/ui"
)

func main() {
	app := ui.NewApp()
	log.Println("ui running.....")
	app.Run()
}
