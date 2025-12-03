package service

import (
	"encoding/json"
	"net/http"

	"github.com/EnesAybeyR/chi-postgresql-jwt.git/database"
	"github.com/EnesAybeyR/chi-postgresql-jwt.git/models"
	"github.com/go-chi/chi/v5"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	database.DB.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(err.Error()))
	}
	database.DB.Create(&user)
	json.NewEncoder(w).Encode(user)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	var user models.User
	database.DB.First(&user, param)
	json.NewEncoder(w).Encode(user)
}
