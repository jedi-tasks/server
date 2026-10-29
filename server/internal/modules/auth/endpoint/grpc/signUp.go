package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/auth/model"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// SignUp регистрация пользователя
func (e *AuthEndpoint) SignUp(ctx context.Context, r *proto.SignUpRequest) (*proto.SignUpResponse, error) {
	res := &proto.SignUpResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoSignUpReq{SignUpRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		return res, err
	}

	// Call service method
	authRes, err := e.authService.SignUp(ctx, req)
	if err != nil {
		return res, err
	}

	// Convert response to proto and return
	return authRes.ConvertToSignUpProto()
}
