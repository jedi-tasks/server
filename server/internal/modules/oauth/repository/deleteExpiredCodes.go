package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/oauth/repository/oauthCodeDDL"
)

// DeleteExpiredCodes удаляет просроченные и использованные коды старше момента before
func (r *OAuthRepository) DeleteExpiredCodes(ctx context.Context, before time.Time) error {
	ctx, span := tracer.Start(ctx, "DeleteExpiredCodes")
	defer span.End()

	return r.db.Exec(ctx, sq.Delete(oauthCodeDDL.Table).Where(sq.Lt{oauthCodeDDL.ColumnExpiresAt: before}))
}
