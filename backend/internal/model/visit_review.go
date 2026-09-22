package model

import "time"

// ReviewStatus values.
const (
	ReviewPending  = "pending"
	ReviewSubmitted = "submitted"
	ReviewOverdue  = "overdue"
)

// VisitReview tracks post-adoption follow-up visits.
type VisitReview struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ApplicationID uint      `gorm:"index;not null" json:"application_id"`
	UserID        uint      `gorm:"index;not null" json:"user_id"`
	OrgID         uint      `gorm:"index" json:"org_id"`
	ScheduledDays int       `json:"scheduled_days"`
	DueDate       time.Time `gorm:"type:date" json:"due_date"`
	Status        string    `gorm:"size:16;default:pending;index" json:"status"`
	Photos        string    `gorm:"type:json" json:"photos"`
	Note          string    `gorm:"type:text" json:"note"`
	CreatedAt     time.Time `json:"created_at"`
}
