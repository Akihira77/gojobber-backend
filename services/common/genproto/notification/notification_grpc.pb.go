// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: notification.proto

package notification

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	NotificationService_UserVerifyingEmail_FullMethodName                      = "/NotificationService/UserVerifyingEmail"
	NotificationService_UserForgotPassword_FullMethodName                      = "/NotificationService/UserForgotPassword"
	NotificationService_UserSucessResetPassword_FullMethodName                 = "/NotificationService/UserSucessResetPassword"
	NotificationService_SendEmailChatNotification_FullMethodName               = "/NotificationService/SendEmailChatNotification"
	NotificationService_SellerHasCompletedAnOrder_FullMethodName               = "/NotificationService/SellerHasCompletedAnOrder"
	NotificationService_SellerRequestDeadlineExtension_FullMethodName          = "/NotificationService/SellerRequestDeadlineExtension"
	NotificationService_SellerCanceledAnOrder_FullMethodName                   = "/NotificationService/SellerCanceledAnOrder"
	NotificationService_BuyerDeadlineExtensionResponse_FullMethodName          = "/NotificationService/BuyerDeadlineExtensionResponse"
	NotificationService_BuyerRefundsAnOrder_FullMethodName                     = "/NotificationService/BuyerRefundsAnOrder"
	NotificationService_NotifySellerOrderHasBeenMade_FullMethodName            = "/NotificationService/NotifySellerOrderHasBeenMade"
	NotificationService_NotifySellerGotAReview_FullMethodName                  = "/NotificationService/NotifySellerGotAReview"
	NotificationService_NotifyBuyerSellerDeliveredOrder_FullMethodName         = "/NotificationService/NotifyBuyerSellerDeliveredOrder"
	NotificationService_NotifyBuyerOrderHasAcknowledged_FullMethodName         = "/NotificationService/NotifyBuyerOrderHasAcknowledged"
	NotificationService_NotifySellerBuyerResponseDeliveredOrder_FullMethodName = "/NotificationService/NotifySellerBuyerResponseDeliveredOrder"
)

// NotificationServiceClient is the client API for NotificationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotificationServiceClient interface {
	// NOTE: From Auth Service
	UserVerifyingEmail(ctx context.Context, in *VerifyingEmailRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UserForgotPassword(ctx context.Context, in *ForgotPasswordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UserSucessResetPassword(ctx context.Context, in *SuccessResetPasswordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// NOTE: From Chat Service
	SendEmailChatNotification(ctx context.Context, in *EmailChatNotificationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// NOTE: From Order Service
	SellerHasCompletedAnOrder(ctx context.Context, in *SellerCompletedAnOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SellerRequestDeadlineExtension(ctx context.Context, in *SellerDeadlineExtensionRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SellerCanceledAnOrder(ctx context.Context, in *SellerCancelOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	BuyerDeadlineExtensionResponse(ctx context.Context, in *BuyerDeadlineExtension, opts ...grpc.CallOption) (*emptypb.Empty, error)
	BuyerRefundsAnOrder(ctx context.Context, in *BuyerRefundsOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	NotifySellerOrderHasBeenMade(ctx context.Context, in *NotifySellerGotAnOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	NotifySellerGotAReview(ctx context.Context, in *NotifySellerGotAReviewRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	NotifyBuyerSellerDeliveredOrder(ctx context.Context, in *NotifyBuyerOrderDeliveredRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	NotifyBuyerOrderHasAcknowledged(ctx context.Context, in *NotifyBuyerOrderAcknowledgeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	NotifySellerBuyerResponseDeliveredOrder(ctx context.Context, in *NotifySellerBuyerResponseDeliveredOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type notificationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationServiceClient(cc grpc.ClientConnInterface) NotificationServiceClient {
	return &notificationServiceClient{cc}
}

func (c *notificationServiceClient) UserVerifyingEmail(ctx context.Context, in *VerifyingEmailRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, NotificationService_UserVerifyingEmail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) UserForgotPassword(ctx context.Context, in *ForgotPasswordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, NotificationService_UserForgotPassword_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) UserSucessResetPassword(ctx context.Context, in *SuccessResetPasswordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, NotificationService_UserSucessResetPassword_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) SendEmailChatNotification(ctx context.Context, in *EmailChatNotificationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, NotificationService_SendEmailChatNotification_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) SellerHasCompletedAnOrder(ctx context.Context, in *SellerCompletedAnOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, NotificationService_SellerHasCompletedAnOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) SellerRequestDeadlineExtension(ctx context.Context, in *SellerDeadlineExtensionRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, NotificationService_SellerRequestDeadlineExtension_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) SellerCanceledAnOrder(ctx context.Context, in *SellerCancelOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, NotificationService_SellerCanceledAnOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) BuyerDeadlineExtensionResponse(ctx context.Context, in *BuyerDeadlineExtension, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, NotificationService_BuyerDeadlineExtensionResponse_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) BuyerRefundsAnOrder(ctx context.Context, in *BuyerRefundsOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, NotificationService_BuyerRefundsAnOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) NotifySellerOrderHasBeenMade(ctx context.Context, in *NotifySellerGotAnOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, NotificationService_NotifySellerOrderHasBeenMade_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) NotifySellerGotAReview(ctx context.Context, in *NotifySellerGotAReviewRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, NotificationService_NotifySellerGotAReview_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) NotifyBuyerSellerDeliveredOrder(ctx context.Context, in *NotifyBuyerOrderDeliveredRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, NotificationService_NotifyBuyerSellerDeliveredOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) NotifyBuyerOrderHasAcknowledged(ctx context.Context, in *NotifyBuyerOrderAcknowledgeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, NotificationService_NotifyBuyerOrderHasAcknowledged_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) NotifySellerBuyerResponseDeliveredOrder(ctx context.Context, in *NotifySellerBuyerResponseDeliveredOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, NotificationService_NotifySellerBuyerResponseDeliveredOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotificationServiceServer is the server API for NotificationService service.
// All implementations must embed UnimplementedNotificationServiceServer
// for forward compatibility.
type NotificationServiceServer interface {
	// NOTE: From Auth Service
	UserVerifyingEmail(context.Context, *VerifyingEmailRequest) (*emptypb.Empty, error)
	UserForgotPassword(context.Context, *ForgotPasswordRequest) (*emptypb.Empty, error)
	UserSucessResetPassword(context.Context, *SuccessResetPasswordRequest) (*emptypb.Empty, error)
	// NOTE: From Chat Service
	SendEmailChatNotification(context.Context, *EmailChatNotificationRequest) (*emptypb.Empty, error)
	// NOTE: From Order Service
	SellerHasCompletedAnOrder(context.Context, *SellerCompletedAnOrderRequest) (*emptypb.Empty, error)
	SellerRequestDeadlineExtension(context.Context, *SellerDeadlineExtensionRequest) (*emptypb.Empty, error)
	SellerCanceledAnOrder(context.Context, *SellerCancelOrderRequest) (*emptypb.Empty, error)
	BuyerDeadlineExtensionResponse(context.Context, *BuyerDeadlineExtension) (*emptypb.Empty, error)
	BuyerRefundsAnOrder(context.Context, *BuyerRefundsOrderRequest) (*emptypb.Empty, error)
	NotifySellerOrderHasBeenMade(context.Context, *NotifySellerGotAnOrderRequest) (*emptypb.Empty, error)
	NotifySellerGotAReview(context.Context, *NotifySellerGotAReviewRequest) (*emptypb.Empty, error)
	NotifyBuyerSellerDeliveredOrder(context.Context, *NotifyBuyerOrderDeliveredRequest) (*emptypb.Empty, error)
	NotifyBuyerOrderHasAcknowledged(context.Context, *NotifyBuyerOrderAcknowledgeRequest) (*emptypb.Empty, error)
	NotifySellerBuyerResponseDeliveredOrder(context.Context, *NotifySellerBuyerResponseDeliveredOrderRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedNotificationServiceServer()
}

// UnimplementedNotificationServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNotificationServiceServer struct{}

func (UnimplementedNotificationServiceServer) UserVerifyingEmail(context.Context, *VerifyingEmailRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserVerifyingEmail not implemented")
}
func (UnimplementedNotificationServiceServer) UserForgotPassword(context.Context, *ForgotPasswordRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserForgotPassword not implemented")
}
func (UnimplementedNotificationServiceServer) UserSucessResetPassword(context.Context, *SuccessResetPasswordRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserSucessResetPassword not implemented")
}
func (UnimplementedNotificationServiceServer) SendEmailChatNotification(context.Context, *EmailChatNotificationRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendEmailChatNotification not implemented")
}
func (UnimplementedNotificationServiceServer) SellerHasCompletedAnOrder(context.Context, *SellerCompletedAnOrderRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SellerHasCompletedAnOrder not implemented")
}
func (UnimplementedNotificationServiceServer) SellerRequestDeadlineExtension(context.Context, *SellerDeadlineExtensionRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SellerRequestDeadlineExtension not implemented")
}
func (UnimplementedNotificationServiceServer) SellerCanceledAnOrder(context.Context, *SellerCancelOrderRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SellerCanceledAnOrder not implemented")
}
func (UnimplementedNotificationServiceServer) BuyerDeadlineExtensionResponse(context.Context, *BuyerDeadlineExtension) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuyerDeadlineExtensionResponse not implemented")
}
func (UnimplementedNotificationServiceServer) BuyerRefundsAnOrder(context.Context, *BuyerRefundsOrderRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuyerRefundsAnOrder not implemented")
}
func (UnimplementedNotificationServiceServer) NotifySellerOrderHasBeenMade(context.Context, *NotifySellerGotAnOrderRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifySellerOrderHasBeenMade not implemented")
}
func (UnimplementedNotificationServiceServer) NotifySellerGotAReview(context.Context, *NotifySellerGotAReviewRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifySellerGotAReview not implemented")
}
func (UnimplementedNotificationServiceServer) NotifyBuyerSellerDeliveredOrder(context.Context, *NotifyBuyerOrderDeliveredRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyBuyerSellerDeliveredOrder not implemented")
}
func (UnimplementedNotificationServiceServer) NotifyBuyerOrderHasAcknowledged(context.Context, *NotifyBuyerOrderAcknowledgeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyBuyerOrderHasAcknowledged not implemented")
}
func (UnimplementedNotificationServiceServer) NotifySellerBuyerResponseDeliveredOrder(context.Context, *NotifySellerBuyerResponseDeliveredOrderRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifySellerBuyerResponseDeliveredOrder not implemented")
}
func (UnimplementedNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {}
func (UnimplementedNotificationServiceServer) testEmbeddedByValue()                             {}

// UnsafeNotificationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificationServiceServer will
// result in compilation errors.
type UnsafeNotificationServiceServer interface {
	mustEmbedUnimplementedNotificationServiceServer()
}

func RegisterNotificationServiceServer(s grpc.ServiceRegistrar, srv NotificationServiceServer) {
	// If the following call pancis, it indicates UnimplementedNotificationServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&NotificationService_ServiceDesc, srv)
}

func _NotificationService_UserVerifyingEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyingEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).UserVerifyingEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_UserVerifyingEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).UserVerifyingEmail(ctx, req.(*VerifyingEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_UserForgotPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForgotPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).UserForgotPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_UserForgotPassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).UserForgotPassword(ctx, req.(*ForgotPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_UserSucessResetPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SuccessResetPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).UserSucessResetPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_UserSucessResetPassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).UserSucessResetPassword(ctx, req.(*SuccessResetPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_SendEmailChatNotification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailChatNotificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).SendEmailChatNotification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_SendEmailChatNotification_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).SendEmailChatNotification(ctx, req.(*EmailChatNotificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_SellerHasCompletedAnOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SellerCompletedAnOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).SellerHasCompletedAnOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_SellerHasCompletedAnOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).SellerHasCompletedAnOrder(ctx, req.(*SellerCompletedAnOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_SellerRequestDeadlineExtension_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SellerDeadlineExtensionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).SellerRequestDeadlineExtension(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_SellerRequestDeadlineExtension_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).SellerRequestDeadlineExtension(ctx, req.(*SellerDeadlineExtensionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_SellerCanceledAnOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SellerCancelOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).SellerCanceledAnOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_SellerCanceledAnOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).SellerCanceledAnOrder(ctx, req.(*SellerCancelOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_BuyerDeadlineExtensionResponse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BuyerDeadlineExtension)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).BuyerDeadlineExtensionResponse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_BuyerDeadlineExtensionResponse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).BuyerDeadlineExtensionResponse(ctx, req.(*BuyerDeadlineExtension))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_BuyerRefundsAnOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BuyerRefundsOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).BuyerRefundsAnOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_BuyerRefundsAnOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).BuyerRefundsAnOrder(ctx, req.(*BuyerRefundsOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_NotifySellerOrderHasBeenMade_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifySellerGotAnOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).NotifySellerOrderHasBeenMade(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_NotifySellerOrderHasBeenMade_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).NotifySellerOrderHasBeenMade(ctx, req.(*NotifySellerGotAnOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_NotifySellerGotAReview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifySellerGotAReviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).NotifySellerGotAReview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_NotifySellerGotAReview_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).NotifySellerGotAReview(ctx, req.(*NotifySellerGotAReviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_NotifyBuyerSellerDeliveredOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyBuyerOrderDeliveredRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).NotifyBuyerSellerDeliveredOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_NotifyBuyerSellerDeliveredOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).NotifyBuyerSellerDeliveredOrder(ctx, req.(*NotifyBuyerOrderDeliveredRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_NotifyBuyerOrderHasAcknowledged_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyBuyerOrderAcknowledgeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).NotifyBuyerOrderHasAcknowledged(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_NotifyBuyerOrderHasAcknowledged_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).NotifyBuyerOrderHasAcknowledged(ctx, req.(*NotifyBuyerOrderAcknowledgeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_NotifySellerBuyerResponseDeliveredOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifySellerBuyerResponseDeliveredOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).NotifySellerBuyerResponseDeliveredOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_NotifySellerBuyerResponseDeliveredOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).NotifySellerBuyerResponseDeliveredOrder(ctx, req.(*NotifySellerBuyerResponseDeliveredOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NotificationService_ServiceDesc is the grpc.ServiceDesc for NotificationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NotificationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "NotificationService",
	HandlerType: (*NotificationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserVerifyingEmail",
			Handler:    _NotificationService_UserVerifyingEmail_Handler,
		},
		{
			MethodName: "UserForgotPassword",
			Handler:    _NotificationService_UserForgotPassword_Handler,
		},
		{
			MethodName: "UserSucessResetPassword",
			Handler:    _NotificationService_UserSucessResetPassword_Handler,
		},
		{
			MethodName: "SendEmailChatNotification",
			Handler:    _NotificationService_SendEmailChatNotification_Handler,
		},
		{
			MethodName: "SellerHasCompletedAnOrder",
			Handler:    _NotificationService_SellerHasCompletedAnOrder_Handler,
		},
		{
			MethodName: "SellerRequestDeadlineExtension",
			Handler:    _NotificationService_SellerRequestDeadlineExtension_Handler,
		},
		{
			MethodName: "SellerCanceledAnOrder",
			Handler:    _NotificationService_SellerCanceledAnOrder_Handler,
		},
		{
			MethodName: "BuyerDeadlineExtensionResponse",
			Handler:    _NotificationService_BuyerDeadlineExtensionResponse_Handler,
		},
		{
			MethodName: "BuyerRefundsAnOrder",
			Handler:    _NotificationService_BuyerRefundsAnOrder_Handler,
		},
		{
			MethodName: "NotifySellerOrderHasBeenMade",
			Handler:    _NotificationService_NotifySellerOrderHasBeenMade_Handler,
		},
		{
			MethodName: "NotifySellerGotAReview",
			Handler:    _NotificationService_NotifySellerGotAReview_Handler,
		},
		{
			MethodName: "NotifyBuyerSellerDeliveredOrder",
			Handler:    _NotificationService_NotifyBuyerSellerDeliveredOrder_Handler,
		},
		{
			MethodName: "NotifyBuyerOrderHasAcknowledged",
			Handler:    _NotificationService_NotifyBuyerOrderHasAcknowledged_Handler,
		},
		{
			MethodName: "NotifySellerBuyerResponseDeliveredOrder",
			Handler:    _NotificationService_NotifySellerBuyerResponseDeliveredOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "notification.proto",
}
