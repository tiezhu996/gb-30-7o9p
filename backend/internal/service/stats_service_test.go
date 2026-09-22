package service

import (
	"testing"

	"github.com/gbadopt/gbadopt/internal/model"
)

func TestComputeStats(t *testing.T) {
	apps := []model.AdoptionApplication{
		{Status: "submitted"},
		{Status: "approved"},
		{Status: "approved"},
		{Status: "rejected"},
	}
	s := ComputeStats(apps)
	if s.Total != 4 {
		t.Errorf("Total = %d, want 4", s.Total)
	}
	if s.ByStatus["approved"] != 2 {
		t.Errorf("approved = %d, want 2", s.ByStatus["approved"])
	}
	if s.Approved != 2 {
		t.Errorf("Approved = %d, want 2", s.Approved)
	}
}
