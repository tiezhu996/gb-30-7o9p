package model

import "time"

// Pet is an adoptable animal published by an organization.
type Pet struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrgID        uint      `gorm:"index;not null" json:"org_id"`
	Name         string    `gorm:"size:64;not null" json:"name"`
	Species      string    `gorm:"size:16;index;not null" json:"species"`
	Breed        string    `gorm:"size:64" json:"breed"`
	Age          int       `json:"age"`
	Gender       string    `gorm:"size:8" json:"gender"`
	Size         string    `gorm:"size:16" json:"size"`
	City         string    `gorm:"size:64;index" json:"city"`
	Description  string    `gorm:"type:text" json:"description"`
	Personality  string    `gorm:"size:255" json:"personality"`
	HealthStatus string    `gorm:"size:128" json:"health_status"`
	Neutered     bool      `json:"neutered"`
	Vaccinated   bool      `json:"vaccinated"`
	ImageURLs    string    `gorm:"type:json" json:"image_urls"`
	Status       string    `gorm:"size:16;default:available;index" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
