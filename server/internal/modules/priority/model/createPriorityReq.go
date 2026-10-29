package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	repoModel "server/internal/modules/priority/repository/model"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type CreatePriorityReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"`        // Идентификатор, генерируется клиентом
	Name      string    `json:"name"`                          // Название
	ColorHex  string    `json:"colorHex"`                      // Цвет
	Marks     string    `json:"marks"`                         // Знаки в строке задачи
	SortOrder float64   `json:"sortOrder"`                     // Порядок важности
	UpdatedAt time.Time `json:"updatedAt" validate:"required"` // Время создания на клиенте
}

// ConvertToRepoReq преобразует запрос в модель репозитория
func (s CreatePriorityReq) ConvertToRepoReq() repoModel.CreatePriorityReq {
	return repoModel.CreatePriorityReq{
		ID:        s.ID,
		UserID:    s.Necessary.UserID,
		Name:      s.Name,
		ColorHex:  s.ColorHex,
		Marks:     s.Marks,
		SortOrder: s.SortOrder,
		UpdatedAt: s.UpdatedAt,
	}
}

// ProtoCreatePriorityReq обертка над proto-запросом
type ProtoCreatePriorityReq struct {
	*proto.CreatePriorityRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoCreatePriorityReq) ConvertToModel() (CreatePriorityReq, error) {
	var res CreatePriorityReq

	if p.CreatePriorityRequest == nil {
		return res, errors.BadRequest.New("CreatePriorityRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return CreatePriorityReq{
		ID:        id,
		Name:      p.Name,
		ColorHex:  p.ColorHex,
		Marks:     p.Marks,
		SortOrder: p.SortOrder,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
