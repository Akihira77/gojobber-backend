package main

import (
	"log"
	"net"

	"github.com/Akihira77/gojobber/services/6-chat/handler"
	"github.com/Akihira77/gojobber/services/6-chat/service"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type gRPCServer struct {
	addr string
}

func NewGRPCServer(addr string) *gRPCServer {
	return &gRPCServer{
		addr: addr,
	}
}

func (s *gRPCServer) Run(db *gorm.DB) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// register our grpc services
	chatSvc := service.NewChatService(db)
	handler.NewChatGRPCHandler(grpcServer, chatSvc)

	log.Println("Starting gRPC server on", s.addr)

	return grpcServer.Serve(lis)
}
