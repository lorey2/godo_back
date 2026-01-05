package database

import (
	"goback/config"
	"goback/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Use the URL from your config package
	dsn := config.DbURL
	
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Fatal will stop the app if the DB isn't reachable
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate using the models we defined
	err = DB.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}
