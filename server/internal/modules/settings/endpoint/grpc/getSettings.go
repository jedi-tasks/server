package grpc

import (
	"context"

	"server/internal/modules/settings/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// GetSettings получение пользовательских настроек
func (e *SettingsEndpoint) GetSettings(ctx context.Context, r *proto.GetSettingsRequest) (*proto.GetSettingsResponse, error) {
	res := &proto.GetSettingsResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoGetSettingsReq{GetSettingsRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Получаем необходимую информацию из контекста
	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return res, err
	}

	// Вызываем метод сервиса
	settings, err := e.settingsService.GetSettings(ctx, req.Necessary.UserID)
	if err != nil {
		return res, err
	}

	// Конвертируем ответ в proto-формат
	res.Settings, err = settings.ConvertToProto()
	return res, err
}
