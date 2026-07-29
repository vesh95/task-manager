package api

import (
	"net/http"

	"github.com/vesh95/task-manager/pkg/authorization"
)

func Auth(password string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if password == "" {
			next(w, r)
			return
		}

		var jwt string
		cookie, err := r.Cookie("token")
		if err == nil {
			jwt = cookie.Value
		}

		valid := authorization.ValidateToken(jwt, password)

		if !valid {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}
