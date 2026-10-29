package mcp

import (
	"context"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/task/model"
	mcpTransport "server/internal/utils/mcp"
	"server/internal/utils/validate"
)

// UpdateTaskIn — аргументы update_task, меняются только переданные поля
type UpdateTaskIn struct {
	ID              string   `json:"id" jsonschema:"Идентификатор задачи"`
	ListID          *string  `json:"listID,omitempty" jsonschema:"Перенос в другой список"`
	SectionID       *string  `json:"sectionID,omitempty" jsonschema:"Раздел, пустая строка — без раздела"`
	Title           *string  `json:"title,omitempty" jsonschema:"Название задачи"`
	Notes           *string  `json:"notes,omitempty" jsonschema:"Описание"`
	PriorityID      *string  `json:"priorityID,omitempty" jsonschema:"Приоритет, пустая строка — «не разобрано»"`
	PlannedDay      *string  `json:"plannedDay,omitempty" jsonschema:"Планируемый срок yyyy-MM-dd, пустая строка — без срока"`
	DeadlineDay     *string  `json:"deadlineDay,omitempty" jsonschema:"Дедлайн yyyy-MM-dd, пустая строка — без дедлайна"`
	DurationMinutes *int32   `json:"durationMinutes,omitempty" jsonschema:"Планируемое время в минутах, 0 — без оценки"`
	ManualOrder     *float64 `json:"manualOrder,omitempty" jsonschema:"Ручной порядок"`
}

// EmptyOut — результат правки без данных
type EmptyOut struct{}

// UpdateTask изменение записи
func (t *TaskTools) UpdateTask(ctx context.Context, _ *sdk.CallToolRequest, in UpdateTaskIn) (*sdk.CallToolResult, EmptyOut, error) {
	var out EmptyOut

	// Собираем запрос из аргументов и контекста
	necessary, err := mcpTransport.Necessary(ctx)
	if err != nil {
		return nil, out, err
	}
	id, err := mcpTransport.ParseUUID("id", in.ID)
	if err != nil {
		return nil, out, err
	}
	listID, err := mcpTransport.ParseOptionalUUID("listID", in.ListID)
	if err != nil {
		return nil, out, err
	}
	sectionID, err := mcpTransport.ParseOptionalUUID("sectionID", in.SectionID)
	if err != nil {
		return nil, out, err
	}
	priorityID, err := mcpTransport.ParseOptionalUUID("priorityID", in.PriorityID)
	if err != nil {
		return nil, out, err
	}
	req := model.UpdateTaskReq{
		Necessary:       necessary,
		ID:              id,
		ListID:          listID,
		SectionID:       sectionID,
		Title:           in.Title,
		Notes:           in.Notes,
		PriorityID:      priorityID,
		PlannedDay:      in.PlannedDay,
		DeadlineDay:     in.DeadlineDay,
		DurationMinutes: in.DurationMinutes,
		ManualOrder:     in.ManualOrder,
		UpdatedAt:       time.Now(),
	}

	// Валидируем запрос
	if err = validate.Validate(ctx, req); err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Вызываем метод сервиса
	if err = t.taskService.UpdateTask(ctx, req); err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}
	return nil, out, nil
}
