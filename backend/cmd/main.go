package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chat_app/backend/logger"
	"chat_app/backend/pkg/storage/database"
	"chat_app/backend/pkg/storage/redis"
	websocket "chat_app/backend/pkg/web-socket"

	"github.com/gorilla/mux"
	"github.com/tinrab/retry"
)

func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Println(err)
		return
	}
	defer logger.Sync()

	redisUrl := "redis://:123456@localhost:6379/0"
	redisCLient, err := redis.NewRedisClient(redisUrl)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Info("redis running")

	scyllaIp := "scylla"
	var scyllaDb database.Repository
	retry.ForeverSleep(
		2*time.Second,
		func(_ int) error {
			scyllaDb, err = database.NewScyllaDb(scyllaIp)
			if err != nil {
				logger.Error(err.Error())
				return err
			}
			logger.Info("scylla running")
			return nil
		},
	)

	ctx, cancel := context.WithCancel(context.Background())

	router := mux.NewRouter()
	manager := websocket.NewManager(ctx, logger, redisCLient, scyllaDb)
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		manager.ServerWS(w, r)
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Server starting on port 8080...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error(err.Error())
		}
	}()

	<-stopChan
	logger.Info("Closing the server....")
	cancel()

	shutdownCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()
	manager.Shutdown(shutdownCtx)
}
