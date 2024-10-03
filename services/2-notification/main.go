package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error reading env file", err)
	}

	// q := queue.NewConnection()
	// go q.ConsumeFromAuthService()
	// q.ConsumeFromChatService()

	grpcServer := NewGRPCServer(os.Getenv("NOTIFICATION_GRPC_PORT"))
	grpcServer.Run()
}
