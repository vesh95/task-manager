package api

import (
	"net/http"
	"strconv"

	"github.com/vesh95/task-manager/pkg/db"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		wrireJson(w, ErrorResponse{err.Error()}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		wrireJson(w, ErrorResponse{err.Error()}, http.StatusBadRequest)
		return
	}

	wrireJson(w, task, http.StatusOK)
}
