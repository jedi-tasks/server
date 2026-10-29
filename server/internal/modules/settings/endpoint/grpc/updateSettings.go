package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/settings/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// UpdateSettings редактирование пользовательских настроек
func (e *SettingsEndpoint) UpdateSettings(ctx context.Context, r *proto.UpdateSettingsRequest) (*proto.UpdateSettingsResponse, error) {
	res := &proto.UpdateSettingsResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoUpdateSettingsReq{UpdateSettingsRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Получаем необходимую информацию из контекста
	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return res, err
	}

	// Валидируем запрос
	if err := validator.Validate(req); err != nil {
		return res, err
	}

	// Вызываем метод сервиса
	return res, e.settingsService.UpdateSettings(ctx, req)
}
