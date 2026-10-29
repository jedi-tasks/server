package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/taskList/model"
	"server/internal/utils/errors"
)

// DeleteTaskList мягко удаляет список. Inbox удалить нельзя, список с неудаленными задачами или
// разделами — тоже: сервер ничего не удаляет каскадом. Удаление старее последнего изменения игнорируется
func (s *TaskListService) DeleteTaskList(ctx context.Context, req model.DeleteTaskListReq) error {
	ctx, span := tracer.Start(ctx, "DeleteTaskList")
	defer span.End()

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {

		// Получаем список для слепка "до" в аудит-логе. Уже удаленный повторно не удаляем
		befores, err := s.taskListRepository.GetTaskLists(ctxTx, model.GetTaskListsReq{ //nolint:exhaustruct
			Necessary: req.Necessary,
			IDs:       []uuid.UUID{req.ID},
		})
		if err != nil || len(befores) == 0 {
			return 0, err
		}

		// Проверяем, что это не Inbox
		if befores[0].IsInbox {
			return 0, errors.BadRequest.New("Inbox can't be deleted").
				WithContextParams(ctxTx).
				WithCustomHumanText("Inbox удалить нельзя")
		}

		// Удаляем список
		deleted, err := s.taskListRepository.DeleteTaskList(ctxTx, req.ID, req.UpdatedAt)
		if err != nil || !deleted {
			return 0, err
		}

		// Проверяем, что на список не ссылаются живые задачи и разделы. Проверка идет после удаления: строка
		// списка уже заблокирована, поэтому параллельное создание задачи или раздела дождется этой транзакции
		if err = s.checkNoLiveContents(ctxTx, req.ID); err != nil {
			return 0, err
		}

		// Фиксируем удаление в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.TaskList,
			Method:   auditLogMethod.Delete,
			EntityID: req.ID.String(),
			Before:   befores[0],
			After:    nil,
			UserID:   req.Necessary.UserID,
			DeviceID: req.Necessary.DeviceID,
		})
	})
}
