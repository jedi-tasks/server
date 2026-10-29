package mcp

import (
	"time"

	"server/internal/modules/tag/model"
)

// TagView — тег в ответе инструмента
type TagView struct {
	ID        string    `json:"id"`        // Идентификатор
	Name      string    `json:"name"`      // Название тега
	UpdatedAt time.Time `json:"updatedAt"` // Дата и время последнего изменения
}

// NewTagView конвертирует модель тега в ответ инструмента
func NewTagView(tag model.Tag) TagView {
	return TagView{ID: tag.ID.String(), Name: tag.Name, UpdatedAt: tag.UpdatedAt}
}

// TaskToTagView — связь задачи и тега в ответе инструмента
type TaskToTagView struct {
	TaskID    string    `json:"taskID"`    // Идентификатор задачи
	TagID     string    `json:"tagID"`     // Идентификатор тега
	UpdatedAt time.Time `json:"updatedAt"` // Дата и время последнего изменения
}

// NewTaskToTagView конвертирует модель связи в ответ инструмента
func NewTaskToTagView(link model.TaskToTag) TaskToTagView {
	return TaskToTagView{TaskID: link.TaskID.String(), TagID: link.TagID.String(), UpdatedAt: link.UpdatedAt}
}
