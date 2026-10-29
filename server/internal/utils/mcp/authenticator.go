package mcp

import (
	"context"
	"strings"
	"sync"

	"github.com/google/uuid"

	"server/internal/enum/osType"
	authModel "server/internal/modules/auth/model"
	userModel "server/internal/modules/user/model"
	"server/internal/utils/contextKeys"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

// DeviceID — идентификатор устройства агента, под ним правки попадают в аудит-лог
const DeviceID = "mcp"

// UserService — поиск пользователя по почте
type UserService interface {
	GetUsers(context.Context, userModel.GetUsersReq) ([]userModel.User, error)
}

// AuthService — создание устройства агента
type AuthService interface {
	EnsureDevice(context.Context, authModel.EnsureDeviceReq) (bool, error)
}

// Authenticator превращает Bearer-токен в пользователя и устройство в контексте запроса
type Authenticator struct {
	tokens      Tokens
	userService UserService
	authService AuthService
	version     string
	users       sync.Map // email → uuid.UUID, пользователи и их устройства уже проверены
}

// NewAuthenticator создает аутентификатор MCP
func NewAuthenticator(tokens Tokens, userService UserService, authService AuthService, version string) *Authenticator {
	return &Authenticator{tokens: tokens, userService: userService, authService: authService, version: version} //nolint:exhaustruct
}

// Authorize проверяет заголовок Authorization и возвращает контекст с UserID и DeviceID, как после
// gRPC-интерсептора авторизации. При первом запросе пользователя создает его устройство mcp
func (a *Authenticator) Authorize(ctx context.Context, authorization string) (context.Context, error) {

	// Достаем токен из "Bearer <токен>"
	token, found := strings.CutPrefix(authorization, "Bearer ")
	if !found || token == "" {
		return ctx, errors.NeedToLogout.New("Нет Bearer-токена")
	}
	email, ok := a.tokens.Lookup(token)
	if !ok {
		return ctx, errors.NeedToLogout.New("Неизвестный токен")
	}

	// Находим пользователя и его устройство, повторно — из памяти
	userID, err := a.userID(ctx, email)
	if err != nil {
		return ctx, err
	}

	ctx = contextKeys.SetUserID(ctx, userID)
	ctx = contextKeys.SetDeviceID(ctx, DeviceID)
	return ctx, nil
}

// userID возвращает идентификатор пользователя по почте и гарантирует, что устройство mcp создано
func (a *Authenticator) userID(ctx context.Context, email string) (uuid.UUID, error) {
	if cached, ok := a.users.Load(email); ok {
		return cached.(uuid.UUID), nil
	}

	users, err := a.userService.GetUsers(ctx, userModel.GetUsersReq{Emails: []string{email}}) //nolint:exhaustruct
	if err != nil {
		return uuid.Nil, err
	}
	if len(users) == 0 {
		return uuid.Nil, errors.NeedToLogout.New("Владелец токена не найден").WithParams("email", email)
	}
	userID := users[0].ID

	// Устройство без локальной копии: гейт синка его не проверяет
	_, err = a.authService.EnsureDevice(ctx, authModel.EnsureDeviceReq{
		UserID:   userID,
		DeviceID: DeviceID,
		Device: userModel.DeviceInformation{ //nolint:exhaustruct
			NameOS:     osType.Server,
			VersionOS:  a.version,
			DeviceName: "MCP",
			ModelName:  "tasks-mcp",
			Stateless:  true,
		},
		Application: userModel.ApplicationInformation{BundleID: "tasks.mcp", Version: a.version, Build: a.version},
	})
	if err != nil {
		return uuid.Nil, err
	}

	a.users.Store(email, userID)
	return userID, nil
}

// Necessary возвращает пользователя и устройство из контекста запроса для сервисных моделей
func Necessary(ctx context.Context) (necessary.NecessaryUserInformation, error) {
	var dest struct {
		Necessary necessary.NecessaryUserInformation
	}
	if err := necessary.ParseNecessary(ctx, &dest); err != nil {
		return dest.Necessary, err
	}
	return dest.Necessary, nil
}
