package interceptor

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc"

	"pkg/jwtManager"
	"server/internal/utils/auth"
	"server/internal/utils/contextKeys"
	"server/internal/utils/errors"
)

type AuthInterceptor struct {
	disableAuthorizationPaths []string
}

func NewAuthInterceptor(accessibleRoutes []string) *AuthInterceptor {
	return &AuthInterceptor{accessibleRoutes}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx, err := interceptor.authorize(ctx, info.FullMethod, req)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Раньше здесь пытались достать accessToken прямо из grpc.ServerStream (а не из
		// распарсенного сообщения) — в проекте не было ни одного streaming-метода, поэтому это
		// ни разу не сработало бы на практике: authorize(..., stream) сериализовал в JSON сам
		// объект стрима, а не тело запроса, поэтому accessToken всегда оказывался пустым.
		// Авторизуем на первом реальном RecvMsg (см. authenticatedServerStream) — раньше interceptor
		// не может знать конкретный Go-тип входящего сообщения. Заодно оборачиваем Context(),
		// иначе обогащённый ctx (UserID/DeviceID) тоже никогда не долетал до хендлера — в отличие
		// от Unary(), server stream не принимает ctx как отдельный параметр.
		return handler(srv, &authenticatedServerStream{
			ServerStream: stream,
			interceptor:  interceptor,
			method:       info.FullMethod,
			ctx:          stream.Context(),
		})
	}
}

// authenticatedServerStream — обёртка над grpc.ServerStream, которая на ПЕРВОМ входящем
// сообщении достаёт accessToken (тем же кодом, что и Unary()) и подменяет Context() на
// обогащённый UserID/DeviceID из токена.
type authenticatedServerStream struct {
	grpc.ServerStream
	interceptor *AuthInterceptor
	method      string
	ctx         context.Context
	authorized  bool
}

func (s *authenticatedServerStream) Context() context.Context {
	return s.ctx
}

func (s *authenticatedServerStream) RecvMsg(m any) error {
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	if s.authorized {
		return nil
	}
	ctx, err := s.interceptor.authorize(s.ctx, s.method, m)
	if err != nil {
		return err
	}
	s.ctx = ctx
	s.authorized = true
	return nil
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string, req any) (context.Context, error) {

	// Игнорируем авторизацию для заранее перечисленных методов
	for _, accessibleRoute := range interceptor.disableAuthorizationPaths {
		if method == accessibleRoute {
			// everyone can access
			return ctx, nil
		}
	}

	type authNecessaryFields struct {
		AccessToken string `json:"accessToken"`
	}

	// Сериализуем запрос в JSON
	requestJson, err := json.Marshal(req)
	if err != nil {
		return ctx, errors.BadRequest.Wrap(err)
	}

	// Десериализуем JSON в структуру с интересующим полем
	var authFields authNecessaryFields
	err = json.Unmarshal(requestJson, &authFields)
	if err != nil {
		return ctx, errors.BadRequest.Wrap(err)
	}

	// Парсим токен
	claims, err := jwtManager.ParseToken[auth.Claims](ctx, authFields.AccessToken, jwtManager.AccessToken)
	if err != nil {

		// Если access-токен истек или некорректен, просим клиента обновить его через refresh-токен,
		// не разлогинивая пользователя сразу
		return ctx, errors.NeedToRefreshToken.Wrap(err)
	}

	// Добавляем в контекст UUID пользователя
	// Устанавливаем данные из токена в контекст
	ctx = contextKeys.SetUserID(ctx, claims.UserID)
	ctx = contextKeys.SetDeviceID(ctx, claims.DeviceID)

	return ctx, nil
}
