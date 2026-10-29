package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/auditLog/repository/auditLogDDL"
	"server/internal/modules/auditLog/repository/model"
)

// CreateAuditLog создает новую запись аудит-лога и возвращает её идентификатор
func (r *AuditLogRepository) CreateAuditLog(ctx context.Context, req model.CreateAuditLogReq) (uint32, error) {
	ctx, span := tracer.Start(ctx, "CreateAuditLog")
	defer span.End()

	// Создаем запись аудит-лога
	return r.db.ExecWithLastInsertID(ctx, sq.Insert(auditLogDDL.Table).
		SetMap(map[string]any{
			auditLogDDL.ColumnEntity:         req.Entity,
			auditLogDDL.ColumnMethod:         req.Method,
			auditLogDDL.ColumnEntityID:       req.EntityID,
			auditLogDDL.ColumnSnapshotBefore: req.SnapshotBefore,
			auditLogDDL.ColumnSnapshotAfter:  req.SnapshotAfter,
			auditLogDDL.ColumnUserID:         req.UserID,
			auditLogDDL.ColumnDeviceID:       req.DeviceID,
		}),
	)
}
