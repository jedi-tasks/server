package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"server/internal/modules/listSection/model"
	"server/internal/modules/listSection/repository/listSectionDDL"
)

// GetListSections возвращает неудаленные записи пользователя по фильтрам
func (r *ListSectionRepository) GetListSections(ctx context.Context, req model.GetListSectionsReq) (res []model.ListSection, err error) {
	ctx, span := tracer.Start(ctx, "GetListSections")
	defer span.End()

	q := sq.
		Select(ddlHelper.SelectAll).
		From(listSectionDDL.Table).
		Where(sq.Eq{
			listSectionDDL.ColumnUserID:    req.Necessary.UserID,
			listSectionDDL.ColumnDeletedAt: nil,
		})

	// Добавляем фильтры
	if len(req.IDs) != 0 {
		q = q.Where(sq.Eq{listSectionDDL.ColumnID: req.IDs})
	}
	if len(req.ListIDs) != 0 {
		q = q.Where(sq.Eq{listSectionDDL.ColumnListID: req.ListIDs})
	}

	return res, r.db.Select(ctx, &res, q)
}
