package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/taskList/model"
	"server/internal/modules/taskList/repository/taskListDDL"
	"server/internal/utils/errors"
)

// UpdateTaskList меняет переданные поля, если изменение новее сохраненного (побеждает последняя
// запись). Возвращает false, если запись не изменилась, потому что на сервере уже более новая версия
func (r *TaskListRepository) UpdateTaskList(ctx context.Context, fields model.UpdateTaskListReq) (updated bool, err error) {
	ctx, span := tracer.Start(ctx, "UpdateTaskList")
	defer span.End()

	updates := make(map[string]any)

	// Добавляем в запрос поля, которые нужно изменить
	if fields.Name != nil {
		updates[taskListDDL.ColumnName] = *fields.Name
	}
	if fields.ColorHex != nil {
		updates[taskListDDL.ColumnColorHex] = *fields.ColorHex
	}
	if fields.Icon != nil {
		updates[taskListDDL.ColumnIcon] = *fields.Icon
	}
	if fields.SortOrder != nil {
		updates[taskListDDL.ColumnSortOrder] = *fields.SortOrder
	}
	if fields.SortMode != nil {
		updates[taskListDDL.ColumnSortMode] = *fields.SortMode
	}

	// Проверяем, что есть поля для обновления
	if len(updates) == 0 {
		return false, errors.BadRequest.New("No fields to update").WithContextParams(ctx)
	}
	updates[taskListDDL.ColumnUpdatedAt] = fields.UpdatedAt

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(taskListDDL.Table).
		SetMap(updates).
		Where(sq.Eq{taskListDDL.ColumnID: fields.ID, taskListDDL.ColumnDeletedAt: nil}).
		Where(sq.Lt{taskListDDL.ColumnUpdatedAt: fields.UpdatedAt}),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
