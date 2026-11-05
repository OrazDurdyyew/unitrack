-- Миграции только вперёд: применённую миграцию менять запрещено, нужна новая.
-- Формат goose: секции Up и Down обязательны.

-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE app_info (
    key        text PRIMARY KEY,
    value      text NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now()
);

INSERT INTO app_info (key, value) VALUES ('schema', 'unitrack');

-- +goose Down
DROP TABLE app_info;
