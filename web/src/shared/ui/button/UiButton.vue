<script setup lang="ts">
import UiSpinner from '../spinner/UiSpinner.vue';

withDefaults(
  defineProps<{
    variant?: 'primary' | 'ghost' | 'danger';
    size?: 'md' | 'sm';
    loading?: boolean;
    disabled?: boolean;
    type?: 'button' | 'submit';
  }>(),
  { variant: 'primary', size: 'md', loading: false, disabled: false, type: 'button' },
);

const emit = defineEmits<{ click: [event: MouseEvent] }>();
</script>

<template>
  <button
    class="ui-button"
    :class="[`ui-button--${variant}`, `ui-button--${size}`, { 'ui-button--loading': loading }]"
    :type="type"
    :disabled="disabled || loading"
    :aria-busy="loading ? 'true' : undefined"
    @click="(e) => !disabled && !loading && emit('click', e)"
  >
    <UiSpinner v-if="loading" class="ui-button__spinner" size="sm" />
    <span class="ui-button__label"><slot /></span>
  </button>
</template>

<style scoped>
.ui-button {
  display: inline-flex;
  gap: var(--space-2);
  align-items: center;
  justify-content: center;
  height: 36px;
  padding: 0 var(--space-4);
  font-size: var(--font-size-md);
  font-weight: 600;
  line-height: 1;
  white-space: nowrap;
  border-width: 1px;
  border-style: solid;
  border-color: transparent;
  border-radius: var(--radius-sm);
  transition:
    background-color 0.12s,
    border-color 0.12s,
    color 0.12s;
}

.ui-button--sm {
  height: 30px;
  padding: 0 var(--space-3);
  font-size: var(--font-size-sm);
}

.ui-button--primary {
  color: var(--color-surface-default);
  background: var(--color-green-600);
  box-shadow: var(--shadow-sm);
}

.ui-button--primary:hover:not(:disabled) {
  background: var(--color-green-700);
}

.ui-button--ghost {
  color: var(--color-ink-700);
  background: var(--color-surface-default);
  border-color: var(--color-line-strong);
}

.ui-button--ghost:hover:not(:disabled) {
  background: var(--color-bg-default);
}

.ui-button--danger {
  color: var(--color-red-700);
  background: var(--color-red-50);
}

.ui-button--danger:hover:not(:disabled) {
  color: var(--color-surface-default);
  background: var(--color-red-700);
}

.ui-button:disabled {
  cursor: not-allowed;
  opacity: 0.55;
}

.ui-button__spinner {
  flex: none;
}
</style>
