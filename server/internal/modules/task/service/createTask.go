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

// CreateTask создает запись
func (s *TaskService) CreateTask(ctx context.Context, req model.CreateTaskReq) error {
	ctx, span := tracer.Start(ctx, "CreateTask")
	defer span.End()

	// Без списка задача попадает в Inbox
	if req.ListID == uuid.Nil {
		inboxID, err := s.taskListService.GetInboxID(ctx, req.Necessary.UserID)
		if err != nil {
			return err
		}
		req.ListID = inboxID
	}

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {

		// Создаем запись. Повтор из офлайн-очереди ничего не создает и не пишет в аудит-лог
		created, err := s.taskRepository.CreateTask(ctxTx, req.ConvertToRepoReq())
		if err != nil || !created {
			return 0, err
		}

		// Проверяем, что связанные записи доступны и не удалены, и блокируем их от удаления
		if err = s.lockReferences(ctxTx, req.Necessary.UserID, &req.ListID, req.SectionID, req.PriorityID); err != nil {
			return 0, err
		}

		// Получаем созданную запись для слепка "после" в аудит-логе
		after, err := slices.FirstWithError(s.taskRepository.GetTasks(ctxTx, model.GetTasksReq{ //nolint:exhaustruct
			Necessary: req.Necessary,
			IDs:       []uuid.UUID{req.ID},
		}))
		if err != nil {
			return 0, err
		}

		// Фиксируем создание в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.Task,
			Method:   auditLogMethod.Create,
			EntityID: req.ID.String(),
			Before:   nil,
			After:    after,
			UserID:   req.Necessary.UserID,
			DeviceID: req.Necessary.DeviceID,
		})
	})
}
