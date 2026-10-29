package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/google/uuid"

	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	"server/internal/modules/transactor"
	userModel "server/internal/modules/user/model"
	userRepository "server/internal/modules/user/repository"
	userRepoModel "server/internal/modules/user/repository/model"
)

var tracer = otel.Tracer("/server/internal/modules/user/service")

var _ UserRepository = new(userRepository.UserRepository)
var _ GeneralRepository = new(transactor.Transactor)

type GeneralRepository interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
}

type UserRepository interface {
	CreateUser(context.Context, userModel.CreateReq) error
	GetUsers(context.Context, userModel.GetUsersReq) ([]userModel.User, error)
	UpdateUser(context.Context, userRepoModel.UpdateUserReq) error

	GetDevices(context.Context, userRepoModel.GetDevicesReq) ([]userModel.Device, error)
	UpdateDevice(context.Context, userRepoModel.UpdateDeviceReq) error

	GetDeviceForUpdate(ctx context.Context, userID uuid.UUID, deviceID string) (userModel.Device, error)
	SetDevicePendingSync(ctx context.Context, userID uuid.UUID, deviceID string, pendingCheckpoint uint32, pendingSyncToken uuid.UUID) error

	ConfirmDeviceSync(ctx context.Context, userID uuid.UUID, deviceID string, pendingSyncToken uuid.UUID) (uint32, error)
	BumpDeviceCheckpoint(ctx context.Context, userID uuid.UUID, deviceID string, auditLogID uint32) error
}

var _ AuditLogService = new(auditLogService.AuditLogService)

// AuditLogService - интерфейс сервиса аудит-лога
type AuditLogService interface {
	TrackMutation(context.Context, auditLogModel.TrackMutationReq) (uint32, error)
}

type UserService struct {
	userRepository    UserRepository
	generalRepository GeneralRepository
	generalSalt       []byte
	auditLogService   AuditLogService
}

func NewUserService(
	userRepository UserRepository,
	generalRepository GeneralRepository,
	generalSalt []byte,
	auditLogService AuditLogService,
) (*UserService, error) {
	s := &UserService{
		userRepository:    userRepository,
		generalRepository: generalRepository,
		generalSalt:       generalSalt,
		auditLogService:   auditLogService,
	}

	return s, nil
}
