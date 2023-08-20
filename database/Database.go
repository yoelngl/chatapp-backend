package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var DB DbInstance

func DBConnect() {
	DATABASE_HOST := os.Getenv("DATABASE_HOST")
	DATABASE_NAME := os.Getenv("DATABASE_NAME")
	DATABASE_USER := os.Getenv("DATABASE_USER")
	DATABASE_PASSWORD := os.Getenv("DATABASE_PASSWORD")
	DATABASE_PORT := os.Getenv("DATABASE_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DATABASE_HOST, DATABASE_USER, DATABASE_PASSWORD, DATABASE_NAME, DATABASE_PORT)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database ", err)
	}

	log.Println("Connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	DB = DbInstance{
		Db: db,
	}
}
