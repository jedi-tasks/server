// Package mcp — инструменты MCP для задач, один к одному с методами TaskEndpoint
package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/task/model"
	taskService "server/internal/modules/task/service"
)

var _ TaskService = new(taskService.TaskService)

// TaskService — интерфейс бизнес-логики задач
type TaskService interface {
	GetTasks(context.Context, model.GetTasksReq) ([]model.Task, error)
	CreateTask(context.Context, model.CreateTaskReq) error
	UpdateTask(context.Context, model.UpdateTaskReq) error
	DeleteTask(context.Context, model.DeleteTaskReq) error
	SetTaskCompleted(context.Context, model.SetTaskCompletedReq) error
}

// TaskTools — обработчики инструментов задач
type TaskTools struct {
	taskService TaskService
}

// Register добавляет инструменты задач на MCP-сервер
func Register(server *sdk.Server, taskService TaskService) {
	tools := &TaskTools{taskService: taskService}
	readOnly := &sdk.ToolAnnotations{ReadOnlyHint: true} //nolint:exhaustruct

	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "get_tasks",
		Description: "Возвращает задачи пользователя. Фильтры по спискам, идентификаторам и выполненности; без фильтров — все задачи, включая выполненные",
		Annotations: readOnly,
	}, tools.GetTasks)
	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "create_task",
		Description: "Создает задачу. Без listID задача попадает в Inbox. Без priorityID задача считается «не разобранной»",
	}, tools.CreateTask)
	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "update_task",
		Description: "Меняет только переданные поля задачи. Сброс: sectionID и priorityID — пустая строка, plannedDay и deadlineDay — пустая строка, durationMinutes — 0",
	}, tools.UpdateTask)
	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "set_task_completed",
		Description: "Отмечает задачу выполненной (completed=true) или возвращает в работу (false). Время выполнения — момент вызова. При выполнении сервер снимает блокеры, которые ссылались на эту задачу как на блокирующую",
	}, tools.SetTaskCompleted)
	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "delete_task",
		Description: "Удаляет задачу. Отклоняется, пока у задачи есть живые связи с тегами или блокеры: сначала remove_tag_from_task и resolve_blocker",
	}, tools.DeleteTask)
}
