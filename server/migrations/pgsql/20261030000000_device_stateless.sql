-- +goose Up
-- Устройство без локальной копии данных (MCP, веб): всегда читает сервер перед записью,
-- поэтому проверка синхронизации перед правками для него не нужна
ALTER TABLE tasks.devices ADD COLUMN is_stateless BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE tasks.devices DROP COLUMN is_stateless;
