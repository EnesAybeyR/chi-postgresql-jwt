package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/EnesAybeyR/chi-postgresql-jwt.git/database"
	"github.com/EnesAybeyR/chi-postgresql-jwt.git/models"
	"github.com/EnesAybeyR/chi-postgresql-jwt.git/service"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type RefreshResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	pwHash, err := service.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	user := models.User{
		Email:    req.Email,
		Password: pwHash,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, "user create error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad requestt", http.StatusBadRequest)
		return
	}
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := service.CheckPassword(user.Password, req.Password); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	access, err := service.GenerateAccessToken(user.Id, user.Email)
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}
	refresh, err := service.GenerateAndStoreRefreshToken(&user)
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}
	resp := LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    time.Now().AddDate(0, 0, 1),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	newAccess, newRefresh, err := service.UseRefreshTokenAndRotate(req.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	resp := RefreshResponse{
		RefreshToken: newRefresh,
		AccessToken:  newAccess,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		RefreshToken string `json:"refresh_token"`
	}
	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	hash := service.HashToken(req.RefreshToken)
	result := database.DB.Model(&models.RefreshToken{}).Where("token_hash = ?", hash).Update("revoked", true)
	if result.Error != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
