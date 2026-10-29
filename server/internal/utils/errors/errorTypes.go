package errors

import (
	"net/http"
	pkgErrors "pkg/errors"
)

var InternalServer = ErrorType{
	ErrorType: pkgErrors.ErrorType{
		HTTPCode:  http.StatusInternalServerError,
		LogAs:     pkgErrors.LogAsError,
		HumanText: "Произошла непредвиденная ошибка",
	},
	Category: CategoryInternal,
}

var BadRequest = ErrorType{
	ErrorType: pkgErrors.ErrorType{
		HTTPCode:  http.StatusBadRequest,
		LogAs:     pkgErrors.LogAsWarning,
		HumanText: "Введены неверные данные",
	},
	Category: CategoryOther,
}

var NeedToLogout = ErrorType{
	ErrorType: pkgErrors.ErrorType{
		HTTPCode:  http.StatusUnauthorized,
		LogAs:     pkgErrors.LogAsWarning,
		HumanText: "Пользователь не авторизован",
	},
	Category: CategoryNeedToLogout,
}

var Forbidden = ErrorType{
	ErrorType: pkgErrors.ErrorType{
		HTTPCode:  http.StatusForbidden,
		LogAs:     pkgErrors.LogAsWarning,
		HumanText: "Доступ запрещен",
	},
	Category: CategoryOther,
}

var BadGateway = ErrorType{
	ErrorType: pkgErrors.ErrorType{
		HTTPCode:  http.StatusBadGateway,
		LogAs:     pkgErrors.LogAsWarning,
		HumanText: "Произошла ошибка на сервере внешнего сервиса",
	},
	Category: CategoryInternal,
}

var NotFound = ErrorType{
	ErrorType: pkgErrors.ErrorType{
		HTTPCode:  http.StatusNotFound,
		LogAs:     pkgErrors.LogAsWarning,
		HumanText: "Данные не найдены",
	},
	Category: CategoryOther,
}

var NeedToSync = ErrorType{
	ErrorType: pkgErrors.ErrorType{
		HTTPCode:  http.StatusConflict,
		LogAs:     pkgErrors.LogAsWarning,
		HumanText: "Требуется синхронизация устройства",
	},
	Category: CategoryNeedToSync,
}

var NeedToRefreshToken = ErrorType{
	ErrorType: pkgErrors.ErrorType{
		HTTPCode:  http.StatusUnauthorized,
		LogAs:     pkgErrors.LogAsWarning,
		HumanText: "Требуется обновить access-токен",
	},
	Category: CategoryNeedToRefreshToken,
}

var ContextCancelled = ErrorType{
	ErrorType: pkgErrors.ErrorType{
		HTTPCode:  http.StatusTeapot,
		LogAs:     pkgErrors.LogAsWarning,
		HumanText: "Вышел таймаут запроса",
	},
	Category: CategoryOther,
}
