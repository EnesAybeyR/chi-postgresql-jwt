package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dbHost := os.Getenv("HOST")
	dbUser := os.Getenv("USER")
	dbPass := os.Getenv("PASSWORD")
	dbName := os.Getenv("DBNAME")
	dbPort := os.Getenv("PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Database Connection Error:", err)
	}
	DB = db
	fmt.Println("Database connected successfully")
}
