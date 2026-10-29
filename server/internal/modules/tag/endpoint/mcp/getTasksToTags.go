package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/tag/model"
	mcpTransport "server/internal/utils/mcp"
)

// GetTasksToTagsIn — аргументы get_tasks_to_tags
type GetTasksToTagsIn struct {
	TaskIDs []string `json:"taskIDs,omitempty" jsonschema:"Фильтр по задачам, пустой — все"`
}

// GetTasksToTagsOut — результат get_tasks_to_tags
type GetTasksToTagsOut struct {
	TasksToTags []TaskToTagView `json:"tasksToTags"` // Связи задач и тегов
}

// GetTasksToTags получение связей задач и тегов
func (t *TagTools) GetTasksToTags(ctx context.Context, _ *sdk.CallToolRequest, in GetTasksToTagsIn) (*sdk.CallToolResult, GetTasksToTagsOut, error) {
	var out GetTasksToTagsOut

	// Собираем запрос из аргументов и контекста
	necessary, err := mcpTransport.Necessary(ctx)
	if err != nil {
		return nil, out, err
	}
	taskIDs, err := mcpTransport.ParseUUIDs("taskIDs", in.TaskIDs)
	if err != nil {
		return nil, out, err
	}

	// Вызываем метод сервиса
	links, err := t.tagService.GetTasksToTags(ctx, model.GetTasksToTagsReq{Necessary: necessary, TaskIDs: taskIDs}) //nolint:exhaustruct
	if err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Конвертируем ответ
	out.TasksToTags = make([]TaskToTagView, 0, len(links))
	for _, link := range links {
		out.TasksToTags = append(out.TasksToTags, NewTaskToTagView(link))
	}
	return nil, out, nil
}
