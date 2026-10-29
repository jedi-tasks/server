package model

import (
	"github.com/google/uuid"

	userModel "server/internal/modules/user/model"
)

// IssueTokensReq — выдача пары токенов уже проверенному пользователю для устройства
type IssueTokensReq struct {
	UserID      uuid.UUID                        `validate:"required"` // Пользователь
	DeviceID    string                           `validate:"required"` // Идентификатор устройства
	Device      userModel.DeviceInformation      // Информация об устройстве
	Application userModel.ApplicationInformation // Информация о приложении
}
