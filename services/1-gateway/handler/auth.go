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

func (ah *AuthHandler) FindGigByID(c *fiber.Ctx) error {
	route := ah.base_url + fmt.Sprintf("/api/v1/auth/search/gig/%s", c.Params("id"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("AUTH - find gig by id error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (ah *AuthHandler) FindGigByQuery(c *fiber.Ctx) error {
	query := fmt.Sprintf("query=%v&max=%v&delivery_time=%v", c.Query("query"), c.QueryInt("max"), c.QueryInt("delivery_time"))
	route := ah.base_url + fmt.Sprintf("/api/v1/auth/search/gig/%s/%s?%s", c.Params("from"), c.Params("size"), query)
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("AUTH - find gig by query error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (ah *AuthHandler) SignIn(c *fiber.Ctx) error {
	route := ah.base_url + fmt.Sprintf("/api/v1/auth/signin")
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

func (ah *AuthHandler) SignUp(c *fiber.Ctx) error {
	route := ah.base_url + fmt.Sprintf("/api/v1/auth/signup")
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
	route := ah.base_url + fmt.Sprintf("/api/v1/auth/user-info")
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
	route := ah.base_url + fmt.Sprintf("/api/v1/auth/refresh-token/%s", c.Params("username"))
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
	route := ah.base_url + fmt.Sprintf("/api/v1/auth/send-verification-email")
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
	route := ah.base_url + fmt.Sprintf("/api/v1/auth/verify-email/%s", c.Params("token"))
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
	route := ah.base_url + fmt.Sprintf("/api/v1/auth/forgot-password/%s", c.Params("email"))
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
	route := ah.base_url + fmt.Sprintf("/api/v1/auth/reset-password/%s", c.Params("token"))
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
	route := ah.base_url + fmt.Sprintf("/api/v1/auth/change-password")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("AUTH - change password error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}
