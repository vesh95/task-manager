package api

import (
	"net/http"
	"strconv"

	"github.com/vesh95/task-manager/pkg/db"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		writeJson(w, ErrorResponse{err.Error()}, http.StatusBadRequest)
		return
	}

	if err = db.DeleteTask(id); err != nil {
		writeJson(w, ErrorResponse{err.Error()}, http.StatusUnprocessableEntity)
		return
	}

	writeJson(w, struct{}{}, http.StatusOK)
}
