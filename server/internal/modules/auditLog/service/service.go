package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/google/uuid"

	"server/internal/modules/auditLog/model"
	"server/internal/modules/auditLog/repository"
	repoModel "server/internal/modules/auditLog/repository/model"
	"server/internal/utils/syncNotifier"
)

var tracer = otel.Tracer("/server/internal/modules/auditLog/service")

// AuditLogService - сервис аудит-лога изменяющих действий пользователей
type AuditLogService struct {
	auditLogRepository AuditLogRepository
	syncNotifier       SyncNotifier
}

var _ AuditLogRepository = new(repository.AuditLogRepository)

// AuditLogRepository - интерфейс репозитория аудит-лога
type AuditLogRepository interface {
	CreateAuditLog(context.Context, repoModel.CreateAuditLogReq) (uint32, error)
	GetAuditLogsSince(ctx context.Context, userID uuid.UUID, sinceID uint32) ([]model.AuditLog, error)
	HasAuditLogsSince(ctx context.Context, userID uuid.UUID, sinceID uint32) (bool, error)
}

var _ SyncNotifier = new(syncNotifier.SyncNotifier)

// SyncNotifier - интерфейс будильника активных SubscribeToSync-стримов (см. пакет syncNotifier)
type SyncNotifier interface {
	Notify(userID uuid.UUID)
}

// NewAuditLogService создает новый сервис аудит-лога
func NewAuditLogService(
	auditLogRepository AuditLogRepository,
	syncNotifier SyncNotifier,
) *AuditLogService {
	return &AuditLogService{
		auditLogRepository: auditLogRepository,
		syncNotifier:       syncNotifier,
	}
}
