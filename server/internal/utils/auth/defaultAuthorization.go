package auth

import (
	"context"
	"net/http"
	"server/internal/utils/errors"

	"go.opentelemetry.io/otel"

	pkgErrors "pkg/errors"
	"pkg/jwtManager"
	"server/internal/utils/contextKeys"
)

const authorizationHeader = "Authorization"

var tracer = otel.Tracer("/server/internal/utils/chain")

func DefaultAuthorization(ctx context.Context, r *http.Request) (context.Context, error) {
	_, span := tracer.Start(ctx, "DefaultAuthorization")
	defer span.End()

	// Пытаемся распарсить токен
	claims, jwtErr := jwtManager.ParseToken[Claims](ctx, r.Header.Get(authorizationHeader), jwtManager.AccessToken)
	if jwtErr != nil && !pkgErrors.Is(jwtErr, jwtManager.ErrTokenExpired) {
		return ctx, errors.Forbidden.Wrap(jwtErr)
	}

	// Валидируем наличие всех полей в токене
	if err := claims.Validate(ctx); err != nil {
		return ctx, err
	}

	// Устанавливаем данные из токена в контекст
	ctx = contextKeys.SetUserID(ctx, claims.UserID)
	ctx = contextKeys.SetDeviceID(ctx, claims.DeviceID)

	return ctx, jwtErr
}
