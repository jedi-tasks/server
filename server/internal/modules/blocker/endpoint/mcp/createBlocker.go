package mcp

import (
	"context"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/blocker/model"
	mcpTransport "server/internal/utils/mcp"
	"server/internal/utils/validate"
)

// CreateBlockerIn — аргументы create_blocker
type CreateBlockerIn struct {
	ID             *string `json:"id,omitempty" jsonschema:"Идентификатор блокера, нет — сервер создаст"`
	TaskID         string  `json:"taskID" jsonschema:"Заблокированная задача"`
	BlockingTaskID *string `json:"blockingTaskID,omitempty" jsonschema:"Блокирующая задача; либо это поле, либо reason"`
	Reason         *string `json:"reason,omitempty" jsonschema:"Текстовая причина; либо это поле, либо blockingTaskID"`
	SinceDay       string  `json:"sinceDay,omitempty" jsonschema:"Дата начала yyyy-MM-dd, нет — сегодня"`
}

// CreatedOut — результат создания записи
type CreatedOut struct {
	ID string `json:"id"` // Идентификатор созданной записи
}

// CreateBlocker создание записи
func (t *BlockerTools) CreateBlocker(ctx context.Context, _ *sdk.CallToolRequest, in CreateBlockerIn) (*sdk.CallToolResult, CreatedOut, error) {
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
	taskID, err := mcpTransport.ParseUUID("taskID", in.TaskID)
	if err != nil {
		return nil, out, err
	}
	blockingTaskID, err := mcpTransport.ParseOptionalUUID("blockingTaskID", in.BlockingTaskID)
	if err != nil {
		return nil, out, err
	}
	now := time.Now()
	req := model.CreateBlockerReq{
		Necessary:      necessary,
		ID:             id,
		TaskID:         taskID,
		BlockingTaskID: blockingTaskID,
		Reason:         in.Reason,
		SinceDay:       in.SinceDay,
		UpdatedAt:      now,
	}
	if req.SinceDay == "" {
		req.SinceDay = now.Format(time.DateOnly)
	}

	// Валидируем запрос
	if err = validate.Validate(ctx, req); err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Вызываем метод сервиса
	if err = t.blockerService.CreateBlocker(ctx, req); err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}
	out.ID = id.String()
	return nil, out, nil
}
