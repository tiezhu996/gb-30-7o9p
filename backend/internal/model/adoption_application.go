package model

import "time"

// AdoptionApplication is an application for adopting a pet.
type AdoptionApplication struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index;not null" json:"user_id"`
	PetID         uint      `gorm:"index;not null" json:"pet_id"`
	OrgID         uint      `gorm:"index;not null" json:"org_id"`
	Questionnaire string    `gorm:"type:json" json:"questionnaire"`
	Status        string    `gorm:"size:32;default:submitted;index" json:"status"`
	WaitlistRank  int       `gorm:"default:0" json:"waitlist_rank"`
	EndReason     string    `gorm:"size:32" json:"end_reason"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
