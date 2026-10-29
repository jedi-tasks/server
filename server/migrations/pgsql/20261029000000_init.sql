-- +goose Up
CREATE SCHEMA IF NOT EXISTS tasks;

-- Пользователи и устройства — как в Coin
CREATE TABLE tasks.users
(
    id            UUID         NOT NULL PRIMARY KEY, -- Идентификатор пользователя
    name          VARCHAR(100) NOT NULL,             -- Имя пользователя
    email         VARCHAR(255) NOT NULL UNIQUE,      -- Почтовый адрес пользователя
    password_hash BYTEA        NOT NULL,             -- Хэш пароля
    password_salt BYTEA        NOT NULL,             -- Соль пароля
    time_create   TIMESTAMPTZ  NOT NULL              -- Дата создания аккаунта
);
COMMENT ON TABLE tasks.users IS 'Пользователи';

CREATE TABLE tasks.devices
(
    id                                UUID         NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                           UUID         NOT NULL REFERENCES tasks.users (id),
    device_id                         VARCHAR(100) NOT NULL,           -- Идентификатор устройства от клиента
    refresh_token                     TEXT         NOT NULL,
    previous_refresh_token            TEXT,                            -- Refresh-токен до последней ротации, валиден до previous_refresh_token_expires_at
    previous_refresh_token_expires_at TIMESTAMPTZ,
    notification_token                VARCHAR,                         -- Токен для системы уведомлений
    application_bundle_id             VARCHAR      NOT NULL,
    application_version               VARCHAR      NOT NULL,
    application_build                 VARCHAR      NOT NULL,
    device_ip_address                 VARCHAR      NOT NULL,
    device_user_agent                 VARCHAR      NOT NULL,
    device_os_name                    VARCHAR      NOT NULL,
    device_os_version                 VARCHAR      NOT NULL,
    device_name                       VARCHAR      NOT NULL,
    device_model_name                 VARCHAR      NOT NULL,
    last_synced_audit_log_id          INT8         NOT NULL DEFAULT 0, -- Последняя подтвержденная (Confirm) запись аудит-лога
    pending_checkpoint                INT8,                            -- Чекпоинт последнего Sync, ожидающий подтверждения
    pending_sync_token                UUID,                            -- Опорный токен последнего Sync
    UNIQUE (user_id, device_id)
);
COMMENT ON TABLE tasks.devices IS 'Устройства пользователей и их чекпоинты синхронизации';

-- Аудит-лог — источник синхронизации, как в Coin
CREATE TABLE tasks.audit_log
(
    id              SERIAL PRIMARY KEY,                           -- Числовой идентификатор записи, он же чекпоинт синка
    entity          VARCHAR(100) NOT NULL,                        -- Название сущности (task, taskList и т.д.)
    method          VARCHAR(20)  NOT NULL,                        -- Метод изменения (create, update, delete)
    entity_id       VARCHAR(100) NOT NULL,                        -- Идентификатор сущности (без FK)
    snapshot_before JSONB,                                        -- Слепок сущности до изменения
    snapshot_after  JSONB,                                        -- Слепок сущности после изменения
    user_id         UUID         NOT NULL,                        -- Пользователь, совершивший действие
    device_id       VARCHAR(255) NOT NULL,                        -- Устройство, с которого совершено действие
    datetime_create TIMESTAMPTZ  NOT NULL DEFAULT now()           -- Дата и время изменения
);
COMMENT ON TABLE tasks.audit_log IS 'Аудит изменяющих действий пользователей';
CREATE INDEX idx_audit_log_user_id_id ON tasks.audit_log (user_id, id);
CREATE INDEX idx_audit_log_entity_entity_id ON tasks.audit_log (entity, entity_id);

-- Данные трекера. У каждой строки id от клиента, владелец, updated_at (последняя запись побеждает)
-- и мягкое удаление deleted_at
CREATE TABLE tasks.task_lists
(
    id         UUID             NOT NULL PRIMARY KEY,
    user_id    UUID             NOT NULL REFERENCES tasks.users (id),
    name       VARCHAR          NOT NULL,
    is_inbox   BOOLEAN          NOT NULL DEFAULT FALSE, -- Системный список Inbox, у пользователя ровно один
    color_hex  VARCHAR(9)       NOT NULL,
    icon       VARCHAR          NOT NULL,               -- SF Symbol иконки списка
    sort_order DOUBLE PRECISION NOT NULL,
    sort_mode  VARCHAR(20)      NOT NULL,               -- Режим сортировки задач: manual, priority, plannedDay, deadline, duration
    created_at TIMESTAMPTZ      NOT NULL,
    updated_at TIMESTAMPTZ      NOT NULL,
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE tasks.task_lists IS 'Списки задач';
CREATE INDEX idx_task_lists_user_id ON tasks.task_lists (user_id);
CREATE UNIQUE INDEX idx_task_lists_one_inbox ON tasks.task_lists (user_id) WHERE is_inbox AND deleted_at IS NULL;

CREATE TABLE tasks.list_sections
(
    id         UUID             NOT NULL PRIMARY KEY,
    user_id    UUID             NOT NULL REFERENCES tasks.users (id),
    list_id    UUID             NOT NULL REFERENCES tasks.task_lists (id),
    name       VARCHAR          NOT NULL,
    sort_order DOUBLE PRECISION NOT NULL,
    updated_at TIMESTAMPTZ      NOT NULL,
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE tasks.list_sections IS 'Разделы внутри списков';
CREATE INDEX idx_list_sections_user_id ON tasks.list_sections (user_id);

CREATE TABLE tasks.priorities
(
    id         UUID             NOT NULL PRIMARY KEY,
    user_id    UUID             NOT NULL REFERENCES tasks.users (id),
    name       VARCHAR          NOT NULL,
    color_hex  VARCHAR(9)       NOT NULL,
    marks      VARCHAR          NOT NULL, -- Знаки в строке задачи, например !!!
    sort_order DOUBLE PRECISION NOT NULL, -- Меньше — важнее
    updated_at TIMESTAMPTZ      NOT NULL,
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE tasks.priorities IS 'Справочник приоритетов';
CREATE INDEX idx_priorities_user_id ON tasks.priorities (user_id);

CREATE TABLE tasks.duration_options
(
    id         UUID        NOT NULL PRIMARY KEY,
    user_id    UUID        NOT NULL REFERENCES tasks.users (id),
    title      VARCHAR     NOT NULL,
    minutes    INT4        NOT NULL, -- Меньше нуля — «непрогнозируемо долго»
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE tasks.duration_options IS 'Справочник длительностей';
CREATE INDEX idx_duration_options_user_id ON tasks.duration_options (user_id);

CREATE TABLE tasks.tasks
(
    id               UUID             NOT NULL PRIMARY KEY,
    user_id          UUID             NOT NULL REFERENCES tasks.users (id),
    list_id          UUID             NOT NULL REFERENCES tasks.task_lists (id),
    section_id       UUID REFERENCES tasks.list_sections (id),
    title            VARCHAR          NOT NULL,
    notes            TEXT             NOT NULL DEFAULT '',
    priority_id      UUID REFERENCES tasks.priorities (id), -- NULL — «не разобрано»
    planned_day      VARCHAR(10),                           -- Планируемый срок, yyyy-MM-dd
    deadline_day     VARCHAR(10),                           -- Дедлайн, yyyy-MM-dd
    duration_minutes INT4,                                  -- Планируемое время в минутах
    manual_order     DOUBLE PRECISION NOT NULL,
    completed_at     TIMESTAMPTZ,
    external_id      VARCHAR,                               -- calendarItemIdentifier напоминания при импорте
    created_at       TIMESTAMPTZ      NOT NULL,
    updated_at       TIMESTAMPTZ      NOT NULL,
    deleted_at       TIMESTAMPTZ
);
COMMENT ON TABLE tasks.tasks IS 'Задачи';
CREATE INDEX idx_tasks_user_id ON tasks.tasks (user_id);
CREATE INDEX idx_tasks_list_id ON tasks.tasks (list_id);

CREATE TABLE tasks.tags
(
    id         UUID        NOT NULL PRIMARY KEY,
    user_id    UUID        NOT NULL REFERENCES tasks.users (id),
    name       VARCHAR     NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE tasks.tags IS 'Теги';
CREATE INDEX idx_tags_user_id ON tasks.tags (user_id);

CREATE TABLE tasks.tasks_to_tags
(
    user_id    UUID        NOT NULL REFERENCES tasks.users (id),
    task_id    UUID        NOT NULL REFERENCES tasks.tasks (id),
    tag_id     UUID        NOT NULL REFERENCES tasks.tags (id),
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (task_id, tag_id)
);
COMMENT ON TABLE tasks.tasks_to_tags IS 'Связи задач и тегов';
CREATE INDEX idx_tasks_to_tags_user_id ON tasks.tasks_to_tags (user_id);

CREATE TABLE tasks.blockers
(
    id               UUID        NOT NULL PRIMARY KEY,
    user_id          UUID        NOT NULL REFERENCES tasks.users (id),
    task_id          UUID        NOT NULL REFERENCES tasks.tasks (id), -- Заблокированная задача
    blocking_task_id UUID REFERENCES tasks.tasks (id),                 -- Блокирующая задача
    reason           VARCHAR,                                          -- Текстовая причина
    since_day        VARCHAR(10) NOT NULL,                             -- Дата начала, yyyy-MM-dd
    resolved_at      TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ NOT NULL,
    deleted_at       TIMESTAMPTZ,
    CHECK ((blocking_task_id IS NULL) <> (reason IS NULL))
);
COMMENT ON TABLE tasks.blockers IS 'Блокеры задач: ссылка на задачу либо текстовая причина';
CREATE INDEX idx_blockers_user_id ON tasks.blockers (user_id);

CREATE TABLE tasks.settings
(
    user_id                UUID        NOT NULL PRIMARY KEY REFERENCES tasks.users (id),
    daily_capacity_minutes INT8        NOT NULL, -- Дневной лимит загрузки в минутах
    today_sort_mode        VARCHAR(20) NOT NULL, -- Режим сортировки экрана «Сегодня»
    updated_at             TIMESTAMPTZ NOT NULL
);
COMMENT ON TABLE tasks.settings IS 'Настройки пользователя, общие для всех его устройств';

-- +goose Down
DROP SCHEMA tasks CASCADE;
