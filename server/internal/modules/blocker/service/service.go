package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	blockerModel "server/internal/modules/blocker/model"
	blockerRepository "server/internal/modules/blocker/repository"
	blockerRepoModel "server/internal/modules/blocker/repository/model"
	taskService "server/internal/modules/task/service"
	"server/internal/modules/transactor"
	userService "server/internal/modules/user/service"
)

var tracer = otel.Tracer("/server/internal/modules/blocker/service")

type BlockerService struct {
	blockerRepository BlockerRepository
	transactor        Transactor
	auditLogService   AuditLogService
	userService       UserService
	taskService       TaskService
}

func NewBlockerService(
	blockerRepository BlockerRepository,
	transactor Transactor,
	auditLogService AuditLogService,
	userService UserService,
	taskService TaskService,
) *BlockerService {
	return &BlockerService{
		blockerRepository: blockerRepository,
		transactor:        transactor,
		auditLogService:   auditLogService,
		userService:       userService,
		taskService:       taskService,
	}
}

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
	WithSyncGate(ctx context.Context, userID uuid.UUID, deviceID string, deviceSyncGate transactor.DeviceSyncGate, auditLogChangeChecker transactor.AuditLogChangeChecker, mutate func(ctxTx context.Context) (auditLogID uint32, err error)) error
}

var _ BlockerRepository = new(blockerRepository.BlockerRepository)

type BlockerRepository interface {
	CreateBlocker(context.Context, blockerRepoModel.CreateBlockerReq) (bool, error)
	DeleteBlocker(ctx context.Context, id uuid.UUID, updatedAt time.Time) (bool, error)
	GetBlockers(context.Context, blockerModel.GetBlockersReq) ([]blockerModel.Blocker, error)

	CheckAccess(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error

	ResolveBlocker(ctx context.Context, id uuid.UUID, resolvedAt time.Time) (bool, error)
}

var _ TaskService = new(taskService.TaskService)

type TaskService interface {
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
	GetDeviceLastSyncedAuditLogIDForUpdate(ctx context.Context, userID uuid.UUID, deviceID string) (uint32, error)
	BumpDeviceCheckpoint(ctx context.Context, userID uuid.UUID, deviceID string, auditLogID uint32) error
}
