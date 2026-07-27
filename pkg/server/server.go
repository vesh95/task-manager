package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	logger *log.Logger
}

func NewServer(addr, port, webDir string, logger *log.Logger) *Server {
	m := http.NewServeMux()
	m.Handle("/", http.FileServer(http.Dir(webDir)))

	s := &http.Server{
		Addr:     fmt.Sprintf("%s:%s", addr, port),
		Handler:  m,
		ErrorLog: logger,
	}

	return &Server{
		server: s,
		logger: logger,
	}
}

func (s *Server) Run() {
	s.logger.Printf("Server started at %s\n", s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil {
		s.logger.Println(err)
	}
}

func (s *Server) Shutdown() {
	ctx, c := context.WithTimeout(context.Background(), 1*time.Second)
	defer c()
	s.server.Shutdown(ctx)
}
