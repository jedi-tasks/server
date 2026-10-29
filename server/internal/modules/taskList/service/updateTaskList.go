package service

import (
	"context"

	"github.com/google/uuid"

	"pkg/slices"
	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/taskList/model"
)

// UpdateTaskList редактирует запись. Изменение старее сохраненного игнорируется
func (s *TaskListService) UpdateTaskList(ctx context.Context, req model.UpdateTaskListReq) error {
	ctx, span := tracer.Start(ctx, "UpdateTaskList")
	defer span.End()

	// Проверяем доступ пользователя к записи
	if err := s.CheckAccess(ctx, req.Necessary.UserID, []uuid.UUID{req.ID}); err != nil {
		return err
	}

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {

		// Получаем запись для слепка "до" в аудит-логе
		before, err := slices.FirstWithError(s.taskListRepository.GetTaskLists(ctxTx, model.GetTaskListsReq{ //nolint:exhaustruct
			Necessary: req.Necessary,
			IDs:       []uuid.UUID{req.ID},
		}))
		if err != nil {
			return 0, err
		}

		// Изменяем запись
		updated, err := s.taskListRepository.UpdateTaskList(ctxTx, req)
		if err != nil || !updated {
			return 0, err
		}

		// Получаем актуальную запись для слепка "после" в аудит-логе
		after, err := slices.FirstWithError(s.taskListRepository.GetTaskLists(ctxTx, model.GetTaskListsReq{ //nolint:exhaustruct
			Necessary: req.Necessary,
			IDs:       []uuid.UUID{req.ID},
		}))
		if err != nil {
			return 0, err
		}

		// Фиксируем изменение в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.TaskList,
			Method:   auditLogMethod.Update,
			EntityID: req.ID.String(),
			Before:   before,
			After:    after,
			UserID:   req.Necessary.UserID,
			DeviceID: req.Necessary.DeviceID,
		})
	})
}
