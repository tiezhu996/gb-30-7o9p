package dto

import "github.com/gbadopt/gbadopt/internal/model"

// ApplicationSubmitRequest submits an adoption application.
type ApplicationSubmitRequest struct {
	PetID         uint   `json:"pet_id" binding:"required"`
	Questionnaire string `json:"questionnaire"`
}

// ApplicationStatusRequest changes application status.
type ApplicationStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// ApplicationReleaseRequest releases a reservation: the selected adopter
// gives up (action=abandon) or the org cancels the reservation (action=cancel).
type ApplicationReleaseRequest struct {
	Action string `json:"action" binding:"required,oneof=abandon cancel"`
}

// ApplicationView enriches an application with waitlist rank, end reason and
// reservation holder information shown on both applicant and org pages.
type ApplicationView struct {
	model.AdoptionApplication
	PetName        string `json:"pet_name"`
	PetStatus      string `json:"pet_status"`
	WaitlistActive bool   `json:"waitlist_active"`
	EndReasonText  string `json:"end_reason_text"`
	ReservedUserID uint   `json:"reserved_user_id"`
	ReservedName   string `json:"reserved_name"`
}
