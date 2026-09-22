package constants

// ApplicationStatus enumerates the adoption application lifecycle.
const (
	AppStatusSubmitted        = "submitted"
	AppStatusOrgReview        = "org_review"
	AppStatusCommunicating    = "communicating"
	AppStatusConfirmed        = "confirmed"
	AppStatusOfflineInterview = "offline_interview"
	AppStatusReserved         = "reserved"
	AppStatusWaitlisted       = "waitlisted"
	AppStatusApproved         = "approved"
	AppStatusRejected         = "rejected"
	AppStatusWithdrawn        = "withdrawn"
	AppStatusCancelled        = "cancelled"
	AppStatusClosed           = "closed"
)

// ApplicationEndReason explains why a terminal application ended.
const (
	AppEndReasonFinalAdoption     = "final_adoption"
	AppEndReasonAdopterGaveUp     = "adopter_gave_up"
	AppEndReasonReservationCancel = "reservation_cancelled"
	AppEndReasonFinalRejection    = "final_rejection"
	AppEndReasonWithdrawn         = "withdrawn"
	AppEndReasonSuperseded        = "superseded"
)

// ValidApplicationStatuses returns all accepted application statuses.
func ValidApplicationStatuses() []string {
	return []string{
		AppStatusSubmitted, AppStatusOrgReview, AppStatusCommunicating,
		AppStatusConfirmed, AppStatusOfflineInterview, AppStatusReserved,
		AppStatusWaitlisted, AppStatusApproved, AppStatusRejected,
		AppStatusWithdrawn, AppStatusCancelled, AppStatusClosed,
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

// IsActiveApplicationStatus reports whether the application can still affect adoption.
func IsActiveApplicationStatus(s string) bool {
	switch s {
	case AppStatusSubmitted, AppStatusOrgReview, AppStatusCommunicating,
		AppStatusConfirmed, AppStatusOfflineInterview, AppStatusReserved,
		AppStatusWaitlisted:
		return true
	default:
		return false
	}
}

// IsReservedApplicationStatus reports whether the applicant holds the one reservation.
func IsReservedApplicationStatus(s string) bool {
	switch s {
	case AppStatusReserved, AppStatusOrgReview, AppStatusCommunicating,
		AppStatusConfirmed, AppStatusOfflineInterview:
		return true
	default:
		return false
	}
}

// IsTerminalApplicationStatus reports whether the application has ended.
func IsTerminalApplicationStatus(s string) bool {
	return !IsActiveApplicationStatus(s)
}

// NextApplicationStatuses returns explicitly requested transitions. Automatic
// waitlist promotion is handled by the service because it updates several rows.
func NextApplicationStatuses(s string) []string {
	reviewRejections := []string{AppStatusRejected}
	switch s {
	case AppStatusSubmitted:
		return append([]string{AppStatusReserved, AppStatusOrgReview, AppStatusWithdrawn}, reviewRejections...)
	case AppStatusOrgReview:
		return append([]string{AppStatusReserved, AppStatusCommunicating, AppStatusCancelled, AppStatusWithdrawn}, reviewRejections...)
	case AppStatusCommunicating:
		return append([]string{AppStatusReserved, AppStatusConfirmed, AppStatusCancelled, AppStatusWithdrawn}, reviewRejections...)
	case AppStatusConfirmed:
		return []string{AppStatusOfflineInterview, AppStatusCancelled, AppStatusWithdrawn, AppStatusRejected}
	case AppStatusOfflineInterview:
		return []string{AppStatusApproved, AppStatusRejected, AppStatusCancelled, AppStatusWithdrawn}
	case AppStatusReserved:
		return []string{
			AppStatusOrgReview, AppStatusCommunicating, AppStatusConfirmed,
			AppStatusOfflineInterview, AppStatusApproved, AppStatusRejected,
			AppStatusCancelled, AppStatusWithdrawn,
		}
	case AppStatusWaitlisted:
		return []string{AppStatusRejected, AppStatusWithdrawn}
	default:
		return nil
	}
}
