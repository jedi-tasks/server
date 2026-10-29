package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/tag/model"
	"server/internal/modules/tag/repository/tagDDL"
	"server/internal/utils/errors"
)

// UpdateTag меняет переданные поля, если изменение новее сохраненного (побеждает последняя
// запись). Возвращает false, если запись не изменилась, потому что на сервере уже более новая версия
func (r *TagRepository) UpdateTag(ctx context.Context, fields model.UpdateTagReq) (updated bool, err error) {
	ctx, span := tracer.Start(ctx, "UpdateTag")
	defer span.End()

	updates := make(map[string]any)

	// Добавляем в запрос поля, которые нужно изменить
	if fields.Name != nil {
		updates[tagDDL.ColumnName] = *fields.Name
	}

	// Проверяем, что есть поля для обновления
	if len(updates) == 0 {
		return false, errors.BadRequest.New("No fields to update").WithContextParams(ctx)
	}
	updates[tagDDL.ColumnUpdatedAt] = fields.UpdatedAt

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(tagDDL.Table).
		SetMap(updates).
		Where(sq.Eq{tagDDL.ColumnID: fields.ID, tagDDL.ColumnDeletedAt: nil}).
		Where(sq.Lt{tagDDL.ColumnUpdatedAt: fields.UpdatedAt}),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
