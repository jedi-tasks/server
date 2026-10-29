package grpc

import (
	"server/internal/modules/sync/model"
	"server/internal/utils/necessary"
	"server/internal/utils/validate"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// SubscribeToSync держит стрим открытым и присылает пустой SyncNotification каждый раз, когда
// для пользователя появились изменения — клиент в ответ сам вызывает обычный Sync/ConfirmSync.
func (e *SyncEndpoint) SubscribeToSync(r *proto.SubscribeToSyncRequest, stream proto.SyncEndpoint_SubscribeToSyncServer) error {
	ctx := stream.Context()

	// Convert proto request to internal model
	req, err := model.ProtoSubscribeToSyncReq{SubscribeToSyncRequest: r}.ConvertToModel()
	if err != nil {
		return err
	}

	// Parse necessary information from context
	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return err
	}

	// Validate request
	if err := validate.Validate(ctx, req); err != nil {
		return err
	}

	return e.syncService.SubscribeToSync(ctx, req.Necessary.UserID, func() error {
		return stream.Send(&proto.SyncNotification{})
	})
}
