# Состояние работ

Обновляется в конце каждой сессии. Дата последнего обновления: 2026-11-06.

## Сделано

- Проектный документ утверждён заказчиком 2026-11-06: `docs/superpowers/specs/2026-11-06-unitrack-design.md`.
- **Этап 0 — каркас (2026-11-06)**, план `docs/superpowers/plans/2026-11-06-stage-0-scaffold.md`, 14 задач:
  - репозиторий, `README`, `CLAUDE.md` (корневой и для `backend/`, `web/`, `infra/`), `.claude/settings.json`;
  - ТЗ в Markdown с картой пробелов (`docs/spec/requirements.md`), 12 ADR (`docs/decisions/`);
  - дизайн-токены `design/tokens.json` → `web/src/shared/ui/tokens.css`, принципы стиля «тихий премиум»;
  - контракт `api/openapi.yaml` (health, ready, Problem) + бандл для генерации;
  - Go: конфигурация из окружения, пул pgx, миграции goose, тестовая PostgreSQL (`TESTDB_URL` или testcontainers), chi-роутер с request id / recover / лог, ошибки RFC 9457, health/ready, graceful shutdown, флаг `-healthcheck`, Dockerfile distroless;
  - веб: Vite 8 + TS 6 strict, FSD с линтером границ (`boundaries/dependencies`), stylelint на токены и БЭМ, Reka UI, 10 компонентов `shared/ui` с тестами, `entities/system`, главная страница со статусом сервера, TS-клиент из контракта;
  - инфраструктура: docker compose (postgres 17, api, nginx со статикой и `/apk/`), `infra/.env.example`, README для администратора;
  - lefthook (pre-commit: gofumpt, golangci-lint, prettier, eslint, stylelint, синхронность токенов и контракта; commit-msg: conventional commits; pre-push: тесты), CI: `backend.yml`, `web.yml`, `contract.yml`.

## В работе

- Ничего. Этап 0 закрыт, ветка `main` опубликована, все три workflow CI (`backend`, `web`, `contract`) зелёные на коммите с правками после финальной ревизии.

## Следующее

- Отложенные мелкие замечания финальной ревизии этапа 0 (не блокируют): `UNITRACK_MIGRATE_ON_START` вне `config.Load`; `Recover` не пробрасывает `http.ErrAbortHandler`; `httpserver.Problem` дублирует `gen.Problem`; nginx `proxy_pass` без `resolver`; версии GitHub Actions и сборка golangci-lint/sqlc из исходников в CI; commit-msg regex не пропускает Merge/Revert; `.gitignore` без `.vscode/*`; `describedBy()` в UiTextField → computed; UiDialog без описания; код ошибки вместо русской строки в `entities/system/api/health.ts`; `shadow.focus` с var() сломает генератор Kotlin (этап 7); `principles.md` не упоминает radius xs/full; `errcheck`/`bodyclose` выключены в тестах целиком; `vitest.config.ts` дублирует `vite.config.ts`; нет автотеста на правило границ FSD. Закрывать по пути на этапе 1, где код затрагивается.

- **Этап 1 — контракт и ядро** (проектный документ §12, первая очередь): `openapi.yaml` для auth, users, structure, students, calendar, reasons, reports; Go-модули auth, users, structure, students, reasons, calendar с тестами; миграции; идемпотентность мутаций (`Idempotency-Key`); `shared/lib/access` на вебе.
- Первым делом — план `docs/superpowers/plans/<дата>-stage-1-core.md` через брейншторм/writing-plans по проектному документу.
- Открытые вопросы к заказчику — проектный документ §14 (роли ректората равны? срок хранения аудита; домен и сертификат; часы сдачи и праздники; Moodle).

## Особенности локальной среды (не правила проекта)

- Docker Hub из этой сети отдаёт образы медленно; локальные тесты бэкенда идут через `TESTDB_URL` на установленной PostgreSQL 16 (`backend/.env`), путь через testcontainers проверен.
- GitHub доступен только через локальный прокси `127.0.0.1:8899`; для репозитория задано `http.https://github.com.proxy`.
