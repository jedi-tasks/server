package service

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"server/internal/modules/blocker/model"
	"server/internal/utils/necessary"
)

// fakeBlockerRepository — подмена репозитория, отдающая заранее заданные блокеры. Остальные методы
// интерфейса в checkCycle не вызываются
type fakeBlockerRepository struct {
	BlockerRepository                 // Неиспользуемые методы, вызов приведет к панике
	blockers          []model.Blocker // Живые блокеры пользователя
}

// GetBlockers возвращает заданные блокеры
func (r fakeBlockerRepository) GetBlockers(context.Context, model.GetBlockersReq) ([]model.Blocker, error) {
	return r.blockers, nil
}

// blocks собирает блокер «blockingTaskID блокирует taskID»
func blocks(taskID, blockingTaskID uuid.UUID) model.Blocker {
	return model.Blocker{TaskID: taskID, BlockingTaskID: &blockingTaskID} //nolint:exhaustruct
}

// TestCheckCycle проверяет поиск циклов в графе блокеров
func TestCheckCycle(t *testing.T) {
	a, b, c, d := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	reason := "ждем ответа"

	testCases := []struct {
		name           string          // Название сценария
		blockers       []model.Blocker // Уже существующие блокеры
		taskID         uuid.UUID       // Заблокированная задача нового блокера
		blockingTaskID uuid.UUID       // Блокирующая задача нового блокера
		wantError      bool            // Ожидается отказ
	}{
		{name: "без блокеров", blockers: nil, taskID: a, blockingTaskID: b, wantError: false},
		{name: "прямой цикл", blockers: []model.Blocker{blocks(b, a)}, taskID: a, blockingTaskID: b, wantError: true},
		{name: "цикл через цепочку", blockers: []model.Blocker{blocks(b, c), blocks(c, a)}, taskID: a, blockingTaskID: b, wantError: true},
		{name: "цепочка без цикла", blockers: []model.Blocker{blocks(b, c), blocks(c, d)}, taskID: a, blockingTaskID: b, wantError: false},
		{name: "общий блокирующий не цикл", blockers: []model.Blocker{blocks(a, c), blocks(b, c)}, taskID: a, blockingTaskID: b, wantError: false},
		{name: "ромб без цикла", blockers: []model.Blocker{blocks(b, c), blocks(b, d), blocks(c, d)}, taskID: a, blockingTaskID: b, wantError: false},
		{name: "текстовая причина не ребро", blockers: []model.Blocker{{TaskID: b, Reason: &reason}}, taskID: a, blockingTaskID: b, wantError: false}, //nolint:exhaustruct
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := &BlockerService{blockerRepository: fakeBlockerRepository{blockers: testCase.blockers}} //nolint:exhaustruct
			err := service.checkCycle(context.Background(), necessary.NecessaryUserInformation{}, testCase.taskID, testCase.blockingTaskID)
			if (err != nil) != testCase.wantError {
				t.Fatalf("ожидался отказ: %v, получено: %v", testCase.wantError, err)
			}
		})
	}
}
