package presenters

// BaseHypermedia is a base hypermedia object to be included in specific hypermedia implementation
type BaseHypermedia struct {
	serverAddr string
	currentURL string
}

// NewBaseHypermedia returns new base hypermedia object to be included in specific hypermedia implementation
func NewBaseHypermedia(serverAddr, currentURL string) *BaseHypermedia {
	return &BaseHypermedia{
		serverAddr: serverAddr,
		currentURL: currentURL,
	}
}

// Self returns 'self' link
func (b BaseHypermedia) Self() string {
	return b.serverAddr + b.currentURL
}

// ServerAddr returns server URL
func (b BaseHypermedia) ServerAddr() string {
	return b.serverAddr
}
