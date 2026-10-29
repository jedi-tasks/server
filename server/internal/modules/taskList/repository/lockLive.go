package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"pkg/slices"
	"server/internal/modules/taskList/repository/taskListDDL"
	"server/internal/utils/errors"
)

// LockLive проверяет, что все записи существуют, принадлежат пользователю и не удалены, и блокирует их
// до конца транзакции: параллельное удаление дождется ее и увидит ссылающиеся на них строки
func (r *TaskListRepository) LockLive(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "LockLive")
	defer span.End()

	ids = slices.GetUniqueFields(ids, func(id uuid.UUID) uuid.UUID { return id })

	// Читаем признак удаления записей с блокировкой FOR SHARE
	var isDeleted []bool
	if err := r.db.Select(ctx, &isDeleted, sq.
		Select(taskListDDL.ColumnDeletedAt+" IS NOT NULL").
		From(taskListDDL.Table).
		Where(sq.Eq{
			taskListDDL.ColumnID:     ids,
			taskListDDL.ColumnUserID: userID,
		}).
		Suffix("FOR SHARE"),
	); err != nil {
		return err
	}

	// Записи нет или она чужая
	if len(isDeleted) != len(ids) {
		return errors.Forbidden.New("You don't have access to some of the requested lists").
			WithContextParams(ctx).
			WithParams("UserID", userID, "IDs", ids)
	}

	// Живая строка не может ссылаться на удаленную
	for _, deleted := range isDeleted {
		if deleted {
			return errors.BadRequest.New("Reference to deleted lists").
				WithContextParams(ctx).
				WithParams("UserID", userID, "IDs", ids).
				WithCustomHumanText("Список удалён — выберите другой список")
		}
	}

	return nil
}
