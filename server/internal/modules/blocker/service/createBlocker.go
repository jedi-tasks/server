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

		// Проверяем, что задачи блокера доступны и не удалены, и блокируем их от удаления до конца транзакции
		taskIDs := []uuid.UUID{req.TaskID}
		if req.BlockingTaskID != nil {
			taskIDs = append(taskIDs, *req.BlockingTaskID)
		}
		if err = s.taskService.LockLive(ctxTx, req.Necessary.UserID, taskIDs); err != nil {
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
