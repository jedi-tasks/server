package service

import (
	"context"

	"github.com/google/uuid"
)

// HasAuditLogsSince легковесно проверяет, есть ли у пользователя изменения после sinceID - используется
// для отсечения мутаций от несинхронизированных устройств
func (s *AuditLogService) HasAuditLogsSince(ctx context.Context, userID uuid.UUID, sinceID uint32) (bool, error) {
	ctx, span := tracer.Start(ctx, "HasAuditLogsSince")
	defer span.End()

	return s.auditLogRepository.HasAuditLogsSince(ctx, userID, sinceID)
}
