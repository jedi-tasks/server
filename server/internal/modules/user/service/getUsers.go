package service

import (
	"context"

	"server/internal/modules/user/model"
)

// GetUsers возвращает всех юзеров по фильтрам
func (s *UserService) GetUsers(ctx context.Context, filters model.GetUsersReq) (users []model.User, err error) {
	ctx, span := tracer.Start(ctx, "GetUsers")
	defer span.End()

	return s.userRepository.GetUsers(ctx, filters)
}
