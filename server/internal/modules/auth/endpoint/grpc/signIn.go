package grpc

import (
	"context"

	"server/internal/modules/auth/model"
	"server/internal/utils/validate"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// SignIn авторизация пользователя по логину и паролю
func (e *AuthEndpoint) SignIn(ctx context.Context, r *proto.SignInRequest) (*proto.SignInResponse, error) {
	res := &proto.SignInResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoSignInReq{SignInRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Validate request
	if err := validate.Validate(ctx, req); err != nil {
		return res, err
	}

	// Call service method
	authRes, err := e.authService.SignIn(ctx, req)
	if err != nil {
		return res, err
	}

	// Convert response to proto and return
	return authRes.ConvertToSignInProto()
}
