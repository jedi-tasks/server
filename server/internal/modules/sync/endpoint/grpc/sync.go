package grpc

import (
	"context"

	"server/internal/modules/sync/model"
	"server/internal/utils/necessary"
	"server/internal/utils/validate"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// Sync атомарно получает все изменения, доступные пользователю по группам счетов, произошедшие
// после последней синхронизации этого устройства
func (e *SyncEndpoint) Sync(ctx context.Context, r *proto.SyncRequest) (*proto.SyncResponse, error) {
	res := new(proto.SyncResponse)

	// Convert proto request to internal model
	req, err := model.ProtoSyncReq{SyncRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Parse necessary information from context
	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return res, err
	}

	// Validate request
	if err := validate.Validate(ctx, req); err != nil {
		return res, err
	}

	// Call service method
	_res, err := e.syncService.Sync(ctx, req)
	if err != nil {
		return res, err
	}

	return _res.ConvertToProto()
}
