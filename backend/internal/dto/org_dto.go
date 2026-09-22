package dto

// OrgRegisterRequest registers an organization.
type OrgRegisterRequest struct {
	Name        string `json:"name" binding:"required,max=128"`
	LicenseURL  string `json:"license_url" binding:"omitempty,max=255"`
	CertType    string `json:"cert_type" binding:"omitempty,max=32"`
	Contact     string `json:"contact" binding:"omitempty,max=64"`
	City        string `json:"city" binding:"omitempty,max=64"`
	Description string `json:"description"`
}

// OrgReviewRequest updates org certification status (admin).
type OrgReviewRequest struct {
	Status string `json:"status" binding:"required"`
}
