// Package mcp — инструменты MCP для блокеров, один к одному с методами BlockerEndpoint
package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/blocker/model"
	blockerService "server/internal/modules/blocker/service"
)

var _ BlockerService = new(blockerService.BlockerService)

// BlockerService — интерфейс бизнес-логики блокеров
type BlockerService interface {
	GetBlockers(context.Context, model.GetBlockersReq) ([]model.Blocker, error)
	CreateBlocker(context.Context, model.CreateBlockerReq) error
	ResolveBlocker(context.Context, model.ResolveBlockerReq) error
}

// BlockerTools — обработчики инструментов блокеров
type BlockerTools struct {
	blockerService BlockerService
}

// Register добавляет инструменты блокеров на MCP-сервер
func Register(server *sdk.Server, blockerService BlockerService) {
	tools := &BlockerTools{blockerService: blockerService}

	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "get_blockers",
		Description: "Возвращает блокеры задач, включая снятые (resolvedAt). Блокер — ссылка на блокирующую задачу либо текстовая причина",
		Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}, //nolint:exhaustruct
	}, tools.GetBlockers)
	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "create_blocker",
		Description: "Блокирует задачу: ровно одно из blockingTaskID (другая задача) и reason (текст). Циклы блокировок отклоняются",
	}, tools.CreateBlocker)
	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "resolve_blocker",
		Description: "Снимает блокер: задача снова не заблокирована, запись о блокере остается",
	}, tools.ResolveBlocker)
}
