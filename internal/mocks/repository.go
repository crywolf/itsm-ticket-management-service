package mocks

import (
	"time"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/types"
)

// FixedClock is used in testing. It implements Clock interface used in mock storages to create timestamp.
type FixedClock struct{}

// Now returns fixed time
func (FixedClock) Now() time.Time {
	tz, err := time.LoadLocation("Europe/Prague")
	if err != nil {
		panic(err)
	}
	return time.Date(2021, 4, 1, 12, 34, 56, 78, tz)
}

// NowFormatted returns time in RFC3339 format
func (c FixedClock) NowFormatted() types.DateTime {
	return types.DateTime(c.Now().Format(time.RFC3339))
}
