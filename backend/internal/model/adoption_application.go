package model

import "time"

// AdoptionApplication is an application for adopting a pet.
type AdoptionApplication struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index;uniqueIndex:uniq_application_user_pet;not null" json:"user_id"`
	PetID         uint      `gorm:"index;uniqueIndex:uniq_application_user_pet;not null" json:"pet_id"`
	OrgID         uint      `gorm:"index;not null" json:"org_id"`
	Questionnaire string    `gorm:"type:json" json:"questionnaire"`
	Status        string    `gorm:"size:32;default:submitted;index" json:"status"`
	EndReason     string    `gorm:"size:64" json:"end_reason"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Derived fields populated by the service; they are not database columns.
	WaitlistPosition int    `gorm:"-" json:"waitlist_position"`
	ReservedUserID   uint   `gorm:"-" json:"reserved_user_id"`
	ReservedUserName string `gorm:"-" json:"reserved_user_name"`
}
