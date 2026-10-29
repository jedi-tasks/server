package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/modules/taskList/repository")

type TaskListRepository struct {
	db *sql.DB
}

func NewTaskListRepository(db *sql.DB) *TaskListRepository {
	return &TaskListRepository{
		db: db,
	}
}
