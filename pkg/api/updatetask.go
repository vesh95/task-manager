package api

import (
	"encoding/json"
	"net/http"

	"github.com/vesh95/task-manager/pkg/db"
)

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	task := db.Task{}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		wrireJson(w, ErrorResponse{err.Error()}, http.StatusBadRequest)
		return
	}
	if err := db.UpdateTask(task); err != nil {
		wrireJson(w, ErrorResponse{err.Error()}, http.StatusInternalServerError)
		return
	}

	wrireJson(w, map[string]string{}, http.StatusOK)
}
