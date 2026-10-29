package auditLogEntity

import (
	"context"

	"server/internal/utils/errors"
)

// AuditLogEntity - тип сущности, изменение которой зафиксировано в аудит-логе
type AuditLogEntity string

// enum:"taskList,listSection,priority,durationOption,task,tag,taskToTag,blocker,settings,user"
const (
	TaskList       = AuditLogEntity("taskList")
	ListSection    = AuditLogEntity("listSection")
	Priority       = AuditLogEntity("priority")
	DurationOption = AuditLogEntity("durationOption")
	Task           = AuditLogEntity("task")
	Tag            = AuditLogEntity("tag")
	TaskToTag      = AuditLogEntity("taskToTag")
	Blocker        = AuditLogEntity("blocker")
	Settings       = AuditLogEntity("settings")
	User           = AuditLogEntity("user")
)

// Validate проверяет, что тип сущности аудит-лога принадлежит известному набору значений
func (e AuditLogEntity) Validate(ctx context.Context) error {
	switch e {
	case TaskList, ListSection, Priority, DurationOption, Task, Tag, TaskToTag, Blocker, Settings, User:
	default:
		return errors.BadRequest.New("Unknown audit log entity").
			WithContextParams(ctx).
			WithParams("entity", e)
	}
	return nil
}
