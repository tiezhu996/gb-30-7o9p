package service

import (
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/model"
	"github.com/gbadopt/gbadopt/internal/repository"
	"github.com/gbadopt/gbadopt/internal/util"
)

// ApplicationService implements the adoption application state machine.
type ApplicationService struct {
	db       *gorm.DB
	repo     *repository.AdoptionApplicationRepository
	petRepo  *repository.PetRepository
	orgRepo  *repository.OrganizationRepository
	userRepo *repository.UserRepository
	logger   *slog.Logger
}

// NewApplicationService creates an ApplicationService.
func NewApplicationService(db *gorm.DB, repo *repository.AdoptionApplicationRepository, petRepo *repository.PetRepository, orgRepo *repository.OrganizationRepository, userRepo *repository.UserRepository, logger *slog.Logger) *ApplicationService {
	return &ApplicationService{db: db, repo: repo, petRepo: petRepo, orgRepo: orgRepo, userRepo: userRepo, logger: logger}
}

// Submit creates an application while the pet is open for adoption. A database
// unique index is the final guard against duplicate or concurrent submissions.
func (s *ApplicationService) Submit(userID, petID uint, questionnaire string) (*model.AdoptionApplication, error) {
	a := &model.AdoptionApplication{
		UserID: userID, PetID: petID,
		Questionnaire: questionnaire, Status: constants.AppStatusSubmitted,
	}
	if a.Questionnaire == "" {
		a.Questionnaire = "{}"
	}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		pet, err := s.petRepo.FindByIDForUpdateTx(tx, petID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("Pet[id=%d] not found", petID))
			}
			return fmt.Errorf("application submit pet find: %w", err)
		}
		if pet.Status == constants.PetStatusAdopted || pet.ReservedUserID != 0 {
			return util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("Application[pet_id=%d] submit failed: pet not open (status=%s)", petID, pet.Status))
		}
		a.OrgID = pet.OrgID
		if pet.Status == constants.PetStatusAvailable {
			if err := s.petRepo.UpdateAdoptionStateTx(tx, pet.ID, constants.PetStatusPending, 0); err != nil {
				return fmt.Errorf("application submit pet update: %w", err)
			}
		}
		if err := s.repo.CreateTx(tx, a); err != nil {
			if errors.Is(err, repository.ErrDuplicate) {
				return util.NewAppError(409, constants.CodeConflict,
					fmt.Sprintf("Application[user_id=%d pet_id=%d] submit failed: already applied", userID, petID))
			}
			return fmt.Errorf("application submit: %w", err)
		}
		return nil
	})
	if err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogAppSubmitFailed, petID), "error", err)
		return nil, err
	}
	s.logger.Info(fmt.Sprintf(constants.LogAppSubmitSuccess, a.ID, petID), "id", a.ID)
	return s.FindByID(a.ID)
}

// FindByID returns one enriched application.
func (s *ApplicationService) FindByID(id uint) (*model.AdoptionApplication, error) {
	a, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("AdoptionApplication[id=%d] not found", id))
		}
		return nil, fmt.Errorf("application find: %w", err)
	}
	if err := s.enrich([]model.AdoptionApplication{*a}); err != nil {
		return nil, err
	}
	return a, nil
}

// UpdateStatus transitions an application along the lifecycle. Selection,
// release, promotion and final adoption all happen in one pet-row transaction.
func (s *ApplicationService) UpdateStatus(userID, id uint, role string, next, requestedReason string) (*model.AdoptionApplication, error) {
	if !constants.IsValidApplicationStatus(next) {
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("Application[id=%d] status=%s invalid", id, next))
	}
	current, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("AdoptionApplication[id=%d] not found", id))
		}
		return nil, fmt.Errorf("application status find: %w", err)
	}
	if err := s.canChangeStatus(userID, role, current, next); err != nil {
		return nil, err
	}
	if !isAllowedTransition(current.Status, next) {
		return nil, util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("AdoptionApplication[id=%d] status change failed: %s -> %s not allowed", id, current.Status, next))
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		pet, err := s.petRepo.FindByIDForUpdateTx(tx, current.PetID)
		if err != nil {
			return fmt.Errorf("application status pet find: %w", err)
		}
		a, err := s.repo.FindByIDForUpdateTx(tx, id)
		if err != nil {
			return fmt.Errorf("application status locked find: %w", err)
		}
		// The pre-lock read can be stale under concurrency; re-validate it.
		if a.Status != current.Status || !isAllowedTransition(a.Status, next) {
			return util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("AdoptionApplication[id=%d] status change failed: %s -> %s not allowed", id, a.Status, next))
		}

		switch next {
		case constants.AppStatusReserved:
			return s.reserveTx(tx, pet, a)
		case constants.AppStatusApproved:
			return s.approveTx(tx, pet, a)
		case constants.AppStatusRejected:
			return s.rejectTx(tx, pet, a)
		case constants.AppStatusWithdrawn:
			if a.Status == constants.AppStatusWaitlisted {
				return s.endWaitlistedTx(tx, a, constants.AppStatusWithdrawn, constants.AppEndReasonWithdrawn)
			}
			return s.releaseFromCandidateTx(tx, pet, a, constants.AppStatusWithdrawn, constants.AppEndReasonAdopterGaveUp, constants.AppEndReasonWithdrawn)
		case constants.AppStatusCancelled:
			if a.Status == constants.AppStatusWaitlisted {
				return util.NewAppError(409, constants.CodeConflict,
					fmt.Sprintf("AdoptionApplication[id=%d] cancellation failed: applicant is not the reserved adopter", id))
			}
			return s.releaseFromCandidateTx(tx, pet, a, constants.AppStatusCancelled, constants.AppEndReasonReservationCancel, constants.AppEndReasonReservationCancel)
		default:
			return s.advanceReviewTx(tx, pet, a, next)
		}
	})
	if err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogAppStatusChangeFailed, id), "error", err)
		return nil, err
	}
	s.logger.Info(fmt.Sprintf(constants.LogAppStatusChanged, id, next), "id", id)
	return s.FindByID(id)
}

// ListByUser returns a user's applications with reservation and waitlist data.
func (s *ApplicationService) ListByUser(userID uint) ([]model.AdoptionApplication, error) {
	items, err := s.repo.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("application list by user: %w", err)
	}
	if err := s.enrich(items); err != nil {
		return nil, err
	}
	return items, nil
}

// ListByOrg returns applications for an org with reservation and waitlist data.
func (s *ApplicationService) ListByOrg(userID uint, status string) ([]model.AdoptionApplication, error) {
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil {
		return nil, util.NewAppError(403, constants.CodeForbidden, "org profile not found")
	}
	items, err := s.repo.ListByOrg(org.ID, status)
	if err != nil {
		return nil, fmt.Errorf("application list by org: %w", err)
	}
	if err := s.enrich(items); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ApplicationService) canChangeStatus(userID uint, role string, a *model.AdoptionApplication, next string) error {
	if role == "org" {
		org, err := s.orgRepo.FindByUserID(userID)
		if err != nil || org.ID != a.OrgID {
			return util.NewAppError(403, constants.CodeForbidden,
				fmt.Sprintf("AdoptionApplication[id=%d] status change failed: user_id=%d not org owner", a.ID, userID))
		}
		return nil
	}
	if role == "user" && a.UserID != userID {
		return util.NewAppError(403, constants.CodeForbidden,
			fmt.Sprintf("AdoptionApplication[id=%d] status change failed: not owner", a.ID))
	}
	return nil
}

func isAllowedTransition(current, next string) bool {
	for _, candidate := range constants.NextApplicationStatuses(current) {
		if candidate == next {
			return true
		}
	}
	return false
}

func mapLifecycleConflict(err error, message string) error {
	if errors.Is(err, repository.ErrConflict) {
		return util.NewAppError(409, constants.CodeConflict,
			"application lifecycle changed concurrently; the operation did not succeed")
	}
	return fmt.Errorf("%s: %w", message, err)
}

func (s *ApplicationService) reserveTx(tx *gorm.DB, pet *model.Pet, selected *model.AdoptionApplication) error {
	if pet.Status == constants.PetStatusAdopted || pet.Status == constants.PetStatusReserved || pet.ReservedUserID != 0 {
		return util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("Pet[id=%d] reserve failed: quota is not open", pet.ID))
	}
	if !isOpenCandidate(*selected) {
		return util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("AdoptionApplication[id=%d] reserve failed: status=%s is not eligible", selected.ID, selected.Status))
	}

	from := selected.Status
	selected.Status = constants.AppStatusReserved
	selected.EndReason = ""
	if err := s.repo.UpdateLifecycleTx(tx, selected, from); err != nil {
		return mapLifecycleConflict(err, "application reserve selected update")
	}
	if err := s.closeOtherActiveTx(tx, pet.ID, selected.ID, constants.AppStatusWaitlisted, ""); err != nil {
		return err
	}
	if err := s.petRepo.UpdateAdoptionStateTx(tx, pet.ID, constants.PetStatusReserved, selected.UserID); err != nil {
		return fmt.Errorf("pet reserve update: %w", err)
	}
	return nil
}

func (s *ApplicationService) advanceReviewTx(tx *gorm.DB, pet *model.Pet, a *model.AdoptionApplication, next string) error {
	if !applicationOwnsReservation(pet, a) {
		return util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("AdoptionApplication[id=%d] review failed: applicant is not the reserved adopter", a.ID))
	}
	from := a.Status
	a.Status = next
	a.EndReason = ""
	if err := s.repo.UpdateLifecycleTx(tx, a, from); err != nil {
		return mapLifecycleConflict(err, "application review update")
	}
	return nil
}

func (s *ApplicationService) approveTx(tx *gorm.DB, pet *model.Pet, a *model.AdoptionApplication) error {
	// Normalize records created by the old flow that left pets in pending.
	if isLegacySoleCandidate(pet, a) {
		pet.Status = constants.PetStatusReserved
		pet.ReservedUserID = a.UserID
	}
	if pet.Status != constants.PetStatusReserved || !applicationOwnsReservation(pet, a) {
		return util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("AdoptionApplication[id=%d] approve failed: applicant is not the reserved adopter", a.ID))
	}
	if a.Status == constants.AppStatusWaitlisted || constants.IsTerminalApplicationStatus(a.Status) {
		return util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("AdoptionApplication[id=%d] approve failed: status=%s", a.ID, a.Status))
	}
	from := a.Status
	a.Status = constants.AppStatusApproved
	a.EndReason = constants.AppEndReasonFinalAdoption
	if err := s.repo.UpdateLifecycleTx(tx, a, from); err != nil {
		return mapLifecycleConflict(err, "application approve update")
	}
	if err := s.closeOtherActiveTx(tx, pet.ID, a.ID, constants.AppStatusClosed, constants.AppEndReasonSuperseded); err != nil {
		return err
	}
	if err := s.petRepo.UpdateAdoptionStateTx(tx, pet.ID, constants.PetStatusAdopted, a.UserID); err != nil {
		return fmt.Errorf("pet approve update: %w", err)
	}
	return nil
}

func (s *ApplicationService) rejectTx(tx *gorm.DB, pet *model.Pet, a *model.AdoptionApplication) error {
	from := a.Status
	ownsSlot := applicationOwnsReservation(pet, a) || isLegacySoleCandidate(pet, a)
	if ownsSlot {
		a.Status = constants.AppStatusRejected
		a.EndReason = constants.AppEndReasonFinalRejection
		if err := s.repo.UpdateLifecycleTx(tx, a, from); err != nil {
			return mapLifecycleConflict(err, "application reject update")
		}
		return s.releaseReservationAndPromoteTx(tx, pet)
	}
	if isOpenCandidate(*a) || a.Status == constants.AppStatusWaitlisted {
		a.Status = constants.AppStatusRejected
		a.EndReason = constants.AppEndReasonFinalRejection
		if err := s.repo.UpdateLifecycleTx(tx, a, from); err != nil {
			return mapLifecycleConflict(err, "application reject update")
		}
		return nil
	}
	return util.NewAppError(409, constants.CodeConflict,
		fmt.Sprintf("AdoptionApplication[id=%d] reject failed: status=%s is not active", a.ID, a.Status))
}

func (s *ApplicationService) releaseFromCandidateTx(tx *gorm.DB, pet *model.Pet, a *model.AdoptionApplication, endedStatus, ownerEndReason, reservationEndReason string) error {
	from := a.Status
	ownsSlot := applicationOwnsReservation(pet, a) || isLegacySoleCandidate(pet, a)
	if ownsSlot {
		a.Status = endedStatus
		a.EndReason = ownerEndReason
		if err := s.repo.UpdateLifecycleTx(tx, a, from); err != nil {
			return mapLifecycleConflict(err, "application release update")
		}
		return s.releaseReservationAndPromoteTx(tx, pet)
	}
	if isOpenCandidate(*a) {
		a.Status = endedStatus
		a.EndReason = reservationEndReason
		if err := s.repo.UpdateLifecycleTx(tx, a, from); err != nil {
			return mapLifecycleConflict(err, "application end update")
		}
		return nil
	}
	return util.NewAppError(409, constants.CodeConflict,
		fmt.Sprintf("AdoptionApplication[id=%d] end failed: status=%s is not active", a.ID, a.Status))
}

func (s *ApplicationService) releaseReservationAndPromoteTx(tx *gorm.DB, pet *model.Pet) error {
	next, err := s.repo.FindFirstWaitlistedForUpdateTx(tx, pet.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return s.petRepo.UpdateAdoptionStateTx(tx, pet.ID, constants.PetStatusAvailable, 0)
		}
		return fmt.Errorf("waitlist promotion find: %w", err)
	}
	next.Status = constants.AppStatusReserved
	next.EndReason = ""
	if err := s.repo.UpdateLifecycleTx(tx, next, constants.AppStatusWaitlisted); err != nil {
		return mapLifecycleConflict(err, "waitlist promotion update")
	}
	if err := s.petRepo.UpdateAdoptionStateTx(tx, pet.ID, constants.PetStatusReserved, next.UserID); err != nil {
		return fmt.Errorf("waitlist promotion pet update: %w", err)
	}
	return nil
}

func (s *ApplicationService) endWaitlistedTx(tx *gorm.DB, a *model.AdoptionApplication, status, reason string) error {
	if a.Status != constants.AppStatusWaitlisted {
		return util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("AdoptionApplication[id=%d] end failed: status=%s is not waitlisted", a.ID, a.Status))
	}
	a.Status = status
	a.EndReason = reason
	if err := s.repo.UpdateLifecycleTx(tx, a, constants.AppStatusWaitlisted); err != nil {
		return mapLifecycleConflict(err, "waitlist application end")
	}
	return nil
}

func (s *ApplicationService) closeOtherActiveTx(tx *gorm.DB, petID, keepApplicationID uint, status, reason string) error {
	var others []model.AdoptionApplication
	if err := tx.Where("pet_id = ? AND id <> ?", petID, keepApplicationID).
		Where("status IN ?", []string{
			constants.AppStatusSubmitted, constants.AppStatusOrgReview,
			constants.AppStatusCommunicating, constants.AppStatusConfirmed,
			constants.AppStatusOfflineInterview, constants.AppStatusReserved,
			constants.AppStatusWaitlisted,
		}).Find(&others).Error; err != nil {
		return fmt.Errorf("other active applications find: %w", err)
	}
	for i := range others {
		from := others[i].Status
		others[i].Status = status
		others[i].EndReason = reason
		if from == constants.AppStatusWaitlisted {
			// Waitlisted applicants already have a rank and are simply closed.
			others[i].EndReason = constants.AppEndReasonSuperseded
		}
		if err := s.repo.UpdateLifecycleTx(tx, &others[i], from); err != nil {
			return mapLifecycleConflict(err, "other application end")
		}
	}
	return nil
}

func isOpenCandidate(a model.AdoptionApplication) bool {
	switch a.Status {
	case constants.AppStatusSubmitted, constants.AppStatusOrgReview,
		constants.AppStatusCommunicating, constants.AppStatusConfirmed,
		constants.AppStatusOfflineInterview:
		return true
	default:
		return false
	}
}

func applicationOwnsReservation(pet *model.Pet, a *model.AdoptionApplication) bool {
	if !constants.IsReservedApplicationStatus(a.Status) {
		return false
	}
	return pet.ReservedUserID == a.UserID ||
		(pet.Status == constants.PetStatusPending && pet.ReservedUserID == 0)
}

// isLegacySoleCandidate recognizes old-flow data: one submitted/review
// application against a pending pet that never had a reserved_user_id.
func isLegacySoleCandidate(pet *model.Pet, a *model.AdoptionApplication) bool {
	return pet.Status == constants.PetStatusPending && pet.ReservedUserID == 0 &&
		(a.Status == constants.AppStatusSubmitted || constants.IsReservedApplicationStatus(a.Status))
}

func (s *ApplicationService) enrich(items []model.AdoptionApplication) error {
	if len(items) == 0 {
		return nil
	}
	petIDs := make([]uint, 0, len(items))
	seen := map[uint]struct{}{}
	for _, a := range items {
		if _, ok := seen[a.PetID]; ok {
			continue
		}
		seen[a.PetID] = struct{}{}
		petIDs = append(petIDs, a.PetID)
	}
	pets, err := s.petRepo.ListByIDs(petIDs)
	if err != nil {
		return fmt.Errorf("application enrichment pets: %w", err)
	}
	petByID := make(map[uint]model.Pet, len(pets))
	for _, p := range pets {
		petByID[p.ID] = p
	}
	userIDs := make([]uint, 0, len(items))
	userSeen := map[uint]struct{}{}
	for _, p := range pets {
		if p.ReservedUserID == 0 {
			continue
		}
		if _, ok := userSeen[p.ReservedUserID]; ok {
			continue
		}
		userSeen[p.ReservedUserID] = struct{}{}
		userIDs = append(userIDs, p.ReservedUserID)
	}
	users, err := s.userRepo.ListByIDs(userIDs)
	if err != nil {
		return fmt.Errorf("application enrichment users: %w", err)
	}
	userByID := make(map[uint]model.User, len(users))
	for _, u := range users {
		userByID[u.ID] = u
	}
	var active []model.AdoptionApplication
	active, err = s.repo.ListActiveByPetIDs(petIDs)
	if err != nil {
		return fmt.Errorf("application enrichment waitlist: %w", err)
	}
	waitlistRank := map[uint]int{}
	for _, a := range active {
		if a.Status == constants.AppStatusWaitlisted {
			waitlistRank[a.ID] = waitlistRankForPet(active, a.PetID, a.ID)
		}
	}
	for i := range items {
		if p, ok := petByID[items[i].PetID]; ok {
			items[i].ReservedUserID = p.ReservedUserID
			if u, ok := userByID[p.ReservedUserID]; ok {
				if u.Nickname != "" {
					items[i].ReservedUserName = u.Nickname
				} else {
					items[i].ReservedUserName = u.Username
				}
			}
		}
		items[i].WaitlistPosition = waitlistRank[items[i].ID]
	}
	return nil
}

func waitlistRankForPet(active []model.AdoptionApplication, petID, applicationID uint) int {
	rank := 0
	for _, a := range active {
		if a.PetID != petID || a.Status != constants.AppStatusWaitlisted {
			continue
		}
		rank++
		if a.ID == applicationID {
			return rank
		}
	}
	return 0
}
