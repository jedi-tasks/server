package mcp

import (
	"context"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/tag/model"
	mcpTransport "server/internal/utils/mcp"
	"server/internal/utils/validate"
)

// CreateTagIn — аргументы create_tag
type CreateTagIn struct {
	ID   *string `json:"id,omitempty" jsonschema:"Идентификатор тега, нет — сервер создаст"`
	Name string  `json:"name" jsonschema:"Название тега"`
}

// CreatedOut — результат создания записи
type CreatedOut struct {
	ID string `json:"id"` // Идентификатор созданной записи
}

// CreateTag создание записи
func (t *TagTools) CreateTag(ctx context.Context, _ *sdk.CallToolRequest, in CreateTagIn) (*sdk.CallToolResult, CreatedOut, error) {
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
	req := model.CreateTagReq{Necessary: necessary, ID: id, Name: in.Name, UpdatedAt: time.Now()}

	// Валидируем запрос
	if err = validate.Validate(ctx, req); err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Вызываем метод сервиса
	if err = t.tagService.CreateTag(ctx, req); err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}
	out.ID = id.String()
	return nil, out, nil
}
