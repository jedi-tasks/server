package grpc

import (
	"context"

	"server/internal/modules/blocker/model"
	blockerService "server/internal/modules/blocker/service"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

var _ BlockerService = new(blockerService.BlockerService)

type BlockerService interface {
	GetBlockers(context.Context, model.GetBlockersReq) ([]model.Blocker, error)
	CreateBlocker(context.Context, model.CreateBlockerReq) error
	DeleteBlocker(context.Context, model.DeleteBlockerReq) error
	ResolveBlocker(context.Context, model.ResolveBlockerReq) error
}

var _ proto.BlockerEndpointServer = new(BlockerEndpoint)

type BlockerEndpoint struct {
	proto.UnsafeBlockerEndpointServer
	blockerService BlockerService
}

func NewBlockerEndpoint(blockerService BlockerService) *BlockerEndpoint {
	return &BlockerEndpoint{
		blockerService: blockerService,
	}
}
