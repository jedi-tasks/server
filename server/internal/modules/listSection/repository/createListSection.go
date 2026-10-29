package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/listSection/repository/listSectionDDL"
	"server/internal/modules/listSection/repository/model"
)

// CreateListSection создает запись. Повтор того же запроса из офлайн-очереди ничего не меняет:
// возвращается false, если запись с таким идентификатором уже есть
func (r *ListSectionRepository) CreateListSection(ctx context.Context, req model.CreateListSectionReq) (created bool, err error) {
	ctx, span := tracer.Start(ctx, "CreateListSection")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Insert(listSectionDDL.Table).
		SetMap(map[string]any{
			listSectionDDL.ColumnID:        req.ID,
			listSectionDDL.ColumnUserID:    req.UserID,
			listSectionDDL.ColumnListID:    req.ListID,
			listSectionDDL.ColumnName:      req.Name,
			listSectionDDL.ColumnSortOrder: req.SortOrder,
			listSectionDDL.ColumnUpdatedAt: req.UpdatedAt,
		}).
		Suffix("ON CONFLICT ("+listSectionDDL.ColumnID+") DO NOTHING"),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
