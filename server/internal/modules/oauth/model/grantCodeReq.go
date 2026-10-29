package model

// GrantCodeReq — выдача кода после отправки формы входа
type GrantCodeReq struct {
	Authorize AuthorizeReq // Проверенные параметры авторизации
	Email     string       // Почта пользователя
	Password  string       // Пароль пользователя
}
