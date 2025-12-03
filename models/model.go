package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password" gorm:"not null"`
}

type RefreshToken struct {
	gorm.Model
	Id        uint      `json:"id"`
	UserId    uint      `json:"userId" gorm:"not null"`
	User      User      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	TokenHash string    `json:"tokenHash" gorm:"not null;index"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"index"`
	Revoked   bool      `json:"revoked" gorm:"default:false"`
}
