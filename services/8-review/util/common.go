package util

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Akihira77/gojobber/services/3-auth/types"
	"github.com/go-playground/validator/v10"
)

func RandomStr(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)

}

func CustomValidationErrors(err error) []types.ErrorResult {
	var errs []types.ErrorResult
	for _, v := range err.(validator.ValidationErrors) {
		var e error
		switch v.Tag() {
		case "required":
			e = fmt.Errorf("Field '%s' cannot be empty", v.Field())
		case "email":
			e = fmt.Errorf("Field '%s' must be a valid email address", v.Field())
		case "eth_addr":
			e = fmt.Errorf("Field '%s' must  be a valid Ethereum address", v.Field())
		case "len":
			e = fmt.Errorf("Field '%s' must be exactly %v characters long", v.Field(), v.Param())
		default:
			e = fmt.Errorf("Field '%s': '%v' must satisfy '%s' '%v' criteria", v.Field(), v.Value(), v.Tag(), v.Param())
		}

		errs = append(errs, types.ErrorResult{
			Field: v.Field(),
			Error: e.Error(),
		})
	}

	return errs
}
