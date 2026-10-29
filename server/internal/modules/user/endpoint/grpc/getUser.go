package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/user/model"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// GetUser получение данных пользователя
func (e *UserEndpoint) GetUser(ctx context.Context, r *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	res := &proto.GetUserResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoGetUserReq{GetUserRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Parse necessary information from context
	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return res, err
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		return res, err
	}

	// Call service method
	user, err := e.userService.GetUser(ctx, req)
	if err != nil {
		return res, err
	}

	// Convert user to proto format
	protoUser, err := user.ConvertToProto()
	if err != nil {
		return res, err
	}

	res.User = protoUser
	return res, nil
}
