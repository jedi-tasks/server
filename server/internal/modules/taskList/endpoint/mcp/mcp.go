// Package mcp — инструменты MCP для списков задач, один к одному с методами TaskListEndpoint
package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/taskList/model"
	taskListService "server/internal/modules/taskList/service"
)

var _ TaskListService = new(taskListService.TaskListService)

// TaskListService — интерфейс бизнес-логики списков
type TaskListService interface {
	GetTaskLists(context.Context, model.GetTaskListsReq) ([]model.TaskList, error)
	CreateTaskList(context.Context, model.CreateTaskListReq) error
	UpdateTaskList(context.Context, model.UpdateTaskListReq) error
	DeleteTaskList(context.Context, model.DeleteTaskListReq) error
}

// TaskListTools — обработчики инструментов списков
type TaskListTools struct {
	taskListService TaskListService
}

// Register добавляет инструменты списков на MCP-сервер
func Register(server *sdk.Server, taskListService TaskListService) {
	tools := &TaskListTools{taskListService: taskListService}

	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "get_task_lists",
		Description: "Возвращает списки задач пользователя. Ровно один список — системный Inbox (isInbox)",
		Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}, //nolint:exhaustruct
	}, tools.GetTaskLists)
	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "create_task_list",
		Description: "Создает список задач. sortMode: manual, priority, plannedDay, deadline, duration",
	}, tools.CreateTaskList)
	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "update_task_list",
		Description: "Меняет только переданные поля списка",
	}, tools.UpdateTaskList)
	sdk.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "delete_task_list",
		Description: "Удаляет список. Inbox удалить нельзя. Отклоняется, пока в списке есть живые задачи или разделы",
	}, tools.DeleteTaskList)
}
