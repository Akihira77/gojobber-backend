package types

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID       string `json:"userId"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	VerifiedUser bool   `json:"verifiedUser"`
}

const (
	NOTIFICATION_SERVICE = "NOTIFICATION_SERVICE"
	AUTH_SERVICE         = "AUTH_SERVICE"
	USER_SERVICE         = "USER_SERVICE"
	GIG_SERVICE          = "GIG_SERVICE"
	CHAT_SERVICE         = "CHAT_SERVICE"
	ORDER_SERVICE        = "ORDER_SERVICE"
	REVIEW_SERVICE       = "REVIEW_SERVICE"
)
