package service

import (
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/dto"
	"github.com/gbadopt/gbadopt/internal/model"
	"github.com/gbadopt/gbadopt/internal/repository"
	"github.com/gbadopt/gbadopt/internal/util"
)

// ApplicationService implements the adoption application state machine,
// including reservation selection, waitlist promotion and slot release.
type ApplicationService struct {
	db       *gorm.DB
	repo     *repository.AdoptionApplicationRepository
	petRepo  *repository.PetRepository
	orgRepo  *repository.OrganizationRepository
	userRepo *repository.UserRepository
	logger   *slog.Logger
}

// NewApplicationService creates an ApplicationService.
func NewApplicationService(
	db *gorm.DB,
	repo *repository.AdoptionApplicationRepository,
	petRepo *repository.PetRepository,
	orgRepo *repository.OrganizationRepository,
	userRepo *repository.UserRepository,
	logger *slog.Logger,
) *ApplicationService {
	return &ApplicationService{db: db, repo: repo, petRepo: petRepo, orgRepo: orgRepo, userRepo: userRepo, logger: logger}
}

// Submit creates an application from a user to a pet. Multiple users may apply
// while the pet remains available; the pet only becomes reserved on selection.
func (s *ApplicationService) Submit(userID, petID uint, questionnaire string) (*dto.ApplicationView, error) {
	pet, err := s.petRepo.FindByID(petID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("Pet[id=%d] not found", petID))
		}
		return nil, fmt.Errorf("application submit pet find: %w", err)
	}
	if pet.Status != constants.PetStatusAvailable {
		return nil, util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("Application[pet_id=%d] submit failed: pet not available (status=%s)", petID, pet.Status))
	}
	if exist, err := s.repo.FindActiveByUserAndPet(userID, petID); err == nil && exist != nil {
		return nil, util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("Application[user_id=%d pet_id=%d] submit failed: active application id=%d", userID, petID, exist.ID))
	}
	a := &model.AdoptionApplication{
		UserID: userID, PetID: petID, OrgID: pet.OrgID,
		Questionnaire: questionnaire, Status: constants.AppStatusSubmitted,
	}
	if a.Questionnaire == "" {
		a.Questionnaire = "{}"
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
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
	return s.loadView(a)
}

// SelectAdopter lets an org choose one open application as the reservation
// holder. The pet becomes reserved and every other in-progress application is
// waitlisted by submission order. Repeated or concurrent selection fails
// without overwriting the existing reservation.
func (s *ApplicationService) SelectAdopter(userID, id uint) (*dto.ApplicationView, error) {
	a, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("AdoptionApplication[id=%d] not found", id))
		}
		return nil, fmt.Errorf("application select find: %w", err)
	}
	if err := s.requireOrgOwner(userID, a.OrgID, id, "select"); err != nil {
		return nil, err
	}
	var promoted *model.AdoptionApplication
	err = s.db.Transaction(func(tx *gorm.DB) error {
		pet, err := s.petRepo.FindByIDForUpdateTx(tx, a.PetID)
		if err != nil {
			return fmt.Errorf("application select pet find: %w", err)
		}
		// pending is accepted for pre-feature data, where the first application
		// moved the pet out of available immediately.
		if pet.Status != constants.PetStatusAvailable && pet.Status != constants.PetStatusPending {
			return util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("AdoptionApplication[id=%d] select failed: pet status=%s not selectable", id, pet.Status))
		}
		target, err := s.repo.FindByIDForUpdateTx(tx, id)
		if err != nil {
			return fmt.Errorf("application select find for update: %w", err)
		}
		if !constants.IsApplicationOpen(target.Status) {
			return util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("AdoptionApplication[id=%d] select failed: status=%s not eligible", id, target.Status))
		}
		actives, err := s.repo.ListActiveByPetForUpdateTx(tx, a.PetID)
		if err != nil {
			return fmt.Errorf("application select list active: %w", err)
		}
		// Conditional update is the single-writer guard: a racing selection on
		// another application changes zero rows and aborts this transaction.
		rows, err := s.repo.UpdateStatusIfTx(tx, id, constants.ApplicationOpenStatuses(),
			constants.AppStatusReserved, "", true)
		if err != nil {
			return fmt.Errorf("application select holder update: %w", err)
		}
		if rows == 0 {
			return util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("AdoptionApplication[id=%d] select failed: reservation already taken", id))
		}
		rank := 1
		for i := range actives {
			cur := actives[i]
			if cur.ID == id {
				continue
			}
			cur.Status = constants.AppStatusWaitlisted
			cur.WaitlistRank = rank
			rank++
			if err := s.repo.UpdateTx(tx, &cur); err != nil {
				return fmt.Errorf("application select waitlist update: %w", err)
			}
		}
		if _, err := s.petRepo.UpdateStatusIfTx(tx, pet.ID,
			[]string{constants.PetStatusAvailable, constants.PetStatusPending}, constants.PetStatusReserved, id); err != nil {
			return fmt.Errorf("application select pet update: %w", err)
		}
		promoted = target
		return nil
	})
	if err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogAppSelectFailed, id), "error", err)
		return nil, err
	}
	s.logger.Info(fmt.Sprintf(constants.LogAppSelected, id, a.PetID), "id", id)
	return s.loadView(promoted)
}

// ReleaseReservation ends a reservation because the holder gave up or the org
// cancelled it. The first waitlisted application is atomically promoted; when
// nobody is waiting the pet returns to available. Only rank-1 can be promoted.
func (s *ApplicationService) ReleaseReservation(userID, id uint, role, action string) (*dto.ApplicationView, error) {
	a, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("AdoptionApplication[id=%d] not found", id))
		}
		return nil, fmt.Errorf("application release find: %w", err)
	}
	switch action {
	case "abandon":
		if role != "user" || a.UserID != userID {
			return nil, util.NewAppError(403, constants.CodeForbidden,
				fmt.Sprintf("AdoptionApplication[id=%d] abandon failed: not holder", id))
		}
	case "cancel":
		if err := s.requireOrgOwner(userID, a.OrgID, id, "cancel reservation"); err != nil {
			return nil, err
		}
	default:
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("AdoptionApplication[id=%d] release failed: action=%s invalid", id, action))
	}
	endReason := constants.AppEndHolderAbandoned
	terminal := constants.AppStatusCancelled
	if action == "cancel" {
		endReason = constants.AppEndOrgCancelled
	}
	var released *model.AdoptionApplication
	err = s.db.Transaction(func(tx *gorm.DB) error {
		pet, err := s.petRepo.FindByIDForUpdateTx(tx, a.PetID)
		if err != nil {
			return fmt.Errorf("application release pet find: %w", err)
		}
		if pet.Status != constants.PetStatusReserved || pet.ReservedApplicationID != id {
			return util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("AdoptionApplication[id=%d] release failed: not current holder of pet %d", id, a.PetID))
		}
		rows, err := s.repo.UpdateStatusIfTx(tx, id,
			constants.ApplicationHolderStatuses(), terminal, endReason, true)
		if err != nil {
			return fmt.Errorf("application release holder update: %w", err)
		}
		if rows == 0 {
			return util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("AdoptionApplication[id=%d] release failed: already released", id))
		}
		waitlist, err := s.repo.ListWaitlistByPetForUpdateTx(tx, a.PetID)
		if err != nil {
			return fmt.Errorf("application release list waitlist: %w", err)
		}
		if len(waitlist) == 0 {
			if _, err := s.petRepo.UpdateStatusIfTx(tx, pet.ID,
				[]string{constants.PetStatusReserved}, constants.PetStatusAvailable, 0); err != nil {
				return fmt.Errorf("application release pet reopen: %w", err)
			}
			released = a
			return nil
		}
		// Only the first-in-line (rank 1) takes the slot.
		next := waitlist[0]
		if next.WaitlistRank != 1 {
			return util.NewAppError(500, constants.CodeInternalError,
				fmt.Sprintf("AdoptionApplication[id=%d] release failed: waitlist head rank=%d", next.ID, next.WaitlistRank))
		}
		promRows, err := s.repo.UpdateStatusIfTx(tx, next.ID,
			[]string{constants.AppStatusWaitlisted}, constants.AppStatusReserved, "", true)
		if err != nil {
			return fmt.Errorf("application release promote: %w", err)
		}
		if promRows == 0 {
			return util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("AdoptionApplication[id=%d] release failed: promotion race", next.ID))
		}
		// Remaining entries shift up one rank.
		for _, w := range waitlist[1:] {
			w.WaitlistRank--
			if err := s.repo.UpdateTx(tx, &w); err != nil {
				return fmt.Errorf("application release rank update: %w", err)
			}
		}
		if _, err := s.petRepo.UpdateStatusIfTx(tx, pet.ID,
			[]string{constants.PetStatusReserved}, constants.PetStatusReserved, next.ID); err != nil {
			return fmt.Errorf("application release pet holder update: %w", err)
		}
		released = a
		return nil
	})
	if err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogAppReleaseFailed, id), "error", err)
		return nil, err
	}
	s.logger.Info(fmt.Sprintf(constants.LogAppReleased, id, action), "id", id)
	return s.loadView(released)
}

// UpdateStatus transitions an application along the review state machine.
// Final approval while the application holds the reservation adopts the pet and
// ends all other applications; final rejection of the holder releases the slot
// to the first waitlist entry. Review steps and permissions are unchanged.
func (s *ApplicationService) UpdateStatus(userID, id uint, role string, next string) (*dto.ApplicationView, error) {
	if !constants.IsValidApplicationStatus(next) {
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("Application[id=%d] status=%s invalid", id, next))
	}
	a, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("AdoptionApplication[id=%d] not found", id))
		}
		return nil, fmt.Errorf("application status find: %w", err)
	}
	if role == "org" {
		if err := s.requireOrgOwner(userID, a.OrgID, id, "status change"); err != nil {
			return nil, err
		}
	} else if role == "user" && a.UserID != userID {
		return nil, util.NewAppError(403, constants.CodeForbidden,
			fmt.Sprintf("AdoptionApplication[id=%d] status change failed: not owner", id))
	}
	allowed := false
	for _, s2 := range constants.NextApplicationStatuses(a.Status) {
		if s2 == next {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("AdoptionApplication[id=%d] status change failed: %s -> %s not allowed", id, a.Status, next))
	}

	var updated *model.AdoptionApplication
	switch {
	case next == constants.AppStatusApproved:
		updated, err = s.approveFinal(a)
	case next == constants.AppStatusRejected:
		updated, err = s.rejectByContext(a)
	default:
		updated, err = s.simpleTransition(a, next)
	}
	if err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogAppStatusChangeFailed, id), "error", err)
		return nil, err
	}
	s.logger.Info(fmt.Sprintf(constants.LogAppStatusChanged, id, next), "id", id)
	return s.loadView(updated)
}

// simpleTransition applies a side-effect-free forward review step.
func (s *ApplicationService) simpleTransition(a *model.AdoptionApplication, next string) (*model.AdoptionApplication, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		return s.forwardStepTx(tx, a, next)
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

// forwardStepTx moves an application one step forward without other effects.
func (s *ApplicationService) forwardStepTx(tx *gorm.DB, a *model.AdoptionApplication, next string) error {
	current, err := s.repo.FindByIDForUpdateTx(tx, a.ID)
	if err != nil {
		return fmt.Errorf("application status find for update: %w", err)
	}
	rows, err := s.repo.UpdateStatusIfTx(tx, a.ID, []string{current.Status}, next, "", false)
	if err != nil {
		return fmt.Errorf("application status update: %w", err)
	}
	if rows == 0 {
		return util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("AdoptionApplication[id=%d] status change failed: state changed concurrently", a.ID))
	}
	a.Status = next
	return nil
}

// approveFinal confirms the final adoption for the current reservation holder:
// the pet becomes adopted and every other in-progress application ends.
func (s *ApplicationService) approveFinal(a *model.AdoptionApplication) (*model.AdoptionApplication, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		pet, err := s.petRepo.FindByIDForUpdateTx(tx, a.PetID)
		if err != nil {
			return fmt.Errorf("application approve pet find: %w", err)
		}
		if pet.Status != constants.PetStatusReserved || pet.ReservedApplicationID != a.ID {
			return util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("AdoptionApplication[id=%d] approve failed: not current holder of pet %d", a.ID, a.PetID))
		}
		current, err := s.repo.FindByIDForUpdateTx(tx, a.ID)
		if err != nil {
			return fmt.Errorf("application approve find for update: %w", err)
		}
		rows, err := s.repo.UpdateStatusIfTx(tx, a.ID, []string{current.Status},
			constants.AppStatusApproved, constants.AppEndFinalAdopted, true)
		if err != nil {
			return fmt.Errorf("application approve update: %w", err)
		}
		if rows == 0 {
			return util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("AdoptionApplication[id=%d] approve failed: state changed concurrently", a.ID))
		}
		// Everyone else (waitlist or stray active applications) ends at once.
		if _, err := s.repo.CloseOtherActiveTx(tx, a.PetID, a.ID,
			constants.AppStatusClosed, constants.AppEndAdoptedByOther); err != nil {
			return fmt.Errorf("application approve close others: %w", err)
		}
		petRows, err := s.petRepo.UpdateStatusIfTx(tx, pet.ID,
			[]string{constants.PetStatusReserved}, constants.PetStatusAdopted, 0)
		if err != nil {
			return fmt.Errorf("application approve pet update: %w", err)
		}
		if petRows == 0 {
			return util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("Pet[id=%d] approve failed: state changed concurrently", pet.ID))
		}
		a.Status = constants.AppStatusApproved
		a.EndReason = constants.AppEndFinalAdopted
		a.WaitlistRank = 0
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

// rejectByContext rejects an application. When it currently holds a reserved
// pet (at any review-pipeline stage), the rejection is final and releases the
// slot to the first waitlist entry; otherwise it is an ordinary rejection.
func (s *ApplicationService) rejectByContext(a *model.AdoptionApplication) (*model.AdoptionApplication, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		pet, err := s.petRepo.FindByIDForUpdateTx(tx, a.PetID)
		if err != nil {
			return fmt.Errorf("application reject pet find: %w", err)
		}
		if pet.Status != constants.PetStatusReserved || pet.ReservedApplicationID != a.ID {
			return s.doSimpleRejectTx(tx, a, constants.AppEndRejected)
		}
		return s.doFinalRejectTx(tx, pet, a)
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

// doSimpleRejectTx records an ordinary rejection (open application period),
// without touching pet status or other applications.
func (s *ApplicationService) doSimpleRejectTx(tx *gorm.DB, a *model.AdoptionApplication, endReason string) error {
	current, err := s.repo.FindByIDForUpdateTx(tx, a.ID)
	if err != nil {
		return fmt.Errorf("application reject find for update: %w", err)
	}
	rows, err := s.repo.UpdateStatusIfTx(tx, a.ID, []string{current.Status},
		constants.AppStatusRejected, endReason, false)
	if err != nil {
		return fmt.Errorf("application reject update: %w", err)
	}
	if rows == 0 {
		return util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("AdoptionApplication[id=%d] reject failed: state changed concurrently", a.ID))
	}
	a.Status = constants.AppStatusRejected
	a.EndReason = endReason
	return nil
}

// doFinalRejectTx rejects the current reservation holder and promotes rank-1.
func (s *ApplicationService) doFinalRejectTx(tx *gorm.DB, pet *model.Pet, a *model.AdoptionApplication) error {
	current, err := s.repo.FindByIDForUpdateTx(tx, a.ID)
	if err != nil {
		return fmt.Errorf("application reject find for update: %w", err)
	}
	rows, err := s.repo.UpdateStatusIfTx(tx, a.ID, []string{current.Status},
		constants.AppStatusRejected, constants.AppEndFinalRejected, true)
	if err != nil {
		return fmt.Errorf("application reject holder update: %w", err)
	}
	if rows == 0 {
		return util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("AdoptionApplication[id=%d] reject failed: already released", a.ID))
	}
	waitlist, err := s.repo.ListWaitlistByPetForUpdateTx(tx, a.PetID)
	if err != nil {
		return fmt.Errorf("application reject list waitlist: %w", err)
	}
	if len(waitlist) == 0 {
		if _, err := s.petRepo.UpdateStatusIfTx(tx, pet.ID,
			[]string{constants.PetStatusReserved}, constants.PetStatusAvailable, 0); err != nil {
			return fmt.Errorf("application reject pet reopen: %w", err)
		}
		a.Status = constants.AppStatusRejected
		a.EndReason = constants.AppEndFinalRejected
		a.WaitlistRank = 0
		return nil
	}
	next := waitlist[0]
	if next.WaitlistRank != 1 {
		return util.NewAppError(500, constants.CodeInternalError,
			fmt.Sprintf("AdoptionApplication[id=%d] reject failed: waitlist head rank=%d", next.ID, next.WaitlistRank))
	}
	promRows, err := s.repo.UpdateStatusIfTx(tx, next.ID,
		[]string{constants.AppStatusWaitlisted}, constants.AppStatusReserved, "", true)
	if err != nil {
		return fmt.Errorf("application reject promote: %w", err)
	}
	if promRows == 0 {
		return util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("AdoptionApplication[id=%d] reject failed: promotion race", next.ID))
	}
	for _, w := range waitlist[1:] {
		w.WaitlistRank--
		if err := s.repo.UpdateTx(tx, &w); err != nil {
			return fmt.Errorf("application reject rank update: %w", err)
		}
	}
	if _, err := s.petRepo.UpdateStatusIfTx(tx, pet.ID,
		[]string{constants.PetStatusReserved}, constants.PetStatusReserved, next.ID); err != nil {
		return fmt.Errorf("application reject pet holder update: %w", err)
	}
	a.Status = constants.AppStatusRejected
	a.EndReason = constants.AppEndFinalRejected
	a.WaitlistRank = 0
	return nil
}

func (s *ApplicationService) requireOrgOwner(userID, orgID, id uint, action string) error {
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil || org.ID != orgID {
		return util.NewAppError(403, constants.CodeForbidden,
			fmt.Sprintf("AdoptionApplication[id=%d] %s failed: user_id=%d not org owner", id, action, userID))
	}
	return nil
}

// ListByUser returns a user's applications with waitlist/reservation context,
// including the holder display name, rank and end reasons.
func (s *ApplicationService) ListByUser(userID uint) ([]dto.ApplicationView, error) {
	items, err := s.repo.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("application list by user: %w", err)
	}
	return s.loadViews(items)
}

// ListByOrg returns applications for an org with waitlist/reservation context.
func (s *ApplicationService) ListByOrg(userID uint, status string) ([]dto.ApplicationView, error) {
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil {
		return nil, util.NewAppError(403, constants.CodeForbidden, "org profile not found")
	}
	items, err := s.repo.ListByOrg(org.ID, status)
	if err != nil {
		return nil, fmt.Errorf("application list by org: %w", err)
	}
	return s.loadViews(items)
}

// loadView enriches a single application.
func (s *ApplicationService) loadView(a *model.AdoptionApplication) (*dto.ApplicationView, error) {
	views, err := s.loadViews([]model.AdoptionApplication{*a})
	if err != nil {
		return nil, err
	}
	if len(views) == 0 {
		return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("AdoptionApplication[id=%d] not found", a.ID))
	}
	return &views[0], nil
}

// loadViews batch-enriches applications with pet and reservation-holder data.
func (s *ApplicationService) loadViews(items []model.AdoptionApplication) ([]dto.ApplicationView, error) {
	views := make([]dto.ApplicationView, 0, len(items))
	if len(items) == 0 {
		return views, nil
	}
	petIDs := make(map[uint]struct{})
	petIDList := make([]uint, 0)
	for _, a := range items {
		if _, ok := petIDs[a.PetID]; !ok {
			petIDs[a.PetID] = struct{}{}
			petIDList = append(petIDList, a.PetID)
		}
	}
	pets, err := s.petRepo.ListByIDs(petIDList)
	if err != nil {
		return nil, fmt.Errorf("application view pets: %w", err)
	}
	petMap := make(map[uint]*model.Pet, len(pets))
	for i := range pets {
		petMap[pets[i].ID] = &pets[i]
	}
	holderIDs := make(map[uint]struct{})
	holderIDList := make([]uint, 0)
	for _, p := range petMap {
		if p.ReservedApplicationID != 0 {
			holderIDs[p.ReservedApplicationID] = struct{}{}
			holderIDList = append(holderIDList, p.ReservedApplicationID)
		}
	}
	holders, err := s.repo.ListByIDs(holderIDList)
	if err != nil {
		return nil, fmt.Errorf("application view holders: %w", err)
	}
	holderAppMap := make(map[uint]*model.AdoptionApplication, len(holders))
	for i := range holders {
		holderAppMap[holders[i].ID] = &holders[i]
	}
	holderUserIDs := make(map[uint]struct{})
	holderUserIDList := make([]uint, 0)
	for _, ha := range holderAppMap {
		if _, ok := holderUserIDs[ha.UserID]; !ok {
			holderUserIDs[ha.UserID] = struct{}{}
			holderUserIDList = append(holderUserIDList, ha.UserID)
		}
	}
	users, err := s.userRepo.ListByIDs(holderUserIDList)
	if err != nil {
		return nil, fmt.Errorf("application view users: %w", err)
	}
	userMap := make(map[uint]*model.User, len(users))
	for i := range users {
		userMap[users[i].ID] = &users[i]
	}
	for _, a := range items {
		v := dto.ApplicationView{AdoptionApplication: a}
		v.EndReasonText = EndReasonText(a.EndReason)
		if a.Status == constants.AppStatusWaitlisted {
			v.WaitlistActive = true
		}
		if p, ok := petMap[a.PetID]; ok {
			v.PetName = p.Name
			v.PetStatus = p.Status
			if p.ReservedApplicationID != 0 {
				if ha, ok2 := holderAppMap[p.ReservedApplicationID]; ok2 {
					v.ReservedUserID = ha.UserID
					if u, ok3 := userMap[ha.UserID]; ok3 {
						v.ReservedName = displayName(u)
					}
				}
			}
		}
		views = append(views, v)
	}
	return views, nil
}

// EndReasonText maps an end-reason code to Chinese text.
func EndReasonText(reason string) string {
	switch reason {
	case constants.AppEndRejected:
		return "审核未通过"
	case constants.AppEndFinalRejected:
		return "预留后最终拒绝"
	case constants.AppEndHolderAbandoned:
		return "获选人放弃"
	case constants.AppEndOrgCancelled:
		return "机构取消预留"
	case constants.AppEndAdoptedByOther:
		return "名额已被其他领养人获得"
	case constants.AppEndFinalAdopted:
		return "最终领养确认"
	default:
		return ""
	}
}

func displayName(u *model.User) string {
	if u.Nickname != "" {
		return u.Nickname
	}
	return u.Username
}
