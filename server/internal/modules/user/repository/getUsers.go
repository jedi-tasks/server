package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	userModel "server/internal/modules/user/model"
	"server/internal/modules/user/repository/userDDL"
)

// GetUsers Возвращает пользователей по фильтрам
func (r *UserRepository) GetUsers(ctx context.Context, filters userModel.GetUsersReq) (user []userModel.User, err error) {
	ctx, span := tracer.Start(ctx, "GetUsers")
	defer span.End()

	filtersEq := make(sq.Eq)

	if len(filters.IDs) > 0 {
		filtersEq[userDDL.ColumnID] = filters.IDs
	}
	if len(filters.Emails) > 0 {
		filtersEq[userDDL.ColumnEmail] = filters.Emails
	}

	return user, r.db.Select(ctx, &user, sq.
		Select(ddlHelper.SelectAll).
		From(userDDL.Table).
		Where(filtersEq),
	)
}
