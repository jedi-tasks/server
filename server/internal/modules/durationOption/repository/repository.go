package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/modules/durationOption/repository")

type DurationOptionRepository struct {
	db *sql.DB
}

func NewDurationOptionRepository(db *sql.DB) *DurationOptionRepository {
	return &DurationOptionRepository{
		db: db,
	}
}
