package api

import (
	"net/http"

	"github.com/vesh95/task-manager/pkg/db"
)

type TasksResponse struct {
	Tasks []db.Task `json:"tasks"`
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	tasks := make([]db.Task, 0)
	tasks, err := db.Tasks(50, search)
	if err != nil {
		wrireJson(w, ErrorResponse{err.Error()}, http.StatusInternalServerError)
		return
	}

	resp := TasksResponse{tasks}
	wrireJson(w, resp, http.StatusOK)
}
