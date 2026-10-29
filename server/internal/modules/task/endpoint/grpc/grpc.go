package grpc

import (
	"context"

	"server/internal/modules/task/model"
	taskService "server/internal/modules/task/service"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

var _ TaskService = new(taskService.TaskService)

type TaskService interface {
	GetTasks(context.Context, model.GetTasksReq) ([]model.Task, error)
	CreateTask(context.Context, model.CreateTaskReq) error
	UpdateTask(context.Context, model.UpdateTaskReq) error
	DeleteTask(context.Context, model.DeleteTaskReq) error
	SetTaskCompleted(context.Context, model.SetTaskCompletedReq) error
}

var _ proto.TaskEndpointServer = new(TaskEndpoint)

type TaskEndpoint struct {
	proto.UnsafeTaskEndpointServer
	taskService TaskService
}

func NewTaskEndpoint(taskService TaskService) *TaskEndpoint {
	return &TaskEndpoint{
		taskService: taskService,
	}
}
