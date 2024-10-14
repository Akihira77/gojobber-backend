package main

import (
	"log"
	"os"

	"github.com/Akihira77/gojobber/services/7-order/handler"
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
	// )
	// db.Debug().AutoMigrate(
	// 	types.DeliveredHistory{},
	// 	types.Order{},
	// )
	// err = types.ApplyDBSetup(db)
	// if err != nil {
	// 	log.Fatalf("Error applying DB setup %v", err)
	// }

	cld := util.NewCloudinary()

	ccs := handler.NewGRPCClients()
	ccs.AddClient("USER_SERVICE", os.Getenv("USER_GRPC_PORT"))
	ccs.AddClient("NOTIFICATION_SERVICE", os.Getenv("NOTIFICATION_GRPC_PORT"))
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	go NewHttpServer(db, cld)

	grpcServer := NewGRPCServer(os.Getenv("ORDER_GRPC_PORT"))
	err = grpcServer.Run(db)
	if err != nil {
		log.Fatalf("Failed running GRPC Server %v", err)
	}
}
