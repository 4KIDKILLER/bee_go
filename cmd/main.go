package main

import (
	"flag"
	"fmt"
	"goserver/pkg/infrastructure/config"
	"goserver/pkg/infrastructure/server"
	"log"
	"net"
	"net/http"
	"os"
)

var systemConf *config.Config

func init() {

	host := flag.String("env", "mac", "服务器地址")

	// 自定义 Usage 函数
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `使用方法: myapp [选项]
选项:
`)
		flag.PrintDefaults()
	}

	flag.Parse()
	fmt.Println("kkkkk")
	if *host != "mac" && *host != "windows" {
		panic("无效的-env参数")
	}

	systemConf = config.NewConfig(*host)

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
