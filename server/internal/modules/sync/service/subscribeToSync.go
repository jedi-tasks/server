package service

import (
	"context"

	"github.com/google/uuid"
)

// SubscribeToSync блокируется и вызывает notify() каждый раз, когда для пользователя появились
// изменения (см. syncNotifier) — до отмены ctx (клиент отключился) или до первой ошибки notify().
// Сам не несёт payload изменений — это чисто сигнал "дёрни Sync" (см. SubscribeToSync в
// sync-endpoint.proto).
func (s *SyncService) SubscribeToSync(ctx context.Context, userID uuid.UUID, notify func() error) error {
	ctx, span := tracer.Start(ctx, "SubscribeToSync")
	defer span.End()

	ch, unsubscribe := s.syncNotifier.Subscribe(userID)
	defer unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
			if err := notify(); err != nil {
				return err
			}
		}
	}
}
