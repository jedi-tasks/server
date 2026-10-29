package service

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	blockerModel "server/internal/modules/blocker/model"
	durationOptionModel "server/internal/modules/durationOption/model"
	listSectionModel "server/internal/modules/listSection/model"
	priorityModel "server/internal/modules/priority/model"
	"server/internal/modules/sync/model"
	tagModel "server/internal/modules/tag/model"
	taskModel "server/internal/modules/task/model"
	taskListModel "server/internal/modules/taskList/model"
	userModel "server/internal/modules/user/model"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

// hydrate раскладывает последние события аудит-лога по сущностям: удалённые объекты собираются
// в списки идентификаторов напрямую, а созданные/изменённые - довычитываются актуальными из своих
// доменных сервисов (а не восстанавливаются из слепка аудит-лога)
func (s *SyncService) hydrate(ctx context.Context, userID uuid.UUID, latest map[entityKey]auditLogModel.AuditLog, res *model.SyncRes) (err error) {
	changed := make(map[auditLogEntity.AuditLogEntity][]uuid.UUID)
	var addedLinks []tagModel.TaskToTag
	var hasUserChange, hasSettingsChange bool

	for key, auditLog := range latest {
		entity := auditLogEntity.AuditLogEntity(key.entity)

		switch entity {
		case auditLogEntity.User:
			// Пользователь не удаляется - метод всегда create/update
			hasUserChange = true
			continue

		case auditLogEntity.Settings:
			// Настройки не удаляются - метод всегда update
			hasSettingsChange = true
			continue

		case auditLogEntity.TaskToTag:
			// Идентификатор связи - "taskID:tagID"
			link, err := parseTaskToTag(key.entityID)
			if err != nil {
				return errors.InternalServer.Wrap(err).WithContextParams(ctx)
			}
			if auditLog.Method == auditLogMethod.Delete {
				link.UpdatedAt = auditLog.DatetimeCreate
				res.RemovedTasksToTags = append(res.RemovedTasksToTags, link)
			} else {
				addedLinks = append(addedLinks, link)
			}
			continue
		}

		id, err := uuid.Parse(key.entityID)
		if err != nil {
			return errors.InternalServer.Wrap(err).WithContextParams(ctx)
		}

		switch entity {
		case auditLogEntity.TaskList:
			if auditLog.Method == auditLogMethod.Delete {
				res.DeletedTaskListIDs = append(res.DeletedTaskListIDs, id)
			} else {
				changed[entity] = append(changed[entity], id)
			}

		case auditLogEntity.ListSection:
			if auditLog.Method == auditLogMethod.Delete {
				res.DeletedListSectionIDs = append(res.DeletedListSectionIDs, id)
			} else {
				changed[entity] = append(changed[entity], id)
			}

		case auditLogEntity.Priority:
			if auditLog.Method == auditLogMethod.Delete {
				res.DeletedPriorityIDs = append(res.DeletedPriorityIDs, id)
			} else {
				changed[entity] = append(changed[entity], id)
			}

		case auditLogEntity.DurationOption:
			if auditLog.Method == auditLogMethod.Delete {
				res.DeletedDurationOptionIDs = append(res.DeletedDurationOptionIDs, id)
			} else {
				changed[entity] = append(changed[entity], id)
			}

		case auditLogEntity.Task:
			if auditLog.Method == auditLogMethod.Delete {
				res.DeletedTaskIDs = append(res.DeletedTaskIDs, id)
			} else {
				changed[entity] = append(changed[entity], id)
			}

		case auditLogEntity.Tag:
			if auditLog.Method == auditLogMethod.Delete {
				res.DeletedTagIDs = append(res.DeletedTagIDs, id)
			} else {
				changed[entity] = append(changed[entity], id)
			}

		case auditLogEntity.Blocker:
			if auditLog.Method == auditLogMethod.Delete {
				res.DeletedBlockerIDs = append(res.DeletedBlockerIDs, id)
			} else {
				changed[entity] = append(changed[entity], id)
			}
		}
	}

	necessaryInfo := necessary.NecessaryUserInformation{UserID: userID} //nolint:exhaustruct

	if ids := changed[auditLogEntity.TaskList]; len(ids) != 0 {
		if res.ChangedTaskLists, err = s.taskListService.GetTaskLists(ctx, taskListModel.GetTaskListsReq{ //nolint:exhaustruct
			Necessary: necessaryInfo,
			IDs:       ids,
		}); err != nil {
			return err
		}
	}

	if ids := changed[auditLogEntity.ListSection]; len(ids) != 0 {
		if res.ChangedListSections, err = s.listSectionService.GetListSections(ctx, listSectionModel.GetListSectionsReq{ //nolint:exhaustruct
			Necessary: necessaryInfo,
			IDs:       ids,
		}); err != nil {
			return err
		}
	}

	if ids := changed[auditLogEntity.Priority]; len(ids) != 0 {
		if res.ChangedPriorities, err = s.priorityService.GetPriorities(ctx, priorityModel.GetPrioritiesReq{ //nolint:exhaustruct
			Necessary: necessaryInfo,
			IDs:       ids,
		}); err != nil {
			return err
		}
	}

	if ids := changed[auditLogEntity.DurationOption]; len(ids) != 0 {
		if res.ChangedDurationOptions, err = s.durationOptionService.GetDurationOptions(ctx, durationOptionModel.GetDurationOptionsReq{ //nolint:exhaustruct
			Necessary: necessaryInfo,
			IDs:       ids,
		}); err != nil {
			return err
		}
	}

	if ids := changed[auditLogEntity.Task]; len(ids) != 0 {
		if res.ChangedTasks, err = s.taskService.GetTasks(ctx, taskModel.GetTasksReq{ //nolint:exhaustruct
			Necessary: necessaryInfo,
			IDs:       ids,
		}); err != nil {
			return err
		}
	}

	if ids := changed[auditLogEntity.Tag]; len(ids) != 0 {
		if res.ChangedTags, err = s.tagService.GetTags(ctx, tagModel.GetTagsReq{ //nolint:exhaustruct
			Necessary: necessaryInfo,
			IDs:       ids,
		}); err != nil {
			return err
		}
	}

	if ids := changed[auditLogEntity.Blocker]; len(ids) != 0 {
		if res.ChangedBlockers, err = s.blockerService.GetBlockers(ctx, blockerModel.GetBlockersReq{ //nolint:exhaustruct
			Necessary: necessaryInfo,
			IDs:       ids,
		}); err != nil {
			return err
		}
	}

	// Связи задач и тегов довычитываем актуальными: связь могла быть отвязана позже в той же транзакции
	if len(addedLinks) != 0 {
		taskIDs := make([]uuid.UUID, 0, len(addedLinks))
		for _, link := range addedLinks {
			taskIDs = append(taskIDs, link.TaskID)
		}
		links, err := s.tagLinkService.GetTasksToTags(ctx, tagModel.GetTasksToTagsReq{ //nolint:exhaustruct
			Necessary: necessaryInfo,
			TaskIDs:   taskIDs,
		})
		if err != nil {
			return err
		}
		wanted := make(map[string]bool, len(addedLinks))
		for _, link := range addedLinks {
			wanted[link.EntityID()] = true
		}
		for _, link := range links {
			if wanted[link.EntityID()] {
				res.AddedTasksToTags = append(res.AddedTasksToTags, link)
			}
		}
	}

	if hasSettingsChange {
		settings, err := s.settingsService.GetSettings(ctx, userID)
		if err != nil {
			return err
		}
		res.ChangedSettings = &settings
	}

	if hasUserChange {
		users, err := s.userService.GetUsers(ctx, userModel.GetUsersReq{ //nolint:exhaustruct
			Necessary: necessaryInfo,
			IDs:       []uuid.UUID{userID},
		})
		if err != nil {
			return err
		}
		if len(users) != 0 {
			res.ChangedUser = &users[0]
		}
	}

	return nil
}

// parseTaskToTag разбирает идентификатор связи задачи и тега из аудит-лога
func parseTaskToTag(entityID string) (tagModel.TaskToTag, error) {
	var link tagModel.TaskToTag

	taskIDString, tagIDString, ok := strings.Cut(entityID, ":")
	if !ok {
		return link, errors.InternalServer.New("Invalid task to tag entity id")
	}

	taskID, err := uuid.Parse(taskIDString)
	if err != nil {
		return link, err
	}

	tagID, err := uuid.Parse(tagIDString)
	if err != nil {
		return link, err
	}

	link.TaskID = taskID
	link.TagID = tagID
	return link, nil
}
