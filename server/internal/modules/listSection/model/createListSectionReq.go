package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	repoModel "server/internal/modules/listSection/repository/model"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type CreateListSectionReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"`        // Идентификатор, генерируется клиентом
	ListID    uuid.UUID `json:"listID" validate:"required"`    // Идентификатор списка
	Name      string    `json:"name"`                          // Название раздела
	SortOrder float64   `json:"sortOrder"`                     // Порядок раздела
	UpdatedAt time.Time `json:"updatedAt" validate:"required"` // Время создания на клиенте
}

// ConvertToRepoReq преобразует запрос в модель репозитория
func (s CreateListSectionReq) ConvertToRepoReq() repoModel.CreateListSectionReq {
	return repoModel.CreateListSectionReq{
		ID:        s.ID,
		UserID:    s.Necessary.UserID,
		ListID:    s.ListID,
		Name:      s.Name,
		SortOrder: s.SortOrder,
		UpdatedAt: s.UpdatedAt,
	}
}

// ProtoCreateListSectionReq обертка над proto-запросом
type ProtoCreateListSectionReq struct {
	*proto.CreateListSectionRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoCreateListSectionReq) ConvertToModel() (CreateListSectionReq, error) {
	var res CreateListSectionReq

	if p.CreateListSectionRequest == nil {
		return res, errors.BadRequest.New("CreateListSectionRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	listID, err := uuid.FromBytes(p.ListID)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}
	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return CreateListSectionReq{
		ID:        id,
		ListID:    listID,
		Name:      p.Name,
		SortOrder: p.SortOrder,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
