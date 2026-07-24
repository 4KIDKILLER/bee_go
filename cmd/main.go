package main

import (
	"goserver/pkg/infrastructure/config"
	"goserver/pkg/server"
	"log"
	"net"
	"net/http"
)

var systemConf *config.Config

func init() {
	systemConf = config.NewConfig()
}

func main() {
	httpServer := server.NewHttpServer(systemConf)

	listener, err := net.Listen("tcp", httpServer.Addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("服务器启动成功:http://localhost%s", httpServer.Addr)

	if err := httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
