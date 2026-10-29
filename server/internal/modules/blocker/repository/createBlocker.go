package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/blocker/repository/blockerDDL"
	"server/internal/modules/blocker/repository/model"
)

// CreateBlocker создает запись. Повтор того же запроса из офлайн-очереди ничего не меняет:
// возвращается false, если запись с таким идентификатором уже есть
func (r *BlockerRepository) CreateBlocker(ctx context.Context, req model.CreateBlockerReq) (created bool, err error) {
	ctx, span := tracer.Start(ctx, "CreateBlocker")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Insert(blockerDDL.Table).
		SetMap(map[string]any{
			blockerDDL.ColumnID:             req.ID,
			blockerDDL.ColumnUserID:         req.UserID,
			blockerDDL.ColumnTaskID:         req.TaskID,
			blockerDDL.ColumnBlockingTaskID: req.BlockingTaskID,
			blockerDDL.ColumnReason:         req.Reason,
			blockerDDL.ColumnSinceDay:       req.SinceDay,
			blockerDDL.ColumnUpdatedAt:      req.UpdatedAt,
		}).
		Suffix("ON CONFLICT ("+blockerDDL.ColumnID+") DO NOTHING"),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
