package grpc

import (
	"context"

	"github.com/google/uuid"

	"server/internal/modules/sync/model"
	syncService "server/internal/modules/sync/service"

	"github.com/jedi-tasks/go-server-grpc/proto"
)

var _ SyncService = new(syncService.SyncService)

type SyncService interface {
	Sync(context.Context, model.SyncReq) (model.SyncRes, error)
	ConfirmSync(context.Context, model.ConfirmSyncReq) error
	SubscribeToSync(ctx context.Context, userID uuid.UUID, notify func() error) error
}

var _ proto.SyncEndpointServer = new(SyncEndpoint)

type SyncEndpoint struct {
	proto.UnsafeSyncEndpointServer
	syncService SyncService
}

func NewSyncEndpoint(syncService SyncService) *SyncEndpoint {
	return &SyncEndpoint{
		syncService: syncService,
	}
}
