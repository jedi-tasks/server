package model

import (
	"server/internal/enum/applicationType"

	"server/internal/utils/errors"

	"github.com/jedi-tasks/go-server-grpc/proto"
)

type GetVersionReq struct {
	ApplicationType applicationType.ApplicationType
}

// ProtoGetVersionReq wrapper for proto request
type ProtoGetVersionReq struct {
	*proto.GetVersionRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoGetVersionReq) ConvertToModel() (res GetVersionReq, err error) {
	if p.GetVersionRequest == nil {
		return res, errors.BadRequest.New("GetVersionRequest is required")
	}

	applicationType, err := applicationType.ProtoApplicationType{ApplicationType: p.ApplicationType}.ConvertToModel()
	if err != nil {
		return res, err
	}

	return GetVersionReq{
		ApplicationType: applicationType,
	}, nil
}
