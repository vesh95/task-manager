package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vesh95/task-manager/pkg/api"
)

type Server struct {
	server *http.Server
	logger *log.Logger
}

func NewServer(addr, port, webDir string, logger *log.Logger) *Server {
	m := http.NewServeMux()
	m.Handle("/", http.FileServer(http.Dir(webDir)))
	m.HandleFunc("/api/nextdate", api.NextDateHandler)
	m.HandleFunc("/api/task", api.Auth(api.TodoPassword, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			api.CreateTaskHandler(w, r)
		case http.MethodGet:
			api.GetTaskHandler(w, r)
		case http.MethodPut:
			api.UpdateTaskHandler(w, r)
		case http.MethodDelete:
			api.DeleteTaskHandler(w, r)
		}
	}))
	m.HandleFunc("/api/tasks", api.Auth(api.TodoPassword, api.GetTasksHandler))
	m.HandleFunc("/api/task/done", api.Auth(api.TodoPassword, api.DoneTaskHandler))
	m.HandleFunc("/api/signin", api.SigninHandler)

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
