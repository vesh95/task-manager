package api

import (
	"encoding/json"
	"net/http"

	"github.com/vesh95/task-manager/pkg/authorization"
)

var TodoPassword string

type SigninRequest struct {
	Password string `json:"password"`
}

type SigninResponse struct {
	Token string `json:"token"`
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var req SigninRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJson(w, ErrorResponse{err.Error()}, http.StatusBadRequest)
		return
	}

	if req.Password == TodoPassword {
		token, err := authorization.CreateToken(req.Password)
		if err != nil {
			writeJson(w, ErrorResponse{err.Error()}, http.StatusInternalServerError)
			return
		}
		writeJson(w, SigninResponse{Token: token}, http.StatusOK)
		return
	}

	writeJson(w, ErrorResponse{"Неверный пароль"}, http.StatusUnauthorized)
}
