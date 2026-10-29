package service

import (
	"context"
	"time"

	"pkg/jwtManager"
	"server/internal/utils/auth"
	"server/internal/utils/errors"

	"github.com/google/uuid"

	"server/internal/modules/auth/model"
	"server/internal/modules/auth/service/utils"
	userRepoModel "server/internal/modules/user/repository/model"
)

// refreshTokenGraceWindow - время, в течение которого предыдущий refresh-токен остаётся валидным
// после ротации. Компенсирует потерю ответа сервера в сети: клиент, не получивший новую пару токенов,
// может повторить запрос старым токеном и получить тот же (уже выпущенный) результат, вместо
// полной потери сессии
const refreshTokenGraceWindow = 5 * time.Minute

// RefreshTokens обновляет токены доступа в базе данных
func (s *AuthService) RefreshTokens(ctx context.Context, req model.RefreshTokensReq) (newTokens model.RefreshTokensRes, err error) {
	ctx, span := tracer.Start(ctx, "RefreshTokens")
	defer span.End()

	// RefreshTokens не проходит через auth-интерцептор (access-токен на этот момент
	// может быть просрочен), поэтому UserID и DeviceID получаем не из контекста,
	// а из самого refresh-токена, переданного клиентом
	claims, err := jwtManager.ParseToken[auth.Claims](ctx, req.Token, jwtManager.RefreshToken)
	if err != nil {
		return newTokens, err
	}

	// Получаем девайс по идентификатору пользователя и девайса
	devices, err := s.userRepository.GetDevices(ctx, userRepoModel.GetDevicesReq{ // nolint:exhaustruct
		DeviceIDs: []string{claims.DeviceID},
		UserIDs:   []uuid.UUID{claims.UserID},
	})
	if err != nil {
		return newTokens, err
	}
	if len(devices) == 0 {
		return newTokens, errors.NeedToLogout.New("Device not found").WithContextParams(ctx).
			WithCustomHumanText("Девайс не найден")
	}
	device := devices[0]

	switch {
	// Переданный токен совпадает с актуальным - обычная ротация
	case req.Token == device.RefreshToken:
		newTokens.Tokens, err = utils.CreatePairTokens(ctx, claims.UserID, claims.DeviceID)
		if err != nil {
			return newTokens, err
		}

		// Атомарно сдвигаем текущий токен в "предыдущий" (с grace-периодом) и записываем новый как актуальный
		if err = s.userRepository.RotateRefreshToken(ctx, claims.UserID, claims.DeviceID, newTokens.Tokens.RefreshToken, refreshTokenGraceWindow); err != nil {
			return newTokens, err
		}

	// Переданный токен - это уже отработавший предыдущий токен в пределах grace-периода: значит,
	// ответ на прошлый рефреш не дошёл до клиента по сети. Не ротируем повторно (иначе повторные
	// ретраи плодили бы новые ротации), а отдаём уже актуальный refresh-токен и свежий access-токен
	case device.PreviousRefreshToken != nil && req.Token == *device.PreviousRefreshToken &&
		device.PreviousRefreshTokenExpiresAt != nil && time.Now().Before(*device.PreviousRefreshTokenExpiresAt):

		pairTokens, err := utils.CreatePairTokens(ctx, claims.UserID, claims.DeviceID)
		if err != nil {
			return newTokens, err
		}
		newTokens.Tokens.AccessToken = pairTokens.AccessToken
		newTokens.Tokens.RefreshToken = device.RefreshToken

	// Переданный refresh-токен не найден ни как актуальный, ни как предыдущий в пределах grace-периода:
	// просим пользователя разлогиниться и авторизоваться заново
	default:
		return newTokens, errors.NeedToLogout.New("Refresh token not found").
			WithContextParams(ctx).
			WithCustomHumanText("Сессия истекла, необходимо авторизоваться заново")
	}

	// Обновляем данные у девайса
	if err = s.userRepository.UpdateDevice(ctx, userRepoModel.UpdateDeviceReq{
		UserID:       claims.UserID,
		DeviceID:     claims.DeviceID,
		RefreshToken: nil,
		DeviceInformation: userRepoModel.UpdateDeviceInformationReq{
			VersionOS: &req.Device.VersionOS,
			UserAgent: &req.Device.UserAgent,
			IPAddress: &req.Device.IPAddress,
		},
		ApplicationInformation: userRepoModel.UpdateApplicationInformationReq{
			BundleID: &req.Application.BundleID,
			Version:  &req.Application.Version,
			Build:    &req.Application.Build,
		},
		NotificationToken: nil,
	}); err != nil {
		return newTokens, err
	}

	// Возвращаем пару токенов клиенту
	return newTokens, nil
}
