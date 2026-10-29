package model

import "github.com/jedi-tasks/go-server-grpc/proto"

type RefreshTokensRes struct {
	Tokens // Токены доступа
}

func (r RefreshTokensRes) ConvertToProto() (*proto.RefreshTokensResponse, error) {
	return &proto.RefreshTokensResponse{
		Error:        nil,
		AccessToken:  &r.Tokens.AccessToken,
		RefreshToken: &r.Tokens.RefreshToken,
	}, nil
}
