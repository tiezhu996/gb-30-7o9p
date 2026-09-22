package model

import "time"

// Organization is a rescue/shelter org certified by the platform.
type Organization struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	Name        string    `gorm:"size:128;not null" json:"name"`
	LicenseURL  string    `gorm:"size:255" json:"license_url"`
	CertType    string    `gorm:"size:32" json:"cert_type"`
	Status      string    `gorm:"size:16;default:pending;index" json:"status"`
	Contact     string    `gorm:"size:64" json:"contact"`
	City        string    `gorm:"size:64;index" json:"city"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
