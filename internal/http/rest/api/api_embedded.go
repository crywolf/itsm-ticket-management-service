package api

import (
	fieldengineer "github.com/KompiTech/itsm-ticket-management-service/internal/domain/field_engineer"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
)

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

// NewEmbeddedFieldEngineer creates new initialized EmbeddedFieldEngineer
func NewEmbeddedFieldEngineer(user fieldengineer.FieldEngineer) *EmbeddedFieldEngineer {
	basicUser := BasicUser{
		UUID:             user.UUID().String(),
		ExternalUserUUID: user.BasicUser.ExternalUserUUID,
		Name:             user.BasicUser.Name,
		Surname:          user.BasicUser.Surname,
		OrgDisplayName:   user.BasicUser.OrgDisplayName,
		OrgName:          user.BasicUser.OrgName,
	}

	return &EmbeddedFieldEngineer{
		FieldEngineerID:       user.UUID().String(),
		BasicUser:             basicUser,
		EmbeddedResourceLinks: &HypermediaLinks{},
	}
}

// EmbeddedFieldEngineer wraps FieldEngineer with hypermedia links
type EmbeddedFieldEngineer struct {
	BasicUser
	FieldEngineerID       string `json:"uuid"`
	EmbeddedResourceLinks `json:"_links"`
}

// UUID returns UUID of the FieldEngineer
func (e EmbeddedFieldEngineer) UUID() string {
	return e.FieldEngineerID
}
