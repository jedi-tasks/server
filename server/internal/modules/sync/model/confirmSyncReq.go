package model

import (
	"github.com/google/uuid"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"

	"github.com/jedi-tasks/go-server-grpc/proto"
)

// ConfirmSyncReq - запрос на подтверждение того, что клиент корректно применил изменения,
// полученные последним вызовом Sync
type ConfirmSyncReq struct {
	Necessary        necessary.NecessaryUserInformation
	PendingSyncToken uuid.UUID // Опорный токен из ответа Sync, который нужно подтвердить
}

// ProtoConfirmSyncReq wrapper for proto request
type ProtoConfirmSyncReq struct {
	*proto.ConfirmSyncRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoConfirmSyncReq) ConvertToModel() (ConfirmSyncReq, error) {
	if p.ConfirmSyncRequest == nil {
		return ConfirmSyncReq{}, errors.BadRequest.New("ConfirmSyncRequest is required")
	}

	pendingSyncToken, err := uuid.FromBytes(p.PendingSyncToken)
	if err != nil {
		return ConfirmSyncReq{}, errors.BadRequest.Wrap(err)
	}

	return ConfirmSyncReq{
		PendingSyncToken: pendingSyncToken,
	}, nil
}
