package model

import (
	"server/internal/utils/errors"
	"server/internal/utils/necessary"

	"github.com/jedi-tasks/go-server-grpc/proto"
)

// SyncReq - запрос на синхронизацию изменений, доступных пользователю по группам счетов, произошедших
// после чекпоинта, хранящегося на клиенте
type SyncReq struct {
	Necessary necessary.NecessaryUserInformation
	SinceID   uint32 // Чекпоинт устройства - идентификатор последней подтвержденной записи аудит-лога
}

// ProtoSyncReq wrapper for proto request
type ProtoSyncReq struct {
	*proto.SyncRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoSyncReq) ConvertToModel() (SyncReq, error) {
	if p.SyncRequest == nil {
		return SyncReq{}, errors.BadRequest.New("SyncRequest is required")
	}

	return SyncReq{
		SinceID: p.SinceID,
	}, nil
}
