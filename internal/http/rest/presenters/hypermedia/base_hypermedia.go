package hypermedia

import "github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"

// BaseHypermediaMapper is a base hypermedia mapping object to be included (via object composition) in specific hypermedia implementation
type BaseHypermediaMapper struct {
	serverAddr string
	currentURL string
	actor      user.Actor
}

// NewBaseHypermedia returns base hypermedia object to be included in specific hypermedia implementation
func NewBaseHypermedia(serverAddr, currentURL string, actor user.Actor) *BaseHypermediaMapper {
	return &BaseHypermediaMapper{
		serverAddr: serverAddr,
		currentURL: currentURL,
		actor:      actor,
	}
}

// SelfLink returns 'self' link URL
func (b BaseHypermediaMapper) SelfLink() string {
	return b.serverAddr + b.currentURL
}

// ServerAddr returns server URL
func (b BaseHypermediaMapper) ServerAddr() string {
	return b.serverAddr
}

// Actor returns user who initiated current API call
func (b BaseHypermediaMapper) Actor() user.Actor {
	return b.actor
}
