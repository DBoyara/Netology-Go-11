package main

import (
	"github.com/DBoyara/Netology-Go-11/pkg/app"
	"github.com/DBoyara/Netology-Go-11/pkg/card"
	"log"
	"net"
	"net/http"
	"os"
)

const defaultPort = "9999"
const defaultHost = "127.0.0.1"

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}

	log.Printf("Server run on http://%s:%s", host, port)

	if err := execute(net.JoinHostPort(host, port)); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(addr string) (err error) {
	cardSvc := card.NewService()
	mux := http.NewServeMux()
	application := app.NewServer(cardSvc, mux)
	application.Init()

	server := &http.Server{
		Addr: addr,
		Handler: application,
	}
	return server.ListenAndServe()
}
