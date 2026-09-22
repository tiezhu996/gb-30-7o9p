package constants

import "testing"

func TestApplicationStatusLifecycle(t *testing.T) {
	if IsApplicationTerminal(AppStatusApproved) != true {
		t.Error("approved must be terminal")
	}
	if IsApplicationTerminal(AppStatusWaitlisted) {
		t.Error("waitlisted must be active")
	}
	if IsApplicationOpen(AppStatusReserved) {
		t.Error("reserved is no longer in the open pipeline")
	}
	if !IsApplicationActive(AppStatusReserved) || !IsApplicationActive(AppStatusWaitlisted) {
		t.Error("reserved and waitlisted must remain active")
	}
}

func TestReservedTransitions(t *testing.T) {
	next := NextApplicationStatuses(AppStatusReserved)
	want := map[string]bool{AppStatusOrgReview: false, AppStatusRejected: false}
	for _, s := range next {
		if _, ok := want[s]; ok {
			want[s] = true
		} else {
			t.Errorf("unexpected transition reserved -> %s", s)
		}
	}
	for s, found := range want {
		if !found {
			t.Errorf("missing transition reserved -> %s", s)
		}
	}
}

func TestWaitlistHasNoSelfTransitions(t *testing.T) {
	// Waitlisted applications only move via first-in-line promotion.
	if got := NextApplicationStatuses(AppStatusWaitlisted); got != nil {
		t.Errorf("waitlisted transitions = %v, want nil", got)
	}
	if got := NextApplicationStatuses(AppStatusApproved); got != nil {
		t.Errorf("approved transitions = %v, want nil", got)
	}
}

func TestReviewPipelineLeadsToApproval(t *testing.T) {
	chain := []string{
		AppStatusReserved, AppStatusOrgReview, AppStatusCommunicating,
		AppStatusConfirmed, AppStatusOfflineInterview, AppStatusApproved,
	}
	for i := 0; i+1 < len(chain); i++ {
		allowed := false
		for _, s := range NextApplicationStatuses(chain[i]) {
			if s == chain[i+1] {
				allowed = true
			}
		}
		if !allowed {
			t.Errorf("%s -> %s must be allowed", chain[i], chain[i+1])
		}
	}
}
