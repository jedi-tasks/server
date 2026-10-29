package grpc

import (
	"context"

	"server/internal/modules/listSection/model"
	listSectionService "server/internal/modules/listSection/service"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

var _ ListSectionService = new(listSectionService.ListSectionService)

type ListSectionService interface {
	GetListSections(context.Context, model.GetListSectionsReq) ([]model.ListSection, error)
	CreateListSection(context.Context, model.CreateListSectionReq) error
	UpdateListSection(context.Context, model.UpdateListSectionReq) error
	DeleteListSection(context.Context, model.DeleteListSectionReq) error
}

var _ proto.ListSectionEndpointServer = new(ListSectionEndpoint)

type ListSectionEndpoint struct {
	proto.UnsafeListSectionEndpointServer
	listSectionService ListSectionService
}

func NewListSectionEndpoint(listSectionService ListSectionService) *ListSectionEndpoint {
	return &ListSectionEndpoint{
		listSectionService: listSectionService,
	}
}
