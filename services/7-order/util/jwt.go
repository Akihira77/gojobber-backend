package util

import (
	"fmt"
	"log"

	"github.com/Akihira77/gojobber/services/7-order/types"
	"github.com/golang-jwt/jwt/v5"
)

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
