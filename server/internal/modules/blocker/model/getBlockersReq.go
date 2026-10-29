package model

import (
	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type GetBlockersReq struct {
	Necessary       necessary.NecessaryUserInformation
	IDs             []uuid.UUID // Идентификаторы
	TaskIDs         []uuid.UUID // Фильтр по заблокированным задачам
	BlockingTaskIDs []uuid.UUID // Фильтр по блокирующим задачам, только для сервера
	OnlyActive      bool        // Только неснятые блокеры, только для сервера
}

// ProtoGetBlockersReq обертка над proto-запросом
type ProtoGetBlockersReq struct {
	*proto.GetBlockersRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoGetBlockersReq) ConvertToModel() (GetBlockersReq, error) {
	var res GetBlockersReq

	if p.GetBlockersRequest == nil {
		return res, errors.BadRequest.New("GetBlockersRequest is required")
	}

	// Конвертируем фильтр по заблокированным задачам
	taskIDs, err := parseUUIDs(p.TaskIDs)
	if err != nil {
		return res, err
	}

	return GetBlockersReq{
		TaskIDs: taskIDs,
	}, nil
}
