package mcp

import (
	"context"
	"time"

	"github.com/google/uuid"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/task/model"
	mcpTransport "server/internal/utils/mcp"
	"server/internal/utils/validate"
)

// CreateTaskIn — аргументы create_task
type CreateTaskIn struct {
	ID              *string  `json:"id,omitempty" jsonschema:"Идентификатор задачи, нет — сервер создаст"`
	ListID          *string  `json:"listID,omitempty" jsonschema:"Идентификатор списка, нет — Inbox"`
	SectionID       *string  `json:"sectionID,omitempty" jsonschema:"Идентификатор раздела"`
	Title           string   `json:"title" jsonschema:"Название задачи"`
	Notes           string   `json:"notes,omitempty" jsonschema:"Описание"`
	PriorityID      *string  `json:"priorityID,omitempty" jsonschema:"Идентификатор приоритета, нет — «не разобрано»"`
	PlannedDay      *string  `json:"plannedDay,omitempty" jsonschema:"Планируемый срок, yyyy-MM-dd"`
	DeadlineDay     *string  `json:"deadlineDay,omitempty" jsonschema:"Дедлайн, yyyy-MM-dd"`
	DurationMinutes *int32   `json:"durationMinutes,omitempty" jsonschema:"Планируемое время в минутах"`
	ManualOrder     *float64 `json:"manualOrder,omitempty" jsonschema:"Ручной порядок, нет — в конец (по времени создания)"`
	ExternalID      *string  `json:"externalID,omitempty" jsonschema:"Идентификатор напоминания при импорте"`
}

// CreatedOut — результат создания записи
type CreatedOut struct {
	ID string `json:"id"` // Идентификатор созданной записи
}

// CreateTask создание записи
func (t *TaskTools) CreateTask(ctx context.Context, _ *sdk.CallToolRequest, in CreateTaskIn) (*sdk.CallToolResult, CreatedOut, error) {
	var out CreatedOut

	// Собираем запрос из аргументов и контекста
	necessary, err := mcpTransport.Necessary(ctx)
	if err != nil {
		return nil, out, err
	}
	id, err := mcpTransport.NewOrParseUUID("id", in.ID)
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
	now := time.Now()
	req := model.CreateTaskReq{
		Necessary:       necessary,
		ID:              id,
		ListID:          uuid.Nil,
		SectionID:       sectionID,
		Title:           in.Title,
		Notes:           in.Notes,
		PriorityID:      priorityID,
		PlannedDay:      in.PlannedDay,
		DeadlineDay:     in.DeadlineDay,
		DurationMinutes: in.DurationMinutes,
		ManualOrder:     float64(now.UnixMilli()),
		ExternalID:      in.ExternalID,
		CreatedAt:       now,
	}
	if listID != nil {
		req.ListID = *listID
	}
	if in.ManualOrder != nil {
		req.ManualOrder = *in.ManualOrder
	}

	// Валидируем запрос
	if err = validate.Validate(ctx, req); err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Вызываем метод сервиса
	if err = t.taskService.CreateTask(ctx, req); err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}
	out.ID = id.String()
	return nil, out, nil
}
