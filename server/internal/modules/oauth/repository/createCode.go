package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/oauth/model"
	"server/internal/modules/oauth/repository/oauthCodeDDL"
)

// CreateCode сохраняет выданный authorization code
func (r *OAuthRepository) CreateCode(ctx context.Context, req model.CreateCodeReq) error {
	ctx, span := tracer.Start(ctx, "CreateCode")
	defer span.End()

	return r.db.Exec(ctx, sq.Insert(oauthCodeDDL.Table).SetMap(map[string]any{
		oauthCodeDDL.ColumnCode:          req.Code,
		oauthCodeDDL.ColumnClientID:      req.ClientID,
		oauthCodeDDL.ColumnUserID:        req.UserID,
		oauthCodeDDL.ColumnRedirectURI:   req.RedirectURI,
		oauthCodeDDL.ColumnCodeChallenge: req.CodeChallenge,
		oauthCodeDDL.ColumnResource:      req.Resource,
		oauthCodeDDL.ColumnExpiresAt:     req.ExpiresAt,
	}))
}
