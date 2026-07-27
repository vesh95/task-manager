package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/vesh95/task-manager/pkg/server"
)

var HTTP_ADDRESS = ""
var HTTP_PORT = "7540"
var webDir = "web"

func main() {
	HTTP_ADDRESS = envOrDefaul("TODO_ADDR", "")
	HTTP_PORT = envOrDefaul("TODO_PORT", "7540")
	logger := log.Default()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	s := server.NewServer(HTTP_ADDRESS, HTTP_PORT, webDir, logger)
	go s.Run()

	<-sig
	s.Shutdown()
}

func envOrDefaul(envName, defaultValue string) string {
	v := os.Getenv(envName)
	if v == "" {
		return defaultValue
	}

	return v
}
