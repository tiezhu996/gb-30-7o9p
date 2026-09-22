package dto

// ReviewCreateRequest schedules a follow-up review.
type ReviewCreateRequest struct {
	ApplicationID uint `json:"application_id" binding:"required"`
	ScheduledDays int  `json:"scheduled_days"`
}

// ReviewSubmitRequest submits a review with photos.
type ReviewSubmitRequest struct {
	Photos string `json:"photos"`
	Note   string `json:"note"`
}
