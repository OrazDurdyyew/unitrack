import { test } from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';

const tokens = JSON.parse(readFileSync(new URL('./tokens.json', import.meta.url), 'utf8'));

test('семейство шрифта в токенах совпадает с тем, что регистрирует @fontsource-variable/inter', () => {
  const css = readFileSync(new URL('../web/node_modules/@fontsource-variable/inter/index.css', import.meta.url), 'utf8');
  const registered = /font-family:\s*'([^']+)'/.exec(css)?.[1];
  assert.ok(registered, 'fontsource не нашёл font-family');
  assert.ok(
    tokens.font.family.base.includes(`'${registered}'`),
    `токен --font-family-base «${tokens.font.family.base}» не содержит '${registered}' — браузер отрисует системный шрифт`,
  );
});

test('все значения токенов — непустые строки', () => {
  const walk = (node, path) => {
    for (const [k, v] of Object.entries(node)) {
      if (v !== null && typeof v === 'object') walk(v, [...path, k]);
      else assert.ok(typeof v === 'string' && v.length > 0, `${[...path, k].join('.')} пуст`);
    }
  };
  walk(tokens, []);
});
