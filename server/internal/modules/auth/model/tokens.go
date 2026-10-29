package model

import "github.com/jedi-tasks/go-server-grpc/proto"

type Tokens struct {
	AccessToken  string `json:"accessToken"`  // Токен доступа
	RefreshToken string `json:"refreshToken"` // Токен восстановления доступа
}

func (t Tokens) ConvertToProto() (*proto.Tokens, error) {
	return &proto.Tokens{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
	}, nil
}
