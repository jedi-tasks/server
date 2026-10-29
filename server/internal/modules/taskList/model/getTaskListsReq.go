package model

import (
	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type GetTaskListsReq struct {
	Necessary necessary.NecessaryUserInformation
	IDs       []uuid.UUID // Идентификаторы

}

// ProtoGetTaskListsReq обертка над proto-запросом
type ProtoGetTaskListsReq struct {
	*proto.GetTaskListsRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoGetTaskListsReq) ConvertToModel() (GetTaskListsReq, error) {
	var res GetTaskListsReq

	if p.GetTaskListsRequest == nil {
		return res, errors.BadRequest.New("GetTaskListsRequest is required")
	}

	return GetTaskListsReq{}, nil
}
