package authorization

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var signKey = "dBBol53EGujOQgZyfEKkJ7Jgd9dU936N"

type TokenClaims struct {
	Payload string `json:"payload"`
	jwt.RegisteredClaims
}

func CreateToken(password string) (string, error) {
	passwordHash := sha256.New().Sum([]byte(password))
	c := TokenClaims{
		Payload: hex.EncodeToString(passwordHash),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 8)),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(signKey))
}

func ValidateToken(jwtString, password string) bool {
	passwordHash := sha256.New().Sum([]byte(password))

	var c TokenClaims
	token, err := jwt.ParseWithClaims(jwtString, &c, func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	})
	if err != nil {
		log.Printf("error while parsing claims: %s", err.Error())
		return false
	}
	return token.Valid && c.Payload == hex.EncodeToString(passwordHash)
}
