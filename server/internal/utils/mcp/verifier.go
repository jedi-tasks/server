package mcp

import (
	"context"
	"fmt"
	"net/http"

	sdkAuth "github.com/modelcontextprotocol/go-sdk/auth"

	"pkg/jwtManager"
	"server/internal/utils/auth"
)

// DeviceIDKey — ключ идентификатора устройства в TokenInfo.Extra
const DeviceIDKey = "deviceID"

// NewTokenVerifier проверяет Bearer-токен как access-JWT сервера: подпись и срок проверяет
// jwtManager, поэтому Expiration в TokenInfo не заполняется (AllowMissingExpiration)
func NewTokenVerifier() sdkAuth.TokenVerifier {
	return func(ctx context.Context, token string, _ *http.Request) (*sdkAuth.TokenInfo, error) {
		claims, err := jwtManager.ParseToken[auth.Claims](ctx, token, jwtManager.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", sdkAuth.ErrInvalidToken, err.Error())
		}
		if err = claims.Validate(ctx); err != nil {
			return nil, fmt.Errorf("%w: %s", sdkAuth.ErrInvalidToken, err.Error())
		}
		return &sdkAuth.TokenInfo{ //nolint:exhaustruct
			UserID: claims.UserID.String(),
			Extra:  map[string]any{DeviceIDKey: claims.DeviceID},
		}, nil
	}
}
