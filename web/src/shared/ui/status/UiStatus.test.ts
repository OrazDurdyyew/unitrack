import { mount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';

import UiStatus from './UiStatus.vue';

describe('UiStatus', () => {
  it('показывает слово и точку с модификатором тона', () => {
    const w = mount(UiStatus, { props: { tone: 'ok', label: 'Утверждён' } });
    expect(w.classes()).toContain('ui-status--ok');
    expect(w.text()).toBe('Утверждён');
    const dot = w.find('.ui-status__dot');
    expect(dot.exists()).toBe(true);
    expect(dot.attributes('aria-hidden')).toBe('true');
  });
});
