package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/utils/errors"
)

// checkNoLiveLinks отклоняет удаление задачи, если к ней привязаны теги или она участвует в блокерах.
// Строка задачи уже заблокирована удалением, поэтому параллельная привязка дождется этой транзакции
func (s *TaskService) checkNoLiveLinks(ctxTx context.Context, taskID uuid.UUID) error {

	// Проверяем связи с тегами
	hasTagLinks, err := s.liveTagLinksChecker.HasLiveLinksForTask(ctxTx, taskID)
	if err != nil {
		return err
	}
	if hasTagLinks {
		return errors.BadRequest.New("Task has live tag links").
			WithContextParams(ctxTx).
			WithParams("TaskID", taskID).
			WithCustomHumanText("К задаче привязаны теги — сначала отвяжите их")
	}

	// Проверяем блокеры, где задача заблокирована или блокирует
	hasBlockers, err := s.liveBlockersChecker.HasLiveBlockersForTask(ctxTx, taskID)
	if err != nil {
		return err
	}
	if hasBlockers {
		return errors.BadRequest.New("Task has live blockers").
			WithContextParams(ctxTx).
			WithParams("TaskID", taskID).
			WithCustomHumanText("Задача участвует в блокерах — сначала удалите их")
	}

	return nil
}
