package model

// EnsureUserReq — пользователь, который должен существовать после запуска сервера
type EnsureUserReq struct {
	Name     string `validate:"required"` // Имя пользователя
	Email    string `validate:"required"` // Электронная почта пользователя
	Password string `validate:"required"` // Пароль пользователя
}
