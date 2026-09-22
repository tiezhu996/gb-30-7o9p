package constants

// OrganizationStatus enumerates org certification states.
const (
	OrgStatusPending  = "pending"
	OrgStatusApproved = "approved"
	OrgStatusRejected = "rejected"
)

// ValidOrganizationStatuses returns all accepted org statuses.
func ValidOrganizationStatuses() []string {
	return []string{OrgStatusPending, OrgStatusApproved, OrgStatusRejected}
}

// IsValidOrganizationStatus reports whether the status is known.
func IsValidOrganizationStatus(s string) bool {
	for _, v := range ValidOrganizationStatuses() {
		if v == s {
			return true
		}
	}
	return false
}
