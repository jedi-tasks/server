package auth

import (
	"context"

	"github.com/google/uuid"

	"server/internal/utils/errors"
)

type Claims struct {
	UserID   uuid.UUID
	DeviceID string
}

func (c *Claims) Validate(ctx context.Context) error {
	if c.UserID == uuid.Nil {
		return errors.InternalServer.New("UserID is empty").WithContextParams(ctx)
	}
	if c.DeviceID == "" {
		return errors.InternalServer.New("DeviceID is empty").WithContextParams(ctx)
	}
	return nil
}
