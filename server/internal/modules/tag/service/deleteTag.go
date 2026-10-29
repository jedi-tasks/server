package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/tag/model"
)

// DeleteTag мягко удаляет запись. Удаление старее последнего изменения игнорируется
func (s *TagService) DeleteTag(ctx context.Context, req model.DeleteTagReq) error {
	ctx, span := tracer.Start(ctx, "DeleteTag")
	defer span.End()

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {
		return s.deleteTag(ctxTx, req)
	})
}

// deleteTag удаляет запись внутри уже открытой транзакции и возвращает идентификатор записи
// аудит-лога (0, если удалять нечего)
func (s *TagService) deleteTag(ctxTx context.Context, req model.DeleteTagReq) (uint32, error) {

	// Получаем запись для слепка "до" в аудит-логе. Уже удаленную повторно не удаляем
	befores, err := s.tagRepository.GetTags(ctxTx, model.GetTagsReq{ //nolint:exhaustruct
		Necessary: req.Necessary,
		IDs:       []uuid.UUID{req.ID},
	})
	if err != nil || len(befores) == 0 {
		return 0, err
	}

	// Удаляем запись
	deleted, err := s.tagRepository.DeleteTag(ctxTx, req.ID, req.UpdatedAt)
	if err != nil || !deleted {
		return 0, err
	}

	// Фиксируем удаление в аудит-логе
	lastAuditLogID, err := s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
		Entity:   auditLogEntity.Tag,
		Method:   auditLogMethod.Delete,
		EntityID: req.ID.String(),
		Before:   befores[0],
		After:    nil,
		UserID:   req.Necessary.UserID,
		DeviceID: req.Necessary.DeviceID,
	})
	if err != nil {
		return 0, err
	}

	// Отвязываем тег от всех задач
	links, err := s.tagRepository.GetTasksToTags(ctxTx, model.GetTasksToTagsReq{
		Necessary: req.Necessary,
		TaskIDs:   nil,
		TagIDs:    []uuid.UUID{req.ID},
	})
	if err != nil {
		return 0, err
	}
	for _, link := range links {
		auditLogID, err := s.removeTagFromTask(ctxTx, req.Necessary, link.TaskID, link.TagID, req.UpdatedAt)
		if err != nil {
			return 0, err
		}
		lastAuditLogID = max(lastAuditLogID, auditLogID)
	}

	return lastAuditLogID, nil
}
