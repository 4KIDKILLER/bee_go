package server

import (
	"goserver/pkg/controller"
	"goserver/pkg/infrastructure/config"
	"net/http"
)

func NewHttpServer(config *config.Config) *http.Server {
	serverConf := config.Server
	mux := http.NewServeMux()
	controller.UserController(mux)

	return &http.Server{
		Handler: mux,
		Addr:    serverConf.Port,
	}
}
