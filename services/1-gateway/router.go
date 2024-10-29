package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Akihira77/gojobber/services/1-gateway/handler"
	"github.com/Akihira77/gojobber/services/1-gateway/types"
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
		fmt.Println("token is not matched with claims:", token.Claims)
		return fiber.NewError(http.StatusUnauthorized, "sign in first")
	}

	c.SetUserContext(context.WithValue(c.UserContext(), "current_user", claims))
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

	authRouter(AUTH_URL, api.Group("/auths"))
	userRouter(USER_URL, api.Group("/users"))
	gigRouter(GIG_URL, api.Group("/gigs"))
	chatRouter(CHAT_URL, api.Group("/chats"))
	orderRouter(ORDER_URL, api.Group("/orders"))
	reviewRouter(REVIEW_URL, api.Group("/reviews"))

	handler.WsUpgrade(api.Use(authOnly))

	app.All("*", func(c *fiber.Ctx) error {
		return c.Status(http.StatusNotFound).SendString("Resource is not found")
	})
}

func authRouter(base_url string, r fiber.Router) {
	ah := handler.NewAuthHandler(base_url)
	r.Get("/health-check", ah.HealthCheck)

	r.Get("/google/:action", ah.AuthWithGoogle)

	r.Get("/signup/google-callback", ah.SignUpWithGoogle)
	r.Post("/signup", ah.SignUp).Name("signup")

	r.Get("/signin/google-callback", ah.SignInWithGoogle)
	r.Post("/signin", ah.SignIn).Name("signin")

	r.Patch("/forgot-password/:email", ah.SendForgotPasswordURL)
	r.Patch("/reset-password/:token", ah.ResetPassword)

	r.Use(authOnly)
	r.Get("/user-info", ah.GetUserInfo)
	r.Get("/refresh-token/:username", ah.RefreshToken)
	r.Post("/send-verification-email", ah.SendVerifyEmailURL)
	r.Patch("/verify-email/:token", ah.VerifyEmail)
	r.Patch("/change-password", ah.ChangePassword)
}

func userRouter(base_url string, r fiber.Router) {
	uh := handler.NewUserHandler(base_url)
	r.Get("/health-check", uh.HealthCheck)

	r.Use(authOnly)
	r.Get("/buyers/my-info", uh.GetMyBuyerInfo)
	r.Get("/buyers/:username", uh.FindBuyerByUsername)

	r.Get("/sellers/my-info", uh.GetMySellerInfo)
	r.Get("/sellers/id/:id", uh.FindSellerByID)
	r.Get("/sellers/username/:username", uh.FindSellerByUsername)
	r.Get("/sellers/random/:count", uh.GetRandomSellers)
	r.Post("/sellers", uh.Create)
	r.Put("/sellers", uh.Update)

}

func gigRouter(base_url string, r fiber.Router) {
	gh := handler.NewGigHandler(base_url)
	r.Get("/health-check", gh.HealthCheck)

	r.Get("/popular", gh.GetPopularGigs).Name("home")
	r.Get("/id/:id", gh.FindGigByID)
	r.Get("/category/:category/:page/:size", gh.FindGigsByCategory)
	r.Get("/similar/:gigId/:page/:size", gh.FindSimilarGigs)
	r.Get("/search/:page/:size", gh.GigQuerySearch)

	r.Use(authOnly)
	r.Get("/sellers/active/:page/:size", gh.FindSellerActiveGigs)
	r.Get("/sellers/inactive/:page/:size", gh.FindSellerInactiveGigs)
	r.Post("", gh.CreateGig)
	r.Put("/:sellerId/:gigId", gh.UpdateGig)
	r.Patch("/update-status/:sellerId/:gigId", gh.ActivateGigStatus)
	r.Delete("/:sellerId/:gigId", gh.DeactivateGigStatus)
}

func chatRouter(base_url string, r fiber.Router) {
	ch := handler.NewChatHandler(base_url)
	r.Get("/health-check", ch.HealthCheck)

	// r.Get("/my-notifications", ch.GetAllMyConversations)
	r.Use(authOnly)
	r.Get("/my-conversations", ch.GetAllMyConversations)
	r.Get("/id/:conversationId", ch.GetMessagesInsideConversation)
	r.Post("", ch.InsertMessage)
	r.Patch("/offer/:messageId/cancel", ch.SellerCancelOffer)
}

func orderRouter(base_url string, r fiber.Router) {
	oh := handler.NewOrderHandler(base_url)
	r.Get("/health-check", oh.HealthCheck)

	r.Use(authOnly)
	r.Get("/:id", oh.FindOrderByID)
	r.Get("/buyer/my-orders", oh.FindOrdersByBuyerID)
	r.Get("/seller/my-orders", oh.FindOrdersBySellerID)
	r.Get("/buyer/my-orders-notifications", oh.FindMyOrdersNotifications)
	r.Post("/stripe/webhook", oh.HandleStripeWebhook)

	r.Post("/payment-intents/create", oh.CreatePaymentIntent)
	//NOTE: JUST FOR TESTING
	r.Post("/payment-intents/:paymentId/confirm", oh.ConfirmPayment)
	r.Post("/stripe/tos-acceptance", oh.StripeTOSAcceptance)

	r.Post("/deadline/extension/:orderId/request", oh.RequestDeadlineExtension)
	r.Post("/deadline/extension/:orderId/response", oh.BuyerDeadlineExtensionResponse)
	r.Post("/:orderId/complete", oh.OrderComplete)
	r.Post("/:orderId/cancel", oh.CancelOrder)
	r.Post("/:orderId/refund", oh.OrderRefund)
	r.Post("/:orderId/acknowledge", oh.AcknowledgeOrder)
	r.Post("/deliver/:orderId", oh.DeliveringOrder)
	r.Post("/deliver/:orderId/response", oh.BuyerResponseForDeliveredOrder)
}

func reviewRouter(base_url string, r fiber.Router) {
	rh := handler.NewReviewHandler(base_url)
	r.Get("/health-check", rh.HealthCheck)

	r.Use(authOnly)
	r.Get("/seller/:sellerId", rh.FindSellerReviews)
	r.Post("", rh.Add)
	r.Patch("/:reviewId", rh.Update)
	r.Delete("/:reviewId", rh.Remove)
}
