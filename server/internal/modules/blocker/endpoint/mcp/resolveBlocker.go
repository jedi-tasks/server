package mcp

import (
	"context"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/blocker/model"
	mcpTransport "server/internal/utils/mcp"
)

// ResolveBlockerIn — аргументы resolve_blocker
type ResolveBlockerIn struct {
	ID string `json:"id" jsonschema:"Идентификатор блокера"`
}

// EmptyOut — результат правки без данных
type EmptyOut struct{}

// ResolveBlocker снимает блокер
func (t *BlockerTools) ResolveBlocker(ctx context.Context, _ *sdk.CallToolRequest, in ResolveBlockerIn) (*sdk.CallToolResult, EmptyOut, error) {
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
	if err = t.blockerService.ResolveBlocker(ctx, model.ResolveBlockerReq{Necessary: necessary, ID: id, UpdatedAt: time.Now()}); err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}
	return nil, out, nil
}
