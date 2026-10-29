package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/modules/listSection/repository")

type ListSectionRepository struct {
	db *sql.DB
}

func NewListSectionRepository(db *sql.DB) *ListSectionRepository {
	return &ListSectionRepository{
		db: db,
	}
}
