package errors

import (
	pkgErrors "pkg/errors"
)

// Category - Категория ошибки, по которой клиент определяет дополнительную логику обработки
type Category int32

const (
	CategoryUnspecified        Category = iota // Категория не определена
	CategoryInternal                            // Внутренняя ошибка сервера
	CategoryNeedToLogout                        // Требуется разлогинить пользователя
	CategoryNeedToSync                          // Требуется синхронизация устройства
	CategoryOther                               // Прочие ошибки
	CategoryNeedToRefreshToken                  // Требуется обновить access-токен
)

// ErrorType - тип ошибки сервиса, расширяющий базовый тип ошибки категорией для клиента
type ErrorType struct {
	pkgErrors.ErrorType
	Category Category // Категория ошибки для клиента
}
