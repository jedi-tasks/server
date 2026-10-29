package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/modules/priority/repository")

type PriorityRepository struct {
	db *sql.DB
}

func NewPriorityRepository(db *sql.DB) *PriorityRepository {
	return &PriorityRepository{
		db: db,
	}
}
