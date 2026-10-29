package service

import (
	"context"

	"server/internal/modules/auth/model"
)

// SignOut удаляет данные девайса пользователя
func (s *AuthService) SignOut(ctx context.Context, req model.SignOutReq) error {
	ctx, span := tracer.Start(ctx, "SignOut")
	defer span.End()

	return s.userRepository.DeleteDevice(ctx, req.Necessary.UserID, req.Necessary.DeviceID)
}
