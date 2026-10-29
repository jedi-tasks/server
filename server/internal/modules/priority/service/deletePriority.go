package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/priority/model"
)

// DeletePriority мягко удаляет запись. Удаление старее последнего изменения игнорируется
func (s *PriorityService) DeletePriority(ctx context.Context, req model.DeletePriorityReq) error {
	ctx, span := tracer.Start(ctx, "DeletePriority")
	defer span.End()

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {
		return s.deletePriority(ctxTx, req)
	})
}

// deletePriority удаляет запись внутри уже открытой транзакции и возвращает идентификатор записи
// аудит-лога (0, если удалять нечего)
func (s *PriorityService) deletePriority(ctxTx context.Context, req model.DeletePriorityReq) (uint32, error) {

	// Получаем запись для слепка "до" в аудит-логе. Уже удаленную повторно не удаляем
	befores, err := s.priorityRepository.GetPriorities(ctxTx, model.GetPrioritiesReq{ //nolint:exhaustruct
		Necessary: req.Necessary,
		IDs:       []uuid.UUID{req.ID},
	})
	if err != nil || len(befores) == 0 {
		return 0, err
	}

	// Удаляем запись
	deleted, err := s.priorityRepository.DeletePriority(ctxTx, req.ID, req.UpdatedAt)
	if err != nil || !deleted {
		return 0, err
	}

	// Фиксируем удаление в аудит-логе
	lastAuditLogID, err := s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
		Entity:   auditLogEntity.Priority,
		Method:   auditLogMethod.Delete,
		EntityID: req.ID.String(),
		Before:   befores[0],
		After:    nil,
		UserID:   req.Necessary.UserID,
		DeviceID: req.Necessary.DeviceID,
	})
	if err != nil || s.priorityDetacher == nil {
		return lastAuditLogID, err
	}

	// Задачи с этим приоритетом становятся «не разобрано»
	auditLogID, err := s.priorityDetacher.DetachPriority(ctxTx, req.Necessary, req.ID, req.UpdatedAt)
	if err != nil {
		return 0, err
	}

	return max(lastAuditLogID, auditLogID), nil
}
