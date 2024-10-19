package database

import (
	"log"

	"github.com/suavelad/go-fibre-api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var DB DbInstance

func ConnectDb() {

	db, err := gorm.Open(sqlite.Open("api.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")

	}

	log.Println("Database connection successfully opened")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running Migrations")
	db.AutoMigrate(&models.User{},&models.Product{},&models.Order{})

	DB = DbInstance{
		Db: db,
	}

}
