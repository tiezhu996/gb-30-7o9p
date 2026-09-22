package model

import "time"

// User represents an account (user/org/admin).
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"size:128;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Nickname     string    `gorm:"size:64" json:"nickname"`
	Avatar       string    `gorm:"size:255" json:"avatar"`
	Phone        string    `gorm:"size:32" json:"phone"`
	Role         string    `gorm:"size:16;default:user;index" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
