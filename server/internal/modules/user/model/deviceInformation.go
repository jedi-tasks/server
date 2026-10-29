package model

import (
	"server/internal/enum/osType"
	"server/internal/utils/errors"

	"github.com/jedi-tasks/go-server-grpc/proto"
)

type DeviceInformation struct {
	NameOS     osType.OSType `json:"nameOS" validate:"required" db:"device_os_name"`       // Название операционной системы
	VersionOS  string        `json:"versionOS" validate:"required" db:"device_os_version"` // Версия операционной системы
	DeviceName string        `json:"deviceName" validate:"required" db:"device_name"`      // Название девайса
	ModelName  string        `json:"modelName" validate:"required" db:"device_model_name"` // Название модели
	IPAddress  string        `json:"-" db:"device_ip_address"`                             // IP-адрес
	UserAgent  string        `json:"-" db:"device_user_agent"`                             // UserAgent
}

type ProtoDeviceInformation struct {
	*proto.DeviceInformation
}

func (p ProtoDeviceInformation) ConvertToModel() (res DeviceInformation, err error) {
	if p.DeviceInformation == nil {
		return res, errors.BadRequest.New("DeviceInformation is required")
	}

	return DeviceInformation{
		NameOS:     osType.OSType(p.NameOS),
		VersionOS:  p.VersionOS,
		DeviceName: p.DeviceName,
		ModelName:  p.ModelName,
		IPAddress:  "",
		UserAgent:  "",
	}, nil
}
