package model

import (
	"time"

	"github.com/google/uuid"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
)

// AuditLog - запись аудит-лога об изменяющем действии пользователя
type AuditLog struct {
	ID             uint32                        `db:"id"`              // Числовой идентификатор записи
	Entity         auditLogEntity.AuditLogEntity `db:"entity"`          // Название сущности (transaction, account и т.д.)
	Method         auditLogMethod.AuditLogMethod `db:"method"`          // Метод изменения (create, update, delete)
	EntityID       string                        `db:"entity_id"`       // Идентификатор сущности (без FK)
	SnapshotBefore []byte                        `db:"snapshot_before"` // Слепок сущности до изменения в формате JSON
	SnapshotAfter  []byte                        `db:"snapshot_after"`  // Слепок сущности после изменения в формате JSON
	UserID         uuid.UUID                     `db:"user_id"`         // Идентификатор пользователя, совершившего действие
	DeviceID       string                        `db:"device_id"`       // Идентификатор устройства, с которого совершено действие
	DatetimeCreate time.Time                     `db:"datetime_create"` // Дата и время изменения
}
