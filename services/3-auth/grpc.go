package main

import (
	"log"
	"net"

	"github.com/Akihira77/gojobber/services/3-auth/handler"
	"github.com/Akihira77/gojobber/services/3-auth/service"
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

	authSvc := service.NewAuthService(db)
	handler.NewAuthGRPCHandler(grpcServer, authSvc)

	log.Println("starting grpc server on", s.addr)

	return grpcServer.Serve(lis)
}
