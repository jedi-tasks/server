package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/priority/repository/model"
	"server/internal/modules/priority/repository/priorityDDL"
)

// CreatePriority создает запись. Повтор того же запроса из офлайн-очереди ничего не меняет:
// возвращается false, если запись с таким идентификатором уже есть
func (r *PriorityRepository) CreatePriority(ctx context.Context, req model.CreatePriorityReq) (created bool, err error) {
	ctx, span := tracer.Start(ctx, "CreatePriority")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Insert(priorityDDL.Table).
		SetMap(map[string]any{
			priorityDDL.ColumnID:        req.ID,
			priorityDDL.ColumnUserID:    req.UserID,
			priorityDDL.ColumnName:      req.Name,
			priorityDDL.ColumnColorHex:  req.ColorHex,
			priorityDDL.ColumnMarks:     req.Marks,
			priorityDDL.ColumnSortOrder: req.SortOrder,
			priorityDDL.ColumnUpdatedAt: req.UpdatedAt,
		}).
		Suffix("ON CONFLICT ("+priorityDDL.ColumnID+") DO NOTHING"),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
