# web — правила для Vue 3 и FSD

Vue 3 + Vite + TypeScript strict, архитектура Feature-Sliced Design (ADR 0005), чистый CSS по БЭМ на токенах, Reka UI (ADR 0006), TanStack Query и zod (ADR 0007).

## Слои (сверху вниз)

```
src/app/        точка входа, провайдеры (i18n, query), роутер, глобальные стили, тестовая настройка
src/pages/      страницы — слайс на маршрут (pages/home)
src/widgets/    крупные составные блоки, переиспользуемые между страницами
src/features/   действия пользователя (отправить отчёт, утвердить, подать запрос)
src/entities/   предметные сущности: api/ (запросы), model/ (TanStack Query, типы), ui/ (карточки, строки)
src/shared/     ui/ (компоненты Ui*), api/ (клиент), lib/, config/, i18n/
```

Слой импортирует только нижележащие. Слайсы одного слоя друг друга не импортируют. Публичный API слайса — `index.ts`, импорт мимо него — ошибка ESLint (`boundaries/dependencies`). Алиас `@/` → `src/`.

## Правила

- Только `<script setup lang="ts">`, Composition API. Пропсы — через `defineProps<…>()`, двусторонние — `defineModel`.
- Никаких `any`. Никаких строк интерфейса в коде и шаблонах — только `t('…')` из `vue-i18n`, словарь `shared/i18n/ru.json` (`vue/no-bare-strings-in-template`).
- Запросы к серверу — только в `entities/*/api` через `api` из `@/shared/api` (сгенерированные типы). Состояние серверных данных — только TanStack Query в `entities/*/model`. Pinia — для клиентского состояния (сессия, настройки интерфейса).
- Стили: `<style scoped>`, классы по БЭМ `block__element--modifier`. Цвета, размеры шрифта, радиусы, тени, семейство шрифта — только `var(--…)` из `shared/ui/tokens.css`; stylelint падает на сыром значении. Новый токен: `design/tokens.json` → `task gen`.
- Компонент больше 200 строк делится. Права и роли — через `shared/lib/access` (появится на этапе 1), не разрозненными условиями.
- Даты — через одну утилиту в `shared/lib/date` с поясом Asia/Ashgabat (этап 2).
- Формы — схема zod + composable `shared/lib/form` (этап 2); vee-validate не используется.

## Как сделать компонент shared/ui

Папка `shared/ui/<name>/Ui<Name>.vue` + `Ui<Name>.test.ts` рядом; корневой класс `ui-<name>`; тексты через i18n (раздел `ui` в `ru.json`); экспорт через `shared/ui/index.ts`. Доступность: подписи связаны с полями (`for/id`), ошибки — `aria-invalid` + `aria-describedby`, статус — точка `aria-hidden` плюс слово.

## Тесты

Vitest + Vue Test Utils, файлы `*.test.ts` рядом с кодом, i18n подключён глобально в `src/app/test-setup.ts`. Серверные данные в тестах страниц — через `vi.mock('@/entities/<name>')`. Playwright для сквозных сценариев — с этапа 2.

## Команды

`task web:install | gen | lint | test | build | dev`. `pnpm lint:fix` — автоисправление ESLint, stylelint, Prettier.
