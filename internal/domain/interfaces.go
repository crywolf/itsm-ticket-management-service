package domain

import (
	"time"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/types"
)

// Clock provides Now method to enable mocking
type Clock interface {
	// Now returns current time
	Now() time.Time

	// NowFormatted returns time in RFC3339 format
	NowFormatted() types.DateTime
}
