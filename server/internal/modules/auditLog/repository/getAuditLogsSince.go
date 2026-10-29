package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"pkg/ddlHelper"
	"server/internal/modules/auditLog/model"
	"server/internal/modules/auditLog/repository/auditLogDDL"
)

// GetAuditLogsSince возвращает все записи аудит-лога пользователя с id больше sinceID, отсортированные
// от старых к новым
func (r *AuditLogRepository) GetAuditLogsSince(ctx context.Context, userID uuid.UUID, sinceID uint32) (auditLogs []model.AuditLog, err error) {
	ctx, span := tracer.Start(ctx, "GetAuditLogsSince")
	defer span.End()

	q := sq.
		Select(ddlHelper.SelectAll).
		From(auditLogDDL.Table).
		Where(sq.Eq{
			auditLogDDL.ColumnUserID: userID,
		}).
		Where(sq.Gt{
			auditLogDDL.ColumnID: sinceID,
		}).
		OrderBy(auditLogDDL.ColumnID)

	return auditLogs, r.db.Select(ctx, &auditLogs, q)
}
