package grpc

import (
	"context"

	"server/internal/modules/user/model"
	userService "server/internal/modules/user/service"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

var _ UserService = new(userService.UserService)

type UserService interface {
	GetUser(context.Context, model.GetUserReq) (model.User, error)
	UpdateUser(context.Context, model.UpdateUserReq) error
}

var _ proto.UserEndpointServer = new(UserEndpoint)

type UserEndpoint struct {
	proto.UnsafeUserEndpointServer
	userService UserService
}

func NewUserEndpoint(userService UserService) *UserEndpoint {
	return &UserEndpoint{
		userService: userService,
	}
}
