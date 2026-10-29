package model

import (
	"github.com/google/uuid"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
)

// CreateAuditLogReq - репозиторная модель запроса на создание записи аудит-лога
type CreateAuditLogReq struct {
	Entity         auditLogEntity.AuditLogEntity // Название сущности (task, taskList и т.д.)
	Method         auditLogMethod.AuditLogMethod // Метод изменения (create, update, delete)
	EntityID       string                        // Идентификатор сущности (без FK)
	SnapshotBefore []byte                        // Слепок сущности до изменения в формате JSON
	SnapshotAfter  []byte                        // Слепок сущности после изменения в формате JSON
	UserID         uuid.UUID                     // Идентификатор пользователя, совершившего действие
	DeviceID       string                        // Идентификатор устройства, с которого совершено действие
}
