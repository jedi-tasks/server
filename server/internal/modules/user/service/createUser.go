package service

import (
	"context"

	"github.com/google/uuid"

	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/user/model"
)

// CreateUser создает нового пользователя
func (s *UserService) CreateUser(ctx context.Context, user model.CreateReq) (id uuid.UUID, err error) {
	ctx, span := tracer.Start(ctx, "CreateUser")
	defer span.End()

	user.ID = uuid.New()

	err = s.generalRepository.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Создаем пользователя
		if err := s.userRepository.CreateUser(ctxTx, user); err != nil {
			return err
		}

		// Получаем созданного пользователя из БД для слепка "после" в аудит-логе (хэш и соль пароля исключены тегом json:"-")
		userAfter, err := slices.FirstWithError(s.userRepository.GetUsers(ctxTx, model.GetUsersReq{ //nolint:exhaustruct
			IDs: []uuid.UUID{user.ID},
		}))
		if err != nil {
			return err
		}

		// Фиксируем создание пользователя в аудит-логе
		_, err = s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.User,
			Method:   auditLogMethod.Create,
			EntityID: user.ID.String(),
			Before:   nil,
			After:    userAfter,
			UserID:   user.ID,
		})
		return err
	})
	if err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}
