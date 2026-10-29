// Package validate — проверка запросов: обязательные поля по тегам (pkg/validator) и собственная
// проверка модели Validate(ctx), если она есть. pkg/validator ждет Validate() без контекста и
// поэтому модели сервера с Validate(ctx) не видел — их проверки не выполнялись
package validate

import (
	"context"

	"pkg/validator"
)

// contextValidator — модель с собственной проверкой, зависящей от контекста
type contextValidator interface {
	Validate(ctx context.Context) error
}

// Validate проверяет обязательные поля и вызывает Validate(ctx) модели
func Validate(ctx context.Context, request any) error {
	if err := validator.Validate(request); err != nil {
		return err
	}
	if v, ok := request.(contextValidator); ok {
		return v.Validate(ctx)
	}
	return nil
}
