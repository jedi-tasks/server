package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/tag/repository/model"
	"server/internal/modules/tag/repository/tagDDL"
)

// CreateTag создает запись. Повтор того же запроса из офлайн-очереди ничего не меняет:
// возвращается false, если запись с таким идентификатором уже есть
func (r *TagRepository) CreateTag(ctx context.Context, req model.CreateTagReq) (created bool, err error) {
	ctx, span := tracer.Start(ctx, "CreateTag")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Insert(tagDDL.Table).
		SetMap(map[string]any{
			tagDDL.ColumnID:        req.ID,
			tagDDL.ColumnUserID:    req.UserID,
			tagDDL.ColumnName:      req.Name,
			tagDDL.ColumnUpdatedAt: req.UpdatedAt,
		}).
		Suffix("ON CONFLICT ("+tagDDL.ColumnID+") DO NOTHING"),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
