package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/durationOption/repository/durationOptionDDL"
	"server/internal/modules/durationOption/repository/model"
)

// CreateDurationOption создает запись. Повтор того же запроса из офлайн-очереди ничего не меняет:
// возвращается false, если запись с таким идентификатором уже есть
func (r *DurationOptionRepository) CreateDurationOption(ctx context.Context, req model.CreateDurationOptionReq) (created bool, err error) {
	ctx, span := tracer.Start(ctx, "CreateDurationOption")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Insert(durationOptionDDL.Table).
		SetMap(map[string]any{
			durationOptionDDL.ColumnID:        req.ID,
			durationOptionDDL.ColumnUserID:    req.UserID,
			durationOptionDDL.ColumnTitle:     req.Title,
			durationOptionDDL.ColumnMinutes:   req.Minutes,
			durationOptionDDL.ColumnUpdatedAt: req.UpdatedAt,
		}).
		Suffix("ON CONFLICT ("+durationOptionDDL.ColumnID+") DO NOTHING"),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
