package model

import (
	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type GetTasksToTagsReq struct {
	Necessary necessary.NecessaryUserInformation
	TaskIDs   []uuid.UUID // Фильтр по задачам
	TagIDs    []uuid.UUID // Фильтр по тегам, только для сервера
}

// ProtoGetTasksToTagsReq обертка над proto-запросом
type ProtoGetTasksToTagsReq struct {
	*proto.GetTasksToTagsRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoGetTasksToTagsReq) ConvertToModel() (GetTasksToTagsReq, error) {
	var res GetTasksToTagsReq

	if p.GetTasksToTagsRequest == nil {
		return res, errors.BadRequest.New("GetTasksToTagsRequest is required")
	}

	taskIDs, err := parseUUIDs(p.TaskIDs)
	if err != nil {
		return res, err
	}

	return GetTasksToTagsReq{
		TaskIDs: taskIDs,
	}, nil
}
