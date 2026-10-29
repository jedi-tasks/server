package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/task/model"
	mcpTransport "server/internal/utils/mcp"
)

// GetTasksIn — аргументы get_tasks
type GetTasksIn struct {
	ListIDs   []string `json:"listIDs,omitempty" jsonschema:"Фильтр по спискам, пустой — все"`
	IDs       []string `json:"ids,omitempty" jsonschema:"Фильтр по идентификаторам задач, пустой — все"`
	Completed *bool    `json:"completed,omitempty" jsonschema:"Фильтр по выполненности, нет — все"`
}

// GetTasksOut — результат get_tasks
type GetTasksOut struct {
	Tasks []TaskView `json:"tasks"` // Задачи
}

// GetTasks получение записей
func (t *TaskTools) GetTasks(ctx context.Context, _ *sdk.CallToolRequest, in GetTasksIn) (*sdk.CallToolResult, GetTasksOut, error) {
	var out GetTasksOut

	// Собираем запрос из аргументов и контекста
	necessary, err := mcpTransport.Necessary(ctx)
	if err != nil {
		return nil, out, err
	}
	listIDs, err := mcpTransport.ParseUUIDs("listIDs", in.ListIDs)
	if err != nil {
		return nil, out, err
	}
	ids, err := mcpTransport.ParseUUIDs("ids", in.IDs)
	if err != nil {
		return nil, out, err
	}

	// Вызываем метод сервиса
	tasks, err := t.taskService.GetTasks(ctx, model.GetTasksReq{Necessary: necessary, IDs: ids, ListIDs: listIDs, Completed: in.Completed})
	if err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Конвертируем ответ
	out.Tasks = make([]TaskView, 0, len(tasks))
	for _, task := range tasks {
		out.Tasks = append(out.Tasks, NewTaskView(task))
	}
	return nil, out, nil
}
