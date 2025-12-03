package database

import (
	"fmt"

	"os"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	Logger, _ := zap.NewDevelopment()
	defer Logger.Sync()
	dbHost := os.Getenv("HOST")
	dbUser := os.Getenv("USER")
	dbPass := os.Getenv("PASSWORD")
	dbName := os.Getenv("DBNAME")
	dbPort := os.Getenv("PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		Logger.Fatal("Database Connection Error: ", zap.Error(err))
	}
	DB = db
	Logger.Info("Database connected successfully")
}
