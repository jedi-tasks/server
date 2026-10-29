package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"server/internal/modules/tag/model"
	"server/internal/modules/tag/repository/tagDDL"
)

// GetTags возвращает неудаленные записи пользователя по фильтрам
func (r *TagRepository) GetTags(ctx context.Context, req model.GetTagsReq) (res []model.Tag, err error) {
	ctx, span := tracer.Start(ctx, "GetTags")
	defer span.End()

	q := sq.
		Select(ddlHelper.SelectAll).
		From(tagDDL.Table).
		Where(sq.Eq{
			tagDDL.ColumnUserID:    req.Necessary.UserID,
			tagDDL.ColumnDeletedAt: nil,
		})

	// Добавляем фильтры
	if len(req.IDs) != 0 {
		q = q.Where(sq.Eq{tagDDL.ColumnID: req.IDs})
	}

	return res, r.db.Select(ctx, &res, q)
}
