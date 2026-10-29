package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/taskList/model"
	mcpTransport "server/internal/utils/mcp"
)

// GetTaskListsIn — аргументы get_task_lists
type GetTaskListsIn struct {
	IDs []string `json:"ids,omitempty" jsonschema:"Фильтр по идентификаторам, пустой — все"`
}

// GetTaskListsOut — результат get_task_lists
type GetTaskListsOut struct {
	TaskLists []TaskListView `json:"taskLists"` // Списки
}

// GetTaskLists получение записей
func (t *TaskListTools) GetTaskLists(ctx context.Context, _ *sdk.CallToolRequest, in GetTaskListsIn) (*sdk.CallToolResult, GetTaskListsOut, error) {
	var out GetTaskListsOut

	// Собираем запрос из аргументов и контекста
	necessary, err := mcpTransport.Necessary(ctx)
	if err != nil {
		return nil, out, err
	}
	ids, err := mcpTransport.ParseUUIDs("ids", in.IDs)
	if err != nil {
		return nil, out, err
	}

	// Вызываем метод сервиса
	lists, err := t.taskListService.GetTaskLists(ctx, model.GetTaskListsReq{Necessary: necessary, IDs: ids})
	if err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Конвертируем ответ
	out.TaskLists = make([]TaskListView, 0, len(lists))
	for _, list := range lists {
		out.TaskLists = append(out.TaskLists, NewTaskListView(list))
	}
	return nil, out, nil
}
