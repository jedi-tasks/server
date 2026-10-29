package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"pkg/slices"
	"server/internal/modules/tag/repository/tagDDL"
	"server/internal/utils/errors"
)

// CheckAccess проверяет, что все записи существуют, не удалены и принадлежат пользователю
func (r *TagRepository) CheckAccess(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "CheckAccess")
	defer span.End()

	ids = slices.GetUniqueFields(ids, func(id uuid.UUID) uuid.UUID { return id })

	var count int
	if err := r.db.Get(ctx, &count, sq.
		Select("COUNT(*)").
		From(tagDDL.Table).
		Where(sq.Eq{
			tagDDL.ColumnID:        ids,
			tagDDL.ColumnUserID:    userID,
			tagDDL.ColumnDeletedAt: nil,
		}),
	); err != nil {
		return err
	}

	if count != len(ids) {
		return errors.Forbidden.New("You don't have access to some of the requested tags").
			WithContextParams(ctx).
			WithParams("UserID", userID, "IDs", ids)
	}

	return nil
}
