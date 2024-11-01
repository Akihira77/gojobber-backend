package util

import (
	"fmt"
	"log"
	"time"

	"github.com/Akihira77/gojobber/services/1-gateway/types"
	"github.com/golang-jwt/jwt/v5"
)

var ServiceID = "API Gateway"
var JWT_EXPIRATION = 1 * time.Hour
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256

func GenerateJWT(secret string) (string, error) {
	claims := jwt.MapClaims{
		"iss": ServiceID,
	}
	token := jwt.NewWithClaims(JWT_SIGNING_METHOD, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println("signinjwt", err)
		return "", fmt.Errorf("error signing jwt")
	}

	return signedToken, nil
}

func VerifyingJWT(secret string, tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signature")
		}

		return []byte(secret), nil
	})

	if err != nil {
		log.Println("verifyingjwt", err)
		return nil, fmt.Errorf("error parsing token")
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return token, nil
}
