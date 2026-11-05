import { mount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';

import UiButton from './UiButton.vue';

describe('UiButton', () => {
  it('ставит БЭМ-модификатор по варианту', () => {
    const w = mount(UiButton, { props: { variant: 'primary' }, slots: { default: 'Сохранить' } });
    expect(w.classes()).toContain('ui-button');
    expect(w.classes()).toContain('ui-button--primary');
  });

  it('в состоянии loading блокируется и помечается aria-busy', () => {
    const w = mount(UiButton, { props: { loading: true }, slots: { default: 'Сохранить' } });
    expect(w.attributes('disabled')).toBeDefined();
    expect(w.attributes('aria-busy')).toBe('true');
  });

  it('не эмитит click, когда disabled', async () => {
    const w = mount(UiButton, { props: { disabled: true }, slots: { default: 'Сохранить' } });
    await w.trigger('click');
    expect(w.emitted('click')).toBeUndefined();
  });

  it('по умолчанию type="button", чтобы не отправлять формы случайно', () => {
    const w = mount(UiButton, { slots: { default: 'Ок' } });
    expect(w.attributes('type')).toBe('button');
  });
});
