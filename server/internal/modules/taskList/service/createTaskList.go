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

// CreateTaskList создает запись
func (s *TaskListService) CreateTaskList(ctx context.Context, req model.CreateTaskListReq) error {
	ctx, span := tracer.Start(ctx, "CreateTaskList")
	defer span.End()

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {

		// Создаем запись. Повтор из офлайн-очереди ничего не создает и не пишет в аудит-лог
		created, err := s.taskListRepository.CreateTaskList(ctxTx, req.ConvertToRepoReq())
		if err != nil || !created {
			return 0, err
		}

		// Получаем созданную запись для слепка "после" в аудит-логе
		after, err := slices.FirstWithError(s.taskListRepository.GetTaskLists(ctxTx, model.GetTaskListsReq{ //nolint:exhaustruct
			Necessary: req.Necessary,
			IDs:       []uuid.UUID{req.ID},
		}))
		if err != nil {
			return 0, err
		}

		// Фиксируем создание в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.TaskList,
			Method:   auditLogMethod.Create,
			EntityID: req.ID.String(),
			Before:   nil,
			After:    after,
			UserID:   req.Necessary.UserID,
			DeviceID: req.Necessary.DeviceID,
		})
	})
}
