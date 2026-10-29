package grpc

import (
	"context"

	"server/internal/modules/priority/model"
	priorityService "server/internal/modules/priority/service"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

var _ PriorityService = new(priorityService.PriorityService)

type PriorityService interface {
	GetPriorities(context.Context, model.GetPrioritiesReq) ([]model.Priority, error)
	CreatePriority(context.Context, model.CreatePriorityReq) error
	UpdatePriority(context.Context, model.UpdatePriorityReq) error
	DeletePriority(context.Context, model.DeletePriorityReq) error
}

var _ proto.PriorityEndpointServer = new(PriorityEndpoint)

type PriorityEndpoint struct {
	proto.UnsafePriorityEndpointServer
	priorityService PriorityService
}

func NewPriorityEndpoint(priorityService PriorityService) *PriorityEndpoint {
	return &PriorityEndpoint{
		priorityService: priorityService,
	}
}
