package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/blocker/model"
	"server/internal/utils/necessary"
)

// ResolveBlocker снимает блокер вручную
func (s *BlockerService) ResolveBlocker(ctx context.Context, req model.ResolveBlockerReq) error {
	ctx, span := tracer.Start(ctx, "ResolveBlocker")
	defer span.End()

	// Проверяем доступ пользователя к блокеру
	if err := s.CheckAccess(ctx, req.Necessary.UserID, []uuid.UUID{req.ID}); err != nil {
		return err
	}

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {
		return s.resolveBlocker(ctxTx, req.Necessary, req.ID, req.UpdatedAt)
	})
}

// resolveBlocker снимает блокер внутри открытой транзакции и фиксирует это в аудит-логе
func (s *BlockerService) resolveBlocker(ctxTx context.Context, necessary necessary.NecessaryUserInformation, id uuid.UUID, resolvedAt time.Time) (uint32, error) {

	// Получаем блокер для слепка "до" в аудит-логе
	before, err := slices.FirstWithError(s.blockerRepository.GetBlockers(ctxTx, model.GetBlockersReq{ //nolint:exhaustruct
		Necessary: necessary,
		IDs:       []uuid.UUID{id},
	}))
	if err != nil {
		return 0, err
	}

	// Снимаем блокер
	resolved, err := s.blockerRepository.ResolveBlocker(ctxTx, id, resolvedAt)
	if err != nil || !resolved {
		return 0, err
	}

	// Получаем блокер для слепка "после" в аудит-логе
	after, err := slices.FirstWithError(s.blockerRepository.GetBlockers(ctxTx, model.GetBlockersReq{ //nolint:exhaustruct
		Necessary: necessary,
		IDs:       []uuid.UUID{id},
	}))
	if err != nil {
		return 0, err
	}

	// Фиксируем изменение в аудит-логе
	return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
		Entity:   auditLogEntity.Blocker,
		Method:   auditLogMethod.Update,
		EntityID: id.String(),
		Before:   before,
		After:    after,
		UserID:   necessary.UserID,
		DeviceID: necessary.DeviceID,
	})
}
