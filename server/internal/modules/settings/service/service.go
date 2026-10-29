package service

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	settingsModel "server/internal/modules/settings/model"
	settingsRepository "server/internal/modules/settings/repository"
	"server/internal/modules/transactor"
	userService "server/internal/modules/user/service"
)

var tracer = otel.Tracer("/server/internal/modules/settings/service")

// Version — версия и сборка сервера
type Version struct {
	Version string
	Build   string
}

type SettingsService struct {
	settingsRepository SettingsRepository
	transactor         Transactor
	auditLogService    AuditLogService
	userService        UserService
	version            Version
}

func NewSettingsService(
	settingsRepository SettingsRepository,
	transactor Transactor,
	auditLogService AuditLogService,
	userService UserService,
	version Version,
) *SettingsService {
	return &SettingsService{
		settingsRepository: settingsRepository,
		transactor:         transactor,
		auditLogService:    auditLogService,
		userService:        userService,
		version:            version,
	}
}

var _ SettingsRepository = new(settingsRepository.SettingsRepository)

type SettingsRepository interface {
	GetSettings(ctx context.Context, userID uuid.UUID) (settingsModel.Settings, bool, error)
	SaveSettings(context.Context, settingsModel.Settings) (bool, error)
}

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithSyncGate(ctx context.Context, userID uuid.UUID, deviceID string, deviceSyncGate transactor.DeviceSyncGate, auditLogChangeChecker transactor.AuditLogChangeChecker, mutate func(ctxTx context.Context) (auditLogID uint32, err error)) error
}

var _ AuditLogService = new(auditLogService.AuditLogService)

// AuditLogService - интерфейс сервиса аудит-лога
type AuditLogService interface {
	TrackMutation(context.Context, auditLogModel.TrackMutationReq) (uint32, error)
	HasAuditLogsSince(ctx context.Context, userID uuid.UUID, sinceID uint32) (bool, error)
}

var _ UserService = new(userService.UserService)

// UserService - интерфейс сервиса пользователей, используется для отсечения мутаций от
// несинхронизированных устройств
type UserService interface {
	GetDeviceLastSyncedAuditLogIDForUpdate(ctx context.Context, userID uuid.UUID, deviceID string) (uint32, error)
	BumpDeviceCheckpoint(ctx context.Context, userID uuid.UUID, deviceID string, auditLogID uint32) error
}
