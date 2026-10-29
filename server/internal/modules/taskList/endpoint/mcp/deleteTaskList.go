package mcp

import (
	"context"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/taskList/model"
	mcpTransport "server/internal/utils/mcp"
)

// DeleteIn — аргументы удаления записи
type DeleteIn struct {
	ID string `json:"id" jsonschema:"Идентификатор записи"`
}

// DeleteTaskList удаление записи
func (t *TaskListTools) DeleteTaskList(ctx context.Context, _ *sdk.CallToolRequest, in DeleteIn) (*sdk.CallToolResult, EmptyOut, error) {
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
	if err = t.taskListService.DeleteTaskList(ctx, model.DeleteTaskListReq{Necessary: necessary, ID: id, UpdatedAt: time.Now()}); err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}
	return nil, out, nil
}
