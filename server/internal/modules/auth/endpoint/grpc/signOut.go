package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/auth/model"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// SignOut выход пользователя из приложения
func (e *AuthEndpoint) SignOut(ctx context.Context, r *proto.SignOutRequest) (*proto.SignOutResponse, error) {
	res := &proto.SignOutResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoSignOutReq{SignOutRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		return res, err
	}

	// Call service method
	err = e.authService.SignOut(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}
