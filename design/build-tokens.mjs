// Генератор CSS-переменных из design/tokens.json.
// Запуск: node design/build-tokens.mjs  → пишет web/src/shared/ui/tokens.css
// Это единственный источник значений цветов, отступов, радиусов, шрифтов и теней
// для веба; Android получает те же токены отдельным генератором (этап 7).

import { readFileSync, writeFileSync, mkdirSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath, pathToFileURL } from 'node:url';

const HEADER = '/* Сгенерировано из design/tokens.json командой `task gen`. Не править руками. */';

/**
 * Рекурсивно превращает вложенный объект токенов в список пар [имя, значение].
 * Имя переменной — путь по ключам через дефис: color.green.600 → --color-green-600.
 * @param {Record<string, unknown>} node
 * @param {string[]} path
 * @returns {Array<[string, string]>}
 */
function flatten(node, path = []) {
  const out = [];
  for (const [key, value] of Object.entries(node)) {
    const next = [...path, key];
    if (value !== null && typeof value === 'object') {
      out.push(...flatten(/** @type {Record<string, unknown>} */ (value), next));
    } else {
      out.push([`--${next.join('-')}`, String(value)]);
    }
  }
  return out;
}

/**
 * @param {Record<string, unknown>} tokens
 * @returns {string} содержимое tokens.css
 */
export function tokensToCss(tokens) {
  const vars = flatten(tokens).sort(([a], [b]) => a.localeCompare(b, 'en'));
  const body = vars.map(([name, value]) => `  ${name}: ${value};`).join('\n');
  return `${HEADER}\n:root {\n${body}\n}\n`;
}

const isMain = process.argv[1] && import.meta.url === pathToFileURL(resolve(process.argv[1])).href;

if (isMain) {
  const here = dirname(fileURLToPath(import.meta.url));
  const source = resolve(here, 'tokens.json');
  const target = resolve(here, '..', 'web', 'src', 'shared', 'ui', 'tokens.css');
  const tokens = JSON.parse(readFileSync(source, 'utf8'));
  mkdirSync(dirname(target), { recursive: true });
  writeFileSync(target, tokensToCss(tokens));
  console.log(`tokens.css: записано ${flatten(tokens).length} переменных → ${target}`);
}
