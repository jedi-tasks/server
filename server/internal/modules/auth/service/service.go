package service

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"

	"github.com/google/uuid"

	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	durationOptionService "server/internal/modules/durationOption/service"
	priorityService "server/internal/modules/priority/service"
	taskListService "server/internal/modules/taskList/service"
	"server/internal/modules/transactor"
	userModel "server/internal/modules/user/model"
	userRepository "server/internal/modules/user/repository"
	userRepoModel "server/internal/modules/user/repository/model"
)

var tracer = otel.Tracer("/server/internal/modules/auth/service")

var _ UserRepository = new(userRepository.UserRepository)
var _ GeneralRepository = new(transactor.Transactor)

type UserRepository interface {
	GetUsers(context.Context, userModel.GetUsersReq) ([]userModel.User, error)
	CreateUser(context.Context, userModel.CreateReq) error

	CreateDevice(context.Context, userModel.Device) error
	DeleteDevice(ctx context.Context, userID uuid.UUID, deviceID string) error
	UpdateDevice(context.Context, userRepoModel.UpdateDeviceReq) error
	GetDevices(context.Context, userRepoModel.GetDevicesReq) ([]userModel.Device, error)
	RotateRefreshToken(ctx context.Context, userID uuid.UUID, deviceID string, newRefreshToken string, graceWindow time.Duration) error
}

type GeneralRepository interface {
	WithinTransaction(ctx context.Context, callback func(ctx context.Context) error) error
}

var _ AuditLogService = new(auditLogService.AuditLogService)

// AuditLogService - интерфейс сервиса аудит-лога
type AuditLogService interface {
	TrackMutation(context.Context, auditLogModel.TrackMutationReq) (uint32, error)
}

// DefaultsCreator - интерфейс сервиса, создающего записи по умолчанию для нового пользователя
// (Inbox, приоритеты и длительности системы «Путь джедая»)
type DefaultsCreator interface {
	CreateDefaults(ctx context.Context, userID uuid.UUID, deviceID string, now time.Time) error
}

var _ DefaultsCreator = new(taskListService.TaskListService)
var _ DefaultsCreator = new(priorityService.PriorityService)
var _ DefaultsCreator = new(durationOptionService.DurationOptionService)

type AuthService struct {
	userRepository    UserRepository
	generalRepository GeneralRepository
	generalSalt       []byte
	signUpEnabled     bool
	auditLogService   AuditLogService
	defaultsCreators  []DefaultsCreator
}

func NewAuthService(
	userRepository UserRepository,
	generalRepository GeneralRepository,
	generalSalt []byte,
	signUpEnabled bool,
	auditLogService AuditLogService,
	defaultsCreators []DefaultsCreator,
) *AuthService {
	return &AuthService{
		userRepository:    userRepository,
		generalRepository: generalRepository,
		generalSalt:       generalSalt,
		signUpEnabled:     signUpEnabled,
		auditLogService:   auditLogService,
		defaultsCreators:  defaultsCreators,
	}
}
