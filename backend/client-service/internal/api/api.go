package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shanto-323/Chat-Server-1/client-service/internal/api/routes"
	"github.com/shanto-323/Chat-Server-1/client-service/internal/cache"
	"github.com/shanto-323/Chat-Server-1/client-service/internal/database"
	"github.com/shanto-323/Chat-Server-1/client-service/util"
)

type Api interface {
	Start() error
}

type gorillaApi struct {
	ipAddr  string
	chche   *cache.RedisService
	service *database.UserService
}

func NewApi(port string, r *cache.RedisService, s *database.UserService) Api {
	return &gorillaApi{
		ipAddr:  fmt.Sprintf(":%s", port),
		chche:   r,
		service: s,
	}
}

func (a *gorillaApi) Start() error {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router = router.PathPrefix("/api/v1/client.service").Subrouter()

	userRouter := router.PathPrefix("/user").Subrouter()
	a.handleUserRoutes(userRouter)

	cacheRouter := router.PathPrefix("/cache").Subrouter()
	a.handleCacheRoutes(cacheRouter)

	return http.ListenAndServe(a.ipAddr, router)
}

func (a *gorillaApi) handleUserRoutes(r *mux.Router) {
	ops := routes.NewUserRoute(a.service)

	r.HandleFunc("/sign.in", util.HandleFunc(ops.SignInHandler)).Methods("POST")
	r.HandleFunc("/sign.up", util.HandleFunc(ops.SignUpHandler)).Methods("POST")
	r.HandleFunc("/remove/{username}", util.HandleFunc(ops.DeleteUserHandler)).Methods("DELETE")
}

func (a *gorillaApi) handleCacheRoutes(r *mux.Router) {
	ops := routes.NewCacheRoute(a.chche)

	r.HandleFunc("/client.up", util.HandleFunc(ops.AddConnectionHandler)).Methods("POST")
	r.HandleFunc("/client.close", util.HandleFunc(ops.RemoveConnectionHandler)).Methods("POST")
	r.HandleFunc("/client.get", util.HandleFunc(ops.GetConnectionHandler)).Methods("GET")
}
