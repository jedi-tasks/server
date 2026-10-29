package errors

import (
	pkgErrors "pkg/errors"
)

// categoryByErrorType содержит соответствие между базовыми типами ошибок и их категориями
var categoryByErrorType = map[pkgErrors.ErrorType]Category{
	InternalServer.ErrorType:     InternalServer.Category,
	BadRequest.ErrorType:         BadRequest.Category,
	NeedToLogout.ErrorType:       NeedToLogout.Category,
	Forbidden.ErrorType:          Forbidden.Category,
	BadGateway.ErrorType:         BadGateway.Category,
	NotFound.ErrorType:           NotFound.Category,
	NeedToSync.ErrorType:         NeedToSync.Category,
	NeedToRefreshToken.ErrorType: NeedToRefreshToken.Category,
	ContextCancelled.ErrorType:   ContextCancelled.Category,
	pkgErrors.Default:            CategoryInternal,
}

// GetCategory возвращает категорию ошибки по ее базовому типу
func GetCategory(errorType pkgErrors.ErrorType) Category {

	// Ищем категорию по типу ошибки
	category, ok := categoryByErrorType[errorType]
	if !ok {
		return CategoryInternal
	}

	return category
}
