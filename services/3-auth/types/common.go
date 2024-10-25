package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID   string `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type RabbitMQResponse[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

type ErrorResult struct {
	Field string `json:"field"`
	Error string `json:"error"`
}
