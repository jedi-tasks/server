package model

import (
	"server/internal/enum/osType"
	userModel "server/internal/modules/user/model"
	"server/internal/utils/errors"

	"github.com/jedi-tasks/go-server-grpc/proto"
)

type SignInReq struct {
	Email       string                           `json:"email" validate:"required" format:"email"` // Электронная почта пользователя
	Password    string                           `json:"password" validate:"required"`             // Пароль пользователя
	Application userModel.ApplicationInformation `json:"application"`                              // Информация о приложении
	Device      userModel.DeviceInformation      `json:"device"`                                   // Информация о девайсе
	DeviceID    string                           `json:"-" validate:"required"`                    // Идентификатор устройства
}

// ProtoSignInReq wrapper for proto request
type ProtoSignInReq struct {
	*proto.SignInRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoSignInReq) ConvertToModel() (res SignInReq, err error) {
	if p.SignInRequest == nil {
		return res, errors.BadRequest.New("SignInRequest is required")
	}

	// Convert application information
	var application userModel.ApplicationInformation
	if p.Application != nil {
		application = userModel.ApplicationInformation{
			BundleID: p.Application.BundleID,
			Version:  p.Application.Version,
			Build:    p.Application.Build,
		}
	}

	// Convert device information
	var device userModel.DeviceInformation
	if p.Device != nil {

		nameOS, err := osType.ProtoOSType{OSType: p.Device.NameOS}.ConvertToModel()
		if err != nil {
			return res, err
		}

		device = userModel.DeviceInformation{
			NameOS:     nameOS,
			VersionOS:  p.Device.VersionOS,
			DeviceName: p.Device.DeviceName,
			ModelName:  p.Device.ModelName,
			IPAddress:  p.Device.IpAddress,
			UserAgent:  "",
			Stateless:  p.Device.Stateless,
		}
	}

	return SignInReq{
		Email:       p.Email,
		Password:    p.Password,
		Application: application,
		Device:      device,
		DeviceID:    p.DeviceID,
	}, nil
}
