package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/task/model"
	"server/internal/modules/task/repository/taskDDL"
	"server/internal/utils/necessary"
)

// DeleteByListID мягко удаляет все задачи списка внутри открытой транзакции удаления списка
func (s *TaskService) DeleteByListID(ctxTx context.Context, necessary necessary.NecessaryUserInformation, listID uuid.UUID, updatedAt time.Time) (lastAuditLogID uint32, err error) {
	ctxTx, span := tracer.Start(ctxTx, "DeleteByListID")
	defer span.End()

	tasks, err := s.taskRepository.GetTasks(ctxTx, model.GetTasksReq{ //nolint:exhaustruct
		Necessary: necessary,
		ListIDs:   []uuid.UUID{listID},
	})
	if err != nil {
		return 0, err
	}

	for _, task := range tasks {
		auditLogID, err := s.deleteTask(ctxTx, model.DeleteTaskReq{
			Necessary: necessary,
			ID:        task.ID,
			UpdatedAt: updatedAt,
		})
		if err != nil {
			return 0, err
		}
		lastAuditLogID = max(lastAuditLogID, auditLogID)
	}

	return lastAuditLogID, nil
}

// DetachSection убирает удаленный раздел из задач внутри открытой транзакции удаления раздела.
// Задачи остаются в списке без раздела
func (s *TaskService) DetachSection(ctxTx context.Context, necessary necessary.NecessaryUserInformation, sectionID uuid.UUID, updatedAt time.Time) (uint32, error) {
	return s.clearReference(ctxTx, necessary, taskDDL.ColumnSectionID, sectionID, updatedAt)
}

// DetachPriority убирает удаленный приоритет из задач внутри открытой транзакции удаления приоритета.
// Задачи становятся «не разобрано»
func (s *TaskService) DetachPriority(ctxTx context.Context, necessary necessary.NecessaryUserInformation, priorityID uuid.UUID, updatedAt time.Time) (uint32, error) {
	return s.clearReference(ctxTx, necessary, taskDDL.ColumnPriorityID, priorityID, updatedAt)
}

// clearReference сбрасывает ссылку у задач и фиксирует каждое изменение в аудит-логе
func (s *TaskService) clearReference(ctxTx context.Context, necessary necessary.NecessaryUserInformation, column string, refID uuid.UUID, updatedAt time.Time) (lastAuditLogID uint32, err error) {
	ctxTx, span := tracer.Start(ctxTx, "clearReference")
	defer span.End()

	// Получаем задачи для слепков "до" в аудит-логе
	befores, err := s.taskRepository.GetTasks(ctxTx, model.GetTasksReq{ //nolint:exhaustruct
		Necessary: necessary,
	})
	if err != nil {
		return 0, err
	}
	beforesMap := slices.ToMap(befores, func(task model.Task) uuid.UUID { return task.ID })

	// Сбрасываем ссылку
	ids, err := s.taskRepository.ClearReference(ctxTx, necessary.UserID, column, refID, updatedAt)
	if err != nil || len(ids) == 0 {
		return 0, err
	}

	// Получаем задачи для слепков "после" в аудит-логе
	afters, err := s.taskRepository.GetTasks(ctxTx, model.GetTasksReq{ //nolint:exhaustruct
		Necessary: necessary,
		IDs:       ids,
	})
	if err != nil {
		return 0, err
	}

	// Фиксируем изменения в аудит-логе
	for _, after := range afters {
		lastAuditLogID, err = s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.Task,
			Method:   auditLogMethod.Update,
			EntityID: after.ID.String(),
			Before:   beforesMap[after.ID],
			After:    after,
			UserID:   necessary.UserID,
			DeviceID: necessary.DeviceID,
		})
		if err != nil {
			return 0, err
		}
	}

	return lastAuditLogID, nil
}
