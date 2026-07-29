package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/vesh95/task-manager/pkg/db"
)

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		writeJson(w, ErrorResponse{err.Error()}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, ErrorResponse{err.Error()}, http.StatusBadRequest)
		return
	}

	if task.Repeat == "" {
		db.DeleteTask(id)
		writeJson(w, struct{}{}, http.StatusOK)
		return
	}

	date, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(w, ErrorResponse{err.Error()}, http.StatusInternalServerError)
	}

	task.Date = date
	if err = db.UpdateTask(task); err != nil {
		writeJson(w, ErrorResponse{err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJson(w, struct{}{}, http.StatusOK)
}
