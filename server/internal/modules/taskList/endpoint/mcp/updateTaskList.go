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

// UpdateTaskListIn — аргументы update_task_list, меняются только переданные поля
type UpdateTaskListIn struct {
	ID        string   `json:"id" jsonschema:"Идентификатор списка"`
	Name      *string  `json:"name,omitempty" jsonschema:"Название списка"`
	ColorHex  *string  `json:"colorHex,omitempty" jsonschema:"Цвет списка"`
	Icon      *string  `json:"icon,omitempty" jsonschema:"SF Symbol иконки списка"`
	SortOrder *float64 `json:"sortOrder,omitempty" jsonschema:"Порядок списка"`
	SortMode  *string  `json:"sortMode,omitempty" jsonschema:"Режим сортировки задач: manual, priority, plannedDay, deadline, duration"`
}

// EmptyOut — результат правки без данных
type EmptyOut struct{}

// UpdateTaskList изменение записи
func (t *TaskListTools) UpdateTaskList(ctx context.Context, _ *sdk.CallToolRequest, in UpdateTaskListIn) (*sdk.CallToolResult, EmptyOut, error) {
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
	var mode *sortMode.SortMode
	if in.SortMode != nil {
		parsed := sortMode.SortMode(*in.SortMode)
		mode = &parsed
	}
	req := model.UpdateTaskListReq{
		Necessary: necessary,
		ID:        id,
		Name:      in.Name,
		ColorHex:  in.ColorHex,
		Icon:      in.Icon,
		SortOrder: in.SortOrder,
		SortMode:  mode,
		UpdatedAt: time.Now(),
	}

	// Валидируем запрос
	if err = validate.Validate(ctx, req); err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Вызываем метод сервиса
	if err = t.taskListService.UpdateTaskList(ctx, req); err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}
	return nil, out, nil
}
