package fieldengineer

import (
	"fmt"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/embedded"
	tsession "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer/time_session"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/types"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
)

// FieldEngineer represents a user who can solve the tickets
type FieldEngineer struct {
	uuid ref.UUID

	BasicUser user.BasicUser

	openTimeSession *tsession.TimeSession

	TimeSessions []ref.UUID

	CreatedUpdated types.CreatedUpdated
}

// UUID getter
func (e FieldEngineer) UUID() ref.UUID {
	return e.uuid
}

// SetUUID returns error if UUID was already set
func (e *FieldEngineer) SetUUID(v ref.UUID) error {
	if !e.uuid.IsZero() {
		return fmt.Errorf("field engineer: cannot set UUID, it was already set (%s)", e.uuid)
	}
	e.uuid = v
	return nil
}

// OpenTimeSession returns open time session if any or nil pointer
func (e FieldEngineer) OpenTimeSession() *tsession.TimeSession {
	return e.openTimeSession
}

// SetOpenTimeSession sets open time session (do not use in the domain, method is used by repository)
func (e *FieldEngineer) SetOpenTimeSession(openTimeSession *tsession.TimeSession) {
	e.openTimeSession = openTimeSession
}

// HasOpenTimeSession returns true if the field engineer has an open time session
func (e FieldEngineer) HasOpenTimeSession() bool {
	return e.openTimeSession != nil
}

// EmbeddedResources returns list of other objects that are 'embedded' in the ticket
func (e FieldEngineer) EmbeddedResources(actor actor.Actor) []embedded.Resource {
	var resources []embedded.Resource
	resources = append(resources, e.CreatedUpdated.EmbeddedResources()...)
	return resources
}
