package memory

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/incident"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/types"
	"github.com/KompiTech/itsm-ticket-management-service/internal/repository"
	"github.com/pkg/errors"
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

	storedInc := Incident{
		ID:               id.String(),
		ExternalID:       inc.ExternalID,
		ShortDescription: inc.ShortDescription,
		Description:      inc.Description,
		State:            incident.StateNew.String(),
		CreatedBy:        inc.GetCreatedBy().String(),
		CreatedAt:        m.Clock.Now().Format(time.RFC3339),
	}
	m.incidents = append(m.incidents, storedInc)

	return id, nil
}

// GetIncident returns the incident with given ID from the repository
func (m *RepositoryMemory) GetIncident(_ context.Context, _ ref.ChannelID, ID ref.UUID) (incident.Incident, error) {
	var inc incident.Incident

	for i := range m.incidents {
		if m.incidents[i].ID == ID.String() {
			si := m.incidents[i] // stored incident
			inc.UUID = ref.UUID(si.ID)
			inc.ExternalID = si.ExternalID
			inc.ShortDescription = si.ShortDescription
			inc.Description = si.Description

			state, err := incident.NewStateFromString(si.State)
			if err != nil {
				return inc, errors.Wrap(err, "error loading incident from the repository")
			}

			err = inc.SetState(state)
			if err != nil {
				return inc, errors.Wrap(err, "error loading incident from the repository")
			}

			err = inc.SetCreatedBy(ref.ExternalUserUUID(si.CreatedBy))
			if err != nil {
				return inc, errors.Wrap(err, "error loading incident from the repository")
			}

			err = inc.SetCreatedAt(types.DateTime(si.CreatedAt))
			if err != nil {
				return inc, errors.Wrap(err, "error loading incident from the repository")
			}

			return inc, nil
		}
	}

	return incident.Incident{}, repository.NewError(ErrNotFound.Error(), http.StatusNotFound)
}

// ListIncidents returns the incidents from the repository
func (m *RepositoryMemory) ListIncidents(_ context.Context, _ ref.ChannelID) ([]incident.Incident, error) {
	panic("implement me")
}
