package model

import (
	"server/internal/utils/errors"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

type GetSettingsReq struct {
	Necessary necessary.NecessaryUserInformation
}

// ProtoGetSettingsReq обертка над proto-запросом
type ProtoGetSettingsReq struct {
	*proto.GetSettingsRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoGetSettingsReq) ConvertToModel() (GetSettingsReq, error) {
	if p.GetSettingsRequest == nil {
		return GetSettingsReq{}, errors.BadRequest.New("GetSettingsRequest is required")
	}
	return GetSettingsReq{}, nil //nolint:exhaustruct
}
