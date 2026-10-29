package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"pkg/ddlHelper"
	"server/internal/modules/oauth/model"
	repoModel "server/internal/modules/oauth/repository/model"
	"server/internal/modules/oauth/repository/oauthClientDDL"
)

// GetClient возвращает клиента по client_id; нет клиента — found=false
func (r *OAuthRepository) GetClient(ctx context.Context, id uuid.UUID) (client model.Client, found bool, err error) {
	ctx, span := tracer.Start(ctx, "GetClient")
	defer span.End()

	var rows []repoModel.Client
	err = r.db.Select(ctx, &rows, sq.Select(ddlHelper.SelectAll).From(oauthClientDDL.Table).Where(sq.Eq{oauthClientDDL.ColumnID: id}))
	if err != nil || len(rows) == 0 {
		return client, false, err
	}

	client, err = rows[0].ConvertToModel()
	return client, err == nil, err
}
