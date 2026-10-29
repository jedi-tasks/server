package service

import (
	"context"

	"github.com/google/uuid"
	"server/internal/modules/user/model"
	"server/internal/utils/errors"
)

// GetUser возвращает данные пользователя
func (s *UserService) GetUser(ctx context.Context, req model.GetUserReq) (model.User, error) {
	ctx, span := tracer.Start(ctx, "GetUser")
	defer span.End()

	// Создаем запрос для получения пользователей
	getUsersReq := model.GetUsersReq{
		Necessary: req.Necessary,
		IDs:       []uuid.UUID{req.Necessary.UserID},
	}

	users, err := s.userRepository.GetUsers(ctx, getUsersReq)
	if err != nil {
		return model.User{}, err
	}

	if len(users) == 0 {
		return model.User{}, errors.NotFound.New("User not found")
	}

	return users[0], nil
}
