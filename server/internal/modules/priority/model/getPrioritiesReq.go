package model

import (
	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type GetPrioritiesReq struct {
	Necessary necessary.NecessaryUserInformation
	IDs       []uuid.UUID // Идентификаторы

}

// ProtoGetPrioritiesReq обертка над proto-запросом
type ProtoGetPrioritiesReq struct {
	*proto.GetPrioritiesRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoGetPrioritiesReq) ConvertToModel() (GetPrioritiesReq, error) {
	var res GetPrioritiesReq

	if p.GetPrioritiesRequest == nil {
		return res, errors.BadRequest.New("GetPrioritiesRequest is required")
	}

	return GetPrioritiesReq{}, nil
}
