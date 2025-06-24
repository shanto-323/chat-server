package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"chat_app/backend"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	manager := backend.NewManager()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		manager.ServerWS(w, r)
	})

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("server starting on port 8080....")
		if err := http.ListenAndServe(":8080", router); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-stopChan
	log.Println("Closing the server....")
}
