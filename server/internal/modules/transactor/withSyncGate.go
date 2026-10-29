package transactor

import (
	"context"

	"github.com/google/uuid"

	"server/internal/utils/errors"
)

// DeviceSyncGate - интерфейс доступа к чекпоинту устройства, необходимый для WithSyncGate
type DeviceSyncGate interface {
	GetDeviceLastSyncedAuditLogIDForUpdate(ctx context.Context, userID uuid.UUID, deviceID string) (uint32, error)
	BumpDeviceCheckpoint(ctx context.Context, userID uuid.UUID, deviceID string, auditLogID uint32) error
}

// AuditLogChangeChecker - интерфейс легковесной проверки наличия изменений после чекпоинта
type AuditLogChangeChecker interface {
	HasAuditLogsSince(ctx context.Context, userID uuid.UUID, sinceID uint32) (bool, error)
}

// WithSyncGate оборачивает мутацию проверкой актуальности синхронизации устройства: блокирует строку
// устройства, отклоняет запрос при наличии несинхронизированных изменений, выполняет переданную
// мутацию и атомарно поднимает чекпоинт устройства до идентификатора записи аудит-лога, созданной
// этой мутацией
func (r *Transactor) WithSyncGate(
	ctx context.Context,
	userID uuid.UUID,
	deviceID string,
	deviceSyncGate DeviceSyncGate,
	auditLogChangeChecker AuditLogChangeChecker,
	mutate func(ctxTx context.Context) (auditLogID uint32, err error),
) error {
	return r.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Блокируем строку устройства до конца транзакции и читаем его чекпоинт
		lastSyncedAuditLogID, err := deviceSyncGate.GetDeviceLastSyncedAuditLogIDForUpdate(ctxTx, userID, deviceID)
		if err != nil {
			return err
		}

		// Отклоняем мутацию, если устройство не синхронизировано с последними изменениями
		hasNewer, err := auditLogChangeChecker.HasAuditLogsSince(ctxTx, userID, lastSyncedAuditLogID)
		if err != nil {
			return err
		}
		if hasNewer {
			return errors.NeedToSync.New("Устройство не синхронизировано с последними изменениями").
				WithContextParams(ctxTx)
		}

		// Выполняем саму мутацию
		auditLogID, err := mutate(ctxTx)
		if err != nil {
			return err
		}

		// Поднимаем чекпоинт устройства до записи собственной мутации
		return deviceSyncGate.BumpDeviceCheckpoint(ctxTx, userID, deviceID, auditLogID)
	})
}
