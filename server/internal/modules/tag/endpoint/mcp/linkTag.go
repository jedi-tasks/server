package mcp

import (
	"context"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/tag/model"
	mcpTransport "server/internal/utils/mcp"
)

// LinkTagIn — аргументы add_tag_to_task и remove_tag_from_task
type LinkTagIn struct {
	TaskID string `json:"taskID" jsonschema:"Идентификатор задачи"`
	TagID  string `json:"tagID" jsonschema:"Идентификатор тега"`
}

// EmptyOut — результат правки без данных
type EmptyOut struct{}

// AddTagToTask привязывает тег к задаче
func (t *TagTools) AddTagToTask(ctx context.Context, _ *sdk.CallToolRequest, in LinkTagIn) (*sdk.CallToolResult, EmptyOut, error) {
	return t.link(ctx, in, t.tagService.AddTagToTask)
}

// RemoveTagFromTask снимает тег с задачи
func (t *TagTools) RemoveTagFromTask(ctx context.Context, _ *sdk.CallToolRequest, in LinkTagIn) (*sdk.CallToolResult, EmptyOut, error) {
	return t.link(ctx, in, t.tagService.RemoveTagFromTask)
}

// link собирает запрос связи и вызывает переданный метод сервиса
func (t *TagTools) link(ctx context.Context, in LinkTagIn, call func(context.Context, model.LinkTagReq) error) (*sdk.CallToolResult, EmptyOut, error) {
	var out EmptyOut

	// Собираем запрос из аргументов и контекста
	necessary, err := mcpTransport.Necessary(ctx)
	if err != nil {
		return nil, out, err
	}
	taskID, err := mcpTransport.ParseUUID("taskID", in.TaskID)
	if err != nil {
		return nil, out, err
	}
	tagID, err := mcpTransport.ParseUUID("tagID", in.TagID)
	if err != nil {
		return nil, out, err
	}

	// Вызываем метод сервиса
	if err = call(ctx, model.LinkTagReq{Necessary: necessary, TaskID: taskID, TagID: tagID, UpdatedAt: time.Now()}); err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}
	return nil, out, nil
}
