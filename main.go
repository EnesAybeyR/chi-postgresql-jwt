package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EnesAybeyR/chi-postgresql-jwt.git/database"
	"github.com/EnesAybeyR/chi-postgresql-jwt.git/models"
	"github.com/EnesAybeyR/chi-postgresql-jwt.git/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.ConnectDB()
	if err := database.DB.AutoMigrate(&models.User{}, &models.RefreshToken{}); err != nil {
		log.Fatalf("db connection error: %v", err)
	}
	fmt.Println("Database Migrated")

	r := routes.GetRoutes()
	fmt.Println("server 8081 portunda calisiyor")
	log.Fatal(http.ListenAndServe(":8081", r))
}
