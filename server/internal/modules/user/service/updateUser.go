package service

import (
	"context"

	"github.com/google/uuid"
	"pkg/passwordManager"
	"pkg/slices"
	"server/internal/utils/errors"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/user/model"
)

// UpdateUser обновляет настройки пользователя
func (s *UserService) UpdateUser(ctx context.Context, req model.UpdateUserReq) error {
	ctx, span := tracer.Start(ctx, "UpdateUser")
	defer span.End()

	// Получаем пользователя для слепка "до" в аудит-логе
	users, err := s.userRepository.GetUsers(ctx, model.GetUsersReq{ //nolint:exhaustruct
		IDs: []uuid.UUID{req.Necessary.UserID},
	})
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return errors.NotFound.New("User not found").WithContextParams(ctx)
	}
	userBefore := users[0]

	return s.generalRepository.WithinTransaction(ctx, func(ctx context.Context) error {

		repoReq := req.ConvertToRepoModel()

		// Если обновляется пароль
		if req.Password != nil {

			if req.OldPassword == nil {
				return errors.BadRequest.New("OldPassword must be filled").WithContextParams(ctx)
			}

			// Сравниваем пришедший пароль и хэш пароля из базы данных
			if err := passwordManager.CompareHashAndPassword(userBefore.PasswordHash, []byte(*req.OldPassword), userBefore.PasswordSalt, s.generalSalt); err != nil {
				return err
			}

			// Генерируем соль для юзера
			userSalt, err := passwordManager.GenerateRandomSalt()
			if err != nil {
				return err
			}

			// Получаем хэш и соль нового пароля
			passwordHash, err := passwordManager.CreateNewPassword([]byte(*req.Password), s.generalSalt, userSalt)
			if err != nil {
				return err
			}

			repoReq.PasswordHash = &passwordHash
			repoReq.PasswordSalt = &userSalt
		}

		if err := s.userRepository.UpdateUser(ctx, repoReq); err != nil {
			return err
		}

		// Получаем актуального пользователя из БД для слепка "после" в аудит-логе (хэш и соль пароля исключены тегом json:"-")
		userAfter, err := slices.FirstWithError(s.userRepository.GetUsers(ctx, model.GetUsersReq{ //nolint:exhaustruct
			IDs: []uuid.UUID{req.Necessary.UserID},
		}))
		if err != nil {
			return err
		}

		// Фиксируем изменение пользователя в аудит-логе
		_, err = s.auditLogService.TrackMutation(ctx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.User,
			Method:   auditLogMethod.Update,
			EntityID: req.Necessary.UserID.String(),
			Before:   userBefore,
			After:    userAfter,
			UserID:   req.Necessary.UserID,
			DeviceID: req.Necessary.DeviceID,
		})
		return err
	})
}
