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
	"server/internal/utils/necessary"
)

var tracer = otel.Tracer("/server/internal/modules/listSection/service")

type ListSectionService struct {
	listSectionRepository ListSectionRepository
	transactor            Transactor
	auditLogService       AuditLogService
	userService           UserService
	taskListService       TaskListService

	sectionDetacher SectionDetacher
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
}

var _ TaskListService = new(taskListService.TaskListService)

type TaskListService interface {
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

// SectionDetacher - интерфейс сервиса задач, убирающего ссылку на удаленную запись внутри открытой транзакции.
// Сервис задач сам зависит от этого сервиса, поэтому подключается после создания через SetSectionDetacher
type SectionDetacher interface {
	DetachSection(ctxTx context.Context, necessary necessary.NecessaryUserInformation, id uuid.UUID, updatedAt time.Time) (uint32, error)
}

// SetSectionDetacher подключает сервис задач
func (s *ListSectionService) SetSectionDetacher(sectionDetacher SectionDetacher) {
	s.sectionDetacher = sectionDetacher
}
