import { mount } from '@vue/test-utils';
import { ref } from 'vue';
import { describe, expect, it, vi } from 'vitest';

import HomePage from './HomePage.vue';

const state = {
  data: ref<{ version: string } | undefined>(undefined),
  isError: ref(false),
  isPending: ref(false),
};

vi.mock('@/entities/system', () => ({
  useHealthQuery: () => state,
}));

describe('HomePage', () => {
  it('показывает версию сервера при успехе', () => {
    state.data.value = { version: '0.1.0' };
    state.isError.value = false;
    const w = mount(HomePage);
    expect(w.text()).toContain('0.1.0');
    expect(w.find('.ui-status--ok').exists()).toBe(true);
  });

  it('показывает статус danger при ошибке', () => {
    state.data.value = undefined;
    state.isError.value = true;
    const w = mount(HomePage);
    expect(w.find('.ui-status--danger').exists()).toBe(true);
  });
});
