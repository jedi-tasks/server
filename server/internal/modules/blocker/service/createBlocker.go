package service

import (
	"context"

	"github.com/google/uuid"

	"pkg/slices"
	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/blocker/model"
)

// CreateBlocker создает запись
func (s *BlockerService) CreateBlocker(ctx context.Context, req model.CreateBlockerReq) error {
	ctx, span := tracer.Start(ctx, "CreateBlocker")
	defer span.End()

	// Проверяем доступ пользователя к связанной записи
	if err := s.taskService.CheckAccess(ctx, req.Necessary.UserID, []uuid.UUID{req.TaskID}); err != nil {
		return err
	}
	// Проверяем доступ пользователя к связанной записи
	if req.BlockingTaskID != nil && *req.BlockingTaskID != uuid.Nil {
		if err := s.taskService.CheckAccess(ctx, req.Necessary.UserID, []uuid.UUID{*req.BlockingTaskID}); err != nil {
			return err
		}
	}

	// Проверяем, что блокер не замыкает цепочку
	if req.BlockingTaskID != nil {
		if err := s.checkCycle(ctx, req.Necessary, req.TaskID, *req.BlockingTaskID); err != nil {
			return err
		}
	}

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {

		// Создаем запись. Повтор из офлайн-очереди ничего не создает и не пишет в аудит-лог
		created, err := s.blockerRepository.CreateBlocker(ctxTx, req.ConvertToRepoReq())
		if err != nil || !created {
			return 0, err
		}

		// Получаем созданную запись для слепка "после" в аудит-логе
		after, err := slices.FirstWithError(s.blockerRepository.GetBlockers(ctxTx, model.GetBlockersReq{ //nolint:exhaustruct
			Necessary: req.Necessary,
			IDs:       []uuid.UUID{req.ID},
		}))
		if err != nil {
			return 0, err
		}

		// Фиксируем создание в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.Blocker,
			Method:   auditLogMethod.Create,
			EntityID: req.ID.String(),
			Before:   nil,
			After:    after,
			UserID:   req.Necessary.UserID,
			DeviceID: req.Necessary.DeviceID,
		})
	})
}
