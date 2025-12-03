package main

import (
	"net/http"

	"github.com/EnesAybeyR/chi-postgresql-jwt.git/database"
	"github.com/EnesAybeyR/chi-postgresql-jwt.git/logger"
	"github.com/EnesAybeyR/chi-postgresql-jwt.git/models"
	"github.com/EnesAybeyR/chi-postgresql-jwt.git/routes"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func init() {

}

func main() {
	logger.InitLogger()
	defer logger.Log.Sync()
	err := godotenv.Load()
	if err != nil {
		logger.Log.Fatal("Error loading .env file")
	}
	database.ConnectDB()
	if err := database.DB.AutoMigrate(&models.User{}, &models.RefreshToken{}); err != nil {
		logger.Log.Error("Db connection error: ", zap.Error(err))
	}
	logger.Log.Info("Database Migrated")

	r := routes.GetRoutes()
	logger.Log.Info("server 8081 portunda calisiyor")
	if err := http.ListenAndServe(":8081", r); err != nil {
		logger.Log.Fatal("server failed to start: ", zap.Error(err))
	}
}
