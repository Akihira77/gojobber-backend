package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	base_url string
}

func NewAuthHandler(base_url string) *AuthHandler {
	return &AuthHandler{
		base_url: base_url,
	}
}

func (ah *AuthHandler) HealthCheck(c *fiber.Ctx) error {
	route := ah.base_url + "/health-check"
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("AUTH - health check error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"response": string(body),
	})

}

func (ah *AuthHandler) SignIn(c *fiber.Ctx) error {
	route := ah.base_url + fmt.Sprintf("/api/v1/auths/signin")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("AUTH - sign in error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	type Response struct {
		Token string `json:"token,omitempty"`
	}

	var res Response
	err := json.Unmarshal(body, &res)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened.")
	}

	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   res.Token,
		Expires: time.Now().Add(1 * time.Hour),
	})

	return c.Status(statusCode).Send(body)
}

func (ah *AuthHandler) SignUp(c *fiber.Ctx) error {
	route := ah.base_url + fmt.Sprintf("/api/v1/auths/signup")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("AUTH - sign up error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	type Response struct {
		Token string `json:"token,omitempty"`
	}

	var res Response
	err := json.Unmarshal(body, &res)
	if err == nil || len(body) > 0 {
		c.Cookie(&fiber.Cookie{
			Name:    "token",
			Value:   res.Token,
			Expires: time.Now().Add(1 * time.Hour),
		})

		return c.Status(statusCode).Send(body)
	}

	return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened.")
}

func (ah *AuthHandler) GetUserInfo(c *fiber.Ctx) error {
	route := ah.base_url + fmt.Sprintf("/api/v1/auths/user-info")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("AUTH - get user info error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (ah *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	route := ah.base_url + fmt.Sprintf("/api/v1/auths/refresh-token/%s", c.Params("username"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("AUTH - refresh token error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (ah *AuthHandler) SendVerifyEmailURL(c *fiber.Ctx) error {
	route := ah.base_url + fmt.Sprintf("/api/v1/auths/send-verification-email")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("AUTH - send verification email error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (ah *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	route := ah.base_url + fmt.Sprintf("/api/v1/auths/verify-email/%s", c.Params("token"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("AUTH - verifying email error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (ah *AuthHandler) SendForgotPasswordURL(c *fiber.Ctx) error {
	route := ah.base_url + fmt.Sprintf("/api/v1/auths/forgot-password/%s", c.Params("email"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("AUTH - send forgot password url error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (ah *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	route := ah.base_url + fmt.Sprintf("/api/v1/auths/reset-password/%s", c.Params("token"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("AUTH - reset password error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (ah *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	route := ah.base_url + fmt.Sprintf("/api/v1/auths/change-password")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("AUTH - change password error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}
