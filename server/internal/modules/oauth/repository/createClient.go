package repository

import (
	"context"
	"encoding/json"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/oauth/model"
	"server/internal/modules/oauth/repository/oauthClientDDL"
	"server/internal/utils/errors"
)

// CreateClient регистрирует клиента
func (r *OAuthRepository) CreateClient(ctx context.Context, req model.CreateClientReq) error {
	ctx, span := tracer.Start(ctx, "CreateClient")
	defer span.End()

	redirectURIs, err := json.Marshal(req.RedirectURIs)
	if err != nil {
		return errors.InternalServer.Wrap(err)
	}

	return r.db.Exec(ctx, sq.Insert(oauthClientDDL.Table).SetMap(map[string]any{
		oauthClientDDL.ColumnID:           req.ID,
		oauthClientDDL.ColumnName:         req.Name,
		oauthClientDDL.ColumnRedirectURIs: redirectURIs,
		oauthClientDDL.ColumnCreatedAt:    req.CreatedAt,
		oauthClientDDL.ColumnLastUsedAt:   req.CreatedAt,
	}))
}
