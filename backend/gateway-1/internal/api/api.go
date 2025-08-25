package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shanto-323/Chat-Server-1/gateway-1/internal/connection"
)

type Api interface {
	Start() error
}

type gorillaApi struct {
	ipAddr  string
	manager *connection.Manager
}

func NewApi(port string, manager *connection.Manager) Api {
	return &gorillaApi{
		ipAddr:  fmt.Sprintf(":%s", port),
		manager: manager,
	}
}

func (a *gorillaApi) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		a.manager.ServerWS(w, r)
	})
	return http.ListenAndServe(a.ipAddr, router)
}
