package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Akihira77/gojobber/services/6-chat/handler"
	"github.com/Akihira77/gojobber/services/6-chat/service"
	"github.com/Akihira77/gojobber/services/6-chat/types"
	"github.com/Akihira77/gojobber/services/6-chat/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const (
	BASE_PATH = "api/v1/chat"
)

func MainRouter(db *gorm.DB, cld *util.Cloudinary, app *fiber.App, grpcServices *handler.GRPCClients) {
	app.Get("health-check", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).SendString("Chat Service is healthy and OK.")
	})

	api := app.Group(BASE_PATH)
	api.Use(verifyGatewayReq)
	api.Use(authOnly)

	cs := service.NewChatService(db)
	ch := handler.NewChatHandler(cld, cs, grpcServices)

	api.Get("/my-conversations", ch.GetAllMyConversations)
	api.Get("/id/:conversationId", ch.GetMessagesInsideConversation)
	api.Post("", ch.InsertMessage)
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
