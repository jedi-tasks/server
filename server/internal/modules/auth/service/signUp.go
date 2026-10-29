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
	"server/internal/modules/auth/service/utils"
	userModel "server/internal/modules/user/model"

	"github.com/google/uuid"
)

// SignUp регистрирует пользователя и возвращает токены доступа
func (s *AuthService) SignUp(ctx context.Context, loginData model.SignUpReq) (accessData model.AuthRes, err error) {
	ctx, span := tracer.Start(ctx, "SignUp")
	defer span.End()

	// Проверяем, что регистрация открыта
	if !s.signUpEnabled {
		return accessData, errors.Forbidden.New("Sign up is disabled").
			WithContextParams(ctx).
			WithCustomHumanText("Регистрация закрыта")
	}

	// Проверяем, есть ли пользователь в бд с таким email
	if _users, err := s.userRepository.GetUsers(ctx, userModel.GetUsersReq{Emails: []string{loginData.Email}}); err != nil { //nolint:exhaustruct
		return accessData, err
	} else if len(_users) != 0 {
		return accessData, errors.Forbidden.New("User with this email is already registered").
			WithContextParams(ctx).
			WithParams("email", loginData.Email).
			WithCustomHumanText("Пользователь с таким email уже зарегистрирован")
	}

	// Генерируем соль для пароля
	userSalt, err := passwordManager.GenerateRandomSalt()
	if err != nil {
		return accessData, err
	}

	// Получаем хэш пароля пользователя
	passwordHash, err := passwordManager.CreateNewPassword([]byte(loginData.Password), s.generalSalt, userSalt)
	if err != nil {
		return accessData, err
	}

	return accessData, s.generalRepository.WithinTransaction(ctx, func(ctx context.Context) error {

		accessData.ID = uuid.New()

		// Создаем пользователя
		datetimeCreate := time.Now()
		err = s.userRepository.CreateUser(ctx, userModel.CreateReq{
			ID:              accessData.ID,
			Name:            loginData.Name,
			Email:           loginData.Email,
			PasswordHash:    passwordHash,
			PasswordSalt:    userSalt,
			TimeCreate:      datetimeCreate,
		})
		if err != nil {
			return err
		}

		// Получаем созданного пользователя из БД для слепка "после" в аудит-логе (хэш и соль пароля исключены тегом json:"-")
		userAfter, err := slices.FirstWithError(s.userRepository.GetUsers(ctx, userModel.GetUsersReq{ //nolint:exhaustruct
			IDs: []uuid.UUID{accessData.ID},
		}))
		if err != nil {
			return err
		}

		// Фиксируем регистрацию пользователя в аудит-логе
		if _, err = s.auditLogService.TrackMutation(ctx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.User,
			Method:   auditLogMethod.Create,
			EntityID: accessData.ID.String(),
			Before:   nil,
			After:    userAfter,
			UserID:   accessData.ID,
			DeviceID: loginData.DeviceID,
		}); err != nil {
			return err
		}

		// Создаем Inbox, приоритеты и длительности по умолчанию
		for _, defaultsCreator := range s.defaultsCreators {
			if err = defaultsCreator.CreateDefaults(ctx, accessData.ID, loginData.DeviceID, datetimeCreate); err != nil {
				return err
			}
		}

		// Создаем пару токенов
		accessData.Tokens, err = utils.CreatePairTokens(ctx, accessData.ID, loginData.DeviceID)
		if err != nil {
			return err
		}

		// Создаем или обновляем девайс пользователя
		err = s.upsertDevice(ctx, userModel.Device{
			DeviceInformation:      loginData.Device,
			NotificationToken:      nil,
			RefreshToken:           accessData.RefreshToken,
			UserID:                 accessData.ID,
			DeviceID:               loginData.DeviceID,
			ApplicationInformation: loginData.Application,
		})
		if err != nil {
			return err
		}

		return nil
	})
}
