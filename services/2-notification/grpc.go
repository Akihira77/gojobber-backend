package main

import (
	"log"
	"net"

	"github.com/Akihira77/gojobber/services/2-notification/handler"
	"github.com/Akihira77/gojobber/services/2-notification/service"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr string
}

func NewGRPCServer(addr string) *gRPCServer {
	return &gRPCServer{
		addr: addr,
	}
}

func (s *gRPCServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// register our grpc services
	notificationSvc := service.NewNotificationService()
	handler.NewNotificationGRPCHandler(grpcServer, notificationSvc)

	log.Println("Starting gRPC server on", s.addr)

	return grpcServer.Serve(lis)
}
