package handler

import (
	"context"
	"log"

	"github.com/Akihira77/gojobber/services/2-notification/service"
	"github.com/Akihira77/gojobber/services/common/genproto/notification"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type NotificationGRPCHandler struct {
	notificationSvc service.NotificationServiceImpl
	notification.UnimplementedNotificationServiceServer
}

func NewNotificationGRPCHandler(grpc *grpc.Server, notificationSvc service.NotificationServiceImpl) {
	gRPCHandler := &NotificationGRPCHandler{
		notificationSvc: notificationSvc,
	}

	notification.RegisterNotificationServiceServer(grpc, gRPCHandler)
}

func (h *NotificationGRPCHandler) UserVerifyingEmail(ctx context.Context, req *notification.VerifyingEmailRequest) (*emptypb.Empty, error) {
	log.Println("Receiving data", req)
	err := h.notificationSvc.UserVerifyingEmail(req.ReceiverEmail, req.HtmlTemplateName, req.VerifyLink)

	if err != nil {
		log.Printf("UserVerifyingEmail for [%s] is error: %v", req.ReceiverEmail, err)
	}
	return nil, err
}

func (h *NotificationGRPCHandler) UserForgotPassword(ctx context.Context, req *notification.ForgotPasswordRequest) (*emptypb.Empty, error) {
	log.Println("Receiving data", req)
	err := h.notificationSvc.UserForgotPassword(req.ReceiverEmail, req.HtmlTemplateName, req.ResetLink, req.Username)

	if err != nil {
		log.Printf("UserForgotPassword for [%s] is error: %v", req.ReceiverEmail, err)
	}
	return nil, err
}

func (h *NotificationGRPCHandler) UserSuccessResetPassword(ctx context.Context, req *notification.SuccessResetPasswordRequest) (*emptypb.Empty, error) {
	log.Println("Receiving data", req)
	err := h.notificationSvc.UserSucessResetPassword(req.ReceiverEmail, req.HtmlTemplateName, req.Username)

	if err != nil {
		log.Printf("UserSuccessResetPassword for [%s] is error: %v", req.ReceiverEmail, err)
	}
	return nil, err
}

func (h *NotificationGRPCHandler) SendEmailChatNotification(ctx context.Context, req *notification.EmailChatNotificationRequest) (*emptypb.Empty, error) {
	log.Println("Receiving data", req)
	err := h.notificationSvc.SendEmailChatNotification(req.ReceiverEmail, req.SenderEmail, req.Message)

	if err != nil {
		log.Printf("SendEmailChatNotificaion for [%s] is error: %v", req.ReceiverEmail, err)
	}
	return nil, err
}

func (h *NotificationGRPCHandler) SellerHasCompletedAnOrder(ctx context.Context, req *notification.SellerCompletedAnOrderRequest) (*emptypb.Empty, error) {
	log.Println("Receiving data", req)
	err := h.notificationSvc.SellerHasCompletedAnOrder(req.ReceiverEmail, req.BuyerEmail, req.OrderId, req.SellerCurrentBalance)

	if err != nil {
		log.Printf("SellerHasCompletedAnOrder for [%s] is error: %v", req.ReceiverEmail, err)
	}
	return nil, err
}

func (h *NotificationGRPCHandler) BuyerDeadlineExtensionResponse(ctx context.Context, req *notification.BuyerDeadlineExtension) (*emptypb.Empty, error) {
	log.Println("Receiving data", req)
	err := h.notificationSvc.BuyerDeadlineExtensionResponse(req.ReceiverEmail, req.Message)

	if err != nil {
		log.Printf("BuyerDeadlineExtensionResponse for [%s] is error: %v", req.ReceiverEmail, err)
	}
	return nil, err
}

func (h *NotificationGRPCHandler) BuyerRefundsAnOrder(ctx context.Context, req *notification.BuyerRefundsOrderRequest) (*emptypb.Empty, error) {
	log.Println("Receiving data", req)
	err := h.notificationSvc.BuyerRefundsAnOrder(req.ReceiverEmail, req.Reason)

	if err != nil {
		log.Printf("BuyerRefundsAnOrder for [%s] is error: %v", req.ReceiverEmail, err)
	}
	return nil, err
}

func (h *NotificationGRPCHandler) SellerCanceledAnOrder(ctx context.Context, req *notification.SellerCancelOrderRequest) (*emptypb.Empty, error) {
	log.Println("Receiving data", req)
	err := h.notificationSvc.SellerCanceledAnOrder(req.ReceiverEmail, req.Reason)

	if err != nil {
		log.Printf("SellerCanceledAnOrder for [%s] is error: %v", req.ReceiverEmail, err)
	}
	return nil, err
}

func (h *NotificationGRPCHandler) NotifySellerOrderHasBeenMade(ctx context.Context, req *notification.NotifySellerGotAnOrderRequest) (*emptypb.Empty, error) {
	log.Println("Receiving data", req)
	err := h.notificationSvc.NotifySellerGotAnOrder(req)

	if err != nil {
		log.Printf("NotifySellerOrderHasBeenMade for [%s] is error: %v", req.ReceiverEmail, err)
	}
	return nil, err
}

func (h *NotificationGRPCHandler) otifySellerGotAReview(ctx context.Context, req *notification.NotifySellerGotAReviewRequest) (*emptypb.Empty, error) {
	log.Println("Receiving data", req)
	err := h.notificationSvc.NotifySellerGotAReview(req)

	if err != nil {
		log.Printf("NotifySellerGotAReview for [%s] is error: %v", req.ReceiverEmail, err)
	}
	return nil, err
}
