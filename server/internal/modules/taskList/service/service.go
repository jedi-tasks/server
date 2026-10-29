package service

import (
	"context"
	userModel "server/internal/modules/user/model"
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
)

var tracer = otel.Tracer("/server/internal/modules/taskList/service")

// TaskListService - сервис списков задач
type TaskListService struct {
	taskListRepository TaskListRepository
	transactor         Transactor
	auditLogService    AuditLogService
	userService        UserService

	liveTasksChecker    LiveTasksChecker
	liveSectionsChecker LiveSectionsChecker
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
	LockLive(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error

	CreateInbox(context.Context, taskListRepoModel.CreateTaskListReq) error
	GetInboxID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}

// LiveTasksChecker - интерфейс сервиса задач, отвечающего, есть ли в списке неудаленные задачи.
// Сервис задач сам зависит от сервиса списков, поэтому подключается после создания через
// SetListContentsCheckers, а не через конструктор
type LiveTasksChecker interface {
	HasLiveTasksInList(ctx context.Context, listID uuid.UUID) (bool, error)
}

// LiveSectionsChecker - интерфейс сервиса разделов, отвечающего, есть ли в списке неудаленные разделы.
// Подключается после создания через SetListContentsCheckers по той же причине
type LiveSectionsChecker interface {
	HasLiveSectionsInList(ctx context.Context, listID uuid.UUID) (bool, error)
}

// SetListContentsCheckers подключает проверки содержимого списка перед его удалением
func (s *TaskListService) SetListContentsCheckers(liveTasksChecker LiveTasksChecker, liveSectionsChecker LiveSectionsChecker) {
	s.liveTasksChecker = liveTasksChecker
	s.liveSectionsChecker = liveSectionsChecker
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
