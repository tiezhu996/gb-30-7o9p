package model

import "time"

// Donation is a donation toward an organization.
type Donation struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index;not null" json:"user_id"`
	OrgID         uint      `gorm:"index;not null" json:"org_id"`
	Amount        float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	TransactionID string    `gorm:"size:64" json:"transaction_id"`
	Status        string    `gorm:"size:16;default:pending" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}
