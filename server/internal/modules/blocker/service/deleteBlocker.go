package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/blocker/model"
)

// DeleteBlocker мягко удаляет запись. Удаление старее последнего изменения игнорируется
func (s *BlockerService) DeleteBlocker(ctx context.Context, req model.DeleteBlockerReq) error {
	ctx, span := tracer.Start(ctx, "DeleteBlocker")
	defer span.End()

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {
		return s.deleteBlocker(ctxTx, req)
	})
}

// deleteBlocker удаляет запись внутри уже открытой транзакции и возвращает идентификатор записи
// аудит-лога (0, если удалять нечего)
func (s *BlockerService) deleteBlocker(ctxTx context.Context, req model.DeleteBlockerReq) (uint32, error) {

	// Получаем запись для слепка "до" в аудит-логе. Уже удаленную повторно не удаляем
	befores, err := s.blockerRepository.GetBlockers(ctxTx, model.GetBlockersReq{ //nolint:exhaustruct
		Necessary: req.Necessary,
		IDs:       []uuid.UUID{req.ID},
	})
	if err != nil || len(befores) == 0 {
		return 0, err
	}

	// Удаляем запись
	deleted, err := s.blockerRepository.DeleteBlocker(ctxTx, req.ID, req.UpdatedAt)
	if err != nil || !deleted {
		return 0, err
	}

	// Фиксируем удаление в аудит-логе
	return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
		Entity:   auditLogEntity.Blocker,
		Method:   auditLogMethod.Delete,
		EntityID: req.ID.String(),
		Before:   befores[0],
		After:    nil,
		UserID:   req.Necessary.UserID,
		DeviceID: req.Necessary.DeviceID,
	})
}
