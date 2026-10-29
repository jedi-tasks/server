package model

import (
	userModel "server/internal/modules/user/model"
	"server/internal/utils/errors"

	"github.com/jedi-tasks/go-server-grpc/proto"
)

type RefreshTokensReq struct {
	Token       string                           `json:"token" validate:"required"` // Токен восстановления доступа
	Application userModel.ApplicationInformation `json:"application"`               // Информация о приложении
	Device      userModel.DeviceInformation      `json:"device"`                    // Информация о девайсе
}

func (req *RefreshTokensReq) ConvertToProto() (res proto.RefreshTokensRequest, err error) {

	nameOS, err := req.Device.NameOS.ConvertToProto()
	if err != nil {
		return res, err
	}

	return proto.RefreshTokensRequest{
		Token: req.Token,
		Application: &proto.ApplicationInformation{
			BundleID: req.Application.BundleID,
			Version:  req.Application.Version,
			Build:    req.Application.Build,
		},
		Device: &proto.DeviceInformation{
			NameOS:     nameOS,
			VersionOS:  req.Device.VersionOS,
			DeviceName: req.Device.DeviceName,
			ModelName:  req.Device.ModelName,
			IpAddress:  req.Device.IPAddress,
		},
	}, nil
}

// ProtoRefreshTokensReq wrapper for proto request
type ProtoRefreshTokensReq struct {
	*proto.RefreshTokensRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoRefreshTokensReq) ConvertToModel() (RefreshTokensReq, error) {
	if p.RefreshTokensRequest == nil {
		return RefreshTokensReq{}, errors.BadRequest.New("RefreshTokensRequest is required")
	}

	// Convert application information
	application, err := userModel.ProtoApplicationInformation{ApplicationInformation: p.Application}.ConvertToModel()
	if err != nil {
		return RefreshTokensReq{}, err
	}

	// Convert device information
	device, err := userModel.ProtoDeviceInformation{DeviceInformation: p.Device}.ConvertToModel()
	if err != nil {
		return RefreshTokensReq{}, err
	}

	return RefreshTokensReq{
		Token:       p.Token,
		Application: application,
		Device:      device,
	}, nil
}
