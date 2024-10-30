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
	"github.com/Akihira77/gojobber/services/7-order/util"
	"github.com/Akihira77/gojobber/services/common/genproto/chat"
	"github.com/Akihira77/gojobber/services/common/genproto/notification"
	"github.com/Akihira77/gojobber/services/common/genproto/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/account"
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

func NewOrderHttpHandler(orderSvc service.OrderServiceImpl, grpcClients *GRPCClients) *OrderHttpHandler {
	return &OrderHttpHandler{
		grpcClient: grpcClients,
		orderSvc:   orderSvc,
		validate:   validator.New(validator.WithRequiredStructEnabled()),
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

	cc, err := oh.grpcClient.GetClient(types.USER_SERVICE)
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
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println("CreatePaymentIntent userInfo", userInfo)
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
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": util.CustomValidationErrors(err),
		})
	}

	cc, err := oh.grpcClient.GetClient(types.USER_SERVICE)
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
		OnBehalfOf:   &s.StripeAccountId,
		Metadata: map[string]string{
			"buyer_id":     userInfo.UserID,
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
		go paymentintent.Cancel(pi.ID, &stripe.PaymentIntentCancelParams{
			CancellationReason: stripe.String("Error saving order data in GoJobber Platform"),
		})
		return fiber.NewError(http.StatusInternalServerError, "Error while processing order")
	}

	go func() {
		if data.MessageID == "" {
			return
		}

		newCtx, canc := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer canc()

		cc, err := oh.grpcClient.GetClient(types.CHAT_SERVICE)
		if err != nil {
			log.Printf("Change Offer Status Error:\n+%v", err)
			return
		}

		chatGrpcClient := chat.NewChatServiceClient(cc)
		_, err = chatGrpcClient.BuyerAcceptedOffer(newCtx, &chat.BuyerAcceptedOfferRequest{
			MessageId: data.MessageID,
		})
		if err != nil {
			log.Printf("Change Offer Status Error:\n+%v", err)
			return
		}
	}()

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"client_secret": pi.ClientSecret,
	})
}

func (oh *OrderHttpHandler) ConfirmPayment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
	}

	o, err := oh.orderSvc.FindOrderByPaymentIntentID(ctx, c.Params("paymentId"))
	if err != nil {
		log.Printf("ConfirmPayment error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Order is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	if o.BuyerID != userInfo.UserID {
		return fiber.ErrForbidden
	}

	result, err := paymentintent.Confirm(o.PaymentIntentID, &stripe.PaymentIntentConfirmParams{
		PaymentMethod: stripe.String("pm_card_visa"),
		ReturnURL:     stripe.String(fmt.Sprintf("%s/orders", os.Getenv("CLIENT_URL"))),
	})

	if err != nil {
		log.Printf("ConfirmPayment error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Order is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	o, err = oh.orderSvc.ChangeOrderStatus(ctx, *o, types.PENDING, fmt.Sprint("Buyer Has Completed Payment Process"))
	if err != nil {
		log.Printf("ConfirmPayment error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"result": result,
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
			cc, err := oh.grpcClient.GetClient(types.NOTIFICATION_SERVICE)
			if err != nil {
				log.Printf("HandleStripeWebhook Error:\n+%v", err)
				return
			}

			notificationGrpcClient := notification.NewNotificationServiceClient(cc)
			_, err = notificationGrpcClient.NotifySellerOrderHasBeenMade(context.TODO(), &notification.NotifySellerGotAnOrderRequest{
				ReceiverEmail: sellerEmail,
				Message:       fmt.Sprintf("You Receive An Order From Buyer [%s]", o.BuyerID),
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

	cc, err := oh.grpcClient.GetClient(types.USER_SERVICE)
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
		cc, err = oh.grpcClient.GetClient(types.NOTIFICATION_SERVICE)
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
			Url:                  fmt.Sprintf("%s/orders/%s", os.Getenv("CLIENT_URL"), o.ID),
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
	ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
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
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": util.CustomValidationErrors(err),
		})
	}

	cc, err := oh.grpcClient.GetClient(types.USER_SERVICE)
	if err != nil {
		log.Printf("SellerCancellingOrder error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while finding order")
	}

	userGrpcClient := user.NewUserServiceClient(cc)
	s, err := userGrpcClient.FindSeller(ctx, &user.FindSellerRequest{
		BuyerId:  userInfo.UserID,
		SellerId: "",
	})
	if err != nil {
		log.Printf("SellerCancellingOrder error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching seller data")
	}

	o, err := oh.orderSvc.FindOrderByIDAndSellerID(ctx, c.Params("orderId"), s.Id)
	if err != nil {
		log.Printf("SellerCancellingOrder error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Order is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while finding order")
	}

	if o.Status != types.PROCESS {
		log.Println("Order is not in PROCESS status, seller cannot cancel this order", o)
		return fiber.NewError(http.StatusBadRequest, "Seller cannot cancel this order because Order Status is not in PROCESS stage")
	}

	refundRes, err := refund.New(&stripe.RefundParams{
		PaymentIntent: &o.PaymentIntentID,
		Reason:        stripe.String(string(stripe.RefundReasonRequestedByCustomer)),
	})
	if err != nil || refundRes.Status != stripe.RefundStatusSucceeded {
		log.Printf("SellerCancellingOrder error:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "Error while processing your refund. Please try again.")
	}

	o, err = oh.orderSvc.ChangeOrderStatus(ctx, *o, types.CANCELED, fmt.Sprintf("Seller Has Canceled This Order With Reason:\n%s", data.Reason))
	if err != nil {
		log.Printf("SellerCancellingOrder error:\n%+v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while updating Order Status")
	}

	//HACK: SEND EMAIL TO BUYER THAT THE SELLER HAS CANCEL THE ORDER
	//IGNORE THE ERROR FROM CODE FLOW
	go func() {
		newCtx, canc := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer canc()

		userGrpcClient := user.NewUserServiceClient(cc)
		b, err := userGrpcClient.FindBuyer(newCtx, &user.FindBuyerRequest{
			BuyerId: o.BuyerID,
		})
		if err != nil {
			log.Printf("SellerCancellingOrder error:\n+%v", err)
			return
		}

		cc, err = oh.grpcClient.GetClient(types.NOTIFICATION_SERVICE)
		if err != nil {
			log.Printf("SellerCancellingOrder error:\n+%v", err)
			return
		}

		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		_, err = notificationGrpcClient.SellerCanceledAnOrder(context.TODO(), &notification.SellerCancelOrderRequest{
			ReceiverEmail: b.Email,
			Url:           fmt.Sprintf("%s/orders/%s", os.Getenv("CLIENT_URL"), o.ID),
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
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": util.CustomValidationErrors(err),
		})
	}

	o, err := oh.orderSvc.FindOrderByID(ctx, c.Params("orderId"))
	if err != nil {
		log.Printf("RequestExtendingDeadline error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Order is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while finding order")
	}

	if o.Status != types.PROCESS {
		log.Println("Order is not in PROCESS status, seller cannot make a deadline extension request for this order", o)
		return fiber.NewError(http.StatusBadRequest, "Seller cannot make a deadline extension request for this order because Order Status is not in PROCESS stage")
	}

	err = oh.orderSvc.RequestDeadlineExtension(ctx, *o, data)
	if err != nil {
		log.Printf("RequestExtendingDeadline error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	//HACK: SEND EMAIL TO BUYER THAT SELLER REQUEST DEADLINE EXTENSION
	// IGNORE ERROR FROM CODE FLOW
	go func() {
		newCtx, canc := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer canc()

		cc, err := oh.grpcClient.GetClient(types.USER_SERVICE)
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

		cc, err = oh.grpcClient.GetClient(types.NOTIFICATION_SERVICE)
		if err != nil {
			log.Printf("RequestExtendingDeadline error:\n+%v", err)
			return
		}

		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		_, err = notificationGrpcClient.SellerRequestDeadlineExtension(context.TODO(), &notification.SellerDeadlineExtensionRequest{
			ReceiverEmail: b.Email,
			Url:           fmt.Sprintf("%s/orders/%s", os.Getenv("CLIENT_URL"), o.ID),
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

	data := new(types.DeadlineExtensionRequest)
	err := c.BodyParser(data)
	if err != nil {
		log.Printf("BuyerDeadlineExtensionResponse error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	err = oh.validate.Struct(data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": util.CustomValidationErrors(err),
		})
	}

	o, err := oh.orderSvc.FindOrderByID(ctx, c.Params("orderId"))
	if err != nil {
		log.Printf("BuyerDeadlineExtensionResponse error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Order is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while finding order")
	}

	_, err = oh.orderSvc.DeadlineExtensionResponse(ctx, *o, types.DeadlineExtensionStatus(data.BuyerResponse), data)
	if err != nil {
		log.Printf("BuyerDeadlineExtensionResponse error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	//HACK: SEND EMAIL TO BUYER THAT SELLER REQUEST DEADLINE EXTENSION
	// IGNORE ERROR FROM CODE FLOW
	go func() {
		newCtx, canc := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer canc()

		cc, err := oh.grpcClient.GetClient(types.USER_SERVICE)
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

		cc, err = oh.grpcClient.GetClient(types.NOTIFICATION_SERVICE)
		if err != nil {
			log.Printf("BuyerDeadlineExtensionResponse error:\n+%v", err)
			return
		}

		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		_, err = notificationGrpcClient.BuyerDeadlineExtensionResponse(context.TODO(), &notification.BuyerDeadlineExtension{
			ReceiverEmail: b.Email,
			Url:           fmt.Sprintf("%s/orders/%s", os.Getenv("CLIENT_URL"), o.ID),
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
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": util.CustomValidationErrors(err),
		})
	}

	o, err := oh.orderSvc.FindOrderByID(ctx, c.Params("orderId"))
	if err != nil {
		log.Printf("BuyerRefundingOrder error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Order is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while finding order")
	}

	if o.Status != types.PENDING {
		return fiber.NewError(http.StatusBadRequest, "Buyer cannot refund this order due to order is not in PENDING stage")
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
			newCtx, canc := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer canc()

			cc, err := oh.grpcClient.GetClient(types.USER_SERVICE)
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

			cc, err = oh.grpcClient.GetClient(types.NOTIFICATION_SERVICE)
			if err != nil {
				log.Printf("BuyerRefundingOrder error:\n+%v", err)
				return
			}

			notificationGrpcClient := notification.NewNotificationServiceClient(cc)
			_, err = notificationGrpcClient.BuyerRefundsAnOrder(context.TODO(), &notification.BuyerRefundsOrderRequest{
				ReceiverEmail: b.Email,
				Url:           fmt.Sprintf("%s/orders/%s", os.Getenv("CLIENT_URL"), o.ID),
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

	data.OrderID = o.ID
	err = oh.validate.Struct(data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": util.CustomValidationErrors(err),
		})
	}

	o, err = oh.orderSvc.DeliveringOrder(ctx, *o, *data)
	if err != nil {
		log.Printf("SellerDeliverOrder Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Error while saving the provided data")
	}

	//HACK: SEND EMAIL TO BUYER THAT SELLER HAS SENT ORDER PROGRESS
	// IGNORE ERROR FROM CODE FLOW
	go func() {
		newCtx, canc := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer canc()

		cc, err := oh.grpcClient.GetClient(types.USER_SERVICE)
		if err != nil {
			log.Printf("SellerDeliveringOrder error:\n+%v", err)
		}

		userGrpcClient := user.NewUserServiceClient(cc)
		b, err := userGrpcClient.FindBuyer(newCtx, &user.FindBuyerRequest{
			BuyerId: o.BuyerID,
		})
		if err != nil {
			log.Printf("SellerDeliveringOrder error:\n+%v", err)
			return
		}

		cc, err = oh.grpcClient.GetClient(types.NOTIFICATION_SERVICE)
		if err != nil {
			log.Printf("SellerDeliveringOrder error:\n+%v", err)
			return
		}

		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		_, err = notificationGrpcClient.NotifyBuyerSellerDeliveredOrder(context.TODO(), &notification.NotifyBuyerOrderDeliveredRequest{
			ReceiverEmail: b.Email,
			Url:           fmt.Sprintf("%s/orders/%s", os.Getenv("CLIENT_URL"), o.ID),
		})
		if err != nil {
			log.Printf("SellerDeliveringOrder error:\n+%v", err)
			return
		}
	}()

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"order": o,
	})
}

func (oh *OrderHttpHandler) SellerAcknowledgeOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 500*time.Millisecond)
	defer cancel()

	o, err := oh.orderSvc.FindOrderByID(ctx, c.Params("orderId"))
	if err != nil {
		log.Printf("SellerAcknowledgeOrder Error:\n+%v", err)
		return fiber.NewError(http.StatusNotFound, "Order is not found")
	}

	o, err = oh.orderSvc.ChangeOrderStatus(ctx, *o, types.PROCESS, fmt.Sprintf("Seller Has Acknowledge This Order. Order Status Change To [%s]", types.PROCESS))
	if err != nil {
		log.Printf("SellerAcknowledgeOrder Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Error while saving the provided data")
	}

	//HACK: SEND EMAIL TO BUYER THAT SELLER HAS ACKNOWLEDGE THIS ORDER AND START WORKING
	// IGNORE ERROR FROM CODE FLOW
	go func() {
		newCtx, canc := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer canc()

		cc, err := oh.grpcClient.GetClient(types.USER_SERVICE)
		if err != nil {
			log.Printf("SellerAcknowledgeOrder error:\n+%v", err)
		}

		userGrpcClient := user.NewUserServiceClient(cc)
		b, err := userGrpcClient.FindBuyer(newCtx, &user.FindBuyerRequest{
			BuyerId: o.BuyerID,
		})
		if err != nil {
			log.Printf("SellerAcknowledgeOrder error:\n+%v", err)
			return
		}

		cc, err = oh.grpcClient.GetClient(types.NOTIFICATION_SERVICE)
		if err != nil {
			log.Printf("SellerAcknowledgeOrder error:\n+%v", err)
			return
		}

		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		_, err = notificationGrpcClient.NotifyBuyerOrderHasAcknowledged(context.TODO(), &notification.NotifyBuyerOrderAcknowledgeRequest{
			ReceiverEmail: b.Email,
			Url:           fmt.Sprintf("%s/orders/%s", os.Getenv("CLIENT_URL"), o.ID),
		})
		if err != nil {
			log.Printf("SellerAcknowledgeOrder error:\n+%v", err)
			return
		}
	}()

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"order": o,
	})
}

func (oh *OrderHttpHandler) BuyerResponseForDeliveredOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 500*time.Millisecond)
	defer cancel()

	o, err := oh.orderSvc.FindOrderByID(ctx, c.Params("orderId"))
	if err != nil {
		log.Printf("BuyerResponseForDeliveredOrder Error:\n+%v", err)
		return fiber.NewError(http.StatusNotFound, "Order is not found")
	}

	data := new(types.BuyerResponseOrderDelivered)
	err = c.BodyParser(data)
	if err != nil {
		log.Printf("BuyerResponseForDeliveredOrder Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Invalid Provided Data")
	}

	err = oh.validate.Struct(data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": util.CustomValidationErrors(err),
		})
	}

	o, err = oh.orderSvc.OrderDeliveredResponse(ctx, *o, data)
	if err != nil {
		log.Printf("BuyerResponseForDeliveredOrder Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Error while saving the provided data")
	}

	//HACK: IGNORE ERROR IN CODE FLOW
	go func() {
		newCtx, canc := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer canc()

		cc, err := oh.grpcClient.GetClient(types.NOTIFICATION_SERVICE)
		if err != nil {
			log.Printf("BuyerResponseForDeliveredOrder Error:\n+%v", err)
			return
		}

		userGrpcClient := user.NewUserServiceClient(cc)
		s, err := userGrpcClient.FindSeller(newCtx, &user.FindSellerRequest{
			SellerId: o.SellerID,
		})
		if err != nil {
			log.Printf("BuyerResponseForDeliveredOrder Error:\n+%v", err)
			return
		}

		cc, err = oh.grpcClient.GetClient(types.NOTIFICATION_SERVICE)
		if err != nil {
			log.Printf("BuyerResponseForDeliveredOrder Error:\n+%v", err)
			return
		}

		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		_, err = notificationGrpcClient.NotifySellerBuyerResponseDeliveredOrder(context.TODO(), &notification.NotifySellerBuyerResponseDeliveredOrderRequest{
			ReceiverEmail: s.Email,
		})
		if err != nil {
			log.Printf("BuyerResponseForDeliveredOrder Error:\n+%v", err)
			return
		}
	}()

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"order": o,
	})
}

func (oh *OrderHttpHandler) FindMyOrdersNotifications(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 500*time.Millisecond)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
	}

	orders, err := oh.orderSvc.FindMyOrderNotifications(ctx, userInfo.UserID)
	if err != nil {
		log.Printf("FindMyOrdersNotifications Error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while finding your orders")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"orders": orders,
	})
}

func (oh *OrderHttpHandler) StripeTOSAcceptance(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 500*time.Millisecond)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
	}

	cc, err := oh.grpcClient.GetClient(types.USER_SERVICE)
	if err != nil {
		log.Printf("StripeTOSAcceptance error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	userGrpcClient := user.NewUserServiceClient(cc)
	s, err := userGrpcClient.FindSeller(ctx, &user.FindSellerRequest{
		BuyerId:  userInfo.UserID,
		SellerId: "",
	})
	if err != nil {
		log.Printf("StripeTOSAcceptance error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while finding seller")
	}

	result, err := account.Update(s.StripeAccountId, &stripe.AccountParams{
		TOSAcceptance: &stripe.AccountTOSAcceptanceParams{
			Date: stripe.Int64(time.Now().Unix()),
			IP:   stripe.String(c.IP()),
		},
	})
	if err != nil {
		log.Printf("StripeTOSAcceptance error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while finding seller")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"result": result,
	})
}
