package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	tagModel "server/internal/modules/tag/model"
	tagRepository "server/internal/modules/tag/repository"
	tagRepoModel "server/internal/modules/tag/repository/model"
	taskService "server/internal/modules/task/service"
	"server/internal/modules/transactor"
	userService "server/internal/modules/user/service"
)

var tracer = otel.Tracer("/server/internal/modules/tag/service")

type TagService struct {
	tagRepository   TagRepository
	transactor      Transactor
	auditLogService AuditLogService
	userService     UserService
	taskService     TaskService
}

func NewTagService(
	tagRepository TagRepository,
	transactor Transactor,
	auditLogService AuditLogService,
	userService UserService,
	taskService TaskService,
) *TagService {
	return &TagService{
		tagRepository:   tagRepository,
		transactor:      transactor,
		auditLogService: auditLogService,
		userService:     userService,
		taskService:     taskService,
	}
}

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
	WithSyncGate(ctx context.Context, userID uuid.UUID, deviceID string, deviceSyncGate transactor.DeviceSyncGate, auditLogChangeChecker transactor.AuditLogChangeChecker, mutate func(ctxTx context.Context) (auditLogID uint32, err error)) error
}

var _ TagRepository = new(tagRepository.TagRepository)

type TagRepository interface {
	CreateTag(context.Context, tagRepoModel.CreateTagReq) (bool, error)
	UpdateTag(context.Context, tagModel.UpdateTagReq) (bool, error)
	DeleteTag(ctx context.Context, id uuid.UUID, updatedAt time.Time) (bool, error)
	GetTags(context.Context, tagModel.GetTagsReq) ([]tagModel.Tag, error)

	CheckAccess(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error

	GetTasksToTags(context.Context, tagModel.GetTasksToTagsReq) ([]tagModel.TaskToTag, error)
	AddTagToTask(ctx context.Context, userID, taskID, tagID uuid.UUID, updatedAt time.Time) (bool, error)
	RemoveTagFromTask(ctx context.Context, taskID, tagID uuid.UUID, updatedAt time.Time) (bool, error)
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
