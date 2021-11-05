package repository

import (
	"io"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	uuidgen "github.com/google/uuid"
	"github.com/pkg/errors"
)

// GenerateUUID returns a random UUID
func GenerateUUID(rand io.Reader) (ref.UUID, error) {
	uuidgen.SetRand(rand)

	uuid, err := uuidgen.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "Could not generate UUID")
	}

	return ref.UUID(uuid.String()), nil
}
