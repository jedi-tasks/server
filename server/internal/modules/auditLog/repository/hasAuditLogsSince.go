package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/auditLog/repository/auditLogDDL"
)

// HasAuditLogsSince легковесно проверяет, есть ли хотя бы одна запись аудит-лога пользователя с id
// больше sinceID, не вычитывая сами записи
func (r *AuditLogRepository) HasAuditLogsSince(ctx context.Context, userID uuid.UUID, sinceID uint32) (bool, error) {
	ctx, span := tracer.Start(ctx, "HasAuditLogsSince")
	defer span.End()

	var rows []int
	if err := r.db.Select(ctx, &rows, sq.
		Select("1").
		From(auditLogDDL.Table).
		Where(sq.Eq{
			auditLogDDL.ColumnUserID: userID,
		}).
		Where(sq.Gt{
			auditLogDDL.ColumnID: sinceID,
		}).
		Limit(1),
	); err != nil {
		return false, err
	}

	return len(rows) > 0, nil
}
