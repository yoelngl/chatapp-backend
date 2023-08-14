package database

import (
	"fmt"
	"log"

	"github.com/chatapp/backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var DB DbInstance

func DBConnect() {
	DATABASE_HOST := config.Config("DATABASE_HOST")
	DATABASE_NAME := config.Config("DATABASE_NAME")
	DATABASE_USER := config.Config("DATABASE_USER")
	DATABASE_PASSWORD := config.Config("DATABASE_PASSWORD")
	DATABASE_PORT := config.Config("DATABASE_PORT")

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
