package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"server/internal/modules/durationOption/model"
	"server/internal/modules/durationOption/repository/durationOptionDDL"
)

// GetDurationOptions возвращает неудаленные записи пользователя по фильтрам
func (r *DurationOptionRepository) GetDurationOptions(ctx context.Context, req model.GetDurationOptionsReq) (res []model.DurationOption, err error) {
	ctx, span := tracer.Start(ctx, "GetDurationOptions")
	defer span.End()

	q := sq.
		Select(ddlHelper.SelectAll).
		From(durationOptionDDL.Table).
		Where(sq.Eq{
			durationOptionDDL.ColumnUserID:    req.Necessary.UserID,
			durationOptionDDL.ColumnDeletedAt: nil,
		})

	// Добавляем фильтры
	if len(req.IDs) != 0 {
		q = q.Where(sq.Eq{durationOptionDDL.ColumnID: req.IDs})
	}

	return res, r.db.Select(ctx, &res, q)
}
