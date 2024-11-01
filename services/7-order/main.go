package main

import (
	"log"
	"os"

	"github.com/Akihira77/gojobber/services/7-order/handler"
	"github.com/Akihira77/gojobber/services/7-order/types"
	"github.com/Akihira77/gojobber/services/7-order/util"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v80"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	db, _ := NewStore()
	// db.Debug().Migrator().DropTable(
	// 	types.DeliveredHistory{},
	// 	types.Order{},
	// 	types.OrderEvent{},
	// )
	// db.Debug().AutoMigrate(
	// 	types.DeliveredHistory{},
	// 	types.Order{},
	// 	types.OrderEvent{},
	// )
	// err = types.ApplyDBSetup(db)
	// if err != nil {
	// 	log.Fatalf("Error applying DB setup %v", err)
	// }

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	cld := util.NewCloudinary()

	ccs := handler.NewGRPCClients()
	ccs.AddClient(types.USER_SERVICE, os.Getenv("USER_GRPC_PORT"))
	ccs.AddClient(types.NOTIFICATION_SERVICE, os.Getenv("NOTIFICATION_GRPC_PORT"))
	ccs.AddClient(types.CHAT_SERVICE, os.Getenv("CHAT_GRPC_PORT"))

	go NewHttpServer(db, cld, ccs)

	grpcServer := NewGRPCServer(os.Getenv("ORDER_GRPC_PORT"))
	err = grpcServer.Run(db)
	if err != nil {
		log.Fatalf("Failed running GRPC Server %v", err)
	}
}
