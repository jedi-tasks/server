package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/task/model"
	"server/internal/modules/task/repository/taskDDL"
	"server/internal/utils/errors"
)

// UpdateTask меняет переданные поля, если изменение новее сохраненного (побеждает последняя
// запись). Возвращает false, если запись не изменилась, потому что на сервере уже более новая версия
func (r *TaskRepository) UpdateTask(ctx context.Context, fields model.UpdateTaskReq) (updated bool, err error) {
	ctx, span := tracer.Start(ctx, "UpdateTask")
	defer span.End()

	updates := make(map[string]any)

	// Добавляем в запрос поля, которые нужно изменить
	if fields.ListID != nil {
		updates[taskDDL.ColumnListID] = *fields.ListID
	}
	if fields.SectionID != nil {
		if *fields.SectionID == uuid.Nil {
			updates[taskDDL.ColumnSectionID] = nil
		} else {
			updates[taskDDL.ColumnSectionID] = *fields.SectionID
		}
	}
	if fields.Title != nil {
		updates[taskDDL.ColumnTitle] = *fields.Title
	}
	if fields.Notes != nil {
		updates[taskDDL.ColumnNotes] = *fields.Notes
	}
	if fields.PriorityID != nil {
		if *fields.PriorityID == uuid.Nil {
			updates[taskDDL.ColumnPriorityID] = nil
		} else {
			updates[taskDDL.ColumnPriorityID] = *fields.PriorityID
		}
	}
	if fields.PlannedDay != nil {
		if *fields.PlannedDay == "" {
			updates[taskDDL.ColumnPlannedDay] = nil
		} else {
			updates[taskDDL.ColumnPlannedDay] = *fields.PlannedDay
		}
	}
	if fields.DeadlineDay != nil {
		if *fields.DeadlineDay == "" {
			updates[taskDDL.ColumnDeadlineDay] = nil
		} else {
			updates[taskDDL.ColumnDeadlineDay] = *fields.DeadlineDay
		}
	}
	if fields.DurationMinutes != nil {
		if *fields.DurationMinutes == 0 {
			updates[taskDDL.ColumnDurationMinutes] = nil
		} else {
			updates[taskDDL.ColumnDurationMinutes] = *fields.DurationMinutes
		}
	}
	if fields.ManualOrder != nil {
		updates[taskDDL.ColumnManualOrder] = *fields.ManualOrder
	}

	// Проверяем, что есть поля для обновления
	if len(updates) == 0 {
		return false, errors.BadRequest.New("No fields to update").WithContextParams(ctx)
	}
	updates[taskDDL.ColumnUpdatedAt] = fields.UpdatedAt

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(taskDDL.Table).
		SetMap(updates).
		Where(sq.Eq{taskDDL.ColumnID: fields.ID, taskDDL.ColumnDeletedAt: nil}).
		Where(sq.Lt{taskDDL.ColumnUpdatedAt: fields.UpdatedAt}),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
