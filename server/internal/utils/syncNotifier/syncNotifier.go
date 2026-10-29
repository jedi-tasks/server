// Package syncNotifier — потокобезопасный in-memory pub/sub: уведомляет активные
// SubscribeToSync-стримы пользователя о том, что для него появились новые изменения, не
// передавая сам payload — получатель сам вызовет обычный Sync/ConfirmSync (см. sync-endpoint.proto,
// SubscribeToSync). Ничего не знает про домен (uuid.UUID — единственный тип), поэтому не создаёт
// циклов импорта между модулями auditLog (пишет уведомления) и sync (читает их).
//
// Одна нода сервера — без внешнего брокера (Redis и т.п.). Если сервер станет многоинстансным,
// это единственное место, которое придётся заменить на pub/sub через шину.
package syncNotifier

import (
	"sync"

	"github.com/google/uuid"
)

type SyncNotifier struct {
	mu   sync.Mutex
	subs map[uuid.UUID]map[chan struct{}]struct{}
}

func New() *SyncNotifier {
	return &SyncNotifier{
		subs: make(map[uuid.UUID]map[chan struct{}]struct{}),
	}
}

// Subscribe регистрирует нового подписчика на уведомления для пользователя. Вызывающий обязан
// вызвать unsubscribe() (обычно defer сразу после Subscribe), когда стрим завершается — иначе
// канал останется висеть в мапе вечно.
func (n *SyncNotifier) Subscribe(userID uuid.UUID) (ch chan struct{}, unsubscribe func()) {
	// Буфер 1 — если получатель ещё не забрал предыдущий сигнал, новый Notify не блокируется и
	// не теряется: "есть новое" уже стоит в канале, второй раз сообщать нечего.
	ch = make(chan struct{}, 1)

	n.mu.Lock()
	if n.subs[userID] == nil {
		n.subs[userID] = make(map[chan struct{}]struct{})
	}
	n.subs[userID][ch] = struct{}{}
	n.mu.Unlock()

	return ch, func() {
		n.mu.Lock()
		defer n.mu.Unlock()
		delete(n.subs[userID], ch)
		if len(n.subs[userID]) == 0 {
			delete(n.subs, userID)
		}
	}
}

// Notify будит все активные стримы пользователя (обычно 0 или 1 — все его открытые сейчас
// устройства). Неблокирующий: если канал уже "полон" (сигнал ещё не забрали), просто пропускаем.
func (n *SyncNotifier) Notify(userID uuid.UUID) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for ch := range n.subs[userID] {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

// NotifyMany — то же самое сразу для набора пользователей (например, всех участников группы
// счетов, в которой произошло изменение).
func (n *SyncNotifier) NotifyMany(userIDs []uuid.UUID) {
	for _, id := range userIDs {
		n.Notify(id)
	}
}
