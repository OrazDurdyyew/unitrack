# План реализации: Этап 0 «Каркас» UniTrack

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Поднять репозиторий UniTrack до состояния, в котором `docker compose up` запускает пустой Go-API и статику Vue, CI и pre-commit зелёные, а любая новая сессия Claude за десять минут понимает проект по `CLAUDE.md`, `docs/STATUS.md` и ADR.

**Architecture:** Монорепозиторий с четырьмя частями (`backend/` Go, `web/` Vue 3 FSD, `infra/` docker + nginx, `android/` — заглушка до этапа 7) и двумя источниками истины над ними: `api/openapi.yaml` (контракт, из него генерируются Go-интерфейсы и TS-клиент) и `design/tokens.json` (дизайн-токены, из них генерируется `tokens.css`). Единые команды через корневой `Taskfile.yml`, проверки через lefthook и GitHub Actions.

**Tech Stack:** Go 1.24, chi v5, pgx v5, goose v3, oapi-codegen v2, testcontainers-go; Vue 3.5, Vite, TypeScript strict, pnpm, ESLint + eslint-plugin-boundaries, stylelint + stylelint-declaration-strict-value, Prettier, Vitest, vue-i18n, vue-router, Pinia, TanStack Query, Reka UI, openapi-typescript + openapi-fetch; Docker Compose, nginx; Taskfile, lefthook.

**Spec:** `docs/superpowers/specs/2026-11-06-unitrack-design.md` (разделы 6, 9, 10, 11, 12 этап 0).

## Global Constraints

- Документация, комментарии, тексты интерфейса, сообщения коммитов — на русском; идентификаторы, имена таблиц и файлов — на английском (спец. §10).
- Conventional commits с префиксами `feat:`, `fix:`, `api:`, `docs:`, `infra:`, `chore:`, `test:`; текст после префикса на русском, в инфинитиве («добавить…»). Без строк атрибуции.
- Go: стандартная библиотека прежде сторонних пакетов; новый пакет в `go.mod` только с записью в ADR; SQL только в `*.sql` для sqlc; `context.Context` первым аргументом; ошибки оборачиваются с контекстом (`fmt.Errorf("…: %w", err)`); глобальных переменных нет, кроме конфигурации (§10).
- Vue: только `<script setup lang="ts">`; никаких `any`; никаких строк интерфейса в коде — только через vue-i18n; никаких стилей вне токенов; компонент больше 200 строк делится; запросы к серверу только в `entities/*/api` (§10).
- Все цвета, отступы, радиусы, шрифты — только через CSS-переменные из `web/src/shared/ui/tokens.css`; stylelint запрещает «сырые» значения (§6.2).
- Внешних сетевых зависимостей во время работы нет: шрифт Inter лежит в репозитории, ничего не грузится с CDN (§3, §6.5).
- Часовой пояс отображения Asia/Ashgabat; хранение в UTC (§3).
- Ничего не удаляется из данных — только статусы и даты (§4.1). В этапе 0 таблиц ещё нет, но правило входит в CLAUDE.md.
- Все даты в документах абсолютные (ГГГГ-ММ-ДД).

## Review Focus

Пять условий, которые спецификация подразумевает, но которые легко сломать в каркасе; каждое закреплено тестом в задаче-владельце:

1. **Токен с «сырым» цветом в компоненте проходит незамеченным.** Ожидание: stylelint падает на `color: #16A964` вне `tokens.css`. Тест — Задача 9, шаг «отрицательная проверка stylelint».
2. **Импорт «снизу вверх» в FSD (например, `shared` импортирует `features`) проходит.** Ожидание: ESLint падает. Тест — Задача 9, шаг «отрицательная проверка boundaries».
3. **Сгенерированный код расходится с `openapi.yaml`, но коммит проходит.** Ожидание: `task gen` + `git diff --exit-code` в CI падает. Тест — Задача 12, workflow `contract.yml`.
4. **API стартует без переменных окружения или с пустым `DATABASE_URL` и молча работает «как-нибудь».** Ожидание: процесс завершается с понятной ошибкой. Тест — Задача 6, `TestLoad_MissingDatabaseURL`.
5. **Readiness отвечает 200 при недоступной базе, и docker compose считает контейнер здоровым.** Ожидание: `/api/v1/ready` отвечает 503. Тест — Задача 7, `TestReady_DatabaseDown`.

---

## Карта файлов

```
unitrack/
├── CLAUDE.md                          # Задача 13
├── README.md                          # Задача 1
├── .editorconfig, .gitattributes      # Задача 1
├── Taskfile.yml                       # Задача 6 (создаётся), 9, 11 (дополняется)
├── lefthook.yml                       # Задача 12
├── .claude/settings.json              # Задача 13
├── .github/workflows/{backend,web,contract}.yml   # Задача 12
├── docs/
│   ├── STATUS.md                      # Задача 1, обновляется в Задаче 14
│   ├── spec/requirements.md           # Задача 2
│   ├── decisions/README.md, 0001…0012 # Задача 3
│   ├── design/style-reference.html    # Задача 1
│   └── design/principles.md           # Задача 4
├── design/
│   ├── tokens.json                    # Задача 4
│   ├── build-tokens.mjs               # Задача 4
│   └── build-tokens.test.mjs          # Задача 4
├── api/
│   ├── openapi.yaml                   # Задача 5
│   └── paths/health.yaml              # Задача 5
├── backend/                           # Задачи 6–8
│   ├── CLAUDE.md                      # Задача 13
│   ├── go.mod, Taskfile.yml, .golangci.yml, Dockerfile
│   ├── cmd/api/main.go
│   ├── internal/platform/config/{config.go,config_test.go}
│   ├── internal/platform/httpserver/{server.go,middleware.go,server_test.go}
│   ├── internal/platform/db/{pool.go,migrate.go,migrations/00001_init.sql}
│   ├── internal/platform/testdb/testdb.go
│   ├── internal/health/{handler.go,handler_test.go}
│   └── internal/api/gen/{api.gen.go,oapi-codegen.yaml,generate.go}
├── web/                               # Задачи 9–10
│   ├── CLAUDE.md                      # Задача 13
│   ├── package.json, vite.config.ts, tsconfig*.json, eslint.config.js, .stylelintrc.json, .prettierrc
│   ├── public/fonts/Inter-*.woff2
│   └── src/{app,pages,widgets,features,entities,shared}/…
├── infra/                             # Задача 11
│   ├── CLAUDE.md                      # Задача 13
│   ├── docker-compose.yml, .env.example
│   ├── nginx/nginx.conf
│   └── web.Dockerfile
└── android/README.md                  # Задача 1 (заглушка)
```

---

### Task 1: Основа репозитория и живой статус

**Files:**
- Create: `README.md`, `.editorconfig`, `.gitattributes`, `docs/STATUS.md`, `docs/design/style-reference.html`, `android/README.md`
- Modify: `.gitignore`

**Interfaces:**
- Produces: `docs/STATUS.md` с тремя разделами `## Сделано`, `## В работе`, `## Следующее` — на эти заголовки ссылается CLAUDE.md (Задача 13) и их обновляет Задача 14.

- [ ] **Шаг 1: Скопировать эталон стиля**

Run: `cp .superpowers/brainstorm/*/content/style-proposal.html docs/design/style-reference.html`
Expected: файл существует, `grep -c "#16A964" docs/design/style-reference.html` ≥ 1.

- [ ] **Шаг 2: Создать `README.md`**

Содержание (на русском): одно предложение о назначении; таблица частей репозитория (`backend/`, `web/`, `infra/`, `android/`, `api/`, `design/`, `docs/`) с одной строкой на каждую; раздел «Быстрый старт» с командами `task setup`, `task dev`, `task check` (появятся в Задачах 6–12); ссылки на `docs/superpowers/specs/2026-11-06-unitrack-design.md`, `docs/STATUS.md`, `docs/decisions/`.

- [ ] **Шаг 3: Создать `docs/STATUS.md`**

Шапка: «Обновляется в конце каждой сессии. Дата последнего обновления: 2026-11-06». Разделы `## Сделано` (пусто, кроме «проектный документ утверждён 2026-11-06»), `## В работе` («Этап 0 — каркас, план `docs/superpowers/plans/2026-11-06-stage-0-scaffold.md`, выполнено задач: 0 из 14»), `## Следующее` («Этап 1 — контракт и ядро»).

- [ ] **Шаг 4: Создать `.editorconfig` и `.gitattributes`**

`.editorconfig`: `root = true`, UTF-8, LF, финальный перевод строки, отступ 2 пробела для `*.{ts,vue,json,yaml,yml,css,md}`, табы для `*.go`, 4 пробела для `*.kt`. `.gitattributes`: `* text=auto eol=lf`, `*.woff2 binary`, `*.png binary`, `*.gen.go linguist-generated=true`, `web/src/shared/api/schema.d.ts linguist-generated=true`.

- [ ] **Шаг 5: Дополнить `.gitignore`**

Добавить секции: Go (`backend/bin/`, `backend/tmp/`), web (`web/node_modules/`, `web/dist/`, `web/coverage/`, `web/playwright-report/`), Android (`android/.gradle/`, `android/**/build/`, `android/local.properties`), IDE (`.idea/`, `.vscode/*` кроме `!.vscode/extensions.json`).

- [ ] **Шаг 6: Заглушка `android/README.md`**

Две строки: «Android-приложение создаётся на этапе 7 (см. проектный документ §12, вторая очередь). До этого папка содержит только этот файл.»

- [ ] **Шаг 7: Проверить и зафиксировать**

Run: `git add -A && git status --short`
Expected: в списке нет `.env`, нет `.superpowers/`.

```bash
git commit -m "chore: создать основу репозитория, статус и эталон стиля"
```

---

### Task 2: ТЗ заказчика в Markdown

**Files:**
- Create: `docs/spec/requirements.md`
- Keep: `docs/eduControl.docx` (исходник не трогать)

**Interfaces:**
- Produces: якоря разделов `## 1`…`## 4.2` в `requirements.md`, на которые ссылаются ADR (Задача 3).

- [ ] **Шаг 1: Извлечь текст из docx во временный файл**

Run (из `/root/projects/unitrack`):
```bash
python3 -I -c "
import zipfile,re,html,sys
x=zipfile.ZipFile('docs/eduControl.docx').read('word/document.xml').decode()
x=re.sub(r'</w:p>','\n',x); x=re.sub(r'<w:tab/>','\t',x)
sys.stdout.write(html.unescape(re.sub(r'<[^>]+>','',x)))
" > /tmp/claude-0/-root-projects-unitrack/1bf6f526-0c9e-49df-b474-e148463760e6/scratchpad/tz.md
wc -l /tmp/claude-0/-root-projects-unitrack/1bf6f526-0c9e-49df-b474-e148463760e6/scratchpad/tz.md
```
Expected: около 170 строк, последняя содержательная строка — про видимость запросов для декана.

- [ ] **Шаг 2: Оформить `docs/spec/requirements.md`**

Шапка: «Источник: `docs/eduControl.docx`, получен 2026-11-06. Это дословный перенос ТЗ заказчика в Markdown без изменений смысла. Документ заказчика обрывается на разделе 4.2 „Видимость“; пробелы закрыты в проектном документе, см. таблицу внизу.» Далее текст ТЗ с исправленной Markdown-разметкой (заголовки `#`, списки `-`, таблица для `absence_reasons` по полям). В конце раздел `## Пробелы ТЗ и где они закрыты` — таблица «Пробел → раздел проектного документа»: типы отчётов и смены → §4.3; роли преподавателя → §1, §4.2; вечерняя проверка → §4.3; возврат отчёта деканом → §4.3; данные студента → §4.1; жизненный цикл групп → §4.1; импорт → §4.6; сводки → §4.7; уведомления → §4.8; видимость запросов для ректората → §4.5; стек → §6; офлайн → §5.

- [ ] **Шаг 3: Проверить**

Run: `grep -c "^#" docs/spec/requirements.md && grep -n "Пробелы ТЗ" docs/spec/requirements.md`
Expected: заголовков ≥ 15; раздел «Пробелы ТЗ» найден.

- [ ] **Шаг 4: Зафиксировать**

```bash
git add docs/spec/requirements.md
git commit -m "docs: перенести ТЗ заказчика в Markdown с картой пробелов"
```

---

### Task 3: Журнал архитектурных решений (ADR)

**Files:**
- Create: `docs/decisions/README.md`, `docs/decisions/0000-template.md`, `docs/decisions/0001-…0012-….md`

**Interfaces:**
- Produces: формат ADR, который CLAUDE.md (Задача 13) обязывает использовать для любого решения вне ТЗ.

- [ ] **Шаг 1: Шаблон и README**

`0000-template.md`: заголовок `# NNNN. Название`, поля `Дата`, `Статус` (предложено / принято / заменено на NNNN), разделы `## Контекст`, `## Решение`, `## Последствия` (плюсы, минусы, что сделать), `## Отклонённые варианты` (с одной строкой «почему нет» на каждый). `README.md`: что такое ADR, когда писать (любое решение, которого нет в `docs/spec/requirements.md`), как нумеровать (следующий номер, ведущие нули до четырёх), как заменять (новый ADR, старому `Статус: заменено на NNNN`), индекс-таблица номер → название → статус.

- [ ] **Шаг 2: Написать двенадцать ADR из проектного документа**

Каждый 15–40 строк, с датой 2026-11-06 и статусом «принято». Контекст и отклонённые варианты берутся из §6 и из истории обсуждения, зафиксированной в проектном документе:

| № | Файл | Решение | Отклонено |
|---|---|---|---|
| 0001 | `0001-go-postgresql-backend.md` | Go + PostgreSQL 17 | Node/NestJS, Kotlin/Spring |
| 0002 | `0002-contract-first-openapi.md` | Рукописный `api/openapi.yaml`, oapi-codegen, генерация TS и Kotlin клиентов | Code-first генерация спеки из кода |
| 0003 | `0003-sqlc-without-orm.md` | pgx + sqlc, goose | GORM, ent |
| 0004 | `0004-modular-monolith.md` | Модульный монолит, три слоя handler/service/repo, интерфейс при второй реализации | Микросервисы, полная clean architecture с портами/адаптерами |
| 0005 | `0005-vue3-fsd-spa.md` | Vue 3 + Vite SPA, FSD с линтером границ | Nuxt 4 |
| 0006 | `0006-pure-css-bem-tokens-reka.md` | Чистый CSS, БЭМ, токены, Reka UI + свой `shared/ui` | Tailwind, SCSS, PrimeVue, shadcn-vue |
| 0007 | `0007-tanstack-query-zod-forms.md` | TanStack Query для серверных данных, Pinia для клиентского состояния, zod + свой composable | Pinia для серверных данных, vee-validate |
| 0008 | `0008-native-kotlin-offline-first.md` | Нативный Kotlin, Room, WorkManager, очередь операций с UUID и идемпотентностью | Flutter, React Native, PWA |
| 0009 | `0009-on-premise-no-external-services.md` | Сервер в институте, docker compose, APK по ссылке, уведомления без Google, Inter в репозитории | Облако, Google Play, FCM |
| 0010 | `0010-report-snapshot-and-unmarked.md` | Состав отчёта — снимок; расхождение с сервером → `unmarked`; нельзя утвердить с `unmarked` | Пересчёт состава на сервере, автоматическое «присутствует» |
| 0011 | `0011-draft-edited-by-author-only.md` | Черновик правит только автор и проверяющий на своей стадии | Совместное редактирование с версионированием |
| 0012 | `0012-russian-docs-english-identifiers.md` | Русский для документов, комментариев, коммитов, интерфейса; английский для идентификаторов | Всё на английском, всё на туркменском |

- [ ] **Шаг 3: Проверить индекс**

Run: `ls docs/decisions/ | wc -l && grep -c "^| 00" docs/decisions/README.md`
Expected: 14 файлов (README, шаблон, 12 ADR); 12 строк в индексе.

- [ ] **Шаг 4: Зафиксировать**

```bash
git add docs/decisions
git commit -m "docs: завести журнал архитектурных решений и первые 12 ADR"
```

---

### Task 4: Дизайн-токены и генератор `tokens.css`

**Files:**
- Create: `design/tokens.json`, `design/build-tokens.mjs`, `design/build-tokens.test.mjs`, `design/package.json`, `docs/design/principles.md`
- Output (генерируется, коммитится): `web/src/shared/ui/tokens.css`

**Interfaces:**
- Produces: CSS-переменные с именами `--color-<group>-<step>`, `--radius-<size>`, `--space-<n>`, `--font-family-base`, `--font-size-<name>`, `--shadow-<name>`; на них ссылаются компоненты Задачи 10 и правило stylelint Задачи 9.
- Produces: `node design/build-tokens.mjs` — команда, которую вызывает `task gen` (Задача 6) и CI (Задача 12).

- [ ] **Шаг 1: Написать `design/tokens.json`**

Структура `{ "color": {...}, "radius": {...}, "space": {...}, "font": {...}, "shadow": {...} }`. Значения — из §6.5 дословно: `color.green.600 = "#16A964"`, `color.green.700 = "#12905A"`, `color.green.100 = "#D4F0E2"`, `color.green.50 = "#EAF7F0"`, `color.ink.900 = "#0B1F17"`, `color.ink.700 = "#2B3A33"`, `color.ink.500 = "#5E6E66"`, `color.ink.400 = "#8A978F"`, `color.line.default = "#E3E9E5"`, `color.line.strong = "#CBD5CF"`, `color.bg.default = "#F4F6F5"`, `color.surface.default = "#FFFFFF"`, `color.amber.700 = "#9A6700"`, `color.amber.50 = "#FFF4DD"`, `color.red.700 = "#B42318"`, `color.red.50 = "#FDECEC"`; `radius.sm = "6px"`, `radius.md = "10px"`, `radius.lg = "14px"`; `space.1…space.8` = 4, 8, 12, 16, 20, 24, 32, 40 px; `font.family.base = "Inter, system-ui, sans-serif"`, `font.size.xs…xl` = 11, 12, 13.5, 15, 18, 22 px (`xs, sm, md, lg, xl, 2xl`); `shadow.sm = "0 1px 2px rgba(11,31,23,.06)"`, `shadow.md = "0 8px 24px -12px rgba(11,31,23,.18)"`.

- [ ] **Шаг 2: Написать падающий тест `design/build-tokens.test.mjs`** (`node:test`, `node:assert/strict`)

```js
import { test } from 'node:test';
import assert from 'node:assert/strict';
import { tokensToCss } from './build-tokens.mjs';

test('вложенные ключи становятся переменными через дефис', () => {
  const css = tokensToCss({ color: { green: { 600: '#16A964' } }, radius: { sm: '6px' } });
  assert.match(css, /--color-green-600: #16A964;/);
  assert.match(css, /--radius-sm: 6px;/);
});

test('результат обёрнут в :root и начинается с предупреждения о генерации', () => {
  const css = tokensToCss({ radius: { sm: '6px' } });
  assert.ok(css.startsWith('/* Сгенерировано из design/tokens.json'));
  assert.match(css, /:root \{[\s\S]*\}/);
});

test('ключи шрифта с несколькими словами не теряют части имени', () => {
  const css = tokensToCss({ font: { family: { base: 'Inter' }, size: { '2xl': '22px' } } });
  assert.match(css, /--font-family-base: Inter;/);
  assert.match(css, /--font-size-2xl: 22px;/);
});
```

- [ ] **Шаг 3: Запустить и убедиться, что падает**

Run: `node --test design/`
Expected: FAIL, `Cannot find module './build-tokens.mjs'`.

- [ ] **Шаг 4: Реализовать `design/build-tokens.mjs`**

Экспорт `tokensToCss(tokens: object): string` — рекурсивный обход, имя переменной = путь через `-`, отступ две пробела, сортировка по имени для стабильного диффа. Когда файл запущен напрямую (`import.meta.url === pathToFileURL(process.argv[1]).href`), читает `design/tokens.json` и пишет `web/src/shared/ui/tokens.css` (создавая папки). `design/package.json`: `{ "type": "module", "private": true }`.

- [ ] **Шаг 5: Запустить тест и генерацию**

Run: `node --test design/ && node design/build-tokens.mjs && grep -c -- "--color-" web/src/shared/ui/tokens.css`
Expected: тесты PASS; счётчик = 16.

- [ ] **Шаг 6: Написать `docs/design/principles.md`**

Одна страница: название стиля «тихий премиум»; ссылка на `style-reference.html`; таблица ролей цветов (действие — green-600, текст — ink-900/700/500/400, ожидание — amber, отказ и отсутствие — red, линии, фон, поверхность); правило «статус = цветная точка + слово»; радиусы 6/10/14; шрифт Inter, цифры табличные (`font-variant-numeric: tabular-nums`); правило «ни одного сырого значения вне `tokens.css`»; как добавить токен (правка `tokens.json` → `task gen` → коммит обоих файлов).

- [ ] **Шаг 7: Зафиксировать**

```bash
git add design docs/design/principles.md web/src/shared/ui/tokens.css
git commit -m "feat: добавить дизайн-токены и генератор tokens.css"
```

---

### Task 5: Контракт API — первая версия

**Files:**
- Create: `api/openapi.yaml`, `api/paths/health.yaml`, `api/components/problem.yaml`, `api/README.md`

**Interfaces:**
- Produces: `GET /api/v1/health` → 200 `HealthResponse { status: "ok", version: string }`; `GET /api/v1/ready` → 200 `ReadyResponse { status: "ok", database: "ok" }` или 503 `ReadyResponse { status: "degraded", database: "unavailable" }`; схема `Problem` (RFC 9457: `type`, `title`, `status`, `detail`, `instance`, `code`). operationId: `getHealth`, `getReady`. Эти имена использует oapi-codegen (Задача 8) и openapi-typescript (Задача 9).

- [ ] **Шаг 1: Написать `api/openapi.yaml`**

`openapi: 3.1.0`; `info.title: UniTrack API`, `info.version: 0.1.0`, описание на русском; `servers: [{ url: /api/v1 }]`; `paths` через `$ref: ./paths/health.yaml#/health` и `#/ready`; `components.schemas.Problem` через `$ref: ./components/problem.yaml#/Problem`; `tags: [{ name: system, description: Служебные операции }]`.

- [ ] **Шаг 2: Написать `api/paths/health.yaml` и `api/components/problem.yaml`**

Схемы `HealthResponse` и `ReadyResponse` объявить inline в `health.yaml` под ключом `components`-подобной секции не нужно — поместить их в `openapi.yaml` → `components.schemas` рядом с `Problem`, а в `health.yaml` ссылаться `$ref: ../openapi.yaml#/components/schemas/HealthResponse`. Все `description` на русском. `status` и `database` — `enum`.

- [ ] **Шаг 3: Проверить валидность спеки**

Run: `cd web 2>/dev/null || true; pnpm dlx @redocly/cli@latest lint /root/projects/unitrack/api/openapi.yaml`
Expected: `Woohoo! Your API description is valid.` (предупреждения о `security` допустимы на этапе 0).

- [ ] **Шаг 4: `api/README.md`**

Правила из §7: версия в пути, контракт сначала, коммит с префиксом `api:` и регенерация клиентов в том же PR, `Idempotency-Key` для мутаций с Android (появится на этапе 1), ошибки по RFC 9457, пагинация курсором. Команда проверки: `task api:lint`.

- [ ] **Шаг 5: Зафиксировать**

```bash
git add api
git commit -m "api: завести контракт OpenAPI с операциями health и ready"
```

---

### Task 6: Go-модуль, конфигурация, Taskfile

**Files:**
- Create: `backend/go.mod`, `backend/Taskfile.yml`, `backend/.golangci.yml`, `backend/internal/platform/config/config.go`, `backend/internal/platform/config/config_test.go`, `Taskfile.yml` (корень)

**Interfaces:**
- Produces: `config.Load(getenv func(string) string) (Config, error)`; `type Config struct { HTTPAddr string; DatabaseURL string; LogLevel slog.Level; Env string }`. Переменные: `UNITRACK_HTTP_ADDR` (по умолчанию `:8080`), `UNITRACK_DATABASE_URL` (обязательна), `UNITRACK_LOG_LEVEL` (`debug|info|warn|error`, по умолчанию `info`), `UNITRACK_ENV` (`development|production`, по умолчанию `production`).
- Produces: корневые команды `task setup`, `task gen`, `task lint`, `task test`, `task check` (= gen + lint + test + проверка чистоты diff), `task dev`; они делегируют в `backend/Taskfile.yml` и далее в `web/Taskfile.yml` (Задача 9), `infra` (Задача 11).

- [ ] **Шаг 1: Инициализировать модуль**

Run: `cd backend && go mod init github.com/OrazDurdyyew/unitrack/backend && go mod edit -go=1.24`
Expected: `backend/go.mod` создан.

- [ ] **Шаг 2: Написать падающий тест `config_test.go`**

```go
func TestLoad_Defaults(t *testing.T) {
	env := map[string]string{"UNITRACK_DATABASE_URL": "postgres://u:p@h/db"}
	cfg, err := config.Load(func(k string) string { return env[k] })
	if err != nil { t.Fatal(err) }
	if cfg.HTTPAddr != ":8080" { t.Errorf("HTTPAddr = %q", cfg.HTTPAddr) }
	if cfg.LogLevel != slog.LevelInfo { t.Errorf("LogLevel = %v", cfg.LogLevel) }
	if cfg.Env != "production" { t.Errorf("Env = %q", cfg.Env) }
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	_, err := config.Load(func(string) string { return "" })
	if err == nil || !strings.Contains(err.Error(), "UNITRACK_DATABASE_URL") { t.Fatalf("ожидалась ошибка про UNITRACK_DATABASE_URL, получено %v", err) }
}

func TestLoad_InvalidLogLevel(t *testing.T) {
	env := map[string]string{"UNITRACK_DATABASE_URL": "postgres://u:p@h/db", "UNITRACK_LOG_LEVEL": "loud"}
	_, err := config.Load(func(k string) string { return env[k] })
	if err == nil { t.Fatal("ожидалась ошибка") }
}
```

- [ ] **Шаг 3: Запустить и убедиться, что падает**

Run: `cd backend && go test ./internal/platform/config/`
Expected: FAIL (пакет не компилируется: `config.Load` не определён).

- [ ] **Шаг 4: Реализовать `config.go`**

Только стандартная библиотека. Ошибки — через `errors.Join` всех найденных проблем, чтобы администратор увидел сразу всё, текст на русском: `"UNITRACK_DATABASE_URL: переменная обязательна"`.

- [ ] **Шаг 5: Запустить тесты**

Run: `cd backend && go test ./internal/platform/config/ -v`
Expected: 3 PASS.

- [ ] **Шаг 6: `.golangci.yml` и `backend/Taskfile.yml`**

`.golangci.yml` (v2 формат): линтеры `govet, errcheck, staticcheck, unused, gofumpt, goimports, revive, gocritic, errorlint, noctx, bodyclose, sqlclosecheck, exhaustive, gosec`; исключить `internal/api/gen/`. `backend/Taskfile.yml` задачи: `tools` (`go install` для `golangci-lint`, `sqlc`, `goose`, `oapi-codegen` с закреплёнными версиями через `go run`-стиль `tool` директивы Go 1.24: добавить их в `go.mod` через `go get -tool`), `gen` (пока только `go generate ./...`), `lint` (`golangci-lint run ./...`), `test` (`go test -race -count=1 ./...`), `build` (`go build -ldflags "-X main.version=$(git describe --tags --always --dirty)" -o bin/api ./cmd/api`), `dev` (`go run ./cmd/api`).

- [ ] **Шаг 7: Корневой `Taskfile.yml`**

`version: '3'`; `includes: { backend: ./backend, web: ./web, infra: ./infra }` (web и infra появятся позже — на этом шаге укажи `optional: true`); задачи `setup` (backend:tools), `gen` (`node design/build-tokens.mjs`, backend:gen), `lint` (backend:lint), `test` (backend:test), `check` (gen → `git diff --exit-code -- web/src/shared/ui/tokens.css backend/internal/api/gen` → lint → test), `dev` (пока backend:dev). Каждая задача с `desc:` на русском.

- [ ] **Шаг 8: Установить инструменты и прогнать проверку**

Run: `task setup && task check`
Expected: завершение с кодом 0; `golangci-lint` отрабатывает без замечаний.

- [ ] **Шаг 9: Зафиксировать**

```bash
git add Taskfile.yml backend
git commit -m "feat: завести Go-модуль, конфигурацию из окружения и Taskfile"
```

---

### Task 7: База данных — пул, миграции, тестовый контейнер

**Files:**
- Create: `backend/internal/platform/db/pool.go`, `backend/internal/platform/db/migrate.go`, `backend/internal/platform/db/migrations/00001_init.sql`, `backend/internal/platform/db/migrate_test.go`, `backend/internal/platform/testdb/testdb.go`, `backend/sqlc.yaml`

**Interfaces:**
- Produces: `db.NewPool(ctx context.Context, url string) (*pgxpool.Pool, error)`; `db.Migrate(ctx context.Context, pool *pgxpool.Pool) error` (goose, SQL-файлы через `embed.FS`); `testdb.New(t testing.TB) *pgxpool.Pool` — поднимает PostgreSQL 17 в testcontainers один раз на пакет (через `sync.Once` в пакете testdb), применяет миграции, на каждый вызов создаёт отдельную схему или базу и возвращает пул к ней; очистка через `t.Cleanup`.
- Consumes: `config.Config.DatabaseURL` (Задача 6).

- [ ] **Шаг 1: Добавить зависимости**

Run: `cd backend && go get github.com/jackc/pgx/v5@latest github.com/pressly/goose/v3@latest github.com/testcontainers/testcontainers-go@latest github.com/testcontainers/testcontainers-go/modules/postgres@latest`
Expected: `go.mod` обновлён. Это пакеты из ADR 0003 (pgx, goose) и §6.1 (testcontainers) — отдельный ADR не нужен.

- [ ] **Шаг 2: Первая миграция**

`00001_init.sql` в формате goose (`-- +goose Up` / `-- +goose Down`): `CREATE EXTENSION IF NOT EXISTS pgcrypto;` и таблица `app_info (key text primary key, value text not null, updated_at timestamptz not null default now())` с строкой `('schema', 'unitrack')`. Down удаляет таблицу. Комментарий в начале файла на русском: правило «миграции только вперёд, изменение применённой миграции запрещено».

- [ ] **Шаг 3: Написать падающий тест `migrate_test.go`**

```go
func TestMigrate_AppliesInitialSchema(t *testing.T) {
	pool := testdb.New(t)
	var value string
	err := pool.QueryRow(context.Background(), `select value from app_info where key = 'schema'`).Scan(&value)
	if err != nil { t.Fatal(err) }
	if value != "unitrack" { t.Errorf("value = %q", value) }
}

func TestMigrate_IsIdempotent(t *testing.T) {
	pool := testdb.New(t)
	if err := db.Migrate(context.Background(), pool); err != nil { t.Fatalf("повторная миграция: %v", err) }
}
```

- [ ] **Шаг 4: Запустить и убедиться, что падает**

Run: `cd backend && go test ./internal/platform/db/`
Expected: FAIL — пакет `testdb` не существует.

- [ ] **Шаг 5: Реализовать `pool.go`, `migrate.go`, `testdb.go`**

`migrate.go`: `goose.SetBaseFS(migrations)`, диалект `postgres`, через `goose.NewProvider` с `database/sql`-обёрткой над pgx (`stdlib.OpenDBFromPool`). `testdb.go`: образ `postgres:17-alpine`, `WithDatabase("unitrack_test")`, ожидание по логу `database system is ready to accept connections` дважды; при `TESTDB_URL` в окружении — использовать его вместо контейнера (для CI без Docker-in-Docker). Изоляция тестов — отдельная база на вызов: `CREATE DATABASE test_<random>` из административного пула, затем пул к ней и `db.Migrate`.

- [ ] **Шаг 6: Запустить тесты**

Run: `cd backend && go test ./internal/platform/db/ -v`
Expected: 2 PASS; первый запуск скачивает образ.

- [ ] **Шаг 7: `sqlc.yaml`**

`version: "2"`, один пакет `queries` → `internal/platform/db/sqlc` (пока без файлов запросов — на этапе 1 запросы появятся по модулям; в `sqlc.yaml` оставить комментарий на русском, где и как добавлять). Добавить `sqlc generate` в `backend/Taskfile.yml` → `gen`, но только если есть `*.sql` в `queries/` (условие `sources`).

- [ ] **Шаг 8: Зафиксировать**

```bash
git add backend
git commit -m "feat: подключить PostgreSQL, миграции goose и тестовый контейнер"
```

---

### Task 8: HTTP-сервер, генерация из OpenAPI, health и ready

**Files:**
- Create: `backend/internal/api/gen/oapi-codegen.yaml`, `backend/internal/api/gen/generate.go`, `backend/internal/api/gen/api.gen.go` (генерируется), `backend/internal/platform/httpserver/server.go`, `backend/internal/platform/httpserver/middleware.go`, `backend/internal/platform/httpserver/problem.go`, `backend/internal/health/handler.go`, `backend/internal/health/handler_test.go`, `backend/cmd/api/main.go`, `backend/Dockerfile`

**Interfaces:**
- Consumes: `api/openapi.yaml` (Задача 5), `config.Load` (Задача 6), `db.NewPool`, `testdb.New` (Задача 7).
- Produces: `httpserver.New(log *slog.Logger, handlers gen.ServerInterface) *http.Server`-подобный конструктор: `httpserver.NewRouter(log *slog.Logger, h gen.ServerInterface) http.Handler` (chi, middleware: request id, recover, структурированный лог запроса, `Content-Type` для ошибок `application/problem+json`); `httpserver.WriteProblem(w http.ResponseWriter, status int, code, title, detail string)`; `health.New(version string, pinger interface{ Ping(context.Context) error }) *health.Handler`, реализующий `gen.ServerInterface`.

- [ ] **Шаг 1: Добавить chi и oapi-codegen**

Run: `cd backend && go get github.com/go-chi/chi/v5@latest && go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest`
Expected: `go.mod` содержит `tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen`.

- [ ] **Шаг 2: Конфигурация генерации**

`oapi-codegen.yaml`: `package: gen`, `generate: { chi-server: true, models: true, embedded-spec: true, strict-server: false }`, `output: api.gen.go`. `generate.go`: `//go:generate go tool oapi-codegen -config oapi-codegen.yaml ../../../../api/openapi.yaml`.

Run: `cd backend && go generate ./... && grep -c "GetHealth\|GetReady" internal/api/gen/api.gen.go`
Expected: ≥ 2; интерфейс `ServerInterface` с методами `GetHealth(w, r)` и `GetReady(w, r)`.

- [ ] **Шаг 3: Написать падающий тест `health/handler_test.go`**

```go
type fakePinger struct{ err error }
func (f fakePinger) Ping(context.Context) error { return f.err }

func TestHealth_ReturnsOkAndVersion(t *testing.T) {
	h := health.New("1.2.3", fakePinger{})
	srv := httptest.NewServer(httpserver.NewRouter(slog.New(slog.DiscardHandler), h))
	defer srv.Close()
	res, _ := http.Get(srv.URL + "/api/v1/health")
	if res.StatusCode != 200 { t.Fatalf("status = %d", res.StatusCode) }
	var body gen.HealthResponse
	json.NewDecoder(res.Body).Decode(&body)
	if body.Status != gen.HealthResponseStatusOk || body.Version != "1.2.3" { t.Errorf("body = %+v", body) }
}

func TestReady_DatabaseUp(t *testing.T) { /* fakePinger{nil} → 200, database == "ok" */ }

func TestReady_DatabaseDown(t *testing.T) { /* fakePinger{errors.New("down")} → 503, status == "degraded", database == "unavailable" */ }

func TestUnknownRoute_ReturnsProblemJSON(t *testing.T) {
	// GET /api/v1/nope → 404, Content-Type application/problem+json, поле code == "not_found"
}
```

- [ ] **Шаг 4: Запустить и убедиться, что падает**

Run: `cd backend && go test ./internal/health/`
Expected: FAIL — пакеты `health`, `httpserver` не определены.

- [ ] **Шаг 5: Реализовать `httpserver` и `health`**

`middleware.go`: `RequestID` (UUID из `crypto/rand`, заголовок `X-Request-Id`), `Recover` (в лог с `slog.Error`, ответ 500 problem), `Logger` (метод, путь, статус, длительность, request id). `problem.go`: структура `Problem` по §7, `WriteProblem`. `server.go`: `NewRouter` монтирует `gen.HandlerWithOptions(h, gen.ChiServerOptions{BaseURL: "/api/v1", BaseRouter: r})`, `r.NotFound` и `r.MethodNotAllowed` пишут problem. `health/handler.go`: `GetReady` вызывает `pinger.Ping` с таймаутом 2 с.

- [ ] **Шаг 6: Запустить тесты**

Run: `cd backend && go test ./... -race`
Expected: все PASS.

- [ ] **Шаг 7: `cmd/api/main.go`**

`var version = "dev"`; `config.Load(os.Getenv)` → при ошибке `slog.Error` + `os.Exit(2)`; логгер JSON с уровнем из конфигурации; `db.NewPool` → `db.Migrate` при старте (флаг `UNITRACK_MIGRATE_ON_START`, по умолчанию `true`) → `health.New(version, pool)` → `http.Server` с таймаутами (`ReadHeaderTimeout: 5s`, `ReadTimeout: 15s`, `WriteTimeout: 30s`, `IdleTimeout: 60s`) → graceful shutdown по `SIGINT`/`SIGTERM` за 10 с.

Run: `cd backend && UNITRACK_DATABASE_URL=postgres://invalid go run ./cmd/api; echo "exit=$?"`
Expected: лог с ошибкой подключения на русском, `exit=1` (не паника).

- [ ] **Шаг 8: `backend/Dockerfile`**

Multi-stage: `golang:1.24-alpine` сборка с `CGO_ENABLED=0` и `-ldflags "-s -w -X main.version=${VERSION}"` (`ARG VERSION=dev`); финальный образ `gcr.io/distroless/static-debian12:nonroot`, `EXPOSE 8080`, `ENTRYPOINT ["/api"]`.

Run: `cd backend && docker build --build-arg VERSION=test -t unitrack-api:test . && docker run --rm -e UNITRACK_DATABASE_URL=postgres://invalid unitrack-api:test; echo "exit=$?"`
Expected: образ собирается; контейнер завершается с кодом 1 и логом ошибки.

- [ ] **Шаг 9: Зафиксировать**

```bash
git add backend
git commit -m "feat: добавить HTTP-сервер на chi, генерацию из OpenAPI и операции health/ready"
```

---

### Task 9: Веб-каркас — Vite, TypeScript, линтеры FSD и токенов

**Files:**
- Create: `web/package.json`, `web/vite.config.ts`, `web/tsconfig.json`, `web/tsconfig.app.json`, `web/tsconfig.node.json`, `web/eslint.config.js`, `web/.stylelintrc.json`, `web/.prettierrc`, `web/Taskfile.yml`, `web/index.html`, `web/src/app/{main.ts,App.vue,router/index.ts,providers/{i18n.ts,query.ts}}`, `web/src/app/styles/{reset.css,base.css}`, `web/src/shared/i18n/ru.json`, `web/src/shared/api/{client.ts,schema.d.ts}`, `web/src/shared/config/env.ts`, `web/public/fonts/*.woff2`, папки `pages/ widgets/ features/ entities/` с `.gitkeep`
- Modify: `Taskfile.yml` (убрать `optional` у `web`), `.gitattributes` при необходимости

**Interfaces:**
- Produces: FSD-слои `app > pages > widgets > features > entities > shared` с алиасом `@/`; правило: слой импортирует только нижележащие, слайсы одного слоя не импортируют друг друга; публичный API слайса — `index.ts`.
- Produces: `shared/api/client.ts` — `export const api = createClient<paths>({ baseUrl: env.apiBaseUrl })` из openapi-fetch; `schema.d.ts` генерируется `openapi-typescript ../api/openapi.yaml`.
- Produces: `task web:gen|lint|test|build|dev`; stylelint-правило «сырые значения запрещены».

- [ ] **Шаг 1: Создать проект**

Run: `cd /root/projects/unitrack && pnpm create vite web --template vue-ts && cd web && pnpm install`
Expected: `web/` с шаблоном Vue + TS.

- [ ] **Шаг 2: Установить зависимости**

Run (в `web/`):
```bash
pnpm add vue-router pinia @tanstack/vue-query vue-i18n zod reka-ui @tanstack/vue-table openapi-fetch
pnpm add -D eslint @eslint/js typescript-eslint eslint-plugin-vue eslint-plugin-boundaries eslint-config-prettier prettier stylelint stylelint-config-standard stylelint-declaration-strict-value stylelint-order vitest @vue/test-utils jsdom vue-tsc openapi-typescript @redocly/cli
```
Expected: `package.json` обновлён; версии закреплены в `pnpm-lock.yaml`.

- [ ] **Шаг 3: TypeScript строгий, алиас, структура**

`tsconfig.app.json`: `strict`, `noUncheckedIndexedAccess`, `noImplicitOverride`, `exactOptionalPropertyTypes`, `verbatimModuleSyntax`, `paths: { "@/*": ["./src/*"] }`. `vite.config.ts`: алиас `@` → `src`, `server.proxy['/api'] → http://localhost:8080`. Удалить шаблонные `HelloWorld.vue`, `style.css`, `assets/`. Создать слои и `.gitkeep`.

- [ ] **Шаг 4: ESLint с границами FSD**

`eslint.config.js` (flat): `typescript-eslint` strict-type-checked, `eslint-plugin-vue` flat/recommended, `eslint-config-prettier`; `boundaries/elements` — шесть слоёв по glob `src/<layer>/**`; `boundaries/element-types`: `app → [pages, widgets, features, entities, shared]`, `pages → [widgets, features, entities, shared]`, `widgets → [features, entities, shared]`, `features → [entities, shared]`, `entities → [shared]`, `shared → [shared]`; `boundaries/no-private` для импортов мимо `index.ts`; `vue/block-lang: script ts`; `vue/component-api-style: ['script-setup']`; `@typescript-eslint/no-explicit-any: error`; `vue/no-bare-strings-in-template: error` (запрет строк интерфейса в шаблонах).

- [ ] **Шаг 5: Отрицательная проверка boundaries**

Создать временный `src/shared/lib/bad.ts` с `import {} from '@/features/x'` (и пустой `src/features/x/index.ts`).

Run: `cd web && pnpm eslint src/shared/lib/bad.ts`
Expected: ошибка `boundaries/element-types`. После проверки удалить оба файла.

- [ ] **Шаг 6: stylelint с запретом сырых значений**

`.stylelintrc.json`: `stylelint-config-standard`, плагины `stylelint-declaration-strict-value`, `stylelint-order`; правило `scale-unlimited/declaration-strict-value: [["/color/", "fill", "stroke", "background", "background-color", "border", "border-color", "box-shadow", "font-family", "font-size", "border-radius"], { "ignoreValues": ["transparent", "inherit", "currentColor", "none", "0", "initial", "unset"] }]`; `selector-class-pattern` — БЭМ `^[a-z]([a-z0-9-]+)?(__([a-z0-9]+-?)+)?(--([a-z0-9]+-?)+){0,2}$`; исключение `src/shared/ui/tokens.css` через `ignoreFiles`. Проверять `src/**/*.{css,vue}` с `postcss-html` синтаксисом (`customSyntax` для `.vue`).

- [ ] **Шаг 7: Отрицательная проверка stylelint**

Временный `src/shared/ui/bad.css` с `.x { color: #16A964; }`.

Run: `cd web && pnpm stylelint src/shared/ui/bad.css`
Expected: ошибка `scale-unlimited/declaration-strict-value`. Удалить файл.

- [ ] **Шаг 8: Шрифт, стили, провайдеры, приложение**

Скачать Inter (woff2, веса 400/500/600/700, латиница + кириллица) в `public/fonts/`; `app/styles/base.css` с `@font-face` и базовыми стилями на токенах (`body { font-family: var(--font-family-base); color: var(--color-ink-900); background: var(--color-bg-default); }`). `app/providers/i18n.ts`: `createI18n({ legacy: false, locale: 'ru', fallbackLocale: 'ru', messages: { ru } })`. `app/providers/query.ts`: `VueQueryPlugin` с `staleTime: 30_000`, `retry: 1`. `app/router/index.ts`: `createRouter` с историей и единственным маршрутом `/` → `pages/home` (появится в Задаче 10; на этом шаге — заглушка компонента внутри `app/`). `shared/config/env.ts`: `apiBaseUrl = import.meta.env.VITE_API_BASE_URL ?? '/api/v1'` с проверкой через zod. `shared/api/client.ts` как в Interfaces. `main.ts` собирает всё; `App.vue` — `<RouterView />`.

- [ ] **Шаг 9: Генерация клиента и `web/Taskfile.yml`**

`package.json scripts`: `gen:api` = `openapi-typescript ../api/openapi.yaml -o src/shared/api/schema.d.ts`, `lint` = `vue-tsc --noEmit && eslint . && stylelint "src/**/*.{css,vue}" && prettier --check .`, `test` = `vitest run`, `build` = `vite build`, `dev` = `vite`. `web/Taskfile.yml` оборачивает их задачами `gen`, `lint`, `test`, `build`, `dev`. Корневой `Taskfile.yml`: `gen` добавляет `web:gen`, `lint` — `web:lint`, `test` — `web:test`, `check` — проверка diff включает `web/src/shared/api/schema.d.ts`.

Run: `task gen && cd web && pnpm lint && pnpm build`
Expected: `schema.d.ts` содержит `"/health"` и `"/ready"`; lint и build без ошибок.

- [ ] **Шаг 10: Зафиксировать**

```bash
git add web Taskfile.yml .gitattributes
git commit -m "feat: завести веб-каркас Vue 3 с FSD, строгим TypeScript и линтерами токенов"
```

---

### Task 10: Базовый набор `shared/ui` и первая страница

**Files:**
- Create: `web/src/shared/ui/{index.ts, button/UiButton.vue, text-field/UiTextField.vue, checkbox/UiCheckbox.vue, select/UiSelect.vue, dialog/UiDialog.vue, tabs/UiTabs.vue, status/UiStatus.vue, empty-state/UiEmptyState.vue, spinner/UiSpinner.vue, page-header/UiPageHeader.vue}`, тесты `*.test.ts` рядом с компонентами, `web/src/entities/system/{index.ts, api/health.ts, model/use-health.ts}`, `web/src/pages/home/{index.ts, HomePage.vue}`, `web/vitest.config.ts`
- Modify: `web/src/app/router/index.ts`, `web/src/shared/i18n/ru.json`

**Interfaces:**
- Consumes: токены (Задача 4), `api` клиент и `paths` (Задача 9), Reka UI.
- Produces: компоненты с префиксом `Ui*`, экспорт только через `shared/ui/index.ts`. Пропсы: `UiButton { variant?: 'primary' | 'ghost' | 'danger'; size?: 'md' | 'sm'; loading?: boolean; type?: 'button' | 'submit' }`; `UiTextField { modelValue: string; label: string; error?: string; hint?: string }` (`defineModel`); `UiCheckbox { modelValue: boolean; label: string }`; `UiSelect<T> { modelValue: T; options: { value: T; label: string }[]; label: string }`; `UiDialog { open: boolean; title: string }` со слотами `default`, `footer`; `UiTabs { modelValue: string; tabs: { value: string; label: string }[] }`; `UiStatus { tone: 'ok' | 'wait' | 'danger' | 'neutral'; label: string }` — точка плюс слово по §6.5; `UiEmptyState { title: string; description?: string }`; `UiSpinner`; `UiPageHeader { title: string; subtitle?: string }` со слотом `actions`.
- Produces: `entities/system`: `getHealth()` через `api.GET('/health')`; `useHealthQuery()` на TanStack Query с ключом `['system', 'health']`.

- [ ] **Шаг 1: `vitest.config.ts`** — `environment: 'jsdom'`, `globals: false`, `include: ['src/**/*.test.ts']`, плагин vue, алиас `@`.

- [ ] **Шаг 2: Падающие тесты для трёх компонентов, задающих правила**

`button/UiButton.test.ts`: рендер с `variant="primary"` даёт класс `ui-button ui-button--primary`; `loading` ставит `disabled` и `aria-busy="true"`; клик по disabled не эмитит. `status/UiStatus.test.ts`: класс `ui-status--ok`, текст равен `label`, внутри есть элемент `.ui-status__dot` с `aria-hidden="true"`. `text-field/UiTextField.test.ts`: `label` связан с `input` через `for/id`; при `error` у `input` `aria-invalid="true"` и `aria-describedby` указывает на элемент с текстом ошибки; ввод эмитит `update:modelValue`.

Run: `cd web && pnpm vitest run`
Expected: FAIL — компоненты не найдены.

- [ ] **Шаг 3: Реализовать десять компонентов**

БЭМ-классы `ui-<name>`, `ui-<name>__<element>`, `ui-<name>--<modifier>`; стили в `<style scoped>` только на токенах; `UiDialog`, `UiSelect`, `UiTabs`, `UiCheckbox` — на примитивах Reka UI (`DialogRoot/…`, `SelectRoot/…`, `TabsRoot/…`, `CheckboxRoot`). Тексты внутри компонентов (например, «Закрыть» у диалога) — через `useI18n` и `ru.json` раздел `ui`. `shared/ui/index.ts` реэкспортирует всё.

- [ ] **Шаг 4: Запустить тесты и линт**

Run: `cd web && pnpm vitest run && pnpm lint`
Expected: тесты PASS; stylelint и eslint без ошибок (ни одного сырого значения, ни одной голой строки в шаблонах).

- [ ] **Шаг 5: `entities/system` и `pages/home`**

`HomePage.vue`: `UiPageHeader` с заголовком из `ru.json` (`pages.home.title = "UniTrack"`, `subtitle = "Система учёта посещаемости"`), ниже `UiStatus` с состоянием сервера: `tone="ok"` и «Сервер доступен, версия {version}» при успехе, `tone="danger"` и «Сервер недоступен» при ошибке, `UiSpinner` во время загрузки. Маршрут `/` → `HomePage` через `pages/home/index.ts`.

Тест `pages/home/HomePage.test.ts`: с замоканным `useHealthQuery` (через `vi.mock('@/entities/system')`) в состоянии успеха рендерится текст с версией `0.1.0`; в состоянии ошибки — класс `ui-status--danger`.

Run: `cd web && pnpm vitest run && pnpm build`
Expected: PASS; `dist/` собран.

- [ ] **Шаг 6: Зафиксировать**

```bash
git add web
git commit -m "feat: добавить базовый набор компонентов shared/ui и главную страницу со статусом сервера"
```

---

### Task 11: Инфраструктура — docker compose, nginx, окружение

**Files:**
- Create: `infra/docker-compose.yml`, `infra/.env.example`, `infra/nginx/nginx.conf`, `infra/web.Dockerfile`, `infra/Taskfile.yml`, `infra/README.md`
- Modify: `Taskfile.yml` (убрать `optional` у `infra`; `dev` → `infra:up`)

**Interfaces:**
- Consumes: `backend/Dockerfile` (Задача 8), сборка `web/` (Задача 10), переменные `UNITRACK_*` (Задача 6).
- Produces: `task infra:up|down|logs|ps`; сервисы `postgres` (образ `postgres:17-alpine`, том `pgdata`, порт не публикуется наружу), `api` (сборка из `../backend`, healthcheck на `/api/v1/ready`, `depends_on: postgres: condition: service_healthy`), `web` (сборка `web.Dockerfile`: `node:24-alpine` → `pnpm build` → `nginx:1.27-alpine` со статикой и `nginx.conf`; порт `${UNITRACK_HTTP_PORT:-8080}:80`).

- [ ] **Шаг 1: `.env.example`**

Переменные с комментариями на русском: `POSTGRES_PASSWORD` (сгенерировать: `openssl rand -base64 24`), `UNITRACK_HTTP_PORT=8080`, `UNITRACK_LOG_LEVEL=info`, `UNITRACK_ENV=production`, `UNITRACK_VERSION=dev`. `UNITRACK_DATABASE_URL` собирается внутри compose из `POSTGRES_PASSWORD`, в `.env` не дублируется.

- [ ] **Шаг 2: `nginx/nginx.conf`**

`server` на 80: `root /usr/share/nginx/html`; `location /api/ { proxy_pass http://api:8080; proxy_set_header X-Request-Id $request_id; proxy_read_timeout 60s; }`; `location /apk/ { alias /srv/apk/; autoindex off; }` (заглушка под этап 9, папка монтируется пустой); `location / { try_files $uri /index.html; }`; `location ~* \.(woff2|js|css)$ { expires 1y; add_header Cache-Control "public, immutable"; }`; `index.html` — `Cache-Control: no-cache`; `gzip on` для текста; заголовки `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `Referrer-Policy: same-origin`. TLS не настраивается в этапе 0 — комментарий со ссылкой на открытый вопрос №4 проектного документа.

- [ ] **Шаг 3: `docker-compose.yml` и `web.Dockerfile`**

Как в Interfaces. `api.environment`: `UNITRACK_DATABASE_URL=postgres://unitrack:${POSTGRES_PASSWORD}@postgres:5432/unitrack?sslmode=disable`, `UNITRACK_HTTP_ADDR=:8080`, `UNITRACK_LOG_LEVEL`, `UNITRACK_ENV`. `api.build.args.VERSION=${UNITRACK_VERSION}`. `api.healthcheck.test`: distroless не содержит curl — healthcheck делает `web` (`wget -qO- http://api:8080/api/v1/ready`) либо для `api` используется `test: ["CMD", "/api", "-healthcheck"]` — **выбрать второй**: добавить в `cmd/api/main.go` флаг `-healthcheck`, который делает HTTP-запрос к собственному `/api/v1/ready` и выходит с кодом 0/1 (стандартная библиотека, 20 строк).

- [ ] **Шаг 4: `infra/Taskfile.yml` и корневой `dev`**

`up` = `docker compose --env-file .env up -d --build`, `down`, `logs` (`-f api`), `ps`. Корневой `task dev` → `infra:up` + подсказка «веб: http://localhost:8080, API: http://localhost:8080/api/v1/health».

- [ ] **Шаг 5: Запустить всё и проверить**

Run:
```bash
cp infra/.env.example infra/.env && sed -i "s/^POSTGRES_PASSWORD=.*/POSTGRES_PASSWORD=$(openssl rand -base64 24 | tr -d '/+=')/" infra/.env
task dev && sleep 10 && docker compose -f infra/docker-compose.yml ps --format '{{.Name}} {{.Status}}'
curl -s http://localhost:8080/api/v1/health; echo
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/
```
Expected: три контейнера `healthy`/`running`; `{"status":"ok","version":"dev"}`; `200`.

- [ ] **Шаг 6: Проверить, что `api` не становится healthy без базы**

Run: `docker compose -f infra/docker-compose.yml stop postgres && sleep 15 && docker compose -f infra/docker-compose.yml ps api --format '{{.Status}}'; docker compose -f infra/docker-compose.yml start postgres`
Expected: статус содержит `unhealthy`.

- [ ] **Шаг 7: `infra/README.md`** — установка на сервер института: требования (Docker, 2 ГБ ОЗУ), шаги (`git clone`, `.env`, `task infra:up`), обновление (`git pull && task infra:up`), где логи, что резервные копии появятся на этапе 6.

- [ ] **Шаг 8: Зафиксировать**

```bash
git add infra Taskfile.yml backend/cmd/api/main.go
git commit -m "infra: собрать docker compose с PostgreSQL, API и nginx со статикой"
```

---

### Task 12: Pre-commit и CI

**Files:**
- Create: `lefthook.yml`, `.github/workflows/backend.yml`, `.github/workflows/web.yml`, `.github/workflows/contract.yml`
- Modify: `Taskfile.yml` (`setup` ставит lefthook и `lefthook install`)

**Interfaces:**
- Consumes: `task` цели из Задач 6, 9, 11.
- Produces: гарантии §11.4: коммит не проходит без lint; CI падает при расхождении генерации с контрактом.

- [ ] **Шаг 1: `lefthook.yml`**

`pre-commit` параллельно: `go` (glob `backend/**/*.go`: `gofumpt -l` → ошибка если вывод непустой; `golangci-lint run --new-from-rev=HEAD` в `backend/`), `web` (glob `web/**/*.{ts,vue,css,json}`: `prettier --check`, `eslint`, `stylelint` по изменённым файлам), `tokens` (glob `design/tokens.json`: `node design/build-tokens.mjs && git diff --exit-code web/src/shared/ui/tokens.css`), `api` (glob `api/**`: `task gen && git diff --exit-code backend/internal/api/gen web/src/shared/api/schema.d.ts`). `commit-msg`: скрипт-проверка префикса `^(feat|fix|api|docs|infra|chore|test|refactor)(\(.+\))?: .+` с сообщением на русском при ошибке. `pre-push`: `task test`.

Run: `task setup && git commit --allow-empty -m "плохое сообщение"; echo "exit=$?"`
Expected: `exit=1`, сообщение о неверном префиксе. Затем: `git commit --allow-empty -m "chore: проверить хук" && git reset --soft HEAD~1`.

- [ ] **Шаг 2: `backend.yml`**

Триггер: push/PR при изменениях в `backend/**`, `api/**`. Шаги: `actions/setup-go` с `go-version-file: backend/go.mod`, `services: postgres:17-alpine` с `TESTDB_URL` (Задача 7, режим без контейнеров), `task backend:lint`, `task backend:test`, `docker build backend`.

- [ ] **Шаг 3: `web.yml`**

Триггер: `web/**`, `api/**`, `design/**`. Шаги: `pnpm/action-setup`, `actions/setup-node` с `node-version-file`, `pnpm install --frozen-lockfile`, `node --test design/`, `task web:lint`, `task web:test`, `task web:build`.

- [ ] **Шаг 4: `contract.yml`**

Триггер: любой push/PR. Шаги: установить Go, Node, task; `task gen`; `git diff --exit-code -- backend/internal/api/gen web/src/shared/api/schema.d.ts web/src/shared/ui/tokens.css`; `pnpm dlx @redocly/cli lint api/openapi.yaml`. Это тест пункта 3 Review Focus.

- [ ] **Шаг 5: Проверить workflow-файлы локально**

Run: `pnpm dlx action-validator@latest .github/workflows/*.yml 2>/dev/null || python3 -I -c "import yaml,glob;[yaml.safe_load(open(f)) for f in glob.glob('.github/workflows/*.yml')];print('yaml ok')"`
Expected: без ошибок разбора.

- [ ] **Шаг 6: Зафиксировать**

```bash
git add lefthook.yml .github Taskfile.yml
git commit -m "chore: настроить pre-commit через lefthook и CI на GitHub Actions"
```

---

### Task 13: CLAUDE.md — корневой и по частям, настройки Claude

**Files:**
- Create: `CLAUDE.md`, `backend/CLAUDE.md`, `web/CLAUDE.md`, `infra/CLAUDE.md`, `.claude/settings.json`

**Interfaces:**
- Consumes: всё созданное выше — команды `task …`, формат ADR, `docs/STATUS.md`.
- Produces: протокол сессии §11 как обязательные инструкции.

- [ ] **Шаг 1: Корневой `CLAUDE.md`** (не более 150 строк)

Разделы в этом порядке:
1. **Что это** — три предложения из §1; ссылка на проектный документ и `docs/spec/requirements.md`.
2. **Протокол сессии** — дословно §11.1: *начало* (прочитать `docs/STATUS.md`, `git log --oneline -10`, открытый план в `docs/superpowers/plans/`), *во время* (решение вне ТЗ → ADR в `docs/decisions/`; изменение `api/openapi.yaml` → `task gen` и коммит с префиксом `api:`), *конец* (`task check` зелёный, `docs/STATUS.md` обновлён, коммит сделан; «сессия без обновлённого STATUS считается незавершённой»).
3. **Команды** — таблица `task setup|dev|gen|lint|test|check` с одной строкой каждая.
4. **Структура репозитория** — дерево из §9 с одной строкой на папку и указанием, что в `backend/`, `web/`, `infra/` лежат свои CLAUDE.md.
5. **Правила кода, общие** — §10 «Общие» дословно, плюс: «Предыдущая реализация eduControl не является образцом; не ссылаться на неё».
6. **Запреты** — список: не удалять данные (только статусы); не добавлять зависимостей без ADR; не писать строки интерфейса в коде; не обходить линтеры (`eslint-disable`, `nolint`) без комментария с причиной; не коммитить `.env`; не использовать внешние сервисы во время работы системы.
7. **Где что решено** — ссылки на ADR-индекс и разделы проектного документа.

- [ ] **Шаг 2: `backend/CLAUDE.md`**

30–50 строк: структура модуля (`internal/<module>/{handler.go,service.go,repo.go,queries.sql}`), правило трёх слоёв и направления зависимостей, «интерфейс при второй реализации», как добавить эндпоинт (контракт → `task gen` → реализовать метод `ServerInterface` → тест через `httptest` → при нужде `testdb.New`), как добавить миграцию (`goose create -dir internal/platform/db/migrations <name> sql`), правила ошибок и логов (§10 Go дословно), что генерируется и не правится руками (`*.gen.go`, `sqlc/`).

- [ ] **Шаг 3: `web/CLAUDE.md`**

30–50 строк: слои FSD и что в каждом лежит; правило публичного API слайса через `index.ts`; §10 Vue дословно; где токены и как добавить; как сделать компонент `shared/ui` (БЭМ-класс `ui-<name>`, тест рядом, экспорт через индекс, тексты через i18n); серверные данные только через TanStack Query в `entities/*/model`, запросы в `entities/*/api` через `api` клиент; формы через zod и `shared/lib/form` (появится на этапе 2).

- [ ] **Шаг 4: `infra/CLAUDE.md`**

20–30 строк: что в compose, как обновлять на сервере, что порт PostgreSQL не публикуется, что секреты только в `.env`, TLS — открытый вопрос, резервные копии — этап 6.

- [ ] **Шаг 5: `.claude/settings.json`**

`permissions.allow`: `Bash(task *)`, `Bash(go test *)`, `Bash(go build *)`, `Bash(go vet *)`, `Bash(golangci-lint *)`, `Bash(pnpm lint*)`, `Bash(pnpm test*)`, `Bash(pnpm build*)`, `Bash(pnpm vitest *)`, `Bash(docker compose *)`, `Bash(git status*)`, `Bash(git diff*)`, `Bash(git log*)`, `Bash(node --test *)`. `permissions.deny`: `Read(./.env)`, `Read(./infra/.env)`, `Bash(git push --force*)`, `Bash(rm -rf *)`.

- [ ] **Шаг 6: Проверить ограничение длины**

Run: `wc -l CLAUDE.md backend/CLAUDE.md web/CLAUDE.md infra/CLAUDE.md`
Expected: корневой ≤ 150, остальные ≤ 60.

- [ ] **Шаг 7: Зафиксировать**

```bash
git add CLAUDE.md backend/CLAUDE.md web/CLAUDE.md infra/CLAUDE.md .claude/settings.json
git commit -m "docs: написать CLAUDE.md с протоколом сессии и правилами для каждой части"
```

---

### Task 14: Полная проверка, статус, публикация

**Files:**
- Modify: `docs/STATUS.md`, `README.md` (если команды изменились)

- [ ] **Шаг 1: Чистая проверка с нуля**

Run: `task check && task dev && sleep 10 && curl -sf http://localhost:8080/api/v1/health && curl -sf -o /dev/null http://localhost:8080/ && task infra:down`
Expected: код 0 на каждой команде.

- [ ] **Шаг 2: Обновить `docs/STATUS.md`**

`## Сделано`: «Этап 0 — каркас (2026-11-xx): репозиторий, ТЗ в Markdown, 12 ADR, токены, контракт health/ready, Go-API с миграциями и тестами на PostgreSQL, веб-каркас FSD с 10 компонентами, docker compose, CI, lefthook, CLAUDE.md». `## В работе`: пусто. `## Следующее`: «Этап 1 — контракт и ядро: auth, users, structure, students, calendar, reasons, reports (см. проектный документ §12); первым делом — план `docs/superpowers/plans/<дата>-stage-1-core.md`». Дата последнего обновления.

- [ ] **Шаг 3: Зафиксировать и опубликовать**

```bash
git add docs/STATUS.md README.md
git commit -m "docs: закрыть этап 0 в STATUS.md"
git push -u origin main
```
Expected: ветка `main` в `https://github.com/OrazDurdyyew/unitrack`; через несколько минут три workflow зелёные (проверить на вкладке Actions).

- [ ] **Шаг 4: Если CI красный** — исправить причину отдельным коммитом с префиксом `fix:`, не отключая проверки. Повторить до зелёного.

---

## Самопроверка плана

**Покрытие спеки (этап 0 из §12):** репозиторий — З1; CLAUDE.md — З13; STATUS — З1, З14; ADR — З3; Taskfile — З6, З9, З11; CI — З12; docker compose с пустым API и статикой — З8, З11; ТЗ в Markdown — З2; токены — З4; набор `shared/ui` — З10. Контракт сначала (§6.1, §7) — З5, З8, З9. Линтеры-гарантии (§11.4) — З9, З12. Пробелы: Playwright (§6.2) — осознанно перенесён на этап 2, когда появятся сквозные сценарии; резервные копии (§6.4) — этап 6 по таблице этапов.

**Согласованность имён:** `config.Load`, `db.NewPool`, `db.Migrate`, `testdb.New`, `httpserver.NewRouter`, `httpserver.WriteProblem`, `health.New`, `gen.ServerInterface`, `api` (openapi-fetch), `useHealthQuery`, `Ui*` — используются одинаково в З6–З10, З13.

**Review Focus → тесты:** 1 → З9 ш.7; 2 → З9 ш.5; 3 → З12 ш.4; 4 → З6 `TestLoad_MissingDatabaseURL`; 5 → З8 `TestReady_DatabaseDown` + З11 ш.6.
