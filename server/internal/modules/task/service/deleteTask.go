package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/task/model"
)

// DeleteTask мягко удаляет запись. Задачу с привязанными тегами или блокерами удалить нельзя. Удаление старее последнего изменения игнорируется
func (s *TaskService) DeleteTask(ctx context.Context, req model.DeleteTaskReq) error {
	ctx, span := tracer.Start(ctx, "DeleteTask")
	defer span.End()

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {
		return s.deleteTask(ctxTx, req)
	})
}

// deleteTask удаляет запись внутри уже открытой транзакции и возвращает идентификатор записи
// аудит-лога (0, если удалять нечего)
func (s *TaskService) deleteTask(ctxTx context.Context, req model.DeleteTaskReq) (uint32, error) {

	// Получаем запись для слепка "до" в аудит-логе. Уже удаленную повторно не удаляем
	befores, err := s.taskRepository.GetTasks(ctxTx, model.GetTasksReq{ //nolint:exhaustruct
		Necessary: req.Necessary,
		IDs:       []uuid.UUID{req.ID},
	})
	if err != nil || len(befores) == 0 {
		return 0, err
	}

	// Удаляем запись
	deleted, err := s.taskRepository.DeleteTask(ctxTx, req.ID, req.UpdatedAt)
	if err != nil || !deleted {
		return 0, err
	}

	// Проверяем, что на задачу не ссылаются живые связи с тегами и блокеры
	if err = s.checkNoLiveLinks(ctxTx, req.ID); err != nil {
		return 0, err
	}

	// Фиксируем удаление в аудит-логе
	return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
		Entity:   auditLogEntity.Task,
		Method:   auditLogMethod.Delete,
		EntityID: req.ID.String(),
		Before:   befores[0],
		After:    nil,
		UserID:   req.Necessary.UserID,
		DeviceID: req.Necessary.DeviceID,
	})
}
