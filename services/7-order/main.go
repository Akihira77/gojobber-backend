package main

import (
	"log"

	"github.com/Akihira77/gojobber/services/7-order/util"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	db, _ := NewStore()
	// db.Migrator().DropTable(
	// 	types.DeliveredHistory{},
	// 	types.Order{},
	// )
	// db.AutoMigrate(
	// 	types.DeliveredHistory{},
	// 	types.Order{},
	// )
	// err = types.ApplyDBSetup(db)
	// if err != nil {
	// 	log.Fatalf("Error applying DB setup %v", err)
	// }

	cld := util.NewCloudinary()

	NewHttpServer(db, cld)
}
