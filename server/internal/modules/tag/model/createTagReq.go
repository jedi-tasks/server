package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	repoModel "server/internal/modules/tag/repository/model"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type CreateTagReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"`        // Идентификатор, генерируется клиентом
	Name      string    `json:"name"`                          // Название тега
	UpdatedAt time.Time `json:"updatedAt" validate:"required"` // Время создания на клиенте
}

// ConvertToRepoReq преобразует запрос в модель репозитория
func (s CreateTagReq) ConvertToRepoReq() repoModel.CreateTagReq {
	return repoModel.CreateTagReq{
		ID:        s.ID,
		UserID:    s.Necessary.UserID,
		Name:      s.Name,
		UpdatedAt: s.UpdatedAt,
	}
}

// ProtoCreateTagReq обертка над proto-запросом
type ProtoCreateTagReq struct {
	*proto.CreateTagRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoCreateTagReq) ConvertToModel() (CreateTagReq, error) {
	var res CreateTagReq

	if p.CreateTagRequest == nil {
		return res, errors.BadRequest.New("CreateTagRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return CreateTagReq{
		ID:        id,
		Name:      p.Name,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
