package model

import (
	"github.com/google/uuid"
	"github.com/jedi-tasks/go-server-grpc/proto"
)

type AuthRes struct {
	Tokens `json:"token"` // Токены доступа
	ID     uuid.UUID      `json:"id"` // Идентификатор пользователя
}

func (a AuthRes) ConvertToSignInProto() (res *proto.SignInResponse, err error) {

	tokens, err := a.Tokens.ConvertToProto()
	if err != nil {
		return nil, err
	}

	return &proto.SignInResponse{
		Id:    a.ID[:],
		Token: tokens,
		Error: nil,
	}, nil
}

func (a AuthRes) ConvertToSignUpProto() (res *proto.SignUpResponse, err error) {

	tokens, err := a.Tokens.ConvertToProto()
	if err != nil {
		return nil, err
	}

	return &proto.SignUpResponse{
		Id:    a.ID[:],
		Token: tokens,
		Error: nil,
	}, nil
}
