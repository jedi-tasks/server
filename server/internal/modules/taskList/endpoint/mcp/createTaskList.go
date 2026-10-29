package mcp

import (
	"context"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/enum/sortMode"
	"server/internal/modules/taskList/model"
	mcpTransport "server/internal/utils/mcp"
	"server/internal/utils/validate"
)

// CreateTaskListIn — аргументы create_task_list
type CreateTaskListIn struct {
	ID        *string  `json:"id,omitempty" jsonschema:"Идентификатор списка, нет — сервер создаст"`
	Name      string   `json:"name" jsonschema:"Название списка"`
	ColorHex  string   `json:"colorHex,omitempty" jsonschema:"Цвет списка, например #3E9BF0"`
	Icon      string   `json:"icon,omitempty" jsonschema:"SF Symbol иконки списка, например list.bullet"`
	SortOrder *float64 `json:"sortOrder,omitempty" jsonschema:"Порядок списка, нет — в конец"`
	SortMode  string   `json:"sortMode,omitempty" jsonschema:"Режим сортировки задач: manual, priority, plannedDay, deadline, duration; нет — manual"`
}

// CreatedOut — результат создания записи
type CreatedOut struct {
	ID string `json:"id"` // Идентификатор созданной записи
}

// CreateTaskList создание записи
func (t *TaskListTools) CreateTaskList(ctx context.Context, _ *sdk.CallToolRequest, in CreateTaskListIn) (*sdk.CallToolResult, CreatedOut, error) {
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
	now := time.Now()
	req := model.CreateTaskListReq{
		Necessary: necessary,
		ID:        id,
		Name:      in.Name,
		ColorHex:  in.ColorHex,
		Icon:      in.Icon,
		SortOrder: float64(now.UnixMilli()),
		SortMode:  sortMode.Manual,
		CreatedAt: now,
	}
	if in.ColorHex == "" {
		req.ColorHex = "#3E9BF0"
	}
	if in.Icon == "" {
		req.Icon = "list.bullet"
	}
	if in.SortOrder != nil {
		req.SortOrder = *in.SortOrder
	}
	if in.SortMode != "" {
		req.SortMode = sortMode.SortMode(in.SortMode)
	}

	// Валидируем запрос
	if err = validate.Validate(ctx, req); err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Вызываем метод сервиса
	if err = t.taskListService.CreateTaskList(ctx, req); err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}
	out.ID = id.String()
	return nil, out, nil
}
