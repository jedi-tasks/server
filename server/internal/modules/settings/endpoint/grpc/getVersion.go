package grpc

import (
	"context"
	"pkg/validator"
	"server/internal/modules/settings/model"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// GetVersion получение текущей версии сервера
func (e *SettingsEndpoint) GetVersion(ctx context.Context, r *proto.GetVersionRequest) (*proto.GetVersionResponse, error) {
	res := new(proto.GetVersionResponse)

	// Convert proto request to internal model
	req, err := model.ProtoGetVersionReq{GetVersionRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		return res, err
	}

	// Call service method
	_res, err := e.settingsService.GetVersion(ctx, req)
	if err != nil {
		return res, err
	}

	return _res.ConvertToProto()
}
