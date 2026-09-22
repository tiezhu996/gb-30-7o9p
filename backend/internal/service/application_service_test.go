package service

import (
	"testing"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/model"
)

func TestReservationOwnership(t *testing.T) {
	pet := &model.Pet{Status: constants.PetStatusReserved, ReservedUserID: 7}
	holder := &model.AdoptionApplication{UserID: 7, Status: constants.AppStatusReserved}
	waiter := &model.AdoptionApplication{UserID: 8, Status: constants.AppStatusWaitlisted}
	if !applicationOwnsReservation(pet, holder) {
		t.Error("reserved application for reserved user must own the slot")
	}
	if applicationOwnsReservation(pet, waiter) {
		t.Error("a waitlisted application must not own the slot")
	}

	// Legacy data from the old pending flow is still recognizable.
	legacyPet := &model.Pet{Status: constants.PetStatusPending}
	legacy := &model.AdoptionApplication{UserID: 9, Status: constants.AppStatusSubmitted}
	if !isLegacySoleCandidate(legacyPet, legacy) {
		t.Error("legacy pending review application should be treated as the sole candidate")
	}
}

func TestWaitlistRanking(t *testing.T) {
	active := []model.AdoptionApplication{
		{ID: 1, PetID: 10, Status: constants.AppStatusReserved},
		{ID: 2, PetID: 10, Status: constants.AppStatusWaitlisted},
		{ID: 3, PetID: 10, Status: constants.AppStatusWaitlisted},
		{ID: 4, PetID: 10, Status: constants.AppStatusWaitlisted},
		{ID: 5, PetID: 20, Status: constants.AppStatusWaitlisted},
	}
	cases := map[uint]int{2: 1, 3: 2, 4: 3, 5: 1}
	for id, want := range cases {
		if got := waitlistRankForPet(active, activeByID(active, id).PetID, id); got != want {
			t.Errorf("application %d rank = %d, want %d", id, got, want)
		}
	}
}

func activeByID(items []model.AdoptionApplication, id uint) model.AdoptionApplication {
	for _, item := range items {
		if item.ID == id {
			return item
		}
	}
	return model.AdoptionApplication{}
}
