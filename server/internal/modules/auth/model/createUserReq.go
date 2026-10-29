package model

// CreateUserReq — данные для создания пользователя вместе с данными по умолчанию
type CreateUserReq struct {
	Name     string // Имя пользователя
	Email    string // Электронная почта пользователя
	Password string // Пароль пользователя
	DeviceID string // Устройство, с которого создается пользователь (для аудит-лога)
}
