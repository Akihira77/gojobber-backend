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
	"strings"
	"time"

	"github.com/Akihira77/gojobber/services/7-order/service"
	"github.com/Akihira77/gojobber/services/7-order/types"
	"github.com/Akihira77/gojobber/services/common/genproto/notification"
	"github.com/Akihira77/gojobber/services/common/genproto/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/paymentintent"
	"github.com/stripe/stripe-go/v80/refund"
	"github.com/stripe/stripe-go/v80/webhook"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (oh *OrderHttpHandler) FindMyOrdersAsBuyer(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 200*time.Millisecond)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
	}

	orders, err := oh.orderSvc.FindOrdersByBuyerID(ctx, userInfo.UserID)
	if err != nil {
		log.Printf("FindMyOrdersAsBuyer error:\n+%v", err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"count":  len(orders),
		"orders": orders,
	})
}

func (oh *OrderHttpHandler) FindMyOrdersAsSeller(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 200*time.Millisecond)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
	}

	cc, err := oh.grpcClient.GetClient("USER_SERVICE")
	if err != nil {
		log.Printf("FindMyOrdersAsSeller error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching orders")
	}

	userGrpcClient := user.NewUserServiceClient(cc)
	s, err := userGrpcClient.FindSeller(ctx, &user.FindSellerRequest{
		BuyerId:  userInfo.UserID,
		SellerId: "",
	})
	if err != nil {
		log.Printf("FindMyOrdersAsSeller error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching seller data")
	}

	orders, err := oh.orderSvc.FindOrdersBySellerID(ctx, s.Id)
	if err != nil {
		log.Printf("FindMyOrdersAsSeller error:\n+%v", err)
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

	userGrpcClient := user.NewUserServiceClient(cc)
	s, err := userGrpcClient.FindSeller(ctx, &user.FindSellerRequest{
		BuyerId:  "",
		SellerId: data.SellerID,
	})
	if err != nil {
		log.Printf("CreatePaymentIntent error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while finding seller related to this gig")
	}

	pi, err := paymentintent.New(&stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(data.Price * 100)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		ReceiptEmail: &userInfo.Email,
		OnBehalfOf:   stripe.String(s.StripeAccountId),
		Metadata: map[string]string{
			"buyer_id":     data.BuyerID,
			"seller_id":    data.SellerID,
			"seller_email": s.Email,
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
	es := os.Getenv("STRIPE_ENDPOINT_SECRET")
	e, err := webhook.ConstructEvent(c.Body(), c.Get("Stripe-Signature"), es)
	if err != nil {
		log.Printf("Error verifying webhook signature: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	switch e.Type {
	case stripe.EventTypePaymentIntentSucceeded:
		go func() {
			ctx, cancel := context.WithTimeout(c.UserContext(), 500*time.Millisecond)
			defer cancel()

			var pi stripe.PaymentIntent
			err := json.Unmarshal(e.Data.Raw, &pi)
			if err != nil {
				log.Printf("Error parsing webhook JSON: %v", err)
				return
			}

			o, err := oh.orderSvc.FindOrderByPaymentIntentID(ctx, pi.ID)
			if err != nil {
				log.Printf("Order did not found:\n+%v", err)
				return
			}

			_, err = oh.orderSvc.ChangeOrderStatus(ctx, *o, types.PENDING, fmt.Sprintf("Buyer Has Paid This Order. Order Status Change To [%s]", types.PENDING))
			if err != nil {
				log.Printf("Changing order status error:\n+%v", err)
				return
			}

			sellerEmail := pi.Metadata["seller_email"]
			cc, err := oh.grpcClient.GetClient("NOTIFICATION_SERVICE")
			if err != nil {
				log.Printf("HandleStripeWebhook Error:\n+%v", err)
				return
			}

			notificationGrpcClient := notification.NewNotificationServiceClient(cc)
			_, err = notificationGrpcClient.NotifySellerOrderHasBeenMade(context.TODO(), &notification.NotifySellerGotAnOrderRequest{
				ReceiverEmail: sellerEmail,
				Message:       fmt.Sprintf("You Recevie An Order From Buyer [%s]", o.BuyerID),
				Detail: &notification.OrderDetail{
					GigTitle:       o.GigTitle,
					GigDescription: o.GigDescription,
					Price:          o.Price,
					ServiceFee:     uint64(o.ServiceFee),
					Deadline:       timestamppb.New(o.Deadline),
				},
			})
			if err != nil {
				log.Printf("HandleStripeWebhook Error:\n+%v", err)
				return
			}
		}()
	default:
		log.Printf("Stripe Unhandled Webhook Event:\n+%v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Unhandled webhook event")
	}

	return c.SendStatus(fiber.StatusOK)
}

func (oh *OrderHttpHandler) BuyerMarkOrderAsComplete(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
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
	s, err := userGrpcClient.UpdateSellerBalance(ctx, &user.UpdateSellerBalanceRequest{
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
			return
		}

		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		_, err = notificationGrpcClient.SellerHasCompletedAnOrder(context.TODO(), &notification.SellerCompletedAnOrderRequest{
			ReceiverEmail:        s.Email,
			BuyerEmail:           userInfo.Email,
			OrderId:              o.ID,
			SellerCurrentBalance: strconv.FormatUint(s.AccountBalance, 10),
		})
		if err != nil {
			log.Printf("OrderComplete error:\n+%v", err)
			return
		}
	}()

	o, err = oh.orderSvc.ChangeOrderStatus(ctx, *o, types.COMPLETED, fmt.Sprintf("Buyer Has Marked This Order As Complete"))
	if err != nil {
		log.Printf("OrderComplete error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while updating Order Status")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"order": o,
	})
}

func (oh *OrderHttpHandler) SellerCancellingOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Please sign-in first")
	}

	type Data struct {
		Reason string `json:"reason" validate:"required"`
	}

	data := new(Data)
	err := c.BodyParser(data)
	if err != nil {
		log.Printf("SellerCancellingOrder error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	err = oh.validate.Struct(data)
	if err != nil {
		log.Printf("SellerCancellingOrder error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	o, err := oh.orderSvc.FindOrderByID(ctx, c.Params("orderId"))
	if err != nil {
		log.Printf("SellerCancellingOrder error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Order is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while finding order")
	}

	o, err = oh.orderSvc.ChangeOrderStatus(ctx, *o, types.CANCELED, fmt.Sprintf("Seller Has Canceled This Order With Reason:\n%s", data.Reason))
	if err != nil {
		log.Printf("SellerCancellingOrder error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while updating Order Status")
	}

	//HACK: SEND EMAIL TO BUYER THAT THE SELLER HAS CANCEL THE ORDER
	//IGNORE THE ERROR FROM CODE FLOW
	go func() {
		newCtx, canc := context.WithTimeout(c.UserContext(), 200*time.Millisecond)
		defer canc()

		_, err := refund.New(&stripe.RefundParams{
			PaymentIntent: &o.PaymentIntentID,
			Reason:        &data.Reason,
		})
		if err != nil {
			log.Printf("SellerCancellingOrder error:\n+%v", err)
			return
		}

		cc, err := oh.grpcClient.GetClient("USER_SERVICE")
		if err != nil {
			log.Printf("SellerCancellingOrder error:\n+%v", err)
			return
		}

		userGrpcClient := user.NewUserServiceClient(cc)
		b, err := userGrpcClient.FindBuyer(newCtx, &user.FindBuyerRequest{
			BuyerId: o.BuyerID,
		})
		if err != nil {
			log.Printf("SellerCancellingOrder error:\n+%v", err)
			return
		}

		cc, err = oh.grpcClient.GetClient("NOTIFICATION_SERVICE")
		if err != nil {
			log.Printf("SellerCancellingOrder error:\n+%v", err)
			return
		}

		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		_, err = notificationGrpcClient.SellerCanceledAnOrder(context.TODO(), &notification.SellerCancelOrderRequest{
			ReceiverEmail: b.Email,
			Reason:        data.Reason,
		})
		if err != nil {
			log.Printf("SellerCancellingOrder error:\n+%v", err)
			return
		}
	}()

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"order": o,
	})
}

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

	//HACK: SEND EMAIL TO BUYER THAT SELLER REQUEST DEADLINE EXTENSION
	// IGNORE ERROR FROM CODE FLOW
	go func() {
		newCtx, canc := context.WithTimeout(c.UserContext(), 200*time.Millisecond)
		defer canc()

		cc, err := oh.grpcClient.GetClient("USER_SERVICE")
		if err != nil {
			log.Printf("RequestExtendingDeadline error:\n+%v", err)
			return
		}

		userGrpcClient := user.NewUserServiceClient(cc)
		b, err := userGrpcClient.FindBuyer(newCtx, &user.FindBuyerRequest{
			BuyerId: o.BuyerID,
		})
		if err != nil {
			log.Printf("RequestExtendingDeadline error:\n+%v", err)
			return
		}

		cc, err = oh.grpcClient.GetClient("NOTIFICATION_SERVICE")
		if err != nil {
			log.Printf("RequestExtendingDeadline error:\n+%v", err)
			return
		}

		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		_, err = notificationGrpcClient.SellerRequestDeadlineExtension(context.TODO(), &notification.SellerDeadlineExtensionRequest{
			ReceiverEmail: b.Email,
			Reason:        fmt.Sprintf("Seller Request Deadline Extension From [%v] To [%v] With Reason:\n%s", o.Deadline, o.Deadline.Add(time.Duration(data.NumberOfDays)), data.Reason),
		})
		if err != nil {
			log.Printf("RequestExtendingDeadline error:\n+%v", err)
			return
		}
	}()

	return c.SendStatus(http.StatusOK)
}

func (oh *OrderHttpHandler) BuyerDeadlineExtensionResponse(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 500*time.Millisecond)
	defer cancel()

	status := strings.ToUpper(c.Query("extension-status", "REJECTED"))

	data := new(types.DeadlineExtensionRequest)
	err := c.BodyParser(data)
	if err != nil {
		log.Printf("BuyerDeadlineExtensionResponse error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	err = oh.validate.Struct(data)
	if err != nil {
		log.Printf("BuyerDeadlineExtensionResponse error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	o, err := oh.orderSvc.FindOrderByID(ctx, c.Params("orderId"))
	if err != nil {
		log.Printf("BuyerDeadlineExtensionResponse error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Order is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while finding order")
	}

	msg, err := oh.orderSvc.DeadlineExtensionResponse(ctx, *o, types.DeadlineExtensionStatus(status), data)
	if err != nil {
		log.Printf("BuyerDeadlineExtensionResponse error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	//HACK: SEND EMAIL TO BUYER THAT SELLER REQUEST DEADLINE EXTENSION
	// IGNORE ERROR FROM CODE FLOW
	go func() {
		newCtx, canc := context.WithTimeout(c.UserContext(), 200*time.Millisecond)
		defer canc()

		cc, err := oh.grpcClient.GetClient("USER_SERVICE")
		if err != nil {
			log.Printf("BuyerDeadlineExtensionResponse error:\n+%v", err)
			return
		}

		userGrpcClient := user.NewUserServiceClient(cc)
		b, err := userGrpcClient.FindBuyer(newCtx, &user.FindBuyerRequest{
			BuyerId: o.BuyerID,
		})
		if err != nil {
			log.Printf("BuyerDeadlineExtensionResponse error:\n+%v", err)
			return
		}

		cc, err = oh.grpcClient.GetClient("NOTIFICATION_SERVICE")
		if err != nil {
			log.Printf("BuyerDeadlineExtensionResponse error:\n+%v", err)
			return
		}

		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		_, err = notificationGrpcClient.BuyerDeadlineExtensionResponse(context.TODO(), &notification.BuyerDeadlineExtension{
			ReceiverEmail: b.Email,
			Message:       msg,
		})
		if err != nil {
			log.Printf("BuyerDeadlineExtensionResponse error:\n+%v", err)
			return
		}
	}()

	return c.SendStatus(http.StatusOK)
}

func (oh *OrderHttpHandler) BuyerRefundingOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 500*time.Millisecond)
	defer cancel()

	type Data struct {
		Reason string `json:"reason" validate:"required"`
	}

	data := new(Data)
	err := c.BodyParser(data)
	if err != nil {
		log.Printf("BuyerRefundingOrder error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	err = oh.validate.Struct(data)
	if err != nil {
		log.Printf("BuyerRefundingOrder error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	o, err := oh.orderSvc.FindOrderByID(ctx, c.Params("orderId"))
	if err != nil {
		log.Printf("BuyerRefundingOrder error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Order is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while finding order")
	}

	result, err := refund.New(&stripe.RefundParams{
		PaymentIntent: &o.PaymentIntentID,
		Reason:        &data.Reason,
	})
	if err != nil {
		log.Printf("BuyerRefundingOrder error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while trying to refund your money")
	}

	switch result.Status {
	case stripe.RefundStatusSucceeded:
		o, err = oh.orderSvc.ChangeOrderStatus(ctx, *o, types.REFUNDED, fmt.Sprintf("Buyer Refunds This Order With Reason:\n%s", data.Reason))
		if err != nil {
			log.Printf("BuyerRefundingOrder error:\n+%v", err)
			return fiber.NewError(http.StatusInternalServerError, "Error while updating order status")
		}

		//HACK: SEND EMAIL TO SELLER THAT BUYER HAS REFUNDS THIS ORDER
		// IGNORE ERROR FROM CODE FLOW
		go func() {
			newCtx, canc := context.WithTimeout(c.UserContext(), 200*time.Millisecond)
			defer canc()

			cc, err := oh.grpcClient.GetClient("USER_SERVICE")
			if err != nil {
				log.Printf("BuyerRefundingOrder error:\n+%v", err)
			}

			userGrpcClient := user.NewUserServiceClient(cc)
			b, err := userGrpcClient.FindBuyer(newCtx, &user.FindBuyerRequest{
				BuyerId: o.BuyerID,
			})
			if err != nil {
				log.Printf("BuyerRefundingOrder error:\n+%v", err)
				return
			}

			cc, err = oh.grpcClient.GetClient("NOTIFICATION_SERVICE")
			if err != nil {
				log.Printf("BuyerRefundingOrder error:\n+%v", err)
				return
			}

			notificationGrpcClient := notification.NewNotificationServiceClient(cc)
			_, err = notificationGrpcClient.BuyerRefundsAnOrder(context.TODO(), &notification.BuyerRefundsOrderRequest{
				ReceiverEmail: b.Email,
				Reason:        data.Reason,
			})
			if err != nil {
				log.Printf("BuyerRefundingOrder error:\n+%v", err)
				return
			}
		}()

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"order": o,
		})
	case stripe.RefundStatusCanceled:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to refund",
		})
	default:
		log.Printf("Unhandled Stripe Refund Status:\n+%v", result)
		return fiber.NewError(http.StatusBadRequest, "Unhandled Stripe Refund Status")
	}
}

// TODO: SEND EMAIL TO BUYER THAT SELLER HAS SENT THE ORDER RESULT
func (oh *OrderHttpHandler) SellerDeliverOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 500*time.Millisecond)
	defer cancel()

	o, err := oh.orderSvc.FindOrderByID(ctx, c.Params("orderId"))
	if err != nil {
		log.Printf("SellerDeliverOrder Error:\n+%v", err)
		return fiber.NewError(http.StatusNotFound, "Order is not found")
	}

	data := new(types.DeliveredHistory)
	err = c.BodyParser(data)
	if err != nil {
		log.Printf("SellerDeliverOrder Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Invalid Provided Data")
	}

	err = oh.validate.Struct(data)
	if err != nil {
		log.Printf("SellerDeliverOrder Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Invalid Provided Data")
	}

	o, err = oh.orderSvc.DeliveringOrder(ctx, *o, *data)
	if err != nil {
		log.Printf("SellerDeliverOrder Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Error while saving the provided data")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"order": o,
	})
}

func (oh *OrderHttpHandler) BuyerResponseForDeliveredOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 500*time.Millisecond)
	defer cancel()

	o, err := oh.orderSvc.FindOrderByID(ctx, c.Params("orderId"))
	if err != nil {
		log.Printf("BuyerNoteForDeliveredOrder Error:\n+%v", err)
		return fiber.NewError(http.StatusNotFound, "Order is not found")
	}

	data := new(types.BuyerResponseOrderDelivered)
	err = c.BodyParser(data)
	if err != nil {
		log.Printf("BuyerNoteForDeliveredOrder Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Invalid Provided Data")
	}

	err = oh.validate.Struct(data)
	if err != nil {
		log.Printf("BuyerNoteForDeliveredOrder Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Invalid Provided Data")
	}

	o, err = oh.orderSvc.OrderDeliveredResponse(ctx, *o, data)
	if err != nil {
		log.Printf("BuyerNoteForDeliveredOrder Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Error while saving the provided data")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"order": o,
	})
}
