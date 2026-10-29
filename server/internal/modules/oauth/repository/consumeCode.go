package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/oauth/model"
	repoModel "server/internal/modules/oauth/repository/model"
	"server/internal/modules/oauth/repository/oauthCodeDDL"
)

// ConsumeCode атомарно помечает код использованным и возвращает его. found=false — кода нет,
// он просрочен или уже обменян: повторный обмен того же кода не пройдет
func (r *OAuthRepository) ConsumeCode(ctx context.Context, code string, now time.Time) (res model.Code, found bool, err error) {
	ctx, span := tracer.Start(ctx, "ConsumeCode")
	defer span.End()

	var rows []repoModel.Code
	err = r.db.Select(ctx, &rows, sq.Update(oauthCodeDDL.Table).
		Set(oauthCodeDDL.ColumnUsedAt, now).
		Where(sq.Eq{oauthCodeDDL.ColumnCode: code, oauthCodeDDL.ColumnUsedAt: nil}).
		Where(sq.Gt{oauthCodeDDL.ColumnExpiresAt: now}).
		Suffix("RETURNING *"))
	if err != nil || len(rows) == 0 {
		return res, false, err
	}
	return rows[0].ConvertToModel(), true, nil
}
