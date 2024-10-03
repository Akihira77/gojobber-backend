package main

import (
	"log"
	"os"

	"github.com/Akihira77/gojobber/services/3-auth/handler"
	"github.com/Akihira77/gojobber/services/3-auth/queue"
	"github.com/Akihira77/gojobber/services/3-auth/types"
	"github.com/Akihira77/gojobber/services/3-auth/util"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, _ := NewStore()
	_ = queue.NewConnection(db)
	cld := util.NewCloudinary()

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "pg_trgm";`)
	// db.Debug().Migrator().DropTable(&types.Auth{})
	db.AutoMigrate(&types.Auth{})
	ccs := handler.NewGRPCClients()
	ccs.AddClient("USER_SERVICE", os.Getenv("USER_GRPC_PORT"))
	ccs.AddClient("NOTIFICATION_SERVICE", os.Getenv("NOTIFICATION_GRPC_PORT"))

	go NewHttpServer(db, cld, ccs)

	grpcServer := NewGRPCServer(os.Getenv("AUTH_GRPC_PORT"))
	err = grpcServer.Run(db)
	if err != nil {
		log.Fatal("Error listen GRPC")
	}
}
