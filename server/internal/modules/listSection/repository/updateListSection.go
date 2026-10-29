package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/listSection/model"
	"server/internal/modules/listSection/repository/listSectionDDL"
	"server/internal/utils/errors"
)

// UpdateListSection меняет переданные поля, если изменение новее сохраненного (побеждает последняя
// запись). Возвращает false, если запись не изменилась, потому что на сервере уже более новая версия
func (r *ListSectionRepository) UpdateListSection(ctx context.Context, fields model.UpdateListSectionReq) (updated bool, err error) {
	ctx, span := tracer.Start(ctx, "UpdateListSection")
	defer span.End()

	updates := make(map[string]any)

	// Добавляем в запрос поля, которые нужно изменить
	if fields.ListID != nil {
		updates[listSectionDDL.ColumnListID] = *fields.ListID
	}
	if fields.Name != nil {
		updates[listSectionDDL.ColumnName] = *fields.Name
	}
	if fields.SortOrder != nil {
		updates[listSectionDDL.ColumnSortOrder] = *fields.SortOrder
	}

	// Проверяем, что есть поля для обновления
	if len(updates) == 0 {
		return false, errors.BadRequest.New("No fields to update").WithContextParams(ctx)
	}
	updates[listSectionDDL.ColumnUpdatedAt] = fields.UpdatedAt

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(listSectionDDL.Table).
		SetMap(updates).
		Where(sq.Eq{listSectionDDL.ColumnID: fields.ID, listSectionDDL.ColumnDeletedAt: nil}).
		Where(sq.Lt{listSectionDDL.ColumnUpdatedAt: fields.UpdatedAt}),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
