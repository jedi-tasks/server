package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

// LinkTagReq — привязка или отвязка тега от задачи
type LinkTagReq struct {
	Necessary necessary.NecessaryUserInformation
	TaskID    uuid.UUID `json:"taskID" validate:"required"`    // Идентификатор задачи
	TagID     uuid.UUID `json:"tagID" validate:"required"`     // Идентификатор тега
	UpdatedAt time.Time `json:"updatedAt" validate:"required"` // Время действия на клиенте
}

// ProtoAddTagToTaskReq обертка над proto-запросом
type ProtoAddTagToTaskReq struct {
	*proto.AddTagToTaskRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoAddTagToTaskReq) ConvertToModel() (LinkTagReq, error) {
	if p.AddTagToTaskRequest == nil {
		return LinkTagReq{}, errors.BadRequest.New("AddTagToTaskRequest is required")
	}
	return convertLinkTagReq(p.TaskID, p.TagID, p.UpdatedAt.AsTime(), p.UpdatedAt != nil)
}

// ProtoRemoveTagFromTaskReq обертка над proto-запросом
type ProtoRemoveTagFromTaskReq struct {
	*proto.RemoveTagFromTaskRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoRemoveTagFromTaskReq) ConvertToModel() (LinkTagReq, error) {
	if p.RemoveTagFromTaskRequest == nil {
		return LinkTagReq{}, errors.BadRequest.New("RemoveTagFromTaskRequest is required")
	}
	return convertLinkTagReq(p.TaskID, p.TagID, p.UpdatedAt.AsTime(), p.UpdatedAt != nil)
}

// convertLinkTagReq разбирает общие поля запросов привязки и отвязки тега
func convertLinkTagReq(taskIDBytes, tagIDBytes []byte, updatedAt time.Time, hasUpdatedAt bool) (LinkTagReq, error) {
	var res LinkTagReq

	taskID, err := uuid.FromBytes(taskIDBytes)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	tagID, err := uuid.FromBytes(tagIDBytes)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	if !hasUpdatedAt {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return LinkTagReq{
		TaskID:    taskID,
		TagID:     tagID,
		UpdatedAt: updatedAt,
	}, nil
}
