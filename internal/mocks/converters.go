package mocks

import (
	"github.com/stretchr/testify/mock"
)

// PaginationParamsMock is a pagination parameters mock
type PaginationParamsMock struct {
	mock.Mock
}

// Page mock
func (p *PaginationParamsMock) Page() uint {
	args := p.Called()
	return args.Get(0).(uint)
}

// ItemsPerPage mock
func (p *PaginationParamsMock) ItemsPerPage() uint {
	args := p.Called()
	return args.Get(0).(uint)
}
