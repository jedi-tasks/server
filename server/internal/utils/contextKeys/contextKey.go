package contextKeys

import (
	"context"

	"github.com/google/uuid"
)

type contextKey int

const (
	deviceIDKey contextKey = iota + 1
	userIDKey
	UserInfoKey
)

func SetDeviceID(ctx context.Context, deviceID string) context.Context {
	return context.WithValue(ctx, deviceIDKey, deviceID)
}

func SetUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetDeviceID(ctx context.Context) *string {
	if v, ok := ctx.Value(deviceIDKey).(string); ok {
		return &v
	}
	return nil
}

func GetUserID(ctx context.Context) *uuid.UUID {
	if v, ok := ctx.Value(userIDKey).(uuid.UUID); ok {
		return &v
	}
	return nil
}
