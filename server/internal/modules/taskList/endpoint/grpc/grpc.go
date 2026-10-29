package grpc

import (
	"context"

	"server/internal/modules/taskList/model"
	taskListService "server/internal/modules/taskList/service"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

var _ TaskListService = new(taskListService.TaskListService)

type TaskListService interface {
	GetTaskLists(context.Context, model.GetTaskListsReq) ([]model.TaskList, error)
	CreateTaskList(context.Context, model.CreateTaskListReq) error
	UpdateTaskList(context.Context, model.UpdateTaskListReq) error
	DeleteTaskList(context.Context, model.DeleteTaskListReq) error
}

var _ proto.TaskListEndpointServer = new(TaskListEndpoint)

type TaskListEndpoint struct {
	proto.UnsafeTaskListEndpointServer
	taskListService TaskListService
}

func NewTaskListEndpoint(taskListService TaskListService) *TaskListEndpoint {
	return &TaskListEndpoint{
		taskListService: taskListService,
	}
}
