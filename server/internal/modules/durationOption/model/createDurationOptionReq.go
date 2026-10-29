package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	repoModel "server/internal/modules/durationOption/repository/model"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type CreateDurationOptionReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"`        // Идентификатор, генерируется клиентом
	Title     string    `json:"title"`                         // Подпись
	Minutes   int32     `json:"minutes"`                       // Минуты, меньше нуля — непрогнозируемо долго
	UpdatedAt time.Time `json:"updatedAt" validate:"required"` // Время создания на клиенте
}

// ConvertToRepoReq преобразует запрос в модель репозитория
func (s CreateDurationOptionReq) ConvertToRepoReq() repoModel.CreateDurationOptionReq {
	return repoModel.CreateDurationOptionReq{
		ID:        s.ID,
		UserID:    s.Necessary.UserID,
		Title:     s.Title,
		Minutes:   s.Minutes,
		UpdatedAt: s.UpdatedAt,
	}
}

// ProtoCreateDurationOptionReq обертка над proto-запросом
type ProtoCreateDurationOptionReq struct {
	*proto.CreateDurationOptionRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoCreateDurationOptionReq) ConvertToModel() (CreateDurationOptionReq, error) {
	var res CreateDurationOptionReq

	if p.CreateDurationOptionRequest == nil {
		return res, errors.BadRequest.New("CreateDurationOptionRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return CreateDurationOptionReq{
		ID:        id,
		Title:     p.Title,
		Minutes:   p.Minutes,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
