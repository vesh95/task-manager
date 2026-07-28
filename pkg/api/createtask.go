package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/vesh95/task-manager/pkg/db"
)

type CreateTaskRequest struct {
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Date    string `json:"date,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type CreateTaskResponse struct {
	ID int64 `json:"id"`
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		wrireJson(w, ErrorResponse{err.Error()}, http.StatusBadRequest)
		return
	}

	validationErr := validateCreateRequest(req)
	if validationErr != nil {
		wrireJson(w, ErrorResponse{validationErr.Error()}, http.StatusUnprocessableEntity)
		return
	}

	now := time.Now()
	task := db.Task{
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	if req.Date == "" {
		task.Date = now.Format(DateFormat)
	} else {
		date, err := time.Parse(DateFormat, req.Date)
		if err != nil {
			wrireJson(w, ErrorResponse{"Неверный формат даты"}, http.StatusBadRequest)
			return
		}

		if date.Before(now) && req.Repeat != "" {
			nd, err := NextDate(now, req.Date, req.Repeat)
			if err != nil {
				wrireJson(w, ErrorResponse{"Неверное правило повторения"}, http.StatusUnprocessableEntity)
				return
			}
			task.Date = nd
		} else {
			task.Date = now.Format(DateFormat)
		}
	}

	fmt.Println(task)
	id, err := db.AddTask(task)

	if err != nil {
		wrireJson(w, ErrorResponse{fmt.Sprintf("Ошибка при создании задачи: %s", err)}, http.StatusInternalServerError)
		return
	}

	wrireJson(w, CreateTaskResponse{id}, http.StatusCreated)
}

func validateCreateRequest(req CreateTaskRequest) error {
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
