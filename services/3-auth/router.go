package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Akihira77/gojobber/services/3-auth/handler"
	"github.com/Akihira77/gojobber/services/3-auth/service"
	"github.com/Akihira77/gojobber/services/3-auth/types"
	"github.com/Akihira77/gojobber/services/3-auth/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const (
	BASE_PATH = "/api/v1/auths"
)

func MainRouter(app *fiber.App, db *gorm.DB, cld *util.Cloudinary, ccs *handler.GRPCClients) {
	app.Get("/health-check", func(c *fiber.Ctx) error {
		return c.SendString("Auth Service is healthy and OK.")
	})

	api := app.Group(BASE_PATH)
	api.Use(verifyGatewayReq)

	as := service.NewAuthService(db)
	ah := handler.NewAuthHttpHandler(as, cld, ccs)

	api.Get("/user-info", authOnly, ah.GetUserInfo)
	api.Post("/signin", ah.SignIn)
	api.Post("/signup", ah.SignUp)
	api.Patch("/forgot-password/:email", ah.SendForgotPasswordURL)
	api.Patch("/reset-password/:token", ah.ResetPassword)

	api.Use(authOnly)

	api.Get("/refresh-token/:username", ah.RefreshToken)
	api.Post("/send-verification-email", ah.SendVerifyEmailURL)
	api.Patch("/verify-email/:token", ah.VerifyEmail)
	api.Patch("/change-password", ah.ChangePassword)
}

func verifyGatewayReq(c *fiber.Ctx) error {
	gatewayToken := c.Get("gatewayToken", "")

	if gatewayToken == "" {
		return fiber.NewError(http.StatusForbidden, "request is not from Gateway")
	}

	GATEWAY_TOKEN := os.Getenv("GATEWAY_TOKEN")

	token, err := jwt.Parse(gatewayToken, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signing method invalid")
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Signing method invalid")
		}

		return []byte(GATEWAY_TOKEN), nil
	})

	if err != nil {
		fmt.Printf("verifyGatewayReq error:\n%+v", err)
		return fiber.NewError(http.StatusForbidden, "invalid gateway token")
	}

	c.Set("gatewayToken", token.Raw)
	return c.Next()
}

func authOnly(c *fiber.Ctx) error {
	tokenStr := c.Cookies("token")
	if tokenStr == "" {
		authHeader := c.Get("Authorization")
		if authHeader == "" || len(strings.Split(authHeader, " ")) == 0 {
			return fiber.NewError(http.StatusUnauthorized, "sign in first")
		}
		tokenStr = strings.Split(authHeader, " ")[1]
	}
	token, err := util.VerifyingJWT(os.Getenv("JWT_SECRET"), tokenStr)
	if err != nil {
		fmt.Printf("authOnly error:\n%+v", err)
		return fiber.NewError(http.StatusUnauthorized, "sign in first")
	}

	claims, ok := token.Claims.(*types.JWTClaims)
	if !ok {
		log.Println("token is not matched with claims: claims", token.Claims)
		return fiber.NewError(http.StatusUnauthorized, "sign in first")
	}

	c.SetUserContext(context.WithValue(c.UserContext(), "current_user", claims))
	return c.Next()
}

func AdminOnly(h func(*fiber.Ctx) error) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		cu, ok := ctx.Value("current_user").(*types.JWTClaims)
		if !ok {
			return fiber.NewError(http.StatusUnauthorized, "sign in first")
		}

		// if strings.ToLower(cu.Role) != "admin" {
		// 	return fiber.ErrForbidden
		// }
		log.Println(cu)
		return c.Next()
	}
}
