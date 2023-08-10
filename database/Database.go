package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var DB DbInstance

func DBConnect() {
	dsn := "host=localhost user=postgres password=password dbname=chatapp port=5432 sslmode=disable"

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
