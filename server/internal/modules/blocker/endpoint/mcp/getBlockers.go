package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/blocker/model"
	mcpTransport "server/internal/utils/mcp"
)

// GetBlockersIn — аргументы get_blockers
type GetBlockersIn struct {
	IDs     []string `json:"ids,omitempty" jsonschema:"Фильтр по идентификаторам блокеров, пустой — все"`
	TaskIDs []string `json:"taskIDs,omitempty" jsonschema:"Фильтр по заблокированным задачам, пустой — все"`
}

// GetBlockersOut — результат get_blockers
type GetBlockersOut struct {
	Blockers []BlockerView `json:"blockers"` // Блокеры
}

// GetBlockers получение записей
func (t *BlockerTools) GetBlockers(ctx context.Context, _ *sdk.CallToolRequest, in GetBlockersIn) (*sdk.CallToolResult, GetBlockersOut, error) {
	var out GetBlockersOut

	// Собираем запрос из аргументов и контекста
	necessary, err := mcpTransport.Necessary(ctx)
	if err != nil {
		return nil, out, err
	}
	ids, err := mcpTransport.ParseUUIDs("ids", in.IDs)
	if err != nil {
		return nil, out, err
	}
	taskIDs, err := mcpTransport.ParseUUIDs("taskIDs", in.TaskIDs)
	if err != nil {
		return nil, out, err
	}

	// Вызываем метод сервиса
	blockers, err := t.blockerService.GetBlockers(ctx, model.GetBlockersReq{Necessary: necessary, IDs: ids, TaskIDs: taskIDs}) //nolint:exhaustruct
	if err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Конвертируем ответ
	out.Blockers = make([]BlockerView, 0, len(blockers))
	for _, blocker := range blockers {
		out.Blockers = append(out.Blockers, NewBlockerView(blocker))
	}
	return nil, out, nil
}
