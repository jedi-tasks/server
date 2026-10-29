// Package mcp — инструменты MCP для тегов, один к одному с методами TagEndpoint
package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/tag/model"
	tagService "server/internal/modules/tag/service"
)

var _ TagService = new(tagService.TagService)

// TagService — интерфейс бизнес-логики тегов
type TagService interface {
	GetTags(context.Context, model.GetTagsReq) ([]model.Tag, error)
	CreateTag(context.Context, model.CreateTagReq) error
	GetTasksToTags(context.Context, model.GetTasksToTagsReq) ([]model.TaskToTag, error)
	AddTagToTask(context.Context, model.LinkTagReq) error
	RemoveTagFromTask(context.Context, model.LinkTagReq) error
}

// TagTools — обработчики инструментов тегов
type TagTools struct {
	tagService TagService
}

// Register добавляет инструменты тегов на MCP-сервер
func Register(server *sdk.Server, tagService TagService) {
	tools := &TagTools{tagService: tagService}
	readOnly := &sdk.ToolAnnotations{ReadOnlyHint: true} //nolint:exhaustruct

	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "get_tags",
		Description: "Возвращает теги пользователя",
		Annotations: readOnly,
	}, tools.GetTags)
	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "create_tag",
		Description: "Создает тег",
	}, tools.CreateTag)
	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "get_tasks_to_tags",
		Description: "Возвращает живые связи задач и тегов. Фильтр по задачам; пустой — все связи",
		Annotations: readOnly,
	}, tools.GetTasksToTags)
	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "add_tag_to_task",
		Description: "Привязывает тег к задаче",
	}, tools.AddTagToTask)
	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "remove_tag_from_task",
		Description: "Снимает тег с задачи",
	}, tools.RemoveTagFromTask)
}
