package constants

// ApplicationStatus enumerates the adoption application state machine.
const (
	AppStatusSubmitted        = "submitted"
	AppStatusOrgReview        = "org_review"
	AppStatusCommunicating    = "communicating"
	AppStatusConfirmed        = "confirmed"
	AppStatusOfflineInterview = "offline_interview"
	AppStatusApproved         = "approved"
	AppStatusRejected         = "rejected"
	AppStatusReserved         = "reserved"
	AppStatusWaitlisted       = "waitlisted"
	AppStatusCancelled        = "cancelled"
	AppStatusClosed           = "closed"
)

// ApplicationEndReason explains why an application ended.
const (
	AppEndRejected        = "rejected"         // 普通审核拒绝（开放期）
	AppEndFinalRejected   = "final_rejected"   // 预留后最终拒绝
	AppEndHolderAbandoned = "holder_abandoned" // 获选人放弃
	AppEndOrgCancelled    = "org_cancelled"    // 机构取消预留
	AppEndAdoptedByOther  = "adopted_by_other" // 名额已被其他领养人获得
	AppEndFinalAdopted    = "final_adopted"    // 最终领养确认
)

// ValidApplicationStatuses returns all accepted application statuses.
func ValidApplicationStatuses() []string {
	return []string{
		AppStatusSubmitted, AppStatusOrgReview, AppStatusCommunicating,
		AppStatusConfirmed, AppStatusOfflineInterview, AppStatusApproved, AppStatusRejected,
		AppStatusReserved, AppStatusWaitlisted, AppStatusCancelled, AppStatusClosed,
	}
}

// IsValidApplicationStatus reports whether a status is known.
func IsValidApplicationStatus(s string) bool {
	for _, v := range ValidApplicationStatuses() {
		if v == s {
			return true
		}
	}
	return false
}

// ApplicationOpenStatuses returns the in-progress statuses before an org
// selects a reservation holder (the open application period).
func ApplicationOpenStatuses() []string {
	return []string{
		AppStatusSubmitted, AppStatusOrgReview, AppStatusCommunicating,
		AppStatusConfirmed, AppStatusOfflineInterview,
	}
}

// ApplicationActiveStatuses returns every non-terminal status, including
// reservation holder and waitlist. Used for uniqueness and cleanup checks.
func ApplicationActiveStatuses() []string {
	return []string{
		AppStatusSubmitted, AppStatusOrgReview, AppStatusCommunicating,
		AppStatusConfirmed, AppStatusOfflineInterview,
		AppStatusReserved, AppStatusWaitlisted,
	}
}

// IsApplicationActive reports whether the application has not ended.
func IsApplicationActive(s string) bool {
	for _, v := range ApplicationActiveStatuses() {
		if v == s {
			return true
		}
	}
	return false
}

// IsApplicationOpen reports whether the application is still in the
// pre-selection pipeline (eligible to be picked or waitlisted).
func IsApplicationOpen(s string) bool {
	for _, v := range ApplicationOpenStatuses() {
		if v == s {
			return true
		}
	}
	return false
}

// ApplicationHolderStatuses returns every status the reservation holder may
// occupy: reserved itself or any review-pipeline step after selection.
func ApplicationHolderStatuses() []string {
	return append([]string{AppStatusReserved}, ApplicationOpenStatuses()...)
}

// IsApplicationHolderStatus reports whether the status belongs to the current
// reservation holder (as opposed to a waitlisted or ended application).
func IsApplicationHolderStatus(s string) bool {
	if s == AppStatusReserved {
		return true
	}
	return IsApplicationOpen(s)
}

// IsApplicationTerminal reports whether the application has ended.
func IsApplicationTerminal(s string) bool {
	switch s {
	case AppStatusApproved, AppStatusRejected, AppStatusCancelled, AppStatusClosed:
		return true
	}
	return false
}

// NextApplicationStatuses returns the allowed forward transitions.
func NextApplicationStatuses(s string) []string {
	switch s {
	case AppStatusSubmitted:
		return []string{AppStatusOrgReview, AppStatusRejected}
	case AppStatusOrgReview:
		return []string{AppStatusCommunicating, AppStatusRejected}
	case AppStatusCommunicating:
		return []string{AppStatusConfirmed, AppStatusRejected}
	case AppStatusConfirmed:
		return []string{AppStatusOfflineInterview}
	case AppStatusOfflineInterview:
		return []string{AppStatusApproved, AppStatusRejected}
	case AppStatusReserved:
		// Holder enters the existing review pipeline; rejecting before review
		// completes is a final rejection that releases the slot to the first
		// waitlist entry. Abandon/cancel use the dedicated release endpoint.
		return []string{AppStatusOrgReview, AppStatusRejected}
	default:
		return nil
	}
}
