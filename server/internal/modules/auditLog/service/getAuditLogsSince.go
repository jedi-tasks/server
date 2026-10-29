package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/modules/auditLog/model"
)

// GetAuditLogsSince возвращает все записи аудит-лога пользователя с id больше sinceID
func (s *AuditLogService) GetAuditLogsSince(ctx context.Context, userID uuid.UUID, sinceID uint32) ([]model.AuditLog, error) {
	ctx, span := tracer.Start(ctx, "GetAuditLogsSince")
	defer span.End()

	return s.auditLogRepository.GetAuditLogsSince(ctx, userID, sinceID)
}
