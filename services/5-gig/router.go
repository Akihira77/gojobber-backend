package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Akihira77/gojobber/services/5-gig/handler"
	svc "github.com/Akihira77/gojobber/services/5-gig/service"
	"github.com/Akihira77/gojobber/services/5-gig/types"
	"github.com/Akihira77/gojobber/services/5-gig/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const BASE_PATH = "/api/v1/gig"

func MainRouter(db *gorm.DB, cld *util.Cloudinary, app *fiber.App, ccs *handler.GRPCClients) {
	app.Get("health-check", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).SendString("Gig Service is healthy and OK.")
	})

	api := app.Group(BASE_PATH)
	api.Use(verifyGatewayReq)

	gigSvc := svc.NewGigService(db)
	gigHandler := handler.NewGigHandler(gigSvc, cld, ccs)

	api.Get("/id/:id", gigHandler.FindGigByID)
	api.Get("/search/:page/:size", gigHandler.GigQuerySearch)
	api.Get("/category/:category/:page/:size", gigHandler.FindGigByCategory)
	api.Get("/popular/:category", gigHandler.GetPopularGigs)
	api.Get("/similar/:gigId/:page/:size", gigHandler.FindSimilarGigs)

	api.Use(authOnly)

	api.Get("/seller/:sellerId/:page/:size", gigHandler.FindSellerGigs)
	api.Get("/seller/inactive/:sellerId/:page/:size", gigHandler.FindSellerInactiveGigs)
	api.Post("", gigHandler.Create)
	api.Put("/:sellerId/:gigId", gigHandler.Update)
	api.Patch("/update-status/:sellerId/:gigId", gigHandler.ActivateGigStatus)
	api.Delete("/:sellerId/:gigId", gigHandler.DeactivateGigStatus)
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
