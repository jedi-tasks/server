package mcp

import (
	"time"

	"server/internal/modules/blocker/model"
	mcpTransport "server/internal/utils/mcp"
)

// BlockerView — блокер в ответе инструмента
type BlockerView struct {
	ID             string     `json:"id"`                       // Идентификатор
	TaskID         string     `json:"taskID"`                   // Заблокированная задача
	BlockingTaskID *string    `json:"blockingTaskID,omitempty"` // Блокирующая задача
	Reason         *string    `json:"reason,omitempty"`         // Текстовая причина
	SinceDay       string     `json:"sinceDay"`                 // Дата начала, yyyy-MM-dd
	ResolvedAt     *time.Time `json:"resolvedAt,omitempty"`     // Дата и время снятия, нет — блокер действует
	UpdatedAt      time.Time  `json:"updatedAt"`                // Дата и время последнего изменения
}

// NewBlockerView конвертирует модель блокера в ответ инструмента
func NewBlockerView(blocker model.Blocker) BlockerView {
	return BlockerView{
		ID:             blocker.ID.String(),
		TaskID:         blocker.TaskID.String(),
		BlockingTaskID: mcpTransport.UUIDString(blocker.BlockingTaskID),
		Reason:         blocker.Reason,
		SinceDay:       blocker.SinceDay,
		ResolvedAt:     blocker.ResolvedAt,
		UpdatedAt:      blocker.UpdatedAt,
	}
}
