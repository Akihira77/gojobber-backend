package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/Akihira77/gojobber/services/6-chat/service"
	"github.com/Akihira77/gojobber/services/6-chat/types"
	"github.com/Akihira77/gojobber/services/common/genproto/chat"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ChatGRPCHandler struct {
	chatSvc service.ChatServiceImpl
	chat.UnimplementedChatServiceServer
}

func NewChatGRPCHandler(grpc *grpc.Server, chatSvc service.ChatServiceImpl) {
	gRPCHandler := &ChatGRPCHandler{
		chatSvc: chatSvc,
	}

	chat.RegisterChatServiceServer(grpc, gRPCHandler)
}

func (ch *ChatGRPCHandler) BuyerAcceptedOffer(ctx context.Context, req *chat.BuyerAcceptedOfferRequest) (*emptypb.Empty, error) {
	log.Println("Receive message", req)

	m, err := ch.chatSvc.FindMessageByID(ctx, req.MessageId)
	if err != nil {
		log.Printf("Message is not found:\n+%v", err)
		return nil, fmt.Errorf("Message is not found")
	}

	err = ch.chatSvc.ChangeOfferStatus(ctx, m, types.ACCEPTED)
	return nil, err
}
