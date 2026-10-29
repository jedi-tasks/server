package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/modules/blocker/model"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

// checkCycle проверяет, что новый блокер «blockingTaskID блокирует taskID» не замыкает цепочку:
// blockingTaskID не должна быть сама заблокирована (прямо или по цепочке) задачей taskID
func (s *BlockerService) checkCycle(ctx context.Context, necessary necessary.NecessaryUserInformation, taskID, blockingTaskID uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "checkCycle")
	defer span.End()

	// Получаем все живые блокеры-задачи пользователя
	blockers, err := s.blockerRepository.GetBlockers(ctx, model.GetBlockersReq{ //nolint:exhaustruct
		Necessary:  necessary,
		OnlyActive: true,
	})
	if err != nil {
		return err
	}

	// Строим граф «задача → чем она заблокирована»
	blockedBy := make(map[uuid.UUID][]uuid.UUID)
	for _, blocker := range blockers {
		if blocker.BlockingTaskID != nil {
			blockedBy[blocker.TaskID] = append(blockedBy[blocker.TaskID], *blocker.BlockingTaskID)
		}
	}

	// Обходим цепочку от блокирующей задачи: если дошли до taskID — это цикл
	visited := map[uuid.UUID]bool{}
	stack := []uuid.UUID{blockingTaskID}
	for len(stack) != 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if current == taskID {
			return errors.BadRequest.New("Blocker creates a cycle").
				WithContextParams(ctx).
				WithParams("TaskID", taskID, "BlockingTaskID", blockingTaskID).
				WithCustomHumanText("Задачи не могут блокировать друг друга по кругу")
		}
		if visited[current] {
			continue
		}
		visited[current] = true
		stack = append(stack, blockedBy[current]...)
	}

	return nil
}
