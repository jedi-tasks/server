package transactor

import (
	"context"
	userModel "server/internal/modules/user/model"

	"github.com/google/uuid"

	"server/internal/utils/errors"
)

// DeviceSyncGate - интерфейс доступа к чекпоинту устройства, необходимый для WithSyncGate
type DeviceSyncGate interface {
	GetDeviceSyncStateForUpdate(ctx context.Context, userID uuid.UUID, deviceID string) (userModel.DeviceSyncState, error)
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

		// Блокируем строку устройства до конца транзакции и читаем его состояние синхронизации
		syncState, err := deviceSyncGate.GetDeviceSyncStateForUpdate(ctxTx, userID, deviceID)
		if err != nil {
			return err
		}

		// Устройство с локальной копией данных должно увидеть чужие изменения до своей правки,
		// иначе оно перезатрет их, не зная о них. Устройство без копии (MCP, веб) читает сервер
		// перед каждой записью, для него проверка не нужна
		if !syncState.IsStateless {
			hasNewer, err := auditLogChangeChecker.HasAuditLogsSince(ctxTx, userID, syncState.LastSyncedAuditLogID)
			if err != nil {
				return err
			}
			if hasNewer {
				return errors.NeedToSync.New("Устройство не синхронизировано с последними изменениями").
					WithContextParams(ctxTx)
			}
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
