package main

import (
	"log"
	"os"

	"github.com/Akihira77/gojobber/services/6-chat/handler"
	"github.com/Akihira77/gojobber/services/6-chat/util"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, _ := NewStore()
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	// db.Migrator().DropTable(
	// 	&types.Conversation{},
	// 	&types.Message{},
	// )
	// db.AutoMigrate(
	// 	&types.Conversation{},
	// 	&types.Message{},
	// )
	// err = types.ApplyDBSetup(db)
	// if err != nil {
	// 	log.Fatal("Error applying DB entities setup")
	// }

	cld := util.NewCloudinary()
	ccs := handler.NewGRPCClients()
	ccs.AddClient("AUTH_SERVICE", os.Getenv("AUTH_GRPC_PORT"))
	ccs.AddClient("USER_SERVICE", os.Getenv("USER_GRPC_PORT"))
	ccs.AddClient("NOTIFICATION_SERVICE", os.Getenv("NOTIFICATION_GRPC_PORT"))

	go NewHttpServer(db, cld, ccs)

	grpcServer := NewGRPCServer(os.Getenv("CHAT_GRPC_PORT"))
	err = grpcServer.Run(db)
	if err != nil {
		log.Fatalf("Error running grpc server:\n+%v", err)
	}
}
