// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.12
// source: notification.proto

package notification

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// INFO: AUTH SERVICE
type VerifyingEmailRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReceiverEmail    string `protobuf:"bytes,1,opt,name=receiverEmail,proto3" json:"receiverEmail,omitempty"`
	HtmlTemplateName string `protobuf:"bytes,2,opt,name=htmlTemplateName,proto3" json:"htmlTemplateName,omitempty"`
	VerifyLink       string `protobuf:"bytes,3,opt,name=verifyLink,proto3" json:"verifyLink,omitempty"`
}

func (x *VerifyingEmailRequest) Reset() {
	*x = VerifyingEmailRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notification_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyingEmailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyingEmailRequest) ProtoMessage() {}

func (x *VerifyingEmailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_notification_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyingEmailRequest.ProtoReflect.Descriptor instead.
func (*VerifyingEmailRequest) Descriptor() ([]byte, []int) {
	return file_notification_proto_rawDescGZIP(), []int{0}
}

func (x *VerifyingEmailRequest) GetReceiverEmail() string {
	if x != nil {
		return x.ReceiverEmail
	}
	return ""
}

func (x *VerifyingEmailRequest) GetHtmlTemplateName() string {
	if x != nil {
		return x.HtmlTemplateName
	}
	return ""
}

func (x *VerifyingEmailRequest) GetVerifyLink() string {
	if x != nil {
		return x.VerifyLink
	}
	return ""
}

type ForgotPasswordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReceiverEmail    string `protobuf:"bytes,1,opt,name=receiverEmail,proto3" json:"receiverEmail,omitempty"`
	HtmlTemplateName string `protobuf:"bytes,2,opt,name=htmlTemplateName,proto3" json:"htmlTemplateName,omitempty"`
	ResetLink        string `protobuf:"bytes,3,opt,name=resetLink,proto3" json:"resetLink,omitempty"`
	Username         string `protobuf:"bytes,4,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *ForgotPasswordRequest) Reset() {
	*x = ForgotPasswordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notification_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ForgotPasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForgotPasswordRequest) ProtoMessage() {}

func (x *ForgotPasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_notification_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForgotPasswordRequest.ProtoReflect.Descriptor instead.
func (*ForgotPasswordRequest) Descriptor() ([]byte, []int) {
	return file_notification_proto_rawDescGZIP(), []int{1}
}

func (x *ForgotPasswordRequest) GetReceiverEmail() string {
	if x != nil {
		return x.ReceiverEmail
	}
	return ""
}

func (x *ForgotPasswordRequest) GetHtmlTemplateName() string {
	if x != nil {
		return x.HtmlTemplateName
	}
	return ""
}

func (x *ForgotPasswordRequest) GetResetLink() string {
	if x != nil {
		return x.ResetLink
	}
	return ""
}

func (x *ForgotPasswordRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

type SuccessResetPasswordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReceiverEmail    string `protobuf:"bytes,1,opt,name=receiverEmail,proto3" json:"receiverEmail,omitempty"`
	HtmlTemplateName string `protobuf:"bytes,2,opt,name=htmlTemplateName,proto3" json:"htmlTemplateName,omitempty"`
	Username         string `protobuf:"bytes,4,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *SuccessResetPasswordRequest) Reset() {
	*x = SuccessResetPasswordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notification_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SuccessResetPasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SuccessResetPasswordRequest) ProtoMessage() {}

func (x *SuccessResetPasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_notification_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SuccessResetPasswordRequest.ProtoReflect.Descriptor instead.
func (*SuccessResetPasswordRequest) Descriptor() ([]byte, []int) {
	return file_notification_proto_rawDescGZIP(), []int{2}
}

func (x *SuccessResetPasswordRequest) GetReceiverEmail() string {
	if x != nil {
		return x.ReceiverEmail
	}
	return ""
}

func (x *SuccessResetPasswordRequest) GetHtmlTemplateName() string {
	if x != nil {
		return x.HtmlTemplateName
	}
	return ""
}

func (x *SuccessResetPasswordRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

// INFO: CHAT SERVICE
type EmailChatNotificationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReceiverEmail string `protobuf:"bytes,1,opt,name=receiverEmail,proto3" json:"receiverEmail,omitempty"`
	SenderEmail   string `protobuf:"bytes,2,opt,name=senderEmail,proto3" json:"senderEmail,omitempty"`
	Message       string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *EmailChatNotificationRequest) Reset() {
	*x = EmailChatNotificationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notification_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmailChatNotificationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmailChatNotificationRequest) ProtoMessage() {}

func (x *EmailChatNotificationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_notification_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmailChatNotificationRequest.ProtoReflect.Descriptor instead.
func (*EmailChatNotificationRequest) Descriptor() ([]byte, []int) {
	return file_notification_proto_rawDescGZIP(), []int{3}
}

func (x *EmailChatNotificationRequest) GetReceiverEmail() string {
	if x != nil {
		return x.ReceiverEmail
	}
	return ""
}

func (x *EmailChatNotificationRequest) GetSenderEmail() string {
	if x != nil {
		return x.SenderEmail
	}
	return ""
}

func (x *EmailChatNotificationRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// INFO: ORDER SERVICE
type SellerCompletedAnOrderRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReceiverEmail        string `protobuf:"bytes,1,opt,name=receiverEmail,proto3" json:"receiverEmail,omitempty"`
	BuyerEmail           string `protobuf:"bytes,2,opt,name=buyerEmail,proto3" json:"buyerEmail,omitempty"`
	OrderId              string `protobuf:"bytes,3,opt,name=orderId,proto3" json:"orderId,omitempty"`
	SellerCurrentBalance string `protobuf:"bytes,4,opt,name=sellerCurrentBalance,proto3" json:"sellerCurrentBalance,omitempty"`
}

func (x *SellerCompletedAnOrderRequest) Reset() {
	*x = SellerCompletedAnOrderRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notification_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SellerCompletedAnOrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SellerCompletedAnOrderRequest) ProtoMessage() {}

func (x *SellerCompletedAnOrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_notification_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SellerCompletedAnOrderRequest.ProtoReflect.Descriptor instead.
func (*SellerCompletedAnOrderRequest) Descriptor() ([]byte, []int) {
	return file_notification_proto_rawDescGZIP(), []int{4}
}

func (x *SellerCompletedAnOrderRequest) GetReceiverEmail() string {
	if x != nil {
		return x.ReceiverEmail
	}
	return ""
}

func (x *SellerCompletedAnOrderRequest) GetBuyerEmail() string {
	if x != nil {
		return x.BuyerEmail
	}
	return ""
}

func (x *SellerCompletedAnOrderRequest) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *SellerCompletedAnOrderRequest) GetSellerCurrentBalance() string {
	if x != nil {
		return x.SellerCurrentBalance
	}
	return ""
}

var File_notification_proto protoreflect.FileDescriptor

var file_notification_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x89, 0x01, 0x0a, 0x15, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x69, 0x6e, 0x67, 0x45,
	0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x72,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69,
	0x6c, 0x12, 0x2a, 0x0a, 0x10, 0x68, 0x74, 0x6d, 0x6c, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x68, 0x74, 0x6d,
	0x6c, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a,
	0x0a, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x4c, 0x69, 0x6e, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x4c, 0x69, 0x6e, 0x6b, 0x22, 0xa3, 0x01,
	0x0a, 0x15, 0x46, 0x6f, 0x72, 0x67, 0x6f, 0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x72, 0x65, 0x63, 0x65, 0x69,
	0x76, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x2a, 0x0a,
	0x10, 0x68, 0x74, 0x6d, 0x6c, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x68, 0x74, 0x6d, 0x6c, 0x54, 0x65, 0x6d,
	0x70, 0x6c, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x73,
	0x65, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65,
	0x73, 0x65, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x22, 0x8b, 0x01, 0x0a, 0x1b, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52,
	0x65, 0x73, 0x65, 0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x45,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x72, 0x65, 0x63, 0x65,
	0x69, 0x76, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x2a, 0x0a, 0x10, 0x68, 0x74, 0x6d,
	0x6c, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x10, 0x68, 0x74, 0x6d, 0x6c, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0x80, 0x01, 0x0a, 0x1c, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x43, 0x68, 0x61, 0x74, 0x4e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x45, 0x6d,
	0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x72, 0x65, 0x63, 0x65, 0x69,
	0x76, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x20, 0x0a, 0x0b, 0x73, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0xb3, 0x01, 0x0a, 0x1d, 0x53, 0x65, 0x6c, 0x6c, 0x65, 0x72, 0x43,
	0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x6e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76,
	0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x72,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1e, 0x0a, 0x0a,
	0x62, 0x75, 0x79, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x62, 0x75, 0x79, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x18, 0x0a, 0x07,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x32, 0x0a, 0x14, 0x73, 0x65, 0x6c, 0x6c, 0x65, 0x72,
	0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x14, 0x73, 0x65, 0x6c, 0x6c, 0x65, 0x72, 0x43, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x74, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x32, 0xa5, 0x03, 0x0a, 0x13, 0x4e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x46, 0x0a, 0x12, 0x55, 0x73, 0x65, 0x72, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79,
	0x69, 0x6e, 0x67, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x16, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66,
	0x79, 0x69, 0x6e, 0x67, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x12, 0x55, 0x73,
	0x65, 0x72, 0x46, 0x6f, 0x72, 0x67, 0x6f, 0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x12, 0x16, 0x2e, 0x46, 0x6f, 0x72, 0x67, 0x6f, 0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x51, 0x0a, 0x17, 0x55, 0x73, 0x65, 0x72, 0x53, 0x75, 0x63, 0x65, 0x73, 0x73,
	0x52, 0x65, 0x73, 0x65, 0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x1c, 0x2e,
	0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x65, 0x74, 0x50, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x54, 0x0a, 0x19, 0x53, 0x65, 0x6e, 0x64, 0x45, 0x6d, 0x61,
	0x69, 0x6c, 0x43, 0x68, 0x61, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1d, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x43, 0x68, 0x61, 0x74, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x55, 0x0a, 0x19, 0x53,
	0x65, 0x6c, 0x6c, 0x65, 0x72, 0x48, 0x61, 0x73, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65,
	0x64, 0x41, 0x6e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x1e, 0x2e, 0x53, 0x65, 0x6c, 0x6c, 0x65,
	0x72, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x6e, 0x4f, 0x72, 0x64, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x41, 0x6b, 0x69, 0x68, 0x69, 0x72, 0x61, 0x37, 0x37, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_notification_proto_rawDescOnce sync.Once
	file_notification_proto_rawDescData = file_notification_proto_rawDesc
)

func file_notification_proto_rawDescGZIP() []byte {
	file_notification_proto_rawDescOnce.Do(func() {
		file_notification_proto_rawDescData = protoimpl.X.CompressGZIP(file_notification_proto_rawDescData)
	})
	return file_notification_proto_rawDescData
}

var file_notification_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_notification_proto_goTypes = []any{
	(*VerifyingEmailRequest)(nil),         // 0: VerifyingEmailRequest
	(*ForgotPasswordRequest)(nil),         // 1: ForgotPasswordRequest
	(*SuccessResetPasswordRequest)(nil),   // 2: SuccessResetPasswordRequest
	(*EmailChatNotificationRequest)(nil),  // 3: EmailChatNotificationRequest
	(*SellerCompletedAnOrderRequest)(nil), // 4: SellerCompletedAnOrderRequest
	(*emptypb.Empty)(nil),                 // 5: google.protobuf.Empty
}
var file_notification_proto_depIdxs = []int32{
	0, // 0: NotificationService.UserVerifyingEmail:input_type -> VerifyingEmailRequest
	1, // 1: NotificationService.UserForgotPassword:input_type -> ForgotPasswordRequest
	2, // 2: NotificationService.UserSucessResetPassword:input_type -> SuccessResetPasswordRequest
	3, // 3: NotificationService.SendEmailChatNotification:input_type -> EmailChatNotificationRequest
	4, // 4: NotificationService.SellerHasCompletedAnOrder:input_type -> SellerCompletedAnOrderRequest
	5, // 5: NotificationService.UserVerifyingEmail:output_type -> google.protobuf.Empty
	5, // 6: NotificationService.UserForgotPassword:output_type -> google.protobuf.Empty
	5, // 7: NotificationService.UserSucessResetPassword:output_type -> google.protobuf.Empty
	5, // 8: NotificationService.SendEmailChatNotification:output_type -> google.protobuf.Empty
	5, // 9: NotificationService.SellerHasCompletedAnOrder:output_type -> google.protobuf.Empty
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_notification_proto_init() }
func file_notification_proto_init() {
	if File_notification_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_notification_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*VerifyingEmailRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_notification_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ForgotPasswordRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_notification_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*SuccessResetPasswordRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_notification_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*EmailChatNotificationRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_notification_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*SellerCompletedAnOrderRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_notification_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_notification_proto_goTypes,
		DependencyIndexes: file_notification_proto_depIdxs,
		MessageInfos:      file_notification_proto_msgTypes,
	}.Build()
	File_notification_proto = out.File
	file_notification_proto_rawDesc = nil
	file_notification_proto_goTypes = nil
	file_notification_proto_depIdxs = nil
}
