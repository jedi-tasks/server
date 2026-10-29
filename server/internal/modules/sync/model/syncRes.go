package model

import (
	"github.com/google/uuid"

	blockerModel "server/internal/modules/blocker/model"
	durationOptionModel "server/internal/modules/durationOption/model"
	listSectionModel "server/internal/modules/listSection/model"
	priorityModel "server/internal/modules/priority/model"
	settingsModel "server/internal/modules/settings/model"
	tagModel "server/internal/modules/tag/model"
	taskModel "server/internal/modules/task/model"
	taskListModel "server/internal/modules/taskList/model"
	userModel "server/internal/modules/user/model"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// SyncRes - результат синхронизации: изменённые и удалённые объекты по каждой сущности
type SyncRes struct {
	PendingCheckpoint uint32     // Новый чекпоинт, который клиент должен подтвердить вызовом ConfirmSync после применения изменений
	PendingSyncToken  *uuid.UUID // Опорный токен этого ответа для ConfirmSync, nil если изменений не было
	HasChanges        bool       // true, если в ответе есть изменения

	ChangedTaskLists   []taskListModel.TaskList // Созданные/изменённые
	DeletedTaskListIDs []uuid.UUID              // Удалённые

	ChangedListSections   []listSectionModel.ListSection // Созданные/изменённые
	DeletedListSectionIDs []uuid.UUID                    // Удалённые

	ChangedPriorities  []priorityModel.Priority // Созданные/изменённые
	DeletedPriorityIDs []uuid.UUID              // Удалённые

	ChangedDurationOptions   []durationOptionModel.DurationOption // Созданные/изменённые
	DeletedDurationOptionIDs []uuid.UUID                          // Удалённые

	ChangedTasks   []taskModel.Task // Созданные/изменённые
	DeletedTaskIDs []uuid.UUID      // Удалённые

	ChangedTags   []tagModel.Tag // Созданные/изменённые
	DeletedTagIDs []uuid.UUID    // Удалённые

	ChangedBlockers   []blockerModel.Blocker // Созданные/изменённые
	DeletedBlockerIDs []uuid.UUID            // Удалённые

	AddedTasksToTags   []tagModel.TaskToTag // Добавленные связи задач и тегов
	RemovedTasksToTags []tagModel.TaskToTag // Удалённые связи задач и тегов

	ChangedSettings *settingsModel.Settings // Изменённые настройки, nil если не менялись
	ChangedUser     *userModel.User         // Изменённые данные пользователя, nil если не менялись
}

// ConvertToProto преобразует SyncRes в proto-формат
func (s *SyncRes) ConvertToProto() (*proto.SyncResponse, error) {
	res := &proto.SyncResponse{ //nolint:exhaustruct
		PendingCheckpoint: s.PendingCheckpoint,
		HasChanges:        s.HasChanges,
	}

	if s.PendingSyncToken != nil {
		res.PendingSyncToken = s.PendingSyncToken[:]
	}

	for _, item := range s.ChangedTaskLists {
		protoItem, err := item.ConvertToProto()
		if err != nil {
			return nil, err
		}
		res.ChangedTaskLists = append(res.ChangedTaskLists, protoItem)
	}
	for _, id := range s.DeletedTaskListIDs {
		res.DeletedTaskListIDs = append(res.DeletedTaskListIDs, id[:])
	}

	for _, item := range s.ChangedListSections {
		protoItem, err := item.ConvertToProto()
		if err != nil {
			return nil, err
		}
		res.ChangedListSections = append(res.ChangedListSections, protoItem)
	}
	for _, id := range s.DeletedListSectionIDs {
		res.DeletedListSectionIDs = append(res.DeletedListSectionIDs, id[:])
	}

	for _, item := range s.ChangedPriorities {
		protoItem, err := item.ConvertToProto()
		if err != nil {
			return nil, err
		}
		res.ChangedPriorities = append(res.ChangedPriorities, protoItem)
	}
	for _, id := range s.DeletedPriorityIDs {
		res.DeletedPriorityIDs = append(res.DeletedPriorityIDs, id[:])
	}

	for _, item := range s.ChangedDurationOptions {
		protoItem, err := item.ConvertToProto()
		if err != nil {
			return nil, err
		}
		res.ChangedDurationOptions = append(res.ChangedDurationOptions, protoItem)
	}
	for _, id := range s.DeletedDurationOptionIDs {
		res.DeletedDurationOptionIDs = append(res.DeletedDurationOptionIDs, id[:])
	}

	for _, item := range s.ChangedTasks {
		protoItem, err := item.ConvertToProto()
		if err != nil {
			return nil, err
		}
		res.ChangedTasks = append(res.ChangedTasks, protoItem)
	}
	for _, id := range s.DeletedTaskIDs {
		res.DeletedTaskIDs = append(res.DeletedTaskIDs, id[:])
	}

	for _, item := range s.ChangedTags {
		protoItem, err := item.ConvertToProto()
		if err != nil {
			return nil, err
		}
		res.ChangedTags = append(res.ChangedTags, protoItem)
	}
	for _, id := range s.DeletedTagIDs {
		res.DeletedTagIDs = append(res.DeletedTagIDs, id[:])
	}

	for _, item := range s.ChangedBlockers {
		protoItem, err := item.ConvertToProto()
		if err != nil {
			return nil, err
		}
		res.ChangedBlockers = append(res.ChangedBlockers, protoItem)
	}
	for _, id := range s.DeletedBlockerIDs {
		res.DeletedBlockerIDs = append(res.DeletedBlockerIDs, id[:])
	}

	for _, link := range s.AddedTasksToTags {
		res.AddedTasksToTags = append(res.AddedTasksToTags, link.ConvertToProto())
	}
	for _, link := range s.RemovedTasksToTags {
		res.RemovedTasksToTags = append(res.RemovedTasksToTags, link.ConvertToProto())
	}

	if s.ChangedSettings != nil {
		protoSettings, err := s.ChangedSettings.ConvertToProto()
		if err != nil {
			return nil, err
		}
		res.ChangedSettings = protoSettings
	}

	if s.ChangedUser != nil {
		protoUser, err := s.ChangedUser.ConvertToProto()
		if err != nil {
			return nil, err
		}
		res.ChangedUser = protoUser
	}

	return res, nil
}
