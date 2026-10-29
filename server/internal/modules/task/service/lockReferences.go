package service

import (
	"context"

	"github.com/google/uuid"
)

// lockReferences внутри открытой транзакции проверяет, что список, раздел и приоритет задачи доступны
// пользователю и не удалены, и блокирует их от удаления. nil — поле не передано, uuid.Nil у раздела и
// приоритета — ссылка сброшена, такие поля не проверяются
func (s *TaskService) lockReferences(ctxTx context.Context, userID uuid.UUID, listID, sectionID, priorityID *uuid.UUID) error {

	// Проверяем список
	if listID != nil {
		if err := s.taskListService.LockLive(ctxTx, userID, []uuid.UUID{*listID}); err != nil {
			return err
		}
	}

	// Проверяем раздел
	if sectionID != nil && *sectionID != uuid.Nil {
		if err := s.listSectionService.LockLive(ctxTx, userID, []uuid.UUID{*sectionID}); err != nil {
			return err
		}
	}

	// Проверяем приоритет
	if priorityID != nil && *priorityID != uuid.Nil {
		if err := s.priorityService.LockLive(ctxTx, userID, []uuid.UUID{*priorityID}); err != nil {
			return err
		}
	}

	return nil
}
