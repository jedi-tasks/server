package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	taskListModel "server/internal/modules/taskList/model"
	taskListRepository "server/internal/modules/taskList/repository"
	taskListRepoModel "server/internal/modules/taskList/repository/model"
	"server/internal/modules/transactor"
	userService "server/internal/modules/user/service"
	"server/internal/utils/necessary"
)

var tracer = otel.Tracer("/server/internal/modules/taskList/service")

type TaskListService struct {
	taskListRepository TaskListRepository
	transactor         Transactor
	auditLogService    AuditLogService
	userService        UserService

	listContentsDeleters []ListContentsDeleter
}

func NewTaskListService(
	taskListRepository TaskListRepository,
	transactor Transactor,
	auditLogService AuditLogService,
	userService UserService,
) *TaskListService {
	return &TaskListService{
		taskListRepository: taskListRepository,
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

var _ TaskListRepository = new(taskListRepository.TaskListRepository)

type TaskListRepository interface {
	CreateTaskList(context.Context, taskListRepoModel.CreateTaskListReq) (bool, error)
	UpdateTaskList(context.Context, taskListModel.UpdateTaskListReq) (bool, error)
	DeleteTaskList(ctx context.Context, id uuid.UUID, updatedAt time.Time) (bool, error)
	GetTaskLists(context.Context, taskListModel.GetTaskListsReq) ([]taskListModel.TaskList, error)

	CheckAccess(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error

	CreateInbox(context.Context, taskListRepoModel.CreateTaskListReq) error
	GetInboxID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}

// ListContentsDeleter - интерфейс сервиса, мягко удаляющего содержимое списка (разделы, задачи) внутри
// открытой транзакции. Возвращает идентификатор последней записи аудит-лога (0, если удалять нечего).
// Сервисы разделов и задач сами зависят от сервиса списков, поэтому подключаются после создания через
// SetListContentsDeleters, а не через конструктор
type ListContentsDeleter interface {
	DeleteByListID(ctxTx context.Context, necessary necessary.NecessaryUserInformation, listID uuid.UUID, updatedAt time.Time) (uint32, error)
}

// SetListContentsDeleters подключает сервисы, удаляющие содержимое списка вместе с ним
func (s *TaskListService) SetListContentsDeleters(deleters ...ListContentsDeleter) {
	s.listContentsDeleters = deleters
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
