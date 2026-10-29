package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"server/internal/modules/tag/model"
	"server/internal/modules/tag/repository/taskToTagDDL"
)

// GetTasksToTags возвращает неудаленные связи задач и тегов пользователя по фильтрам
func (r *TagRepository) GetTasksToTags(ctx context.Context, req model.GetTasksToTagsReq) (res []model.TaskToTag, err error) {
	ctx, span := tracer.Start(ctx, "GetTasksToTags")
	defer span.End()

	q := sq.
		Select(ddlHelper.SelectAll).
		From(taskToTagDDL.Table).
		Where(sq.Eq{
			taskToTagDDL.ColumnUserID:    req.Necessary.UserID,
			taskToTagDDL.ColumnDeletedAt: nil,
		})

	// Добавляем фильтры
	if len(req.TaskIDs) != 0 {
		q = q.Where(sq.Eq{taskToTagDDL.ColumnTaskID: req.TaskIDs})
	}
	if len(req.TagIDs) != 0 {
		q = q.Where(sq.Eq{taskToTagDDL.ColumnTagID: req.TagIDs})
	}

	return res, r.db.Select(ctx, &res, q)
}
