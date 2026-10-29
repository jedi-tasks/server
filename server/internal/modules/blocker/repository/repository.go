package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/modules/blocker/repository")

type BlockerRepository struct {
	db *sql.DB
}

func NewBlockerRepository(db *sql.DB) *BlockerRepository {
	return &BlockerRepository{
		db: db,
	}
}
