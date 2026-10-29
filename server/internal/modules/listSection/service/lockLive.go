package service

import (
	"context"

	"github.com/google/uuid"
)

// LockLive внутри открытой транзакции проверяет, что на записи можно сослаться (они принадлежат
// пользователю и не удалены), и блокирует их от удаления до конца транзакции
func (s *ListSectionService) LockLive(ctxTx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	ctxTx, span := tracer.Start(ctxTx, "LockLive")
	defer span.End()

	return s.listSectionRepository.LockLive(ctxTx, userID, ids)
}
