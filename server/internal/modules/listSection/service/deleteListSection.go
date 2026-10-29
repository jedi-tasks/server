package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/listSection/model"
)

// DeleteListSection мягко удаляет запись. Удаление старее последнего изменения игнорируется
func (s *ListSectionService) DeleteListSection(ctx context.Context, req model.DeleteListSectionReq) error {
	ctx, span := tracer.Start(ctx, "DeleteListSection")
	defer span.End()

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {
		return s.deleteListSection(ctxTx, req)
	})
}

// deleteListSection удаляет запись внутри уже открытой транзакции и возвращает идентификатор записи
// аудит-лога (0, если удалять нечего)
func (s *ListSectionService) deleteListSection(ctxTx context.Context, req model.DeleteListSectionReq) (uint32, error) {

	// Получаем запись для слепка "до" в аудит-логе. Уже удаленную повторно не удаляем
	befores, err := s.listSectionRepository.GetListSections(ctxTx, model.GetListSectionsReq{ //nolint:exhaustruct
		Necessary: req.Necessary,
		IDs:       []uuid.UUID{req.ID},
	})
	if err != nil || len(befores) == 0 {
		return 0, err
	}

	// Удаляем запись
	deleted, err := s.listSectionRepository.DeleteListSection(ctxTx, req.ID, req.UpdatedAt)
	if err != nil || !deleted {
		return 0, err
	}

	// Фиксируем удаление в аудит-логе
	lastAuditLogID, err := s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
		Entity:   auditLogEntity.ListSection,
		Method:   auditLogMethod.Delete,
		EntityID: req.ID.String(),
		Before:   befores[0],
		After:    nil,
		UserID:   req.Necessary.UserID,
		DeviceID: req.Necessary.DeviceID,
	})
	if err != nil || s.sectionDetacher == nil {
		return lastAuditLogID, err
	}

	// Задачи раздела остаются в списке без раздела
	auditLogID, err := s.sectionDetacher.DetachSection(ctxTx, req.Necessary, req.ID, req.UpdatedAt)
	if err != nil {
		return 0, err
	}

	return max(lastAuditLogID, auditLogID), nil
}
