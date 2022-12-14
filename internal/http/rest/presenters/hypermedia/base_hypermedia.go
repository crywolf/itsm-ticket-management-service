package hypermedia

import (
	"net/url"

	"github.com/crywolf/itsm-ticket-management-service/internal/domain/user/actor"
)

// BaseHypermediaMapper is a base hypermedia mapping object to be included (via object composition) in specific hypermedia implementation
type BaseHypermediaMapper struct {
	serverAddr string
	currentURL *url.URL
	actor      actor.Actor
}

// NewBaseHypermedia returns base hypermedia object to be included in specific hypermedia implementation
func NewBaseHypermedia(serverAddr string, currentURL *url.URL, actor actor.Actor) *BaseHypermediaMapper {
	return &BaseHypermediaMapper{
		serverAddr: serverAddr,
		currentURL: currentURL,
		actor:      actor,
	}
}

// SelfLink returns 'self' link URL
func (b BaseHypermediaMapper) SelfLink() string {
	return b.serverAddr + b.currentURL.String()
}

// RequestURL returns current URL
func (b BaseHypermediaMapper) RequestURL() *url.URL {
	return b.currentURL
}

// ServerAddr returns server URL
func (b BaseHypermediaMapper) ServerAddr() string {
	return b.serverAddr
}

// Actor returns user who initiated current API call
func (b BaseHypermediaMapper) Actor() actor.Actor {
	return b.actor
}
