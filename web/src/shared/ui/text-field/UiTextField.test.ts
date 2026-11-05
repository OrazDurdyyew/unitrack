import { mount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';

import UiTextField from './UiTextField.vue';

describe('UiTextField', () => {
  it('связывает подпись с полем через for/id', () => {
    const w = mount(UiTextField, { props: { modelValue: '', label: 'Логин' } });
    const input = w.find('input');
    const label = w.find('label');
    expect(input.attributes('id')).toBeTruthy();
    expect(label.attributes('for')).toBe(input.attributes('id'));
  });

  it('при ошибке помечает поле aria-invalid и ссылается на текст ошибки', () => {
    const w = mount(UiTextField, {
      props: { modelValue: '', label: 'Логин', error: 'Обязательное поле' },
    });
    const input = w.find('input');
    expect(input.attributes('aria-invalid')).toBe('true');
    const describedBy = input.attributes('aria-describedby');
    expect(describedBy).toBeTruthy();
    expect(w.find(`#${describedBy ?? ''}`).text()).toBe('Обязательное поле');
  });

  it('эмитит update:modelValue при вводе', async () => {
    const w = mount(UiTextField, { props: { modelValue: '', label: 'Логин' } });
    await w.find('input').setValue('oraz');
    expect(w.emitted('update:modelValue')?.[0]).toEqual(['oraz']);
  });
});
