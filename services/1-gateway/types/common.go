package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID       string `json:"userId"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	VerifiedUser bool   `json:"verifiedUser"`
}

type GoogleUserData struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

type SignUpParams struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Country        string `json:"country"`
	ProfilePicture string `json:"profilePicture"`
}

type SignInParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
