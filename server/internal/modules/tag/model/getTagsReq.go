package model

import (
	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type GetTagsReq struct {
	Necessary necessary.NecessaryUserInformation
	IDs       []uuid.UUID // Идентификаторы

}

// ProtoGetTagsReq обертка над proto-запросом
type ProtoGetTagsReq struct {
	*proto.GetTagsRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoGetTagsReq) ConvertToModel() (GetTagsReq, error) {
	var res GetTagsReq

	if p.GetTagsRequest == nil {
		return res, errors.BadRequest.New("GetTagsRequest is required")
	}

	return GetTagsReq{}, nil
}
