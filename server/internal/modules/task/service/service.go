package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	listSectionService "server/internal/modules/listSection/service"
	priorityService "server/internal/modules/priority/service"
	taskModel "server/internal/modules/task/model"
	taskRepository "server/internal/modules/task/repository"
	taskRepoModel "server/internal/modules/task/repository/model"
	taskListService "server/internal/modules/taskList/service"
	"server/internal/modules/transactor"
	userService "server/internal/modules/user/service"
	"server/internal/utils/necessary"
)

var tracer = otel.Tracer("/server/internal/modules/task/service")

type TaskService struct {
	taskRepository     TaskRepository
	transactor         Transactor
	auditLogService    AuditLogService
	userService        UserService
	taskListService    TaskListService
	listSectionService ListSectionService
	priorityService    PriorityService

	blockerResolver BlockerResolver
}

func NewTaskService(
	taskRepository TaskRepository,
	transactor Transactor,
	auditLogService AuditLogService,
	userService UserService,
	taskListService TaskListService,
	listSectionService ListSectionService,
	priorityService PriorityService,
) *TaskService {
	return &TaskService{
		taskRepository:     taskRepository,
		transactor:         transactor,
		auditLogService:    auditLogService,
		userService:        userService,
		taskListService:    taskListService,
		listSectionService: listSectionService,
		priorityService:    priorityService,
	}
}

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
	WithSyncGate(ctx context.Context, userID uuid.UUID, deviceID string, deviceSyncGate transactor.DeviceSyncGate, auditLogChangeChecker transactor.AuditLogChangeChecker, mutate func(ctxTx context.Context) (auditLogID uint32, err error)) error
}

var _ TaskRepository = new(taskRepository.TaskRepository)

type TaskRepository interface {
	CreateTask(context.Context, taskRepoModel.CreateTaskReq) (bool, error)
	UpdateTask(context.Context, taskModel.UpdateTaskReq) (bool, error)
	DeleteTask(ctx context.Context, id uuid.UUID, updatedAt time.Time) (bool, error)
	GetTasks(context.Context, taskModel.GetTasksReq) ([]taskModel.Task, error)

	CheckAccess(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error

	SetTaskCompleted(ctx context.Context, id uuid.UUID, completedAt *time.Time, updatedAt time.Time) (bool, error)
	ClearReference(ctx context.Context, userID uuid.UUID, column string, refID uuid.UUID, updatedAt time.Time) ([]uuid.UUID, error)
}

// BlockerResolver - интерфейс сервиса блокеров, снимающего блокеры выполненной задачи внутри открытой
// транзакции. Сервис блокеров сам зависит от сервиса задач, поэтому подключается после создания через
// SetBlockerResolver, а не через конструктор
type BlockerResolver interface {
	ResolveByBlockingTask(ctxTx context.Context, necessary necessary.NecessaryUserInformation, blockingTaskID uuid.UUID, resolvedAt time.Time) (uint32, error)
}

// SetBlockerResolver подключает сервис блокеров
func (s *TaskService) SetBlockerResolver(blockerResolver BlockerResolver) {
	s.blockerResolver = blockerResolver
}

var _ TaskListService = new(taskListService.TaskListService)

type TaskListService interface {
	CheckAccess(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
	GetInboxID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}

var _ ListSectionService = new(listSectionService.ListSectionService)

type ListSectionService interface {
	CheckAccess(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
}

var _ PriorityService = new(priorityService.PriorityService)

type PriorityService interface {
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
