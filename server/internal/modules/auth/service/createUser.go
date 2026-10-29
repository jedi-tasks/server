package service

import (
	"context"
	"time"

	"pkg/passwordManager"
	"pkg/slices"
	"server/internal/utils/errors"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/auth/model"
	userModel "server/internal/modules/user/model"

	"github.com/google/uuid"
)

// createUser создает пользователя, фиксирует его в аудит-логе и создает Inbox, приоритеты и длительности.
// Вызывается внутри транзакции
func (s *AuthService) createUser(ctx context.Context, req model.CreateUserReq) (userID uuid.UUID, err error) {

	// Проверяем, есть ли пользователь в бд с таким email
	if _users, err := s.userRepository.GetUsers(ctx, userModel.GetUsersReq{Emails: []string{req.Email}}); err != nil { //nolint:exhaustruct
		return userID, err
	} else if len(_users) != 0 {
		return userID, errors.Forbidden.New("User with this email is already registered").
			WithContextParams(ctx).
			WithParams("email", req.Email).
			WithCustomHumanText("Пользователь с таким email уже зарегистрирован")
	}

	// Генерируем соль для пароля
	userSalt, err := passwordManager.GenerateRandomSalt()
	if err != nil {
		return userID, err
	}

	// Получаем хэш пароля пользователя
	passwordHash, err := passwordManager.CreateNewPassword([]byte(req.Password), s.generalSalt, userSalt)
	if err != nil {
		return userID, err
	}

	// Создаем пользователя
	userID = uuid.New()
	datetimeCreate := time.Now()
	err = s.userRepository.CreateUser(ctx, userModel.CreateReq{
		ID:           userID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		PasswordSalt: userSalt,
		TimeCreate:   datetimeCreate,
	})
	if err != nil {
		return userID, err
	}

	// Получаем созданного пользователя из БД для слепка "после" в аудит-логе (хэш и соль пароля исключены тегом json:"-")
	userAfter, err := slices.FirstWithError(s.userRepository.GetUsers(ctx, userModel.GetUsersReq{ //nolint:exhaustruct
		IDs: []uuid.UUID{userID},
	}))
	if err != nil {
		return userID, err
	}

	// Фиксируем создание пользователя в аудит-логе
	if _, err = s.auditLogService.TrackMutation(ctx, auditLogModel.TrackMutationReq{
		Entity:   auditLogEntity.User,
		Method:   auditLogMethod.Create,
		EntityID: userID.String(),
		Before:   nil,
		After:    userAfter,
		UserID:   userID,
		DeviceID: req.DeviceID,
	}); err != nil {
		return userID, err
	}

	// Создаем Inbox, приоритеты и длительности по умолчанию
	for _, defaultsCreator := range s.defaultsCreators {
		if err = defaultsCreator.CreateDefaults(ctx, userID, req.DeviceID, datetimeCreate); err != nil {
			return userID, err
		}
	}

	return userID, nil
}
