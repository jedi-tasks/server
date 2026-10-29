package model

import (
	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type GetTasksReq struct {
	Necessary necessary.NecessaryUserInformation
	IDs       []uuid.UUID // Идентификаторы
	ListIDs   []uuid.UUID // Фильтр по спискам
	Completed *bool       // Фильтр по выполненности, nil — все
}

// ProtoGetTasksReq обертка над proto-запросом
type ProtoGetTasksReq struct {
	*proto.GetTasksRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoGetTasksReq) ConvertToModel() (GetTasksReq, error) {
	var res GetTasksReq

	if p.GetTasksRequest == nil {
		return res, errors.BadRequest.New("GetTasksRequest is required")
	}

	// Конвертируем фильтр по спискам
	listIDs, err := parseUUIDs(p.ListIDs)
	if err != nil {
		return res, err
	}

	// Конвертируем фильтр по идентификаторам
	ids, err := parseUUIDs(p.Ids)
	if err != nil {
		return res, err
	}

	return GetTasksReq{
		IDs:       ids,
		ListIDs:   listIDs,
		Completed: p.Completed,
	}, nil
}
