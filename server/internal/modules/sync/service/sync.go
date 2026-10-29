package service

import (
	"context"

	"github.com/google/uuid"

	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/sync/model"
)

// entityKey - уникальный идентификатор объекта в рамках сущности
type entityKey struct {
	entity   string
	entityID string
}

// Sync возвращает изменения, доступные пользователю по группам счетов, произошедшие после чекпоинта,
// переданного клиентом (sinceID). Чекпоинт хранится на стороне клиента - Sync не сохраняет его сам,
// а лишь выдаёт pendingCheckpoint/pendingSyncToken, которые клиент должен подтвердить вызовом ConfirmSync
// после того, как корректно применит полученные изменения у себя
func (s *SyncService) Sync(ctx context.Context, req model.SyncReq) (res model.SyncRes, err error) {
	ctx, span := tracer.Start(ctx, "Sync")
	defer span.End()

	err = s.transactor.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Получаем изменения, произошедшие после чекпоинта клиента, в доступных пользователю группах счетов
		auditLogs, err := s.auditLogService.GetAuditLogsSince(ctxTx, req.Necessary.UserID, req.SinceID)
		if err != nil {
			return err
		}

		res.PendingCheckpoint = req.SinceID

		// Изменений нет - клиенту нечего подтверждать, чекпоинт остаётся прежним
		if len(auditLogs) == 0 {
			return nil
		}

		// Для каждого объекта оставляем только последнее событие - записи идут по возрастанию id,
		// поэтому при перезаписи в мапе остается самая свежая запись
		latest := make(map[entityKey]auditLogModel.AuditLog, len(auditLogs))
		for _, auditLog := range auditLogs {
			latest[entityKey{entity: string(auditLog.Entity), entityID: auditLog.EntityID}] = auditLog

			if auditLog.ID > res.PendingCheckpoint {
				res.PendingCheckpoint = auditLog.ID
			}
		}

		// Раскладываем последние события по сущностям, вычитывая актуальные данные из доменных сервисов
		if err := s.hydrate(ctxTx, req.Necessary.UserID, latest, &res); err != nil {
			return err
		}

		res.HasChanges = true

		// Выдаём опорный токен, который клиент подтвердит вызовом ConfirmSync после применения изменений
		pendingSyncToken := uuid.New()
		res.PendingSyncToken = &pendingSyncToken

		return s.userService.SetDevicePendingSync(ctxTx, req.Necessary.UserID, req.Necessary.DeviceID, res.PendingCheckpoint, pendingSyncToken)
	})
	if err != nil {
		return model.SyncRes{}, err //nolint:exhaustruct
	}

	return res, nil
}
