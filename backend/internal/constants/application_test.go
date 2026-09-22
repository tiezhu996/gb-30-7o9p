package constants

import "testing"

func TestWaitlistTransitions(t *testing.T) {
	next := NextApplicationStatuses(AppStatusWaitlisted)
	for _, want := range []string{AppStatusRejected, AppStatusWithdrawn} {
		found := false
		for _, s := range next {
			if s == want {
				found = true
			}
		}
		if !found {
			t.Errorf("waitlisted -> %s missing, got %v", want, next)
		}
	}
	if IsActiveApplicationStatus(AppStatusWaitlisted) != true {
		t.Error("waitlisted should be active")
	}
	if IsTerminalApplicationStatus(AppStatusClosed) != true {
		t.Error("closed should be terminal")
	}
	if IsReservedApplicationStatus(AppStatusReserved) != true {
		t.Error("reserved should own a reservation slot")
	}
	if IsReservedApplicationStatus(AppStatusWaitlisted) {
		t.Error("waitlisted must not own a reservation slot")
	}
}

func TestSelectionFlow(t *testing.T) {
	// submitted can be selected directly, and rejected is the only other terminal path.
	next := NextApplicationStatuses(AppStatusSubmitted)
	if !contains(next, AppStatusReserved) {
		t.Errorf("submitted must be selectable for reservation, got %v", next)
	}
	// A reservation holder may confirm final adoption, give up, or be cancelled.
	next = NextApplicationStatuses(AppStatusReserved)
	for _, want := range []string{AppStatusApproved, AppStatusWithdrawn, AppStatusCancelled, AppStatusRejected} {
		if !contains(next, want) {
			t.Errorf("reserved must allow -> %s, got %v", want, next)
		}
	}
	// Terminal states cannot transition further.
	if got := NextApplicationStatuses(AppStatusApproved); got != nil {
		t.Errorf("approved must be terminal, got %v", got)
	}
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
