package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	listSectionModel "server/internal/modules/listSection/model"
	listSectionRepository "server/internal/modules/listSection/repository"
	listSectionRepoModel "server/internal/modules/listSection/repository/model"
	taskListService "server/internal/modules/taskList/service"
	"server/internal/modules/transactor"
	userService "server/internal/modules/user/service"
)

var tracer = otel.Tracer("/server/internal/modules/listSection/service")

// ListSectionService - сервис разделов списков
type ListSectionService struct {
	listSectionRepository ListSectionRepository
	transactor            Transactor
	auditLogService       AuditLogService
	userService           UserService
	taskListService       TaskListService

	liveTasksChecker LiveTasksChecker
}

func NewListSectionService(
	listSectionRepository ListSectionRepository,
	transactor Transactor,
	auditLogService AuditLogService,
	userService UserService,
	taskListService TaskListService,
) *ListSectionService {
	return &ListSectionService{
		listSectionRepository: listSectionRepository,
		transactor:            transactor,
		auditLogService:       auditLogService,
		userService:           userService,
		taskListService:       taskListService,
	}
}

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
	WithSyncGate(ctx context.Context, userID uuid.UUID, deviceID string, deviceSyncGate transactor.DeviceSyncGate, auditLogChangeChecker transactor.AuditLogChangeChecker, mutate func(ctxTx context.Context) (auditLogID uint32, err error)) error
}

var _ ListSectionRepository = new(listSectionRepository.ListSectionRepository)

type ListSectionRepository interface {
	CreateListSection(context.Context, listSectionRepoModel.CreateListSectionReq) (bool, error)
	UpdateListSection(context.Context, listSectionModel.UpdateListSectionReq) (bool, error)
	DeleteListSection(ctx context.Context, id uuid.UUID, updatedAt time.Time) (bool, error)
	GetListSections(context.Context, listSectionModel.GetListSectionsReq) ([]listSectionModel.ListSection, error)

	CheckAccess(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
	LockLive(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
	HasLiveSectionsInList(ctx context.Context, listID uuid.UUID) (bool, error)
}

var _ TaskListService = new(taskListService.TaskListService)

type TaskListService interface {
	LockLive(ctxTx context.Context, userID uuid.UUID, ids []uuid.UUID) error
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

// LiveTasksChecker - интерфейс сервиса задач, отвечающего, есть ли в разделе неудаленные задачи.
// Сервис задач сам зависит от этого сервиса, поэтому подключается после создания через SetLiveTasksChecker
type LiveTasksChecker interface {
	HasLiveTasksInSection(ctx context.Context, sectionID uuid.UUID) (bool, error)
}

// SetLiveTasksChecker подключает проверку задач раздела перед его удалением
func (s *ListSectionService) SetLiveTasksChecker(liveTasksChecker LiveTasksChecker) {
	s.liveTasksChecker = liveTasksChecker
}
