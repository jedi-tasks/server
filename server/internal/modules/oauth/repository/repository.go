// Package repository — хранение OAuth-клиентов и authorization code
package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/modules/oauth/repository")

// OAuthRepository — репозиторий OAuth
type OAuthRepository struct {
	db *sql.DB
}

// NewOAuthRepository создает репозиторий OAuth
func NewOAuthRepository(db *sql.DB) *OAuthRepository {
	return &OAuthRepository{db: db}
}
