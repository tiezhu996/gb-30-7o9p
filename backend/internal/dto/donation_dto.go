package dto

// DonateRequest creates a donation.
type DonateRequest struct {
	OrgID  uint    `json:"org_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

// UsageCreateRequest publishes a donation usage record.
type UsageCreateRequest struct {
	DonationID uint    `json:"donation_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	UsageDesc  string  `json:"usage_desc" binding:"required,max=512"`
	ProofURL   string  `json:"proof_url" binding:"omitempty,max=255"`
}
