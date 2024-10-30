package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Akihira77/gojobber/services/1-gateway/config"
	"github.com/Akihira77/gojobber/services/1-gateway/types"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
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

func (ah *AuthHandler) googleCallback(code string, cfg oauth2.Config) (types.GoogleUserData, error) {
	token, err := cfg.Exchange(context.Background(), code)
	if err != nil {
		return types.GoogleUserData{}, err
	}

	agent := fiber.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	_, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return types.GoogleUserData{}, err
	}

	var userData types.GoogleUserData
	err = json.Unmarshal(body, &userData)
	if err != nil {
		log.Printf("Google Auth Error:\n%+v", err)
		return types.GoogleUserData{}, err
	}

	return userData, nil
}

func (ah *AuthHandler) AuthWithGoogle(c *fiber.Ctx) error {
	action := c.Params("action")
	log.Println(action)

	var url string
	switch action {
	case "signin":
		url = config.GoogleOAuthConfig["signin"].AuthCodeURL("randomstate")
		break
	case "signup":
		url = config.GoogleOAuthConfig["signup"].AuthCodeURL("randomstate")
		break
	default:
		return c.Status(http.StatusBadRequest).SendString("Invalid Auth Action")
	}

	return c.Status(fiber.StatusSeeOther).Redirect(url)
}

func (ah *AuthHandler) SignInWithGoogle(c *fiber.Ctx) error {
	cfg, err := config.GetGoogleOAuthConfig("signin")
	if err != nil {
		log.Printf("Google Signup Error:\n%+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	userData, err := ah.googleCallback(c.Query("code"), *cfg)
	if err != nil {
		log.Printf("Google Signin Error:\n%+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	body, err := json.Marshal(types.SignInParams{
		Username: userData.Email,
	})
	if err != nil {
		log.Printf("Google Signin Error:\n%+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	c.Method(http.MethodPost)
	c.Request().SetBody(body)
	c.SetUserContext(context.WithValue(c.UserContext(), "via", "google"))
	return ah.SignIn(c)
}

func (ah *AuthHandler) SignIn(c *fiber.Ctx) error {
	via := c.UserContext().Value("via").(string)
	var route string
	switch via {
	case "google":
		route = ah.base_url + fmt.Sprintf("/api/v1/auths/signin?via=google")
		break
	default:
		route = ah.base_url + fmt.Sprintf("/api/v1/auths/signin")
		break
	}

	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("AUTH - sign in error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	if statusCode >= 400 {
		fmt.Println("AUTH - sign in error", errs)
		return c.Status(statusCode).Send(body)
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

	// return c.Status(statusCode).Send(body)
	return c.RedirectToRoute("home", fiber.Map{}, statusCode)
}

func (ah *AuthHandler) SignUpWithGoogle(c *fiber.Ctx) error {
	cfg, err := config.GetGoogleOAuthConfig("signup")
	if err != nil {
		log.Printf("Google Signup Error:\n%+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	userData, err := ah.googleCallback(c.Query("code"), *cfg)
	if err != nil {
		log.Printf("Google Signup Error:\n%+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	body, err := json.Marshal(types.SignUpParams{
		Username:       "",
		Email:          userData.Email,
		Password:       "",
		Country:        "",
		ProfilePicture: userData.Picture,
	})
	if err != nil {
		log.Printf("Google Signup Error:\n%+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	c.Response().SetBody(body)
	return c.Redirect(fmt.Sprintf("%s/signup", os.Getenv("CLIENT_URL")))
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

	if statusCode >= 400 {
		fmt.Println("AUTH - sign up error", errs)
		return c.Status(statusCode).Send(body)
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

	// return c.Status(statusCode).Send(body)
	return c.RedirectToRoute("home", fiber.Map{}, statusCode)
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
