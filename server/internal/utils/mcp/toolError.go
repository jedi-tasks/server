package mcp

import (
	"fmt"

	pkgErrors "pkg/errors"
)

// ToolError превращает ошибку сервиса в текст для агента: человеческий текст типа ошибки и
// причина. Агенту важна причина — по ней он решает, что делать дальше (например, снять блокер
// перед удалением задачи)
func ToolError(err error) error {
	castErr := pkgErrors.CastError(err)
	return fmt.Errorf("%s: %s", castErr.ErrorType.HumanText, castErr.Err.Error())
}
