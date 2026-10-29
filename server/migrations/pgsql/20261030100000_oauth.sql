-- +goose Up
-- OAuth 2.1 для MCP-клиентов: сервер сам выдает токены после входа пользователя

-- Клиенты, зарегистрировавшиеся динамически (RFC 7591). Секрета у них нет: защита — вход пользователя
CREATE TABLE tasks.oauth_clients
(
    id            UUID        NOT NULL PRIMARY KEY, -- client_id
    name          VARCHAR     NOT NULL,             -- client_name из регистрации, показывается на форме входа
    redirect_uris JSONB       NOT NULL,             -- Разрешенные redirect_uri, JSON-массив строк
    created_at    TIMESTAMPTZ NOT NULL,
    last_used_at  TIMESTAMPTZ NOT NULL              -- Последняя авторизация, для чистки заброшенных клиентов
);
COMMENT ON TABLE tasks.oauth_clients IS 'OAuth-клиенты MCP, зарегистрированные динамически';

-- Одноразовые authorization code с PKCE. Живут 10 минут, после обмена помечаются used_at
CREATE TABLE tasks.oauth_codes
(
    code           VARCHAR     NOT NULL PRIMARY KEY,                                -- Случайный код
    client_id      UUID        NOT NULL REFERENCES tasks.oauth_clients (id),
    user_id        UUID        NOT NULL REFERENCES tasks.users (id),
    redirect_uri   TEXT        NOT NULL,                                            -- Должен совпасть при обмене
    code_challenge VARCHAR     NOT NULL,                                            -- PKCE S256: BASE64URL(SHA256(code_verifier))
    resource       TEXT        NOT NULL,                                            -- Канонический URL MCP, для которого выдан код
    expires_at     TIMESTAMPTZ NOT NULL,
    used_at        TIMESTAMPTZ
);
COMMENT ON TABLE tasks.oauth_codes IS 'Одноразовые authorization code OAuth с PKCE';
CREATE INDEX idx_oauth_codes_expires_at ON tasks.oauth_codes (expires_at);

-- +goose Down
DROP TABLE tasks.oauth_codes;
DROP TABLE tasks.oauth_clients;
