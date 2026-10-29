package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/durationOption/model"
	"server/internal/modules/durationOption/repository/durationOptionDDL"
	"server/internal/utils/errors"
)

// UpdateDurationOption меняет переданные поля, если изменение новее сохраненного (побеждает последняя
// запись). Возвращает false, если запись не изменилась, потому что на сервере уже более новая версия
func (r *DurationOptionRepository) UpdateDurationOption(ctx context.Context, fields model.UpdateDurationOptionReq) (updated bool, err error) {
	ctx, span := tracer.Start(ctx, "UpdateDurationOption")
	defer span.End()

	updates := make(map[string]any)

	// Добавляем в запрос поля, которые нужно изменить
	if fields.Title != nil {
		updates[durationOptionDDL.ColumnTitle] = *fields.Title
	}
	if fields.Minutes != nil {
		updates[durationOptionDDL.ColumnMinutes] = *fields.Minutes
	}

	// Проверяем, что есть поля для обновления
	if len(updates) == 0 {
		return false, errors.BadRequest.New("No fields to update").WithContextParams(ctx)
	}
	updates[durationOptionDDL.ColumnUpdatedAt] = fields.UpdatedAt

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(durationOptionDDL.Table).
		SetMap(updates).
		Where(sq.Eq{durationOptionDDL.ColumnID: fields.ID, durationOptionDDL.ColumnDeletedAt: nil}).
		Where(sq.Lt{durationOptionDDL.ColumnUpdatedAt: fields.UpdatedAt}),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
