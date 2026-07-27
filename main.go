package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/vesh95/task-manager/pkg/server"
	"github.com/vesh95/task-manager/pkg/server/db"
)

var (
	HTTP_ADDRESS,
	HTTP_PORT,
	TODO_DBFILE,
	webDir string
)

func main() {
	HTTP_ADDRESS = envOrDefaul("TODO_ADDR", "")
	HTTP_PORT = envOrDefaul("TODO_PORT", "7540")
	TODO_DBFILE = envOrDefaul("TODO_DBFILE", "scheduler.db")
	webDir = envOrDefaul("WEB_DIR", "web")
	logger := log.Default()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	err := db.Init(TODO_DBFILE)
	if err != nil {
		log.Fatalf("error while connecting database: %s", err)
	}

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
