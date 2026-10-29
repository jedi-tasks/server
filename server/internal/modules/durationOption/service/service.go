package service

import (
	"context"
	userModel "server/internal/modules/user/model"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	durationOptionModel "server/internal/modules/durationOption/model"
	durationOptionRepository "server/internal/modules/durationOption/repository"
	durationOptionRepoModel "server/internal/modules/durationOption/repository/model"
	"server/internal/modules/transactor"
	userService "server/internal/modules/user/service"
)

var tracer = otel.Tracer("/server/internal/modules/durationOption/service")

type DurationOptionService struct {
	durationOptionRepository DurationOptionRepository
	transactor               Transactor
	auditLogService          AuditLogService
	userService              UserService
}

func NewDurationOptionService(
	durationOptionRepository DurationOptionRepository,
	transactor Transactor,
	auditLogService AuditLogService,
	userService UserService,
) *DurationOptionService {
	return &DurationOptionService{
		durationOptionRepository: durationOptionRepository,
		transactor:               transactor,
		auditLogService:          auditLogService,
		userService:              userService,
	}
}

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
	WithSyncGate(ctx context.Context, userID uuid.UUID, deviceID string, deviceSyncGate transactor.DeviceSyncGate, auditLogChangeChecker transactor.AuditLogChangeChecker, mutate func(ctxTx context.Context) (auditLogID uint32, err error)) error
}

var _ DurationOptionRepository = new(durationOptionRepository.DurationOptionRepository)

type DurationOptionRepository interface {
	CreateDurationOption(context.Context, durationOptionRepoModel.CreateDurationOptionReq) (bool, error)
	UpdateDurationOption(context.Context, durationOptionModel.UpdateDurationOptionReq) (bool, error)
	DeleteDurationOption(ctx context.Context, id uuid.UUID, updatedAt time.Time) (bool, error)
	GetDurationOptions(context.Context, durationOptionModel.GetDurationOptionsReq) ([]durationOptionModel.DurationOption, error)

	CheckAccess(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
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
	GetDeviceSyncStateForUpdate(ctx context.Context, userID uuid.UUID, deviceID string) (userModel.DeviceSyncState, error)
	BumpDeviceCheckpoint(ctx context.Context, userID uuid.UUID, deviceID string, auditLogID uint32) error
}
