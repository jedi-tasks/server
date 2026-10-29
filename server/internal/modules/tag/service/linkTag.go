package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/tag/model"
	"server/internal/utils/necessary"
)

// AddTagToTask привязывает тег к задаче
func (s *TagService) AddTagToTask(ctx context.Context, req model.LinkTagReq) error {
	ctx, span := tracer.Start(ctx, "AddTagToTask")
	defer span.End()

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {

		// Привязываем тег. Повтор из офлайн-очереди ничего не меняет и не пишет в аудит-лог
		added, err := s.tagRepository.AddTagToTask(ctxTx, req.Necessary.UserID, req.TaskID, req.TagID, req.UpdatedAt)
		if err != nil || !added {
			return 0, err
		}

		// Проверяем, что задача и тег доступны и не удалены, и блокируем их от удаления до конца транзакции
		if err = s.taskService.LockLive(ctxTx, req.Necessary.UserID, []uuid.UUID{req.TaskID}); err != nil {
			return 0, err
		}
		if err = s.LockLive(ctxTx, req.Necessary.UserID, []uuid.UUID{req.TagID}); err != nil {
			return 0, err
		}

		// Получаем связь для слепка "после" в аудит-логе
		after, err := slices.FirstWithError(s.tagRepository.GetTasksToTags(ctxTx, model.GetTasksToTagsReq{
			Necessary: req.Necessary,
			TaskIDs:   []uuid.UUID{req.TaskID},
			TagIDs:    []uuid.UUID{req.TagID},
		}))
		if err != nil {
			return 0, err
		}

		// Фиксируем привязку в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.TaskToTag,
			Method:   auditLogMethod.Create,
			EntityID: after.EntityID(),
			Before:   nil,
			After:    after,
			UserID:   req.Necessary.UserID,
			DeviceID: req.Necessary.DeviceID,
		})
	})
}

// RemoveTagFromTask отвязывает тег от задачи
func (s *TagService) RemoveTagFromTask(ctx context.Context, req model.LinkTagReq) error {
	ctx, span := tracer.Start(ctx, "RemoveTagFromTask")
	defer span.End()

	// Проверяем доступ пользователя к задаче
	if err := s.taskService.CheckAccess(ctx, req.Necessary.UserID, []uuid.UUID{req.TaskID}); err != nil {
		return err
	}

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {
		return s.removeTagFromTask(ctxTx, req.Necessary, req.TaskID, req.TagID, req.UpdatedAt)
	})
}

// removeTagFromTask отвязывает тег внутри открытой транзакции и фиксирует это в аудит-логе
func (s *TagService) removeTagFromTask(ctxTx context.Context, necessary necessary.NecessaryUserInformation, taskID, tagID uuid.UUID, updatedAt time.Time) (uint32, error) {

	// Получаем связь для слепка "до" в аудит-логе. Уже удаленную повторно не удаляем
	befores, err := s.tagRepository.GetTasksToTags(ctxTx, model.GetTasksToTagsReq{
		Necessary: necessary,
		TaskIDs:   []uuid.UUID{taskID},
		TagIDs:    []uuid.UUID{tagID},
	})
	if err != nil || len(befores) == 0 {
		return 0, err
	}

	// Отвязываем тег
	removed, err := s.tagRepository.RemoveTagFromTask(ctxTx, taskID, tagID, updatedAt)
	if err != nil || !removed {
		return 0, err
	}

	// Фиксируем отвязку в аудит-логе
	return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
		Entity:   auditLogEntity.TaskToTag,
		Method:   auditLogMethod.Delete,
		EntityID: befores[0].EntityID(),
		Before:   befores[0],
		After:    nil,
		UserID:   necessary.UserID,
		DeviceID: necessary.DeviceID,
	})
}
