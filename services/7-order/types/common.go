package types

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID   string `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
