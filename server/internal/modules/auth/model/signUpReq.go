package model

import (
	userModel "server/internal/modules/user/model"
	"server/internal/utils/errors"

	"github.com/jedi-tasks/go-server-grpc/proto"
)

type SignUpReq struct {
	Name        string                           `json:"name" validate:"required"`                 // Имя пользователя
	Email       string                           `json:"email" validate:"required" format:"email"` // Электронная почта пользователя
	Password    string                           `json:"password" validate:"required"`             // Пароль пользователя
	Application userModel.ApplicationInformation `json:"application"`                              // Информация о приложении
	Device      userModel.DeviceInformation      `json:"device"`                                   // Информация о девайсе
	DeviceID    string                           `json:"-" validate:"required"`                    // Идентификатор устройства
}

// ProtoSignUpReq wrapper for proto request
type ProtoSignUpReq struct {
	*proto.SignUpRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoSignUpReq) ConvertToModel() (res SignUpReq, err error) {
	if p.SignUpRequest == nil {
		return res, errors.BadRequest.New("SignUpRequest is required")
	}

	// Convert application information
	application, err := userModel.ProtoApplicationInformation{ApplicationInformation: p.Application}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Convert device information
	device, err := userModel.ProtoDeviceInformation{DeviceInformation: p.Device}.ConvertToModel()
	if err != nil {
		return res, err
	}

	return SignUpReq{
		Name:        p.Name,
		Email:       p.Email,
		Password:    p.Password,
		Application: application,
		Device:      device,
		DeviceID:    p.DeviceID,
	}, nil
}
