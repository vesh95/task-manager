package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/vesh95/task-manager/pkg/db"
)

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req db.Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		wrireJson(w, ErrorResponse{err.Error()}, http.StatusBadRequest)
		return
	}

	if err := validateUpdateTaskRequest(req); err != nil {
		wrireJson(w, ErrorResponse{err.Error()}, http.StatusUnprocessableEntity)
		return
	}

	date, err := recalculateReplaceDate(time.Now(), req.Date, req.Repeat)
	if err != nil {
		wrireJson(w, ErrorResponse{err.Error()}, http.StatusUnprocessableEntity)
		return
	}

	req.Date = date
	err = db.UpdateTask(req)

	if err != nil {
		wrireJson(w, ErrorResponse{err.Error()}, http.StatusInternalServerError)
		return
	}

	wrireJson(w, struct{}{}, http.StatusOK)
}

func validateUpdateTaskRequest(req db.Task) error {
	if req.Title == "" {
		return errors.New("Не указан заголовок задачи")
	}

	if len(req.Title) > 512 {
		return errors.New("Длина заголовка более 512 символов")
	}

	if req.Date != "" {
		_, err := time.Parse(DateFormat, req.Date)
		if err != nil {
			return fmt.Errorf("Неверный формат даты")
		}
	}

	return nil
}
