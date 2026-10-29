package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/priority/model"
	"server/internal/modules/priority/repository/priorityDDL"
	"server/internal/utils/errors"
)

// UpdatePriority меняет переданные поля, если изменение новее сохраненного (побеждает последняя
// запись). Возвращает false, если запись не изменилась, потому что на сервере уже более новая версия
func (r *PriorityRepository) UpdatePriority(ctx context.Context, fields model.UpdatePriorityReq) (updated bool, err error) {
	ctx, span := tracer.Start(ctx, "UpdatePriority")
	defer span.End()

	updates := make(map[string]any)

	// Добавляем в запрос поля, которые нужно изменить
	if fields.Name != nil {
		updates[priorityDDL.ColumnName] = *fields.Name
	}
	if fields.ColorHex != nil {
		updates[priorityDDL.ColumnColorHex] = *fields.ColorHex
	}
	if fields.Marks != nil {
		updates[priorityDDL.ColumnMarks] = *fields.Marks
	}
	if fields.SortOrder != nil {
		updates[priorityDDL.ColumnSortOrder] = *fields.SortOrder
	}

	// Проверяем, что есть поля для обновления
	if len(updates) == 0 {
		return false, errors.BadRequest.New("No fields to update").WithContextParams(ctx)
	}
	updates[priorityDDL.ColumnUpdatedAt] = fields.UpdatedAt

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(priorityDDL.Table).
		SetMap(updates).
		Where(sq.Eq{priorityDDL.ColumnID: fields.ID, priorityDDL.ColumnDeletedAt: nil}).
		Where(sq.Lt{priorityDDL.ColumnUpdatedAt: fields.UpdatedAt}),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
