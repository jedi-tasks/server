package model

import (
	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type GetDurationOptionsReq struct {
	Necessary necessary.NecessaryUserInformation
	IDs       []uuid.UUID // Идентификаторы

}

// ProtoGetDurationOptionsReq обертка над proto-запросом
type ProtoGetDurationOptionsReq struct {
	*proto.GetDurationOptionsRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoGetDurationOptionsReq) ConvertToModel() (GetDurationOptionsReq, error) {
	var res GetDurationOptionsReq

	if p.GetDurationOptionsRequest == nil {
		return res, errors.BadRequest.New("GetDurationOptionsRequest is required")
	}

	return GetDurationOptionsReq{}, nil
}
