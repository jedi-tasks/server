package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/tag/model"
	"server/internal/modules/tag/repository/taskToTagDDL"
	"server/internal/utils/errors"
)

// DeleteTag мягко удаляет запись. Тег, привязанный к задачам, удалить нельзя. Удаление старее последнего изменения игнорируется
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

	// Проверяем, что на тег не ссылаются живые связи с задачами. Строка тега уже заблокирована
	// удалением, поэтому параллельная привязка дождется этой транзакции
	hasLinks, err := s.tagRepository.HasLiveLinks(ctxTx, taskToTagDDL.ColumnTagID, req.ID)
	if err != nil {
		return 0, err
	}
	if hasLinks {
		return 0, errors.BadRequest.New("Tag has live task links").
			WithContextParams(ctxTx).
			WithParams("TagID", req.ID).
			WithCustomHumanText("Тег привязан к задачам — сначала отвяжите его")
	}

	// Фиксируем удаление в аудит-логе
	return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
		Entity:   auditLogEntity.Tag,
		Method:   auditLogMethod.Delete,
		EntityID: req.ID.String(),
		Before:   befores[0],
		After:    nil,
		UserID:   req.Necessary.UserID,
		DeviceID: req.Necessary.DeviceID,
	})
}
