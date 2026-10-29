package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"server/internal/modules/user/model"
	"server/internal/utils/validate"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// UpdateUser редактирование пользователя
func (e *UserEndpoint) UpdateUser(ctx context.Context, r *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	res := &proto.UpdateUserResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoUpdateUserReq{UpdateUserRequest: r}.ConvertToModel()
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
	err = e.userService.UpdateUser(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}
