package api

// EmbeddedResourceLinks can append 'self' link to its '_links' object
type EmbeddedResourceLinks interface {
	AppendSelfLink(url string)
}
