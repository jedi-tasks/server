package grpc

import (
	"context"

	"server/internal/modules/auth/model"
	authService "server/internal/modules/auth/service"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

var _ AuthService = new(authService.AuthService)

type AuthService interface {
	SignIn(context.Context, model.SignInReq) (model.AuthRes, error)
	SignUp(context.Context, model.SignUpReq) (model.AuthRes, error)
	SignOut(context.Context, model.SignOutReq) error
	RefreshTokens(context.Context, model.RefreshTokensReq) (model.RefreshTokensRes, error)
}

var _ proto.AuthEndpointServer = new(AuthEndpoint)

type AuthEndpoint struct {
	proto.UnsafeAuthEndpointServer
	authService AuthService
}

func NewAuthEndpoint(authService AuthService) *AuthEndpoint {
	return &AuthEndpoint{
		authService: authService,
	}
}
