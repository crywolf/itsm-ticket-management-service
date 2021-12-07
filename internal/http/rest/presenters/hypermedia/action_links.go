package hypermedia

import "github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/api"

// ActionLinks maps allowed resource domain actions to hypermedia action links
type ActionLinks struct {
	mapper *BaseHypermediaMapper
	data   map[string]api.Link
}

// NewActionLinks creates new empty action links
func NewActionLinks(mapper *BaseHypermediaMapper) ActionLinks {
	return ActionLinks{
		mapper: mapper,
		data:   map[string]api.Link{},
	}
}

// Add adds new action link
func (l *ActionLinks) Add(domainAction, linkName, route string) {
	l.data[domainAction] = api.Link{
		Name: linkName,
		Href: l.mapper.ServerAddr() + route}
}

// Get reruns action link
func (l ActionLinks) Get(domainAction string) api.Link {
	return l.data[domainAction]
}
