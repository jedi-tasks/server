package service

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	blockerModel "server/internal/modules/blocker/model"
	blockerService "server/internal/modules/blocker/service"
	durationOptionModel "server/internal/modules/durationOption/model"
	durationOptionService "server/internal/modules/durationOption/service"
	listSectionModel "server/internal/modules/listSection/model"
	listSectionService "server/internal/modules/listSection/service"
	priorityModel "server/internal/modules/priority/model"
	priorityService "server/internal/modules/priority/service"
	settingsModel "server/internal/modules/settings/model"
	settingsService "server/internal/modules/settings/service"
	tagModel "server/internal/modules/tag/model"
	tagService "server/internal/modules/tag/service"
	taskModel "server/internal/modules/task/model"
	taskService "server/internal/modules/task/service"
	taskListModel "server/internal/modules/taskList/model"
	taskListService "server/internal/modules/taskList/service"
	"server/internal/modules/transactor"
	userModel "server/internal/modules/user/model"
	userService "server/internal/modules/user/service"
	"server/internal/utils/syncNotifier"
)

var tracer = otel.Tracer("/server/internal/modules/sync/service")

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
}

var _ UserService = new(userService.UserService)

// UserService - интерфейс сервиса пользователей
type UserService interface {
	SetDevicePendingSync(ctx context.Context, userID uuid.UUID, deviceID string, pendingCheckpoint uint32, pendingSyncToken uuid.UUID) error
	ConfirmDeviceSync(ctx context.Context, userID uuid.UUID, deviceID string, pendingSyncToken uuid.UUID) error
	GetUsers(ctx context.Context, filters userModel.GetUsersReq) ([]userModel.User, error)
}

var _ AuditLogService = new(auditLogService.AuditLogService)

// AuditLogService - интерфейс сервиса аудит-лога
type AuditLogService interface {
	GetAuditLogsSince(ctx context.Context, userID uuid.UUID, sinceID uint32) ([]auditLogModel.AuditLog, error)
}

var _ TaskListService = new(taskListService.TaskListService)

// TaskListService - интерфейс сервиса сущности taskList
type TaskListService interface {
	GetTaskLists(context.Context, taskListModel.GetTaskListsReq) ([]taskListModel.TaskList, error)
}

var _ ListSectionService = new(listSectionService.ListSectionService)

// ListSectionService - интерфейс сервиса сущности listSection
type ListSectionService interface {
	GetListSections(context.Context, listSectionModel.GetListSectionsReq) ([]listSectionModel.ListSection, error)
}

var _ PriorityService = new(priorityService.PriorityService)

// PriorityService - интерфейс сервиса сущности priority
type PriorityService interface {
	GetPriorities(context.Context, priorityModel.GetPrioritiesReq) ([]priorityModel.Priority, error)
}

var _ DurationOptionService = new(durationOptionService.DurationOptionService)

// DurationOptionService - интерфейс сервиса сущности durationOption
type DurationOptionService interface {
	GetDurationOptions(context.Context, durationOptionModel.GetDurationOptionsReq) ([]durationOptionModel.DurationOption, error)
}

var _ TaskService = new(taskService.TaskService)

// TaskService - интерфейс сервиса сущности task
type TaskService interface {
	GetTasks(context.Context, taskModel.GetTasksReq) ([]taskModel.Task, error)
}

var _ TagService = new(tagService.TagService)

// TagService - интерфейс сервиса сущности tag
type TagService interface {
	GetTags(context.Context, tagModel.GetTagsReq) ([]tagModel.Tag, error)
}

var _ BlockerService = new(blockerService.BlockerService)

// BlockerService - интерфейс сервиса сущности blocker
type BlockerService interface {
	GetBlockers(context.Context, blockerModel.GetBlockersReq) ([]blockerModel.Blocker, error)
}

var _ TagLinkService = new(tagService.TagService)

// TagLinkService - интерфейс сервиса тегов для связей задач и тегов
type TagLinkService interface {
	GetTasksToTags(context.Context, tagModel.GetTasksToTagsReq) ([]tagModel.TaskToTag, error)
}

var _ SettingsService = new(settingsService.SettingsService)

// SettingsService - интерфейс сервиса настроек
type SettingsService interface {
	GetSettings(ctx context.Context, userID uuid.UUID) (settingsModel.Settings, error)
}

var _ SyncNotifier = new(syncNotifier.SyncNotifier)

// SyncNotifier - интерфейс будильника активных SubscribeToSync-стримов (см. пакет syncNotifier)
type SyncNotifier interface {
	Subscribe(userID uuid.UUID) (ch chan struct{}, unsubscribe func())
}

// SyncService - сервис синхронизации изменений между устройствами
type SyncService struct {
	transactor            Transactor
	userService           UserService
	auditLogService       AuditLogService
	taskListService       TaskListService
	listSectionService    ListSectionService
	priorityService       PriorityService
	durationOptionService DurationOptionService
	taskService           TaskService
	tagService            TagService
	blockerService        BlockerService
	tagLinkService        TagLinkService
	settingsService       SettingsService
	syncNotifier          SyncNotifier
}

// NewSyncService создает новый сервис синхронизации
func NewSyncService(
	transactor Transactor,
	userService UserService,
	auditLogService AuditLogService,
	taskListService TaskListService,
	listSectionService ListSectionService,
	priorityService PriorityService,
	durationOptionService DurationOptionService,
	taskService TaskService,
	tagService TagService,
	blockerService BlockerService,
	tagLinkService TagLinkService,
	settingsService SettingsService,
	syncNotifier SyncNotifier,
) *SyncService {
	return &SyncService{
		transactor:            transactor,
		userService:           userService,
		auditLogService:       auditLogService,
		taskListService:       taskListService,
		listSectionService:    listSectionService,
		priorityService:       priorityService,
		durationOptionService: durationOptionService,
		taskService:           taskService,
		tagService:            tagService,
		blockerService:        blockerService,
		tagLinkService:        tagLinkService,
		settingsService:       settingsService,
		syncNotifier:          syncNotifier,
	}
}
