package mcp

import (
	"time"

	"server/internal/modules/taskList/model"
)

// TaskListView — список в ответе инструмента
type TaskListView struct {
	ID        string    `json:"id"`        // Идентификатор
	Name      string    `json:"name"`      // Название списка
	IsInbox   bool      `json:"isInbox"`   // Системный список Inbox
	ColorHex  string    `json:"colorHex"`  // Цвет списка
	Icon      string    `json:"icon"`      // SF Symbol иконки списка
	SortOrder float64   `json:"sortOrder"` // Порядок списка
	SortMode  string    `json:"sortMode"`  // Режим сортировки задач в списке
	CreatedAt time.Time `json:"createdAt"` // Дата и время создания
	UpdatedAt time.Time `json:"updatedAt"` // Дата и время последнего изменения
}

// NewTaskListView конвертирует модель списка в ответ инструмента
func NewTaskListView(list model.TaskList) TaskListView {
	return TaskListView{
		ID:        list.ID.String(),
		Name:      list.Name,
		IsInbox:   list.IsInbox,
		ColorHex:  list.ColorHex,
		Icon:      list.Icon,
		SortOrder: list.SortOrder,
		SortMode:  string(list.SortMode),
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}
}
