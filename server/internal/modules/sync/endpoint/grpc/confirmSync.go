package grpc

import (
	"context"

	"server/internal/modules/sync/model"
	"server/internal/utils/necessary"
	"server/internal/utils/validate"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// ConfirmSync подтверждает, что клиент корректно применил изменения, полученные последним вызовом Sync
func (e *SyncEndpoint) ConfirmSync(ctx context.Context, r *proto.ConfirmSyncRequest) (*proto.ConfirmSyncResponse, error) {
	res := new(proto.ConfirmSyncResponse)

	// Convert proto request to internal model
	req, err := model.ProtoConfirmSyncReq{ConfirmSyncRequest: r}.ConvertToModel()
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
	return res, e.syncService.ConfirmSync(ctx, req)
}
