package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/utils/errors"
)

// checkNoLiveContents отклоняет удаление списка, если в нем остались неудаленные задачи или разделы
func (s *TaskListService) checkNoLiveContents(ctxTx context.Context, listID uuid.UUID) error {

	// Проверяем задачи списка
	hasTasks, err := s.liveTasksChecker.HasLiveTasksInList(ctxTx, listID)
	if err != nil {
		return err
	}
	if hasTasks {
		return errors.BadRequest.New("Task list has live tasks").
			WithContextParams(ctxTx).
			WithParams("ListID", listID).
			WithCustomHumanText("В списке есть задачи — сначала удалите их или перенесите в другой список")
	}

	// Проверяем разделы списка
	hasSections, err := s.liveSectionsChecker.HasLiveSectionsInList(ctxTx, listID)
	if err != nil {
		return err
	}
	if hasSections {
		return errors.BadRequest.New("Task list has live sections").
			WithContextParams(ctxTx).
			WithParams("ListID", listID).
			WithCustomHumanText("В списке есть разделы — сначала удалите их")
	}

	return nil
}
