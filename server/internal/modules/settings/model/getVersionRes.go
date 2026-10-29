package model

import (
	"github.com/jedi-tasks/go-server-grpc/proto"
)

type GetVersionRes struct {
	Version Version
}

// ConvertToModel converts proto request to internal model
func (p GetVersionRes) ConvertToProto() (res *proto.GetVersionResponse, err error) {

	version, err := p.Version.ConvertToProto()
	if err != nil {
		return res, err
	}

	return &proto.GetVersionResponse{
		Version: version,
		Error:   nil,
	}, nil
}
