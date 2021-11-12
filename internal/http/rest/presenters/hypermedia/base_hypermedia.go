package hypermedia

// BaseHypermediaMapper is a base hypermedia mapping object to be included (via object composition) in specific hypermedia implementation
type BaseHypermediaMapper struct {
	serverAddr string
	currentURL string
}

// NewBaseHypermedia returns base hypermedia object to be included in specific hypermedia implementation
func NewBaseHypermedia(serverAddr, currentURL string) *BaseHypermediaMapper {
	return &BaseHypermediaMapper{
		serverAddr: serverAddr,
		currentURL: currentURL,
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
