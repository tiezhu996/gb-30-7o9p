package dto

// ApplicationSubmitRequest submits an adoption application.
type ApplicationSubmitRequest struct {
	PetID         uint   `json:"pet_id" binding:"required"`
	Questionnaire string `json:"questionnaire"`
}

// ApplicationStatusRequest changes application status.
type ApplicationStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
