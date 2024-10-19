package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, _ := NewStore()

	// db.Debug().Exec(`CREATE EXTENSION IF NOT EXISTS tsm_system_rows;`)
	// db.Debug().Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	// err = db.
	// 	Debug().
	// 	Migrator().
	// 	DropTable(
	// 		&types.Buyer{},
	// 		&types.Seller{},
	// 		&types.Skill{},
	// 		&types.SellerSkill{},
	// 		&types.Language{},
	// 		&types.SellerLanguage{},
	// 		&types.Experience{},
	// 		&types.Certificate{},
	// 		&types.Education{},
	// 	)
	// if err != nil {
	// 	log.Fatal("Droping table error")
	// }
	// db.
	// 	Debug().
	// 	AutoMigrate(
	// 		&types.Buyer{},
	// 		&types.Seller{},
	// 		&types.Skill{},
	// 		&types.SellerSkill{},
	// 		&types.Language{},
	// 		&types.SellerLanguage{},
	// 		&types.Experience{},
	// 		&types.Certificate{},
	// 		&types.Education{},
	// 	)
	// if err = types.ApplyDBSetup(db); err != nil {
	// 	log.Fatal(err)
	// }

	go NewHttpServer(db)

	grpcServer := NewGRPCServer(os.Getenv("USER_GRPC_PORT"))
	err = grpcServer.Run(db)
	if err != nil {
		log.Fatal("Error listen GRPC")
	}
}
