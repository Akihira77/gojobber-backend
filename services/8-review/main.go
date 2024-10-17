package main

import (
	"log"
	"os"

	"github.com/Akihira77/gojobber/services/8-review/handler"
	"github.com/Akihira77/gojobber/services/8-review/types"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error reading environment")
	}

	db, _ := NewStore()
	err = db.
		Debug().
		Migrator().
		DropTable(&types.Review{})
	if err != nil {
		log.Fatalf("Error while droping table")
	}

	err = db.
		Debug().
		AutoMigrate(&types.Review{})
	if err != nil {
		log.Fatalf("Error while migrating schema to database")
	}

	err = types.ApplyDBSetup(db)
	if err != nil {
		log.Fatalf("Error while applying other db setup")
	}

	ccs := handler.NewGRPCClients()
	ccs.AddClient("USER_SERVICE", os.Getenv("USER_GRPC_PORT"))
	ccs.AddClient("NOTIFICATION_SERVICE", os.Getenv("NOTIFICATION_GRPC_PORT"))

	NewHttpServer(db)
}
