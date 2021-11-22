package memory

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/types"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository"
)

// Clock provides Now method to enable mocking
type Clock interface {
	Now() time.Time
}

// RepositoryMemory keeps data in memory
type RepositoryMemory struct {
	Rand      io.Reader
	Clock     Clock
	incidents []Incident
}

// AddIncident adds the given incident to the repository
func (m *RepositoryMemory) AddIncident(_ context.Context, _ ref.ChannelID, inc incident.Incident) (ref.UUID, error) {
	id, err := repository.GenerateUUID(m.Rand)
	if err != nil {
		log.Fatal(err)
	}

	now := m.Clock.Now().Format(time.RFC3339)

	storedInc := Incident{
		ID:               id.String(),
		Number:           inc.Number,
		ExternalID:       inc.ExternalID,
		ShortDescription: inc.ShortDescription,
		Description:      inc.Description,
		State:            incident.StateNew.String(),
		CreatedBy:        inc.CreatedUpdated.CreatedBy().String(),
		CreatedAt:        now,
		UpdatedBy:        inc.CreatedUpdated.UpdatedBy().String(),
		UpdatedAt:        now,
	}
	m.incidents = append(m.incidents, storedInc)

	return id, nil
}

// GetIncident returns the incident with given ID from the repository
func (m *RepositoryMemory) GetIncident(_ context.Context, _ ref.ChannelID, ID ref.UUID) (incident.Incident, error) {
	var inc incident.Incident
	var err error

	for i := range m.incidents {
		if m.incidents[i].ID == ID.String() {
			storedInc := m.incidents[i]

			inc, err = m.convertStoredToDomainIncident(storedInc)
			if err != nil {
				return incident.Incident{}, err
			}

			return inc, nil
		}
	}

	return incident.Incident{}, domain.WrapErrorf(ErrNotFound, domain.ErrorCodeNotFound, "repo GetIncident")
}

// ListIncidents returns the list of incidents from the repository
func (m *RepositoryMemory) ListIncidents(_ context.Context, _ ref.ChannelID) ([]incident.Incident, error) {
	var list []incident.Incident

	for _, storedInc := range m.incidents {
		inc, err := m.convertStoredToDomainIncident(storedInc)
		if err != nil {
			return nil, err
		}

		list = append(list, inc)
	}

	return list, nil
}

func (m RepositoryMemory) convertStoredToDomainIncident(storedInc Incident) (incident.Incident, error) {
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
