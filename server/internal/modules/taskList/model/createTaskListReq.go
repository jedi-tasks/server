package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/enum/sortMode"
	repoModel "server/internal/modules/taskList/repository/model"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type CreateTaskListReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID         `json:"id" validate:"required"`        // Идентификатор, генерируется клиентом
	Name      string            `json:"name"`                          // Название списка
	ColorHex  string            `json:"colorHex"`                      // Цвет списка
	Icon      string            `json:"icon"`                          // SF Symbol иконки списка
	SortOrder float64           `json:"sortOrder"`                     // Порядок списка
	SortMode  sortMode.SortMode `json:"sortMode"`                      // Режим сортировки задач в списке
	CreatedAt time.Time         `json:"createdAt" validate:"required"` // Время создания на клиенте
}

// ConvertToRepoReq преобразует запрос в модель репозитория
func (s CreateTaskListReq) ConvertToRepoReq() repoModel.CreateTaskListReq {
	return repoModel.CreateTaskListReq{
		ID:        s.ID,
		UserID:    s.Necessary.UserID,
		Name:      s.Name,
		ColorHex:  s.ColorHex,
		Icon:      s.Icon,
		SortOrder: s.SortOrder,
		SortMode:  s.SortMode,
		CreatedAt: s.CreatedAt,
	}
}

// ProtoCreateTaskListReq обертка над proto-запросом
type ProtoCreateTaskListReq struct {
	*proto.CreateTaskListRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoCreateTaskListReq) ConvertToModel() (CreateTaskListReq, error) {
	var res CreateTaskListReq

	if p.CreateTaskListRequest == nil {
		return res, errors.BadRequest.New("CreateTaskListRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	sortModeValue, err := sortMode.ProtoSortMode{SortMode: p.SortMode}.ConvertToModel()
	if err != nil {
		return res, err
	}
	if p.CreatedAt == nil {
		return res, errors.BadRequest.New("CreatedAt is required")
	}

	return CreateTaskListReq{
		ID:        id,
		Name:      p.Name,
		ColorHex:  p.ColorHex,
		Icon:      p.Icon,
		SortOrder: p.SortOrder,
		SortMode:  sortModeValue,
		CreatedAt: p.CreatedAt.AsTime(),
	}, nil
}
