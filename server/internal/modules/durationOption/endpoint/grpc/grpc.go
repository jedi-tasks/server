package grpc

import (
	"context"

	"server/internal/modules/durationOption/model"
	durationOptionService "server/internal/modules/durationOption/service"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

var _ DurationOptionService = new(durationOptionService.DurationOptionService)

type DurationOptionService interface {
	GetDurationOptions(context.Context, model.GetDurationOptionsReq) ([]model.DurationOption, error)
	CreateDurationOption(context.Context, model.CreateDurationOptionReq) error
	UpdateDurationOption(context.Context, model.UpdateDurationOptionReq) error
	DeleteDurationOption(context.Context, model.DeleteDurationOptionReq) error
}

var _ proto.DurationOptionEndpointServer = new(DurationOptionEndpoint)

type DurationOptionEndpoint struct {
	proto.UnsafeDurationOptionEndpointServer
	durationOptionService DurationOptionService
}

func NewDurationOptionEndpoint(durationOptionService DurationOptionService) *DurationOptionEndpoint {
	return &DurationOptionEndpoint{
		durationOptionService: durationOptionService,
	}
}
