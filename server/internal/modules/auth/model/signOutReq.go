package model

import (
	"github.com/jedi-tasks/go-server-grpc/proto"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type SignOutReq struct {
	Necessary necessary.NecessaryUserInformation
}

// ProtoSignOutReq wrapper for proto request
type ProtoSignOutReq struct {
	*proto.SignOutRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoSignOutReq) ConvertToModel() (SignOutReq, error) {
	if p.SignOutRequest == nil {
		return SignOutReq{}, errors.BadRequest.New("SignOutRequest is required")
	}

	return SignOutReq{}, nil
}
