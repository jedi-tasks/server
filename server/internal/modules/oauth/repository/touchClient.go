package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/oauth/repository/oauthClientDDL"
)

// TouchClient отмечает, что клиент только что авторизовался
func (r *OAuthRepository) TouchClient(ctx context.Context, id uuid.UUID, now time.Time) error {
	ctx, span := tracer.Start(ctx, "TouchClient")
	defer span.End()

	return r.db.Exec(ctx, sq.Update(oauthClientDDL.Table).
		Set(oauthClientDDL.ColumnLastUsedAt, now).
		Where(sq.Eq{oauthClientDDL.ColumnID: id}))
}
