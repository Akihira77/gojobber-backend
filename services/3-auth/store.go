package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewStore() (*gorm.DB, string) {
	db_user := os.Getenv("DB_USERNAME")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", db_user, db_password, db_name, db_port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		log.Fatalf("Error connecting to Postgres DB\n%+v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error setting connection pool db\n%+v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	log.Println("Success connect to Postgres DB")
	return db, dsn
}
