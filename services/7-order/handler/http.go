package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Akihira77/gojobber/services/7-order/service"
	"github.com/Akihira77/gojobber/services/7-order/types"
	"github.com/Akihira77/gojobber/services/common/genproto/notification"
	"github.com/Akihira77/gojobber/services/common/genproto/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/paymentintent"
	"github.com/stripe/stripe-go/v80/webhook"
	"gorm.io/gorm"
)

type OrderHttpHandler struct {
	orderSvc   service.OrderServiceImpl
	grpcClient *GRPCClients
	validate   *validator.Validate
}

func NewOrderHttpHandler(orderSvc service.OrderServiceImpl) *OrderHttpHandler {
	return &OrderHttpHandler{
		orderSvc: orderSvc,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (oh *OrderHttpHandler) FindOrderByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 200*time.Millisecond)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
	}

	o, err := oh.orderSvc.FindOrderByID(ctx, c.Params("id"))
	if err != nil {
		log.Printf("FindOrderById error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Order is not found")
		}

		return fiber.NewError(http.StatusInternalServerError, "Error while finding order")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"order": o,
	})
}

func (oh *OrderHttpHandler) FindOrdersByBuyerID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 200*time.Millisecond)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
	}

	orders, err := oh.orderSvc.FindOrdersByBuyerID(ctx, userInfo.UserID)
	if err != nil {
		log.Printf("FindOrderByBuyerID error:\n+%v", err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"count":  len(orders),
		"orders": orders,
	})
}

func (oh *OrderHttpHandler) FindOrdersBySellerID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 200*time.Millisecond)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
	}

	cc, err := oh.grpcClient.GetClient("USER_SERVICE")
	if err != nil {
		log.Printf("FindOrdersBySellerID error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching orders")
	}

	sellerGrpcClient := user.NewUserServiceClient(cc)
	sellerInfo, err := sellerGrpcClient.FindSeller(ctx, &user.FindSellerRequest{
		BuyerId:  userInfo.UserID,
		SellerId: "",
	})
	if err != nil {
		log.Printf("FindOrdersBySellerID error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching seller data")
	}

	orders, err := oh.orderSvc.FindOrdersBySellerID(ctx, sellerInfo.Id)
	if err != nil {
		log.Printf("FindOrdersBySellerID error:\n+%v", err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"count":  len(orders),
		"orders": orders,
	})
}

// NOTE: SAVE ORDER DATA WHEN BUYER INTENTS TO PAY
// AND SET PAYMENT SUCCEED WEBHOOKS FROM STRIPE TO HANDLE CHANGING ORDER STATUS ON BACKEND
func (oh *OrderHttpHandler) CreatePaymentIntent(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 500*time.Millisecond)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
	}

	data := new(types.CreateOrderDTO)
	err := c.BodyParser(data)
	if err != nil {
		log.Printf("CreatePaymentIntent error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	err = oh.validate.Struct(data)
	if err != nil {
		log.Printf("CreatePaymentIntent error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	cc, err := oh.grpcClient.GetClient("USER_SERVICE")
	if err != nil {
		log.Printf("CreatePaymentIntent error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while validating gig")
	}

	sellerGrpcClient := user.NewUserServiceClient(cc)
	sellerInfo, err := sellerGrpcClient.FindSeller(ctx, &user.FindSellerRequest{
		BuyerId:  "",
		SellerId: data.SellerID,
	})

	pi, err := paymentintent.New(&stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(data.Price * 100)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		ReceiptEmail: &userInfo.Email,
		OnBehalfOf:   stripe.String(sellerInfo.StripeAccountId),
		Metadata: map[string]string{
			"buyer_id":  data.BuyerID,
			"seller_id": data.SellerID,
		},
	})
	if err != nil {
		log.Printf("CreatePaymentIntent error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while processing payment")
	}

	data.BuyerID = userInfo.UserID
	data.PaymentIntentID = pi.ID
	data.StripeClientSecret = pi.ClientSecret
	_, err = oh.orderSvc.SaveOrder(ctx, data)
	if err != nil {
		log.Printf("CreatePaymentIntent error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while processing order")
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"client_secret": pi.ClientSecret,
	})
}

// NOTE: TEST FIRST TO DETERMINE IF IT NEEDS TO USING GOROUTINE
func (oh *OrderHttpHandler) HandleStripeWebhook(c *fiber.Ctx) error {
	endpointSecret := os.Getenv("STRIPE_ENDPOINT_SECRET")
	event, err := webhook.ConstructEvent(c.Body(), c.Get("Stripe-Signature"), endpointSecret)
	if err != nil {
		log.Printf("Error verifying webhook signature: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	switch event.Type {
	case stripe.EventTypePaymentIntentSucceeded:
		ctx, cancel := context.WithTimeout(c.UserContext(), 200*time.Millisecond)
		defer cancel()

		var pi stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &pi)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Unexpected error happened. Please try again.")
		}

		o, err := oh.orderSvc.FindOrderByPaymentIntentID(ctx, pi.ID)
		if err != nil {
			log.Printf("Order did not found:\n+%v", err)
			return fiber.NewError(fiber.StatusNotFound, "Order is invalid")
		}

		err = oh.orderSvc.ChangeOrderStatus(ctx, *o, types.PENDING)
		if err != nil {
			log.Printf("Changing order status error:\n+%v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Unexpected error happened. Please try again.")
		}
	default:
		log.Printf("Stripe Invalid Webhook Event:\n+%v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Unknown webhook event")
	}

	return c.SendStatus(fiber.StatusOK)
}

func (oh *OrderHttpHandler) BuyerMarkOrderAsComplete(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Please sign-in first")
	}

	o, err := oh.orderSvc.FindOrderByID(ctx, c.Params("orderId"))
	if err != nil {
		log.Printf("OrderComplete error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Order is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while finding order")
	}

	cc, err := oh.grpcClient.GetClient("USER_SERVICE")
	if err != nil {
		log.Printf("OrderComplete error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	userGrpcClient := user.NewUserServiceClient(cc)
	sellerInfo, err := userGrpcClient.UpdateSellerBalance(ctx, &user.UpdateSellerBalanceRequest{
		SellerId: o.SellerID,
		Amount:   o.Price,
	})
	if err != nil {
		log.Printf("OrderComplete error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching seller data")
	}

	//HACK: IF THERE ARE ERRORS FOR SENDING EMAIL NOTIFICATION
	// THEN IT SHOULDN'T AFFECT CODE FLOW
	go func() {
		cc, err = oh.grpcClient.GetClient("NOTIFICATION_SERVICE")
		if err != nil {
			log.Printf("OrderComplete error:\n+%v", err)
		}

		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		_, err = notificationGrpcClient.SellerHasCompletedAnOrder(context.TODO(), &notification.SellerCompletedAnOrderRequest{
			ReceiverEmail:        sellerInfo.Email,
			BuyerEmail:           userInfo.Email,
			OrderId:              o.ID,
			SellerCurrentBalance: strconv.FormatUint(sellerInfo.AccountBalance, 10),
		})
		if err != nil {
			log.Printf("OrderComplete error:\n+%v", err)
		}
	}()

	err = oh.orderSvc.ChangeOrderStatus(ctx, *o, types.COMPLETED)
	if err != nil {
		log.Printf("OrderComplete error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while updating Order Status")
	}

	return c.SendStatus(http.StatusOK)
}

// TODO: REFACTORE
func (oh *OrderHttpHandler) SellerCancellingOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Please sign-in first")
	}

	o, err := oh.orderSvc.FindOrderByID(ctx, c.Params("orderId"))
	if err != nil {
		log.Printf("CancelOrder error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Order is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while finding order")
	}

	err = oh.orderSvc.ChangeOrderStatus(ctx, *o, types.CANCELED)
	if err != nil {
		log.Printf("CancelOrder error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while updating Order Status")
	}

	return c.SendStatus(http.StatusOK)
}

// TODO: IMPLEMENT
func (oh *OrderHttpHandler) RequestDeadlineExtension(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 500*time.Millisecond)
	defer cancel()

	data := new(types.DeadlineExtensionRequest)
	err := c.BodyParser(data)
	if err != nil {
		log.Printf("RequestExtendingDeadline error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	err = oh.validate.Struct(data)
	if err != nil {
		log.Printf("RequestExtendingDeadline error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	o, err := oh.orderSvc.FindOrderByID(ctx, c.Params("orderId"))
	if err != nil {
		log.Printf("RequestExtendingDeadline error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Order is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while finding order")
	}

	err = oh.orderSvc.RequestDeadlineExtension(ctx, *o, data)
	if err != nil {
		log.Printf("RequestExtendingDeadline error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	go func() {
		cc, err := oh.grpcClient.GetClient("USER_SERVICE")
		if err != nil {
			log.Printf("RequestExtendingDeadline error:\n+%v", err)
		}

		userGrpcClient := user.NewUserServiceClient(cc)
		b, err := userGrpcClient.FindBuyer(ctx, &user.FindBuyerRequest{
			BuyerId: o.BuyerID,
		})
		if err != nil {
			log.Printf("RequestExtendingDeadline error:\n+%v", err)
		}

		cc, err = oh.grpcClient.GetClient("NOTIFICATION_SERVICE")
		if err != nil {
			log.Printf("RequestExtendingDeadline error:\n+%v", err)
		}

		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		_, err = notificationGrpcClient.SellerRequestDeadlineExtension(context.TODO(), &notification.SellerDeadlineExtensionRequest{
			ReceiverEmail: b.Email,
			Message:       fmt.Sprintf("Seller Requesting To Extend The Deadline To Become [%v] With Reason:\n%s", o.Deadline.Add(time.Duration(data.NumberOfDays)), data.Reason),
		})
		if err != nil {
			log.Printf("RequestExtendingDeadline error:\n+%v", err)
		}
	}()

	return c.SendStatus(http.StatusOK)
}

// TODO: IMPLEMENT
func (oh *OrderHttpHandler) ApproveDeadlineExtension(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 500*time.Millisecond)
	defer cancel()

	o, err := oh.orderSvc.FindOrderByID(ctx, c.Params("orderId"))
	if err != nil {
		log.Printf("RequestExtendingDeadline error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Order is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while finding order")
	}

	err = oh.orderSvc.ApproveDeadlineExtension(ctx, *o, data)
	if err != nil {
		log.Printf("RequestExtendingDeadline error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	return c.SendStatus(http.StatusOK)
}
