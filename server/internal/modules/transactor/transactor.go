package transactor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/modules/transactor/service")

type Transactor struct {
	db *sql.DB
}

func NewTransactor(db *sql.DB) *Transactor {
	return &Transactor{
		db: db,
	}
}

// WithinTransaction принимает коллбэк, который будет выполнен в рамках транзакции
func (r *Transactor) WithinTransaction(ctx context.Context, callback func(ctx context.Context) error) error {
	ctx, span := tracer.Start(ctx, "WithinTransaction")
	defer span.End()

	// Запускаем транзакцию
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// Запускаем коллбэк
	err = callback(sql.InjectTx(ctx, tx))
	if err != nil {
		// Если произошла ошибка, откатываем изменения
		_ = tx.Rollback()
		return err
	}
	// Если ошибок нет, подтверждаем изменения
	return tx.Commit()
}
