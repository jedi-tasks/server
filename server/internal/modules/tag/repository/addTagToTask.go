package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/tag/repository/taskToTagDDL"
)

// AddTagToTask привязывает тег к задаче или восстанавливает удаленную связь, если привязка новее
// последнего изменения связи. Возвращает false, если связь уже есть или на сервере более новая версия
func (r *TagRepository) AddTagToTask(ctx context.Context, userID, taskID, tagID uuid.UUID, updatedAt time.Time) (added bool, err error) {
	ctx, span := tracer.Start(ctx, "AddTagToTask")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Insert(taskToTagDDL.Table).
		SetMap(map[string]any{
			taskToTagDDL.ColumnUserID:    userID,
			taskToTagDDL.ColumnTaskID:    taskID,
			taskToTagDDL.ColumnTagID:     tagID,
			taskToTagDDL.ColumnUpdatedAt: updatedAt,
		}).
		Suffix(
			"ON CONFLICT ("+taskToTagDDL.ColumnTaskID+", "+taskToTagDDL.ColumnTagID+") DO UPDATE SET "+
				taskToTagDDL.ColumnDeletedAt+" = NULL, "+
				taskToTagDDL.ColumnUpdatedAt+" = EXCLUDED."+taskToTagDDL.ColumnUpdatedAt+" "+
				"WHERE "+taskToTagDDL.Table+"."+taskToTagDDL.ColumnDeletedAt+" IS NOT NULL "+
				"AND "+taskToTagDDL.Table+"."+taskToTagDDL.ColumnUpdatedAt+" < EXCLUDED."+taskToTagDDL.ColumnUpdatedAt,
		),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
