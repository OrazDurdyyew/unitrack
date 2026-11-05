# backend — правила для Go

Модульный монолит на Go (ADR 0001, 0004). Стандартная библиотека прежде сторонних пакетов; новый пакет в `go.mod` — только с ADR.

## Структура

```
cmd/api/main.go                  точка входа: конфиг → пул БД → миграции → HTTP-сервер → graceful shutdown; флаг -healthcheck
internal/api/gen/                сгенерировано из api/openapi.bundle.yaml (oapi-codegen) — руками не править
internal/platform/config         настройки из переменных UNITRACK_*
internal/platform/db             пул pgx, миграции goose (migrations/*.sql, встроены в бинарник)
internal/platform/httpserver     роутер chi, middleware (request id, recover, лог), ошибки RFC 9457
internal/platform/testdb         настоящая PostgreSQL для тестов (TESTDB_URL или testcontainers)
internal/<module>/               предметный модуль: handler.go, service.go, repo.go, queries.sql, *_test.go
```

Модули (появляются по этапам): `health`, `auth`, `users`, `structure`, `students`, `calendar`, `reasons`, `reports`, `changerequests`, `imports`, `analytics`, `notifications`, `audit`, `sync`, `apk`.

## Три слоя, зависимости вниз

- **handler** — разбор HTTP, проверка прав, вызов сервиса, ответ. Не знает о SQL.
- **service** — бизнес-правила, транзакции, запись в аудит. Не знает о HTTP.
- **repo** — запросы через sqlc. Не знает о правилах.

Модуль не читает чужие таблицы напрямую — только через интерфейс чужого сервиса. Интерфейс появляется при второй реализации или для подмены в тесте, не заранее.

## Как добавить эндпоинт

1. Описать операцию в `api/paths/<module>.yaml`, подключить в `api/openapi.yaml`, `task gen`.
2. Реализовать метод `gen.ServerInterface` в `internal/<module>/handler.go`; собрать обработчики в `cmd/api/main.go`.
3. Тест через `httptest` + `httpserver.NewRouter`; для базы — `testdb.New(t)`.
4. Ошибки клиенту — `httpserver.WriteProblem(w, r, status, code, title, detail)`, `code` в snake_case.

## Как добавить миграцию

`"$(go env GOPATH)/bin/goose" -dir internal/platform/db/migrations create <name> sql`. Применённую миграцию не менять. Данные не удаляются — только статусы и даты.

## Правила

- `context.Context` — первый аргумент. Ошибки оборачивать: `fmt.Errorf("что делали: %w", err)`, тексты на русском. Не глотать ошибки; `_ =` только с комментарием почему.
- Логи — `slog` с полями, не `fmt.Println`. Персональные данные в логи не попадают.
- SQL только в `*.sql` для sqlc (`sqlc.yaml`). Транзакции открывает сервисный слой.
- Глобальных переменных нет, кроме `version` в `main`.
- Тесты: `go test -race`, настоящая PostgreSQL, без моков SQL. Таблица в тесте — одно поведение на тест.

## Команды

`task backend:gen | lint | test | build | dev`. Локальные переменные — `backend/.env` (см. `.env.example`): `TESTDB_URL` для тестов без Docker, `UNITRACK_*` для `dev`.
