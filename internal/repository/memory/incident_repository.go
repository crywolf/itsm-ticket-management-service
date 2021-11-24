package memory

import (
	"context"
	"io"
	"log"
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
	Now() time.Time
}

// IncidentRepositoryMemory keeps data in memory
type IncidentRepositoryMemory struct {
	Rand      io.Reader
	Clock     Clock
	incidents []Incident
	timelogs  map[string]Timelog
}

// AddIncident adds the given incident to the repository
func (r *IncidentRepositoryMemory) AddIncident(_ context.Context, _ ref.ChannelID, inc incident.Incident) (ref.UUID, error) {
	id, err := repository.GenerateUUID(r.Rand)
	if err != nil {
		log.Fatal(err)
	}

	now := r.Clock.Now().Format(time.RFC3339)

	var timelogUUIDs []string
	for _, tmlg := range inc.Timelogs {
		timelogUUIDs = append(timelogUUIDs, tmlg.String())
	}

	storedInc := Incident{
		ID:               id.String(),
		Number:           inc.Number,
		ExternalID:       inc.ExternalID,
		ShortDescription: inc.ShortDescription,
		Description:      inc.Description,
		State:            incident.StateNew.String(),
		Timelogs:         timelogUUIDs,
		CreatedBy:        inc.CreatedUpdated.CreatedBy().String(),
		CreatedAt:        now,
		UpdatedBy:        inc.CreatedUpdated.UpdatedBy().String(),
		UpdatedAt:        now,
	}
	r.incidents = append(r.incidents, storedInc)

	return id, nil
}

// GetIncident returns the incident with given ID from the repository
func (r *IncidentRepositoryMemory) GetIncident(_ context.Context, _ ref.ChannelID, ID ref.UUID) (incident.Incident, error) {
	var inc incident.Incident
	var err error

	for i := range r.incidents {
		if r.incidents[i].ID == ID.String() {
			storedInc := r.incidents[i]

			inc, err = r.convertStoredToDomainIncident(storedInc)
			if err != nil {
				return incident.Incident{}, err
			}

			return inc, nil
		}
	}

	return incident.Incident{}, domain.WrapErrorf(ErrNotFound, domain.ErrorCodeNotFound, "repo GetIncident")
}

// ListIncidents returns the list of incidents from the repository
func (r *IncidentRepositoryMemory) ListIncidents(_ context.Context, _ ref.ChannelID) ([]incident.Incident, error) {
	var list []incident.Incident

	for _, storedInc := range r.incidents {
		inc, err := r.convertStoredToDomainIncident(storedInc)
		if err != nil {
			return nil, err
		}

		list = append(list, inc)
	}

	return list, nil
}

func (r IncidentRepositoryMemory) convertStoredToDomainIncident(storedInc Incident) (incident.Incident, error) {
	var inc incident.Incident
	errMsg := "error loading incident from the repository"

	err := inc.SetUUID(ref.UUID(storedInc.ID))
	if err != nil {
		return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg)
	}

	inc.Number = storedInc.Number
	inc.ExternalID = storedInc.ExternalID
	inc.ShortDescription = storedInc.ShortDescription
	inc.Description = storedInc.Description

	// load and set open timelog if any
	var timelogUUIDs []ref.UUID
	for _, tmlg := range storedInc.Timelogs {
		timelogUUIDs = append(timelogUUIDs, ref.UUID(tmlg))
	}

	inc.Timelogs = timelogUUIDs

	for _, timelogID := range storedInc.Timelogs {
		storedTmlg := r.timelogs[timelogID]
		if storedTmlg.Work > 0 { // timelog is open
			openTmlg := &timelog.Timelog{
				Remote:       storedTmlg.Remote,
				Work:         storedTmlg.Work,
				VisitSummary: storedTmlg.VisitSummary,
			}

			err = openTmlg.CreatedUpdated.SetCreated(ref.ExternalUserUUID(storedTmlg.CreatedBy), types.DateTime(storedTmlg.CreatedAt))
			if err != nil {
				return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg)
			}

			err = openTmlg.CreatedUpdated.SetUpdated(ref.ExternalUserUUID(storedTmlg.UpdatedBy), types.DateTime(storedTmlg.UpdatedAt))
			if err != nil {
				return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg)
			}

			inc.SetOpenTimelog(openTmlg)
			break
		}
	}

	state, err := incident.NewStateFromString(storedInc.State)
	if err != nil {
		return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg)
	}

	err = inc.SetState(state)
	if err != nil {
		return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg)
	}

	err = inc.CreatedUpdated.SetCreated(ref.ExternalUserUUID(storedInc.CreatedBy), types.DateTime(storedInc.CreatedAt))
	if err != nil {
		return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg)
	}

	err = inc.CreatedUpdated.SetUpdated(ref.ExternalUserUUID(storedInc.UpdatedBy), types.DateTime(storedInc.UpdatedAt))
	if err != nil {
		return incident.Incident{}, domain.WrapErrorf(err, domain.ErrorCodeUnknown, errMsg)
	}

	return inc, nil
}