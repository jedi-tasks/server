package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/auditLog/repository")

// AuditLogRepository - репозиторий аудит-лога
type AuditLogRepository struct {
	db *sql.DB
}

// NewAuditLogRepository создает новый репозиторий аудит-лога
func NewAuditLogRepository(db *sql.DB) *AuditLogRepository {
	return &AuditLogRepository{
		db: db,
	}
}
