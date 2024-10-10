// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.12
// source: user.proto

package user

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SaveBuyerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Username       string                 `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Email          string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Country        string                 `protobuf:"bytes,4,opt,name=country,proto3" json:"country,omitempty"`
	ProfilePicture string                 `protobuf:"bytes,5,opt,name=profilePicture,proto3" json:"profilePicture,omitempty"`
	IsSeller       bool                   `protobuf:"varint,6,opt,name=isSeller,proto3" json:"isSeller,omitempty"`
	CreatedAt      *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
}

func (x *SaveBuyerRequest) Reset() {
	*x = SaveBuyerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveBuyerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveBuyerRequest) ProtoMessage() {}

func (x *SaveBuyerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveBuyerRequest.ProtoReflect.Descriptor instead.
func (*SaveBuyerRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{0}
}

func (x *SaveBuyerRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SaveBuyerRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *SaveBuyerRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *SaveBuyerRequest) GetCountry() string {
	if x != nil {
		return x.Country
	}
	return ""
}

func (x *SaveBuyerRequest) GetProfilePicture() string {
	if x != nil {
		return x.ProfilePicture
	}
	return ""
}

func (x *SaveBuyerRequest) GetIsSeller() bool {
	if x != nil {
		return x.IsSeller
	}
	return false
}

func (x *SaveBuyerRequest) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type SaveBuyerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *SaveBuyerResponse) Reset() {
	*x = SaveBuyerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveBuyerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveBuyerResponse) ProtoMessage() {}

func (x *SaveBuyerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveBuyerResponse.ProtoReflect.Descriptor instead.
func (*SaveBuyerResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{1}
}

func (x *SaveBuyerResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *SaveBuyerResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type FindSellerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SellerId string `protobuf:"bytes,1,opt,name=sellerId,proto3" json:"sellerId,omitempty"`
	BuyerId  string `protobuf:"bytes,2,opt,name=buyerId,proto3" json:"buyerId,omitempty"`
}

func (x *FindSellerRequest) Reset() {
	*x = FindSellerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindSellerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindSellerRequest) ProtoMessage() {}

func (x *FindSellerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindSellerRequest.ProtoReflect.Descriptor instead.
func (*FindSellerRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{2}
}

func (x *FindSellerRequest) GetSellerId() string {
	if x != nil {
		return x.SellerId
	}
	return ""
}

func (x *FindSellerRequest) GetBuyerId() string {
	if x != nil {
		return x.BuyerId
	}
	return ""
}

type FindSellerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id               string          `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FullName         string          `protobuf:"bytes,2,opt,name=fullName,proto3" json:"fullName,omitempty"`
	RatingsCount     int64           `protobuf:"varint,3,opt,name=ratingsCount,proto3" json:"ratingsCount,omitempty"`
	RatingSum        int64           `protobuf:"varint,4,opt,name=ratingSum,proto3" json:"ratingSum,omitempty"`
	RatingCategories *RatingCategory `protobuf:"bytes,5,opt,name=ratingCategories,proto3" json:"ratingCategories,omitempty"`
}

func (x *FindSellerResponse) Reset() {
	*x = FindSellerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindSellerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindSellerResponse) ProtoMessage() {}

func (x *FindSellerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindSellerResponse.ProtoReflect.Descriptor instead.
func (*FindSellerResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{3}
}

func (x *FindSellerResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *FindSellerResponse) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

func (x *FindSellerResponse) GetRatingsCount() int64 {
	if x != nil {
		return x.RatingsCount
	}
	return 0
}

func (x *FindSellerResponse) GetRatingSum() int64 {
	if x != nil {
		return x.RatingSum
	}
	return 0
}

func (x *FindSellerResponse) GetRatingCategories() *RatingCategory {
	if x != nil {
		return x.RatingCategories
	}
	return nil
}

type RatingCategory struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Five  int32 `protobuf:"varint,1,opt,name=five,proto3" json:"five,omitempty"`
	Four  int32 `protobuf:"varint,2,opt,name=four,proto3" json:"four,omitempty"`
	Three int32 `protobuf:"varint,3,opt,name=three,proto3" json:"three,omitempty"`
	Two   int32 `protobuf:"varint,4,opt,name=two,proto3" json:"two,omitempty"`
	One   int32 `protobuf:"varint,5,opt,name=one,proto3" json:"one,omitempty"`
}

func (x *RatingCategory) Reset() {
	*x = RatingCategory{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RatingCategory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RatingCategory) ProtoMessage() {}

func (x *RatingCategory) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RatingCategory.ProtoReflect.Descriptor instead.
func (*RatingCategory) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{4}
}

func (x *RatingCategory) GetFive() int32 {
	if x != nil {
		return x.Five
	}
	return 0
}

func (x *RatingCategory) GetFour() int32 {
	if x != nil {
		return x.Four
	}
	return 0
}

func (x *RatingCategory) GetThree() int32 {
	if x != nil {
		return x.Three
	}
	return 0
}

func (x *RatingCategory) GetTwo() int32 {
	if x != nil {
		return x.Two
	}
	return 0
}

func (x *RatingCategory) GetOne() int32 {
	if x != nil {
		return x.One
	}
	return 0
}

var File_user_proto protoreflect.FileDescriptor

var file_user_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xec, 0x01,
	0x0a, 0x10, 0x53, 0x61, 0x76, 0x65, 0x42, 0x75, 0x79, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x26,
	0x0a, 0x0e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x69, 0x63, 0x74, 0x75, 0x72, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x50,
	0x69, 0x63, 0x74, 0x75, 0x72, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x73, 0x53, 0x65, 0x6c, 0x6c,
	0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x53, 0x65, 0x6c, 0x6c,
	0x65, 0x72, 0x12, 0x38, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x47, 0x0a, 0x11,
	0x53, 0x61, 0x76, 0x65, 0x42, 0x75, 0x79, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x49, 0x0a, 0x11, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x65, 0x6c,
	0x6c, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65,
	0x6c, 0x6c, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65,
	0x6c, 0x6c, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x62, 0x75, 0x79, 0x65, 0x72, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x62, 0x75, 0x79, 0x65, 0x72, 0x49, 0x64,
	0x22, 0xbf, 0x01, 0x0a, 0x12, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x65, 0x6c, 0x6c, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x72, 0x61, 0x74, 0x69, 0x6e,
	0x67, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x61, 0x74, 0x69, 0x6e,
	0x67, 0x53, 0x75, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x72, 0x61, 0x74, 0x69,
	0x6e, 0x67, 0x53, 0x75, 0x6d, 0x12, 0x3b, 0x0a, 0x10, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x43,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0f, 0x2e, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79,
	0x52, 0x10, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69,
	0x65, 0x73, 0x22, 0x72, 0x0a, 0x0e, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x43, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x69, 0x76, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x66, 0x69, 0x76, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x6f, 0x75, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x66, 0x6f, 0x75, 0x72, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x68, 0x72, 0x65, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x74, 0x68, 0x72,
	0x65, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x77, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x03, 0x74, 0x77, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x6f, 0x6e, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x03, 0x6f, 0x6e, 0x65, 0x32, 0x80, 0x01, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x38, 0x0a, 0x0d, 0x53, 0x61, 0x76, 0x65, 0x42, 0x75,
	0x79, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x11, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x42, 0x75,
	0x79, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x53, 0x61, 0x76,
	0x65, 0x42, 0x75, 0x79, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x37, 0x0a, 0x0a, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x65, 0x6c, 0x6c, 0x65, 0x72, 0x12, 0x12,
	0x2e, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x65, 0x6c, 0x6c, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x13, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x65, 0x6c, 0x6c, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x22, 0x5a, 0x20, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x41, 0x6b, 0x69, 0x68, 0x69, 0x72, 0x61, 0x37,
	0x37, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_user_proto_rawDescOnce sync.Once
	file_user_proto_rawDescData = file_user_proto_rawDesc
)

func file_user_proto_rawDescGZIP() []byte {
	file_user_proto_rawDescOnce.Do(func() {
		file_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_proto_rawDescData)
	})
	return file_user_proto_rawDescData
}

var file_user_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_user_proto_goTypes = []any{
	(*SaveBuyerRequest)(nil),      // 0: SaveBuyerRequest
	(*SaveBuyerResponse)(nil),     // 1: SaveBuyerResponse
	(*FindSellerRequest)(nil),     // 2: FindSellerRequest
	(*FindSellerResponse)(nil),    // 3: FindSellerResponse
	(*RatingCategory)(nil),        // 4: RatingCategory
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
}
var file_user_proto_depIdxs = []int32{
	5, // 0: SaveBuyerRequest.createdAt:type_name -> google.protobuf.Timestamp
	4, // 1: FindSellerResponse.ratingCategories:type_name -> RatingCategory
	0, // 2: UserService.SaveBuyerData:input_type -> SaveBuyerRequest
	2, // 3: UserService.FindSeller:input_type -> FindSellerRequest
	1, // 4: UserService.SaveBuyerData:output_type -> SaveBuyerResponse
	3, // 5: UserService.FindSeller:output_type -> FindSellerResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_user_proto_init() }
func file_user_proto_init() {
	if File_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_user_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*SaveBuyerRequest); i {
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
		file_user_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*SaveBuyerResponse); i {
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
		file_user_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*FindSellerRequest); i {
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
		file_user_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*FindSellerResponse); i {
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
		file_user_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*RatingCategory); i {
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
			RawDescriptor: file_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_proto_goTypes,
		DependencyIndexes: file_user_proto_depIdxs,
		MessageInfos:      file_user_proto_msgTypes,
	}.Build()
	File_user_proto = out.File
	file_user_proto_rawDesc = nil
	file_user_proto_goTypes = nil
	file_user_proto_depIdxs = nil
}
