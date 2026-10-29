package model

import (
	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type GetListSectionsReq struct {
	Necessary necessary.NecessaryUserInformation
	IDs       []uuid.UUID // Идентификаторы
	ListIDs   []uuid.UUID // Фильтр по спискам
}

// ProtoGetListSectionsReq обертка над proto-запросом
type ProtoGetListSectionsReq struct {
	*proto.GetListSectionsRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoGetListSectionsReq) ConvertToModel() (GetListSectionsReq, error) {
	var res GetListSectionsReq

	if p.GetListSectionsRequest == nil {
		return res, errors.BadRequest.New("GetListSectionsRequest is required")
	}

	// Конвертируем фильтр по спискам
	listIDs, err := parseUUIDs(p.ListIDs)
	if err != nil {
		return res, err
	}

	return GetListSectionsReq{
		ListIDs: listIDs,
	}, nil
}
