package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/taskList/repository/taskListDDL"
)

// GetInboxID возвращает идентификатор системного списка Inbox пользователя
func (r *TaskListRepository) GetInboxID(ctx context.Context, userID uuid.UUID) (id uuid.UUID, err error) {
	ctx, span := tracer.Start(ctx, "GetInboxID")
	defer span.End()

	return id, r.db.Get(ctx, &id, sq.
		Select(taskListDDL.ColumnID).
		From(taskListDDL.Table).
		Where(sq.Eq{
			taskListDDL.ColumnUserID:    userID,
			taskListDDL.ColumnIsInbox:   true,
			taskListDDL.ColumnDeletedAt: nil,
		}),
	)
}
