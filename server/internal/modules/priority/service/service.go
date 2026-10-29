package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	priorityModel "server/internal/modules/priority/model"
	priorityRepository "server/internal/modules/priority/repository"
	priorityRepoModel "server/internal/modules/priority/repository/model"
	"server/internal/modules/transactor"
	userService "server/internal/modules/user/service"
)

var tracer = otel.Tracer("/server/internal/modules/priority/service")

// PriorityService - сервис справочника приоритетов
type PriorityService struct {
	priorityRepository PriorityRepository
	transactor         Transactor
	auditLogService    AuditLogService
	userService        UserService

	liveTasksChecker LiveTasksChecker
}

func NewPriorityService(
	priorityRepository PriorityRepository,
	transactor Transactor,
	auditLogService AuditLogService,
	userService UserService,
) *PriorityService {
	return &PriorityService{
		priorityRepository: priorityRepository,
		transactor:         transactor,
		auditLogService:    auditLogService,
		userService:        userService,
	}
}

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
	WithSyncGate(ctx context.Context, userID uuid.UUID, deviceID string, deviceSyncGate transactor.DeviceSyncGate, auditLogChangeChecker transactor.AuditLogChangeChecker, mutate func(ctxTx context.Context) (auditLogID uint32, err error)) error
}

var _ PriorityRepository = new(priorityRepository.PriorityRepository)

type PriorityRepository interface {
	CreatePriority(context.Context, priorityRepoModel.CreatePriorityReq) (bool, error)
	UpdatePriority(context.Context, priorityModel.UpdatePriorityReq) (bool, error)
	DeletePriority(ctx context.Context, id uuid.UUID, updatedAt time.Time) (bool, error)
	GetPriorities(context.Context, priorityModel.GetPrioritiesReq) ([]priorityModel.Priority, error)

	CheckAccess(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
	LockLive(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
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

// LiveTasksChecker - интерфейс сервиса задач, отвечающего, есть ли неудаленные задачи с приоритетом.
// Сервис задач сам зависит от этого сервиса, поэтому подключается после создания через SetLiveTasksChecker
type LiveTasksChecker interface {
	HasLiveTasksWithPriority(ctx context.Context, priorityID uuid.UUID) (bool, error)
}

// SetLiveTasksChecker подключает проверку задач с приоритетом перед его удалением
func (s *PriorityService) SetLiveTasksChecker(liveTasksChecker LiveTasksChecker) {
	s.liveTasksChecker = liveTasksChecker
}
