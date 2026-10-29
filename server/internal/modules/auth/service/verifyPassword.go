package service

import (
	"context"

	"github.com/google/uuid"

	"pkg/passwordManager"
	userModel "server/internal/modules/user/model"
	"server/internal/utils/errors"
)

// VerifyPassword проверяет пароль пользователя по почте и возвращает его идентификатор
func (s *AuthService) VerifyPassword(ctx context.Context, email, password string) (uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "VerifyPassword")
	defer span.End()

	// Получаем пользователя по email
	users, err := s.userRepository.GetUsers(ctx, userModel.GetUsersReq{Emails: []string{email}}) //nolint:exhaustruct
	if err != nil {
		return uuid.Nil, err
	}
	if len(users) == 0 {
		return uuid.Nil, errors.NotFound.New("User not found").
			WithContextParams(ctx).
			WithCustomHumanText("Пользователь не найден")
	}
	user := users[0]

	// Сравниваем пришедший пароль и хэш пароля из базы данных
	if err = passwordManager.CompareHashAndPassword(user.PasswordHash, []byte(password), user.PasswordSalt, s.generalSalt); err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}
