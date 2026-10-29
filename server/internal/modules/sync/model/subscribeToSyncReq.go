package model

import (
	"server/internal/utils/errors"
	"server/internal/utils/necessary"

	"github.com/jedi-tasks/go-server-grpc/proto"
)

// SubscribeToSyncReq - запрос на подписку на уведомления о новых изменениях (см. TrackMutation/
// syncNotifier) — сам не несёт ничего, кроме access token'а, распознаваемого интерцептором
type SubscribeToSyncReq struct {
	Necessary necessary.NecessaryUserInformation
}

// ProtoSubscribeToSyncReq wrapper for proto request
type ProtoSubscribeToSyncReq struct {
	*proto.SubscribeToSyncRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoSubscribeToSyncReq) ConvertToModel() (SubscribeToSyncReq, error) {
	if p.SubscribeToSyncRequest == nil {
		return SubscribeToSyncReq{}, errors.BadRequest.New("SubscribeToSyncRequest is required")
	}
	return SubscribeToSyncReq{}, nil
}
