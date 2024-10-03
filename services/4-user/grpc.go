package main

import (
	"log"
	"net"

	"github.com/Akihira77/gojobber/services/4-user/handler"
	"github.com/Akihira77/gojobber/services/4-user/service"
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
	buyerSvc := service.NewBuyerService(db)
	sellerSvc := service.NewSellerService(db)
	handler.NewUserGRPCHandler(grpcServer, buyerSvc, sellerSvc)

	log.Println("Starting gRPC server on", s.addr)

	return grpcServer.Serve(lis)
}
