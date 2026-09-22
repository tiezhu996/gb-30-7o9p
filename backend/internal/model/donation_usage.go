package model

import "time"

// DonationUsage is a published usage record of donated funds.
type DonationUsage struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrgID      uint      `gorm:"index;not null" json:"org_id"`
	DonationID uint      `gorm:"index" json:"donation_id"`
	Amount     float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	UsageDesc  string    `gorm:"size:512;not null" json:"usage_desc"`
	ProofURL   string    `gorm:"size:255" json:"proof_url"`
	CreatedAt  time.Time `json:"created_at"`
}
