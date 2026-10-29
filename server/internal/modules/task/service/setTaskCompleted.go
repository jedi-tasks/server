package service

import (
	"context"

	"github.com/google/uuid"

	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/task/model"
)

// SetTaskCompleted выполняет задачу или возвращает ее в работу. При выполнении снимаются блокеры,
// ссылающиеся на эту задачу. Действие старее последнего изменения задачи игнорируется
func (s *TaskService) SetTaskCompleted(ctx context.Context, req model.SetTaskCompletedReq) error {
	ctx, span := tracer.Start(ctx, "SetTaskCompleted")
	defer span.End()

	// Проверяем доступ пользователя к задаче
	if err := s.CheckAccess(ctx, req.Necessary.UserID, []uuid.UUID{req.ID}); err != nil {
		return err
	}

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {

		// Получаем задачу для слепка "до" в аудит-логе
		before, err := slices.FirstWithError(s.taskRepository.GetTasks(ctxTx, model.GetTasksReq{ //nolint:exhaustruct
			Necessary: req.Necessary,
			IDs:       []uuid.UUID{req.ID},
		}))
		if err != nil {
			return 0, err
		}

		// Выставляем или снимаем время выполнения
		completedAt := &req.UpdatedAt
		if !req.Completed {
			completedAt = nil
		}
		updated, err := s.taskRepository.SetTaskCompleted(ctxTx, req.ID, completedAt, req.UpdatedAt)
		if err != nil || !updated {
			return 0, err
		}

		// Получаем актуальную задачу для слепка "после" в аудит-логе
		after, err := slices.FirstWithError(s.taskRepository.GetTasks(ctxTx, model.GetTasksReq{ //nolint:exhaustruct
			Necessary: req.Necessary,
			IDs:       []uuid.UUID{req.ID},
		}))
		if err != nil {
			return 0, err
		}

		// Фиксируем изменение в аудит-логе
		lastAuditLogID, err := s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.Task,
			Method:   auditLogMethod.Update,
			EntityID: req.ID.String(),
			Before:   before,
			After:    after,
			UserID:   req.Necessary.UserID,
			DeviceID: req.Necessary.DeviceID,
		})
		if err != nil || !req.Completed || s.blockerResolver == nil {
			return lastAuditLogID, err
		}

		// Снимаем блокеры, ссылающиеся на выполненную задачу
		auditLogID, err := s.blockerResolver.ResolveByBlockingTask(ctxTx, req.Necessary, req.ID, req.UpdatedAt)
		if err != nil {
			return 0, err
		}

		return max(lastAuditLogID, auditLogID), nil
	})
}
