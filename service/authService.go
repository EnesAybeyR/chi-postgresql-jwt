package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/EnesAybeyR/chi-postgresql-jwt.git/database"
	"github.com/EnesAybeyR/chi-postgresql-jwt.git/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte(os.Getenv("JWTKEY"))

func HashPassword(pw string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(h), err
}
func CheckPassword(hash, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}

func GenerateAccessToken(userId uint, email string) (string, error) {
	expDay, _ := strconv.Atoi("1")
	if expDay == 0 {
		expDay = 1
	}
	claims := jwt.MapClaims{
		"sub":   userId,
		"email": email,
		"exp":   time.Now().AddDate(0, 0, expDay).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func generateRefreshTokenString() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h[:])
}

func GenerateAndStoreRefreshToken(user *models.User) (string, error) {
	tokenStr, err := generateRefreshTokenString()
	if err != nil {
		return "", err
	}
	hash := HashToken(tokenStr)
	expDay, _ := strconv.Atoi("15")
	if expDay == 0 {
		expDay = 15
	}

	rt := models.RefreshToken{
		UserId:    user.Id,
		TokenHash: hash,
		ExpiresAt: time.Now().AddDate(0, 0, expDay),
		Revoked:   false,
	}
	if err := database.DB.Create(&rt).Error; err != nil {
		return "", err
	}
	return tokenStr, nil
}

func UseRefreshTokenAndRotate(tokenStr string) (newAccess string, newRefresh string, err error) {
	hash := HashToken(tokenStr)
	var rt models.RefreshToken
	if err := database.DB.Preload("User").Where("token_hash = ?", hash).First(&rt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", errors.New("invalid refresh token")
		}
		return "", "", err
	}
	if rt.Revoked || rt.ExpiresAt.Before(time.Now()) {
		return "", "", errors.New("refresh token revoked or expired")
	}
	var user models.User
	if err := database.DB.First(&user, rt.UserId).Error; err != nil {
		return "", "", err
	}
	if err := database.DB.Model(&rt).Update("revoked", true).Error; err != nil {
		return "", "", err
	}
	accessToken, err := GenerateAccessToken(user.Id, user.Email)
	if err != nil {
		return "", "", err
	}
	newRT, err := GenerateAndStoreRefreshToken(&user)
	if err != nil {
		return "", "", err
	}
	return accessToken, newRT, nil
}
