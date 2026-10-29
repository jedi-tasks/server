package grpc

import (
	"context"

	"server/internal/modules/tag/model"
	tagService "server/internal/modules/tag/service"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

var _ TagService = new(tagService.TagService)

type TagService interface {
	GetTags(context.Context, model.GetTagsReq) ([]model.Tag, error)
	CreateTag(context.Context, model.CreateTagReq) error
	UpdateTag(context.Context, model.UpdateTagReq) error
	DeleteTag(context.Context, model.DeleteTagReq) error
	GetTasksToTags(context.Context, model.GetTasksToTagsReq) ([]model.TaskToTag, error)
	AddTagToTask(context.Context, model.LinkTagReq) error
	RemoveTagFromTask(context.Context, model.LinkTagReq) error
}

var _ proto.TagEndpointServer = new(TagEndpoint)

type TagEndpoint struct {
	proto.UnsafeTagEndpointServer
	tagService TagService
}

func NewTagEndpoint(tagService TagService) *TagEndpoint {
	return &TagEndpoint{
		tagService: tagService,
	}
}
