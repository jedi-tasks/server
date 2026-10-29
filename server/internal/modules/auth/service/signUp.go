package service

import (
	"context"

	"server/internal/utils/errors"

	"server/internal/modules/auth/model"
	"server/internal/modules/auth/service/utils"
	userModel "server/internal/modules/user/model"
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

	return accessData, s.generalRepository.WithinTransaction(ctx, func(ctx context.Context) error {

		// Создаем пользователя с данными по умолчанию
		accessData.ID, err = s.createUser(ctx, model.CreateUserReq{
			Name:     loginData.Name,
			Email:    loginData.Email,
			Password: loginData.Password,
			DeviceID: loginData.DeviceID,
		})
		if err != nil {
			return err
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
