package memory

import (
	"context"
	"io"
	"time"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident/timelog"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/types"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository"
)

// Clock provides Now method to enable mocking
type Clock interface {
	// Now returns current time
	Now() time.Time

	// NowFormatted returns time in RFC3339 format
	NowFormatted() types.DateTime
}

// IncidentRepositoryMemory keeps data in memory
type IncidentRepositoryMemory struct {
	basicUserRepository     repository.BasicUserRepository
	fieldEngineerRepository repository.FieldEngineerRepository
	Rand                    io.Reader
	clock                   Clock
	incidents               []Incident
	timelogs                map[string]Timelog
}

// NewIncidentRepositoryMemory returns new initialized repository
func NewIncidentRepositoryMemory(clock Clock, basicUserRepo repository.BasicUserRepository, fieldEngineerRepository repository.FieldEngineerRepository) *IncidentRepositoryMemory {
	return &IncidentRepositoryMemory{
		basicUserRepository:     basicUserRepo,
		fieldEngineerRepository: fieldEngineerRepository,
		clock:                   clock,
		timelogs:                make(map[string]Timelog),
	}
}

// AddIncident adds the given incident to the repository
func (r *IncidentRepositoryMemory) AddIncident(_ context.Context, _ ref.ChannelID, inc incident.Incident) (ref.UUID, error) {
	now := r.clock.NowFormatted().String()

	incidentID, err := repository.GenerateUUID(r.Rand)
	if err != nil {
		return ref.UUID(""), err
	}

	feUUID := ""
	if inc.FieldEngineerID != nil {
		feUUID = inc.FieldEngineerID.String()
	}

	storedInc := Incident{
		ID:               incidentID.String(),
		Number:           inc.Number,
		ExternalID:       inc.ExternalID,
		ShortDescription: inc.ShortDescription,
		Description:      inc.Description,
		FieldEngineerID:  feUUID,
		State:            incident.StateNew.String(),
		CreatedBy:        inc.CreatedUpdated.CreatedByID().String(),
		CreatedAt:        now,
		UpdatedBy:        inc.CreatedUpdated.UpdatedByID().String(),
		UpdatedAt:        now,
	}
	r.incidents = append(r.incidents, storedInc)

	return incidentID, nil
}

// UpdateIncident updates the given incident in the repository
func (r *IncidentRepositoryMemory) UpdateIncident(_ context.Context, _ ref.ChannelID, inc incident.Incident) (ref.UUID, error) {
	var err error
	now := r.clock.NowFormatted().String()

	if inc.HasOpenTimelog() {
		openTimelog := inc.OpenTimelog()
		createdAt := openTimelog.CreatedUpdated.CreatedAt().String()
		updatedAt := openTimelog.CreatedUpdated.UpdatedAt().String()

		timelogID := openTimelog.UUID()
		if timelogID.IsZero() { // newly opened => set new UUID
			timelogID, err = repository.GenerateUUID(r.Rand)
			if err != nil {
				return ref.UUID(""), err
			}

			createdAt = now
			updatedAt = now

			inc.Timelogs = append(inc.Timelogs, timelogID)
		}

		storedTimelog := Timelog{
			ID:           timelogID.String(),
			Remote:       openTimelog.Remote,
			Work:         openTimelog.Work,
			VisitSummary: openTimelog.VisitSummary,
			CreatedBy:    openTimelog.CreatedUpdated.CreatedByID().String(),
			CreatedAt:    createdAt,
			UpdatedBy:    openTimelog.CreatedUpdated.UpdatedByID().String(),
			UpdatedAt:    updatedAt,
		}

		r.timelogs[storedTimelog.ID] = storedTimelog
	}

	var timelogUUIDs []string
	for _, timelogID := range inc.Timelogs {
		timelogUUIDs = append(timelogUUIDs, timelogID.String())
	}

	feUUID := ""
	if inc.FieldEngineerID != nil {
		feUUID = inc.FieldEngineerID.String()
	}

	storedInc := Incident{
		ID:               inc.UUID().String(),
		Number:           inc.Number,
		ExternalID:       inc.ExternalID,
		ShortDescription: inc.ShortDescription,
		Description:      inc.Description,
		FieldEngineerID:  feUUID,
		State:            incident.StateNew.String(),
		Timelogs:         timelogUUIDs,
		CreatedBy:        inc.CreatedUpdated.CreatedByID().String(),
		CreatedAt:        inc.CreatedUpdated.CreatedAt().String(),
		UpdatedBy:        inc.CreatedUpdated.UpdatedByID().String(),
		UpdatedAt:        now,
	}

	for i := range r.incidents {
		if r.incidents[i].ID == inc.UUID().String() {
			r.incidents[i] = storedInc
			return inc.UUID(), nil
		}
	}

	return inc.UUID(), domain.WrapErrorf(ErrNotFound, domain.ErrorCodeNotFound, "error updating incident in repository")
}

// GetIncident returns the incident with given ID from the repository
func (r *IncidentRepositoryMemory) GetIncident(ctx context.Context, channelID ref.ChannelID, ID ref.UUID) (incident.Incident, error) {
	var inc incident.Incident
	var err error

	for i := range r.incidents {
		if r.incidents[i].ID == ID.String() {
			storedInc := r.incidents[i]

			inc, err = r.convertStoredToDomainIncident(ctx, channelID, storedInc)
			if err != nil {
				return incident.Incident{}, err
			}

			return inc, nil
		}
	}

	return incident.Incident{}, domain.WrapErrorf(ErrNotFound, domain.ErrorCodeNotFound, "error loading incident from repository")
}

// ListIncidents returns the list of incidents from the repository
func (r *IncidentRepositoryMemory) ListIncidents(ctx context.Context, channelID ref.ChannelID, page, itemsPerPage uint) (repository.IncidentList, error) {
	var list []incident.Incident

	total := len(r.incidents)

	pagination := repository.NewPagination(total, page, itemsPerPage)

	firstElementIndex := pagination.FirstElementIndex
	lastElementIndex := pagination.LastElementIndex

	var perPageList []Incident
	if total > 0 {
		perPageList = r.incidents[firstElementIndex : lastElementIndex+1]
	}

	for _, storedInc := range perPageList {
		inc, err := r.convertStoredToDomainIncident(ctx, channelID, storedInc)
		if err != nil {
			return repository.IncidentList{}, err
		}

		list = append(list, inc)
	}

	incidentList := repository.IncidentList{
		Result:     list,
		Pagination: pagination,
	}
	return incidentList, nil
}

func (r IncidentRepositoryMemory) convertStoredToDomainIncident(ctx context.Context, channelID ref.ChannelID, storedInc Incident) (incident.Incident, error) {
	var inc incident.Incident
	errMsg := "error loading incident from repository (%s)"

	err := inc.SetUUID(ref.UUID(storedInc.ID))
	if err != nil {
		return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedInc.ID")
	}

	inc.Number = storedInc.Number
	inc.ExternalID = storedInc.ExternalID
	inc.ShortDescription = storedInc.ShortDescription
	inc.Description = storedInc.Description

	if storedInc.FieldEngineerID != "" {
		feUUID := ref.UUID(storedInc.FieldEngineerID)
		inc.FieldEngineerID = &feUUID
	}

	// load and set open timelog if any
	var timelogUUIDs []ref.UUID
	for _, timelogID := range storedInc.Timelogs {
		timelogUUIDs = append(timelogUUIDs, ref.UUID(timelogID))
	}

	inc.Timelogs = timelogUUIDs

	for _, timelogID := range storedInc.Timelogs {
		storedTimelog := r.timelogs[timelogID]

		if storedTimelog.Work == 0 { // timelog is open
			openTimelog := &timelog.Timelog{
				Remote:       storedTimelog.Remote,
				Work:         storedTimelog.Work,
				VisitSummary: storedTimelog.VisitSummary,
			}

			createdByUser, err := r.basicUserRepository.GetBasicUser(ctx, channelID, ref.UUID(storedTimelog.CreatedBy))
			if err != nil {
				return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedTimelog.createdBy")
			}

			err = openTimelog.CreatedUpdated.SetCreated(createdByUser, types.DateTime(storedTimelog.CreatedAt))
			if err != nil {
				return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedTimelog.createdAt")
			}

			updatedByUser, err := r.basicUserRepository.GetBasicUser(ctx, channelID, ref.UUID(storedTimelog.UpdatedBy))
			if err != nil {
				return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedTimelog.UpdatedBy")
			}

			err = openTimelog.CreatedUpdated.SetUpdated(updatedByUser, types.DateTime(storedTimelog.UpdatedAt))
			if err != nil {
				return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedTimelog.UpdatedAt")
			}

			inc.SetOpenTimelog(openTimelog)
			break
		}
	}

	state, err := incident.NewStateFromString(storedInc.State)
	if err != nil {
		return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "state")
	}

	err = inc.SetState(state)
	if err != nil {
		return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "state")
	}

	createdByUser, err := r.basicUserRepository.GetBasicUser(ctx, channelID, ref.UUID(storedInc.CreatedBy))
	if err != nil {
		return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedInc.CreatedBy")
	}

	err = inc.CreatedUpdated.SetCreated(createdByUser, types.DateTime(storedInc.CreatedAt))
	if err != nil {
		return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedInc.CreatedAt")
	}

	updatedByUser, err := r.basicUserRepository.GetBasicUser(ctx, channelID, ref.UUID(storedInc.UpdatedBy))
	if err != nil {
		return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedInc.UpdatedBy")
	}

	err = inc.CreatedUpdated.SetUpdated(updatedByUser, types.DateTime(storedInc.UpdatedAt))
	if err != nil {
		return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg, "storedInc.UpdatedAt")
	}

	return inc, nil
}
