package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Akihira77/gojobber/services/1-gateway/handler"
	"github.com/Akihira77/gojobber/services/1-gateway/util"
	"github.com/gofiber/fiber/v2"
)

var (
	BASE_PATH = "/api/v1/gateway"
)

func generateGatewayToken(c *fiber.Ctx) error {
	gatewaySecret := os.Getenv("GATEWAY_TOKEN")
	gatewayToken, err := util.SigningJWT(gatewaySecret)
	if err != nil {
		fmt.Printf("generating gateway token error:\n%+v", err)
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened.")
	}

	ctx := context.WithValue(c.UserContext(), "gatewayToken", gatewayToken)
	c.SetUserContext(ctx)
	return c.Next()
}

func MainRouter(app *fiber.App) {
	app.Get("/health-check", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).SendString("API Gateway Service is health and OK!")
	})

	AUTH_URL := os.Getenv("AUTH_URL")
	USER_URL := os.Getenv("USER_URL")
	GIG_URL := os.Getenv("GIG_URL")
	CHAT_URL := os.Getenv("CHAT_URL")
	ORDER_URL := os.Getenv("ORDER_URL")
	REVIEW_URL := os.Getenv("REVIEW_URL")
	api := app.Group(BASE_PATH)
	api.Use(generateGatewayToken)

	handler.WsUpgrade(api)

	authRouter(AUTH_URL, api.Group("/auths"))
	userRouter(USER_URL, api.Group("/users"))
	gigRouter(GIG_URL, api.Group("/gigs"))
	chatRouter(CHAT_URL, api.Group("/chats"))
	orderRouter(ORDER_URL, api.Group("/orders"))
	reviewRouter(REVIEW_URL, api.Group("/reviews"))

	app.All("*", func(c *fiber.Ctx) error {
		return c.Status(http.StatusNotFound).SendString("Resource is not found")
	})
}

func authRouter(base_url string, r fiber.Router) {
	ah := handler.NewAuthHandler(base_url)
	r.Get("/health-check", ah.HealthCheck)

	//INFO: AUTH ROUTER
	r.Get("/user-info", ah.GetUserInfo)
	r.Get("/refresh-token/:username", ah.RefreshToken)
	r.Post("/signin", ah.SignIn)
	r.Post("/signup", ah.SignUp)
	r.Post("/send-verification-email", ah.SendVerifyEmailURL)
	r.Patch("/verify-email/:token", ah.VerifyEmail)
	r.Patch("/forgot-password/:email", ah.SendForgotPasswordURL)
	r.Patch("/reset-password/:token", ah.ResetPassword)
	r.Patch("/change-password", ah.ChangePassword)
}

func userRouter(base_url string, r fiber.Router) {
	uh := handler.NewUserHandler(base_url)
	r.Get("/health-check", uh.HealthCheck)

	//INFO: BUYER ROUTER
	r.Get("/buyer/my-info", uh.GetMyBuyerInfo)
	r.Get("/buyer/:username", uh.FindBuyerByUsername)

	//INFO: SELLER ROUTER
	r.Get("/seller/my-info", uh.GetMySellerInfo)
	r.Get("/seller/id/:id", uh.FindSellerByID)
	r.Get("/seller/username/:username", uh.FindSellerByUsername)
	r.Get("/seller/random/:count", uh.GetRandomSellers)
	r.Post("/seller", uh.Create)
	r.Put("/seller", uh.Update)

}

func gigRouter(base_url string, r fiber.Router) {
	gh := handler.NewGigHandler(base_url)
	r.Get("/health-check", gh.HealthCheck)

	//INFO: GIG ROUTE
	r.Get("/id/:id", gh.FindGigByID)
	r.Get("/seller/:sellerId/:page/:size", gh.FindSellerActiveGigs)
	r.Get("/seller/inactive/:sellerId/:page/:size", gh.FindSellerInactiveGigs)
	r.Get("/category/:category/:page/:size", gh.FindGigsByCategory)
	r.Get("/popular/:category", gh.GetPopularGigs)
	r.Get("/similar/:gigId/:page/:size", gh.FindSimilarGigs)
	r.Get("/search/:page/:size", gh.GigQuerySearch)
	r.Post("", gh.CreateGig)
	r.Put("/:sellerId/:gigId", gh.UpdateGig)
	r.Patch("/update-status/:sellerId/:gigId", gh.ActivateGigStatus)
	r.Delete("/:sellerId/:gigId", gh.DeactivateGigStatus)
}

func chatRouter(base_url string, r fiber.Router) {
	ch := handler.NewChatHandler(base_url)
	r.Get("/health-check", ch.HealthCheck)

	//INFO: CHAT ROUTER
	// r.Get("/my-notifications", ch.GetAllMyConversations)
	r.Get("/my-conversations", ch.GetAllMyConversations)
	r.Get("/id/:conversationId", ch.GetMessagesInsideConversation)
	r.Post("", ch.InsertMessage)

}

func orderRouter(base_url string, r fiber.Router) {
	oh := handler.NewOrderHandler(base_url)
	r.Get("/health-check", oh.HealthCheck)
	r.Post("/payment-intent/create", oh.CreatePaymentIntent)
	r.Post("/payment-intent/:paymentId/confirm", oh.ConfirmPayment)
}

func reviewRouter(base_url string, r fiber.Router) {
	rh := handler.NewReviewHandler(base_url)
	r.Get("/health-check", rh.HealthCheck)
}
