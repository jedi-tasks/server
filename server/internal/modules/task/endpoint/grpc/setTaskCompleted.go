package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/task/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// SetTaskCompleted выполнение задачи или возврат в работу
func (e *TaskEndpoint) SetTaskCompleted(ctx context.Context, r *proto.SetTaskCompletedRequest) (*proto.SetTaskCompletedResponse, error) {
	res := &proto.SetTaskCompletedResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoSetTaskCompletedReq{SetTaskCompletedRequest: r}.ConvertToModel()
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
	return res, e.taskService.SetTaskCompleted(ctx, req)
}
