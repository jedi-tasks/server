package mcp

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	pkgErrors "pkg/errors"
	"pkg/log"
)

// instructions — что агенту нужно знать об API. Только семантика полей, без правил работы с
// задачами: они принадлежат пользователю и приходят из его инструкций
const instructions = `Трекер задач Tasks. Инструменты повторяют методы серверного API один к одному.
Идентификаторы — UUID строками, новые идентификаторы создает агент. Даты дней — yyyy-MM-dd, время — RFC 3339.
Сервер ничего не меняет скрыто: удаление записи, на которую ссылаются живые строки, отклоняется с текстом, что мешает.
Правка принимается, только если она новее сохраненной (последняя запись побеждает).`

// NewServer создает MCP-сервер без инструментов: их добавляют модули
func NewServer(version string) *sdk.Server {
	server := sdk.NewServer(
		&sdk.Implementation{Name: "tasks", Title: "Tasks", Version: version}, //nolint:exhaustruct
		&sdk.ServerOptions{Instructions: instructions},                       //nolint:exhaustruct
	)
	server.AddReceivingMiddleware(logCalls)
	return server
}

// logCalls пишет каждый вызов инструмента: имя, длительность, ошибка. Аргументы не пишем — в них
// заметки задач пользователя
func logCalls(next sdk.MethodHandler) sdk.MethodHandler {
	return func(ctx context.Context, method string, req sdk.Request) (sdk.Result, error) {
		started := time.Now()
		result, err := next(ctx, method, req)

		name := method
		if params, ok := req.GetParams().(*sdk.CallToolParamsRaw); ok {
			name = fmt.Sprintf("%s %s", method, params.Name)
		}
		duration := time.Since(started).Round(time.Millisecond)

		switch {
		case err != nil:
			log.Warning(fmt.Sprintf("MCP %s за %s: %s", name, duration, pkgErrors.CastError(err).Error()))
		case isToolError(result):
			log.Warning(fmt.Sprintf("MCP %s за %s: инструмент вернул ошибку", name, duration))
		default:
			log.Info(fmt.Sprintf("MCP %s за %s", name, duration))
		}
		return result, err
	}
}

// isToolError — результат вызова инструмента с флагом ошибки
func isToolError(result sdk.Result) bool {
	callResult, ok := result.(*sdk.CallToolResult)
	return ok && callResult.IsError
}
