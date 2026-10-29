package mcp

import (
	"context"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/task/model"
	mcpTransport "server/internal/utils/mcp"
)

// SetTaskCompletedIn — аргументы set_task_completed
type SetTaskCompletedIn struct {
	ID        string `json:"id" jsonschema:"Идентификатор задачи"`
	Completed bool   `json:"completed" jsonschema:"true — выполнить, false — вернуть в работу"`
}

// SetTaskCompleted выполнение задачи или возврат в работу
func (t *TaskTools) SetTaskCompleted(ctx context.Context, _ *sdk.CallToolRequest, in SetTaskCompletedIn) (*sdk.CallToolResult, EmptyOut, error) {
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

	// Вызываем метод сервиса
	err = t.taskService.SetTaskCompleted(ctx, model.SetTaskCompletedReq{Necessary: necessary, ID: id, Completed: in.Completed, UpdatedAt: time.Now()})
	if err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}
	return nil, out, nil
}
