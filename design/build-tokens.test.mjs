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
