package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Akihira77/gojobber/services/7-order/service"
	"github.com/Akihira77/gojobber/services/7-order/types"
	"github.com/Akihira77/gojobber/services/common/genproto/notification"
	"github.com/Akihira77/gojobber/services/common/genproto/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/paymentintent"
	"gorm.io/gorm"
)

type OrderHttpHandler struct {
	orderSvc   service.OrderServiceImpl
	grpcClient *GRPCClients
	stripeKey  string
	validate   *validator.Validate
}

func NewOrderHttpHandler(orderSvc service.OrderServiceImpl) *OrderHttpHandler {
	return &OrderHttpHandler{
		orderSvc:  orderSvc,
		stripeKey: os.Getenv("STRIPE_SECRET_KEY"),
		validate:  validator.New(validator.WithRequiredStructEnabled()),
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

func (oh *OrderHttpHandler) CreatePaymentIntent(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 200*time.Millisecond)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
	}

	data := new(types.CreatePaymentIntentDTO)
	err := c.BodyParser(data)
	if err != nil || data.Amount <= 0 {
		log.Printf("ConfirmPayment error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
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

	stripe.Key = oh.stripeKey
	paymentIntent, err := paymentintent.New(&stripe.PaymentIntentParams{
		Amount:   stripe.Int64(data.Amount * 100),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		ReceiptEmail: &userInfo.Email,
		OnBehalfOf:   stripe.String(sellerInfo.StripeAccountId),
		// TransferData: &stripe.PaymentIntentTransferDataParams{
		// 	Amount:      stripe.Int64(data.Amount * (1 - (25 / 100))),
		// 	Destination: stripe.String(sellerInfo.Id),
		// },
	})
	if err != nil {
		log.Printf("CreatePaymentIntent error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while processing order")
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"id":            paymentIntent.ID,
		"client_secret": paymentIntent.ClientSecret,
	})
}

func (oh *OrderHttpHandler) ConfirmPayment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
	}

	data := new(types.CreateOrderDTO)
	err := c.BodyParser(data)
	if err != nil {
		log.Printf("ConfirmPayment error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	err = oh.validate.Struct(data)
	if err != nil {
		log.Printf("ConfirmPayment error:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	stripe.Key = oh.stripeKey
	buyerPaymentIntent, err := paymentintent.Get(c.Params("paymentId"), &stripe.PaymentIntentParams{
		ClientSecret: &data.ClientSecret,
	})
	if err != nil {
		log.Printf("ConfirmPayment error:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "Error while processing order. Please try again!")
	}

	if buyerPaymentIntent.Status == stripe.PaymentIntentStatusSucceeded {
		return c.Status(http.StatusOK).SendString("This Payment Is Already Succeded")
	} else if buyerPaymentIntent.Status != stripe.PaymentIntentStatusRequiresPaymentMethod {
		log.Printf("Buyer Payment Status is not RequiresPaymentMethod")
		return fiber.NewError(http.StatusBadRequest, "There is something wrong with your payment. Please check and try again!")
	}

	result, err := paymentintent.Confirm(c.Params("paymentId"), &stripe.PaymentIntentConfirmParams{
		PaymentMethod: stripe.String("pm_card_visa"),
		ReceiptEmail:  &userInfo.Email,
		ReturnURL:     stripe.String(fmt.Sprintf("%s/order/buyer/my-orders", os.Getenv("CLIENT_URL"))),
	})
	if err != nil {
		log.Printf("ConfirmPayment error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while confirming payment")
	}

	data.BuyerID = userInfo.UserID
	data.PaymentIntentID = result.ID
	o, err := oh.orderSvc.CreateOrder(ctx, data)
	if err != nil {
		log.Printf("ConfirmPayment error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while confirming payment")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"order": o,
	})
}

func (oh *OrderHttpHandler) OrderComplete(c *fiber.Ctx) error {
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

	//HACK: IF THERE IS ERROR FOR SENDING EMAIL NOTIFICATION
	// THEN IT SHOULD'NT AFFECT THE CODE FLOW
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
			SellerCurrentBalance: string(sellerInfo.AccountBalance),
		})
		if err != nil {
			log.Printf("OrderComplete error:\n+%v", err)
		}
	}()

	err = oh.orderSvc.ChangeOrderStatus(ctx, o, types.COMPLETED)
	if err != nil {
		log.Printf("OrderComplete error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while updating Order Status")
	}

	return c.SendStatus(http.StatusNoContent)
}

// TODO: REFACTORE
func (oh *OrderHttpHandler) CancelOrder(c *fiber.Ctx) error {
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

	err = oh.orderSvc.ChangeOrderStatus(ctx, o, types.CANCELED)
	if err != nil {
		log.Printf("CancelOrder error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while updating Order Status")
	}

	return c.SendStatus(http.StatusNoContent)
}

// TODO: IMPLEMENT
func (oh *OrderHttpHandler) RequestExtendingDeadline(c *fiber.Ctx) error {
	return nil
}

// TODO: IMPLEMENT
func (oh *OrderHttpHandler) ExtendDeadline(c *fiber.Ctx) error {
	return nil
}
