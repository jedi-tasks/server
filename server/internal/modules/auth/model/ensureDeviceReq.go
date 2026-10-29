package model

import (
	"github.com/google/uuid"

	userModel "server/internal/modules/user/model"
)

// EnsureDeviceReq — устройство пользователя, которое должно существовать без входа по паролю
type EnsureDeviceReq struct {
	UserID      uuid.UUID                        `validate:"required"` // Идентификатор пользователя
	DeviceID    string                           `validate:"required"` // Идентификатор устройства
	Device      userModel.DeviceInformation      // Информация об устройстве
	Application userModel.ApplicationInformation // Информация о приложении
}
