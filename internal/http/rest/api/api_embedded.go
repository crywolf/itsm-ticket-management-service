package api

import "github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"

// EmbeddedResources contain information about embedded objects
type EmbeddedResources map[string]interface{}

// NewEmbeddedBasicUser creates new initialized EmbeddedBasicUser
func NewEmbeddedBasicUser(user user.BasicUser) *EmbeddedBasicUser {
	basicUser := BasicUser{
		UUID:             user.UUID().String(),
		ExternalUserUUID: user.ExternalUserUUID,
		Name:             user.Name,
		Surname:          user.Surname,
		OrgDisplayName:   user.OrgDisplayName,
		OrgName:          user.OrgName,
	}

	return &EmbeddedBasicUser{
		BasicUser:             basicUser,
		EmbeddedResourceLinks: &HypermediaLinks{},
	}
}

// EmbeddedBasicUser wraps BasicUser with hypermedia links
type EmbeddedBasicUser struct {
	BasicUser
	EmbeddedResourceLinks `json:"_links"`
}

// UUID returns UUID of the BasicUser
func (e EmbeddedBasicUser) UUID() string {
	return e.BasicUser.UUID
}
