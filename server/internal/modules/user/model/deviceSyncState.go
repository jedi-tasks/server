package model

// DeviceSyncState — состояние синхронизации устройства для проверки перед правкой
type DeviceSyncState struct {
	LastSyncedAuditLogID uint32 // Идентификатор последней подтвержденной записи аудит-лога
	IsStateless          bool   // Устройство без локальной копии данных: проверка синхронизации не нужна
}
