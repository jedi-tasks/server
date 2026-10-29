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

// UpdateTask редактирует запись. Изменение старее сохраненного игнорируется
func (s *TaskService) UpdateTask(ctx context.Context, req model.UpdateTaskReq) error {
	ctx, span := tracer.Start(ctx, "UpdateTask")
	defer span.End()

	// Проверяем доступ пользователя к записи
	if err := s.CheckAccess(ctx, req.Necessary.UserID, []uuid.UUID{req.ID}); err != nil {
		return err
	}

	// Проверяем доступ пользователя к связанной записи
	if req.ListID != nil && *req.ListID != uuid.Nil {
		if err := s.taskListService.CheckAccess(ctx, req.Necessary.UserID, []uuid.UUID{*req.ListID}); err != nil {
			return err
		}
	}
	// Проверяем доступ пользователя к связанной записи
	if req.SectionID != nil && *req.SectionID != uuid.Nil {
		if err := s.listSectionService.CheckAccess(ctx, req.Necessary.UserID, []uuid.UUID{*req.SectionID}); err != nil {
			return err
		}
	}
	// Проверяем доступ пользователя к связанной записи
	if req.PriorityID != nil && *req.PriorityID != uuid.Nil {
		if err := s.priorityService.CheckAccess(ctx, req.Necessary.UserID, []uuid.UUID{*req.PriorityID}); err != nil {
			return err
		}
	}

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {

		// Получаем запись для слепка "до" в аудит-логе
		before, err := slices.FirstWithError(s.taskRepository.GetTasks(ctxTx, model.GetTasksReq{ //nolint:exhaustruct
			Necessary: req.Necessary,
			IDs:       []uuid.UUID{req.ID},
		}))
		if err != nil {
			return 0, err
		}

		// Изменяем запись
		updated, err := s.taskRepository.UpdateTask(ctxTx, req)
		if err != nil || !updated {
			return 0, err
		}

		// Получаем актуальную запись для слепка "после" в аудит-логе
		after, err := slices.FirstWithError(s.taskRepository.GetTasks(ctxTx, model.GetTasksReq{ //nolint:exhaustruct
			Necessary: req.Necessary,
			IDs:       []uuid.UUID{req.ID},
		}))
		if err != nil {
			return 0, err
		}

		// Фиксируем изменение в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.Task,
			Method:   auditLogMethod.Update,
			EntityID: req.ID.String(),
			Before:   before,
			After:    after,
			UserID:   req.Necessary.UserID,
			DeviceID: req.Necessary.DeviceID,
		})
	})
}
