# Tasks — server

Сервер трекера задач jedi-tasks. Собран из сервера Coin: те же слои (`endpoint/service/repository/model`),
авторизация (логин и пароль, JWT), аудит-лог и синк. Контекст предметной области — в репо `ios` (`CONTEXT.md`),
контракты — в [contracts](https://github.com/jedi-tasks/contracts).

## Состав

- `server/internal/modules/*` — модули: `auth`, `user`, `settings`, `taskList`, `listSection`, `priority`,
  `durationOption`, `task`, `tag` (вместе со связями задач и тегов), `blocker`, `auditLog`, `sync`, `transactor`.
- `server/migrations/pgsql` — миграции goose, применяются при старте.
- `server/cmd/smoke` — смоук-тест запущенного сервера.
- `pkg` — сабмодуль [bonavi/pkg](https://github.com/bonavi/pkg), общая библиотека с Coin.

## Синк

- Клиент работает офлайн и отправляет мутации обычными Create/Update/Delete из своей очереди.
- `WithSyncGate`, как в Coin: мутация от устройства, которое не скачало последние изменения, отклоняется с
  `NeedToSync` — клиент делает `Sync`, `ConfirmSync` и повторяет очередь.
- Каждая мутация несёт `updatedAt` с клиента. Запись меняется, только если изменение новее сохранённого
  (побеждает последняя запись). Повторный `Create` с тем же id ничего не делает — повтор из очереди безопасен.
- Каждое применённое изменение пишется в аудит-лог; `Sync` отдаёт всё после чекпоинта устройства.
- Каскады тоже идут через аудит-лог: удаление списка удаляет его задачи и разделы, удаление раздела или
  приоритета снимает их с задач, удаление тега отвязывает его от задач, выполнение задачи снимает блокеры.

## Новый пользователь

Регистрация открыта только при `AUTH_SIGN_UP_ENABLED=true`. При регистрации сервер создаёт Inbox,
три приоритета «Пути джедая» и справочник длительностей — клиент получает их первым синком.

## Локальный запуск

```bash
docker compose -f localEnvironment/docker-compose.yml up -d
```

Переменные окружения — в `server/internal/config/config.go` и `pkg/database/pgsql`. Пока CI контрактов не
опубликовал `go-server-grpc`, сервер собирается на локальной генерации: в `server/go.mod` стоит
`replace github.com/jedi-tasks/go-server-grpc => ../../contracts/go-server-grpc-local`
(`make build-go` в `contracts`).
