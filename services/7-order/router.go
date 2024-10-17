package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Akihira77/gojobber/services/7-order/handler"
	"github.com/Akihira77/gojobber/services/7-order/service"
	"github.com/Akihira77/gojobber/services/7-order/types"
	"github.com/Akihira77/gojobber/services/7-order/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const (
	BASE_PATH = "api/v1/orders"
)

func MainRouter(db *gorm.DB, cld *util.Cloudinary, app *fiber.App) {
	app.Get("health-check", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).SendString("Order Service is healthy and OK.")
	})

	api := app.Group(BASE_PATH)
	api.Use(verifyGatewayReq)
	api.Use(authOnly)

	os := service.NewOrderService(db)
	oh := handler.NewOrderHttpHandler(os)

	api.Get("/:id", oh.FindOrderByID)
	api.Get("/buyer/my-orders", oh.FindMyOrdersAsBuyer)
	api.Get("/seller/my-orders", oh.FindMyOrdersAsSeller)
	api.Post("/payment-intents/create", oh.CreatePaymentIntent)
	api.Post("/stripe-webhook", oh.HandleStripeWebhook)
	api.Post("/:orderId/refund", oh.BuyerRefundingOrder)
	api.Post("/:orderId/complete", oh.BuyerMarkOrderAsComplete)
	api.Post("/:orderId/cancel", oh.SellerCancellingOrder)
	api.Post("/deadline/extension/:orderId/request", oh.RequestDeadlineExtension)
	api.Post("/deadline/extension/:orderId/response", oh.BuyerDeadlineExtensionResponse)
	api.Post("/deliver/:orderId", oh.SellerDeliverOrder)
	api.Post("/deliver/:orderId/response", oh.BuyerResponseForDeliveredOrder)
	api.Get("/buyer/my-orders-notifications", oh.FindMyOrdersNotifications)
	// api.Patch("/buyer/my-orders-notification/reads", oh.MarkReadsMyOrderNotifications)

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
		if authHeader == "" || len(strings.Split(authHeader, " ")) <= 1 {
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
	log.Println(claims)
	if !ok {
		log.Println("token is not matched with claims")
		return fiber.NewError(http.StatusUnauthorized, "sign in first")
	}

	c.SetUserContext(context.WithValue(c.UserContext(), "current_user", claims))
	return c.Next()
}
