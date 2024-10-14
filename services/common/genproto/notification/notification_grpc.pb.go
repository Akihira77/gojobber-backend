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
	NotificationService_UserVerifyingEmail_FullMethodName             = "/NotificationService/UserVerifyingEmail"
	NotificationService_UserForgotPassword_FullMethodName             = "/NotificationService/UserForgotPassword"
	NotificationService_UserSucessResetPassword_FullMethodName        = "/NotificationService/UserSucessResetPassword"
	NotificationService_SendEmailChatNotification_FullMethodName      = "/NotificationService/SendEmailChatNotification"
	NotificationService_SellerHasCompletedAnOrder_FullMethodName      = "/NotificationService/SellerHasCompletedAnOrder"
	NotificationService_SellerRequestDeadlineExtension_FullMethodName = "/NotificationService/SellerRequestDeadlineExtension"
)

// NotificationServiceClient is the client API for NotificationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotificationServiceClient interface {
	// From Auth Service
	UserVerifyingEmail(ctx context.Context, in *VerifyingEmailRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UserForgotPassword(ctx context.Context, in *ForgotPasswordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UserSucessResetPassword(ctx context.Context, in *SuccessResetPasswordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// From Chat Service
	SendEmailChatNotification(ctx context.Context, in *EmailChatNotificationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// From Order Service
	SellerHasCompletedAnOrder(ctx context.Context, in *SellerCompletedAnOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SellerRequestDeadlineExtension(ctx context.Context, in *SellerDeadlineExtensionRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
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

// NotificationServiceServer is the server API for NotificationService service.
// All implementations must embed UnimplementedNotificationServiceServer
// for forward compatibility.
type NotificationServiceServer interface {
	// From Auth Service
	UserVerifyingEmail(context.Context, *VerifyingEmailRequest) (*emptypb.Empty, error)
	UserForgotPassword(context.Context, *ForgotPasswordRequest) (*emptypb.Empty, error)
	UserSucessResetPassword(context.Context, *SuccessResetPasswordRequest) (*emptypb.Empty, error)
	// From Chat Service
	SendEmailChatNotification(context.Context, *EmailChatNotificationRequest) (*emptypb.Empty, error)
	// From Order Service
	SellerHasCompletedAnOrder(context.Context, *SellerCompletedAnOrderRequest) (*emptypb.Empty, error)
	SellerRequestDeadlineExtension(context.Context, *SellerDeadlineExtensionRequest) (*emptypb.Empty, error)
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
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "notification.proto",
}
