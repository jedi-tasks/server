package mcp

import (
	"time"

	mcpTransport "server/internal/utils/mcp"

	"server/internal/modules/task/model"
)

// TaskView — задача в ответе инструмента: идентификаторы строками, необязательные поля опущены
type TaskView struct {
	ID              string     `json:"id"`                        // Идентификатор
	ListID          string     `json:"listID"`                    // Идентификатор списка
	SectionID       *string    `json:"sectionID,omitempty"`       // Идентификатор раздела
	Title           string     `json:"title"`                     // Название задачи
	Notes           string     `json:"notes,omitempty"`           // Описание
	PriorityID      *string    `json:"priorityID,omitempty"`      // Идентификатор приоритета, нет — «не разобрано»
	PlannedDay      *string    `json:"plannedDay,omitempty"`      // Планируемый срок, yyyy-MM-dd
	DeadlineDay     *string    `json:"deadlineDay,omitempty"`     // Дедлайн, yyyy-MM-dd
	DurationMinutes *int32     `json:"durationMinutes,omitempty"` // Планируемое время в минутах
	ManualOrder     float64    `json:"manualOrder"`               // Ручной порядок
	CompletedAt     *time.Time `json:"completedAt,omitempty"`     // Дата и время выполнения
	ExternalID      *string    `json:"externalID,omitempty"`      // Идентификатор напоминания при импорте
	CreatedAt       time.Time  `json:"createdAt"`                 // Дата и время создания
	UpdatedAt       time.Time  `json:"updatedAt"`                 // Дата и время последнего изменения
}

// NewTaskView конвертирует модель задачи в ответ инструмента
func NewTaskView(task model.Task) TaskView {
	return TaskView{
		ID:              task.ID.String(),
		ListID:          task.ListID.String(),
		SectionID:       mcpTransport.UUIDString(task.SectionID),
		Title:           task.Title,
		Notes:           task.Notes,
		PriorityID:      mcpTransport.UUIDString(task.PriorityID),
		PlannedDay:      task.PlannedDay,
		DeadlineDay:     task.DeadlineDay,
		DurationMinutes: task.DurationMinutes,
		ManualOrder:     task.ManualOrder,
		CompletedAt:     task.CompletedAt,
		ExternalID:      task.ExternalID,
		CreatedAt:       task.CreatedAt,
		UpdatedAt:       task.UpdatedAt,
	}
}
