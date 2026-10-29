package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/tag/model"
	mcpTransport "server/internal/utils/mcp"
)

// GetTagsIn — аргументы get_tags
type GetTagsIn struct {
	IDs []string `json:"ids,omitempty" jsonschema:"Фильтр по идентификаторам, пустой — все"`
}

// GetTagsOut — результат get_tags
type GetTagsOut struct {
	Tags []TagView `json:"tags"` // Теги
}

// GetTags получение записей
func (t *TagTools) GetTags(ctx context.Context, _ *sdk.CallToolRequest, in GetTagsIn) (*sdk.CallToolResult, GetTagsOut, error) {
	var out GetTagsOut

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
	tags, err := t.tagService.GetTags(ctx, model.GetTagsReq{Necessary: necessary, IDs: ids})
	if err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Конвертируем ответ
	out.Tags = make([]TagView, 0, len(tags))
	for _, tag := range tags {
		out.Tags = append(out.Tags, NewTagView(tag))
	}
	return nil, out, nil
}
