package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	userRepoModel "server/internal/modules/user/repository/model"
	"server/internal/modules/user/repository/userDDL"
	"server/internal/utils/errors"
)

// UpdateUser редактирует пользователя
func (r *UserRepository) UpdateUser(ctx context.Context, fields userRepoModel.UpdateUserReq) error {
	ctx, span := tracer.Start(ctx, "UpdateUser")
	defer span.End()

	// Изменяем поля пользователя
	updates := make(map[string]any)

	// Добавляем в запрос поля, которые нужно изменить
	if fields.Email != nil {
		updates[userDDL.ColumnEmail] = *fields.Email
	}
	if fields.Name != nil {
		updates[userDDL.ColumnName] = *fields.Name
	}
	if fields.PasswordHash != nil && fields.PasswordSalt != nil {
		updates[userDDL.ColumnPasswordHash] = *fields.PasswordHash
		updates[userDDL.ColumnPasswordSalt] = *fields.PasswordSalt
	}

	if len(updates) == 0 {
		return errors.BadRequest.New("No fields to update").
			WithContextParams(ctx)
	}

	// Обновляем пользователя
	return r.db.Exec(ctx, sq.
		Update(userDDL.Table).
		SetMap(updates).
		Where(sq.Eq{userDDL.ColumnID: fields.ID}),
	)
}
