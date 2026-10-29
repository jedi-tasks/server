package auditLogMethod

import (
	"context"

	"server/internal/utils/errors"
)

// AuditLogMethod - тип изменяющего действия, зафиксированного в аудит-логе
type AuditLogMethod string

// enum:"create,update,delete"
const (
	Create = AuditLogMethod("create")
	Update = AuditLogMethod("update")
	Delete = AuditLogMethod("delete")
)

// Validate проверяет, что метод аудит-лога принадлежит известному набору значений
func (m AuditLogMethod) Validate(ctx context.Context) error {
	switch m {
	case Create, Update, Delete:
	default:
		return errors.BadRequest.New("Unknown audit log method").
			WithContextParams(ctx).
			WithParams("method", m)
	}
	return nil
}
