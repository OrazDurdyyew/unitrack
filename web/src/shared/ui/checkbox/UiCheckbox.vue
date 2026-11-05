<script setup lang="ts">
import { CheckboxIndicator, CheckboxRoot } from 'reka-ui';
import { useId } from 'vue';

defineProps<{ label: string; disabled?: boolean }>();

const model = defineModel<boolean>({ required: true });
const id = useId();
</script>

<template>
  <div class="ui-checkbox">
    <CheckboxRoot :id="id" v-model="model" class="ui-checkbox__box" :disabled="disabled">
      <CheckboxIndicator class="ui-checkbox__indicator">
        <svg viewBox="0 0 16 16" aria-hidden="true" class="ui-checkbox__icon">
          <path d="M3.5 8.5l2.8 2.8L12.5 5" />
        </svg>
      </CheckboxIndicator>
    </CheckboxRoot>
    <label class="ui-checkbox__label" :for="id">{{ label }}</label>
  </div>
</template>

<style scoped>
.ui-checkbox {
  display: inline-flex;
  gap: var(--space-2);
  align-items: center;
}

.ui-checkbox__box {
  display: grid;
  flex: none;
  place-items: center;
  width: 18px;
  height: 18px;
  padding: 0;
  background: var(--color-surface-default);
  border: 1px solid var(--color-line-strong);
  border-radius: var(--radius-xs);
}

.ui-checkbox__box[data-state='checked'] {
  background: var(--color-green-600);
  border-color: var(--color-green-600);
}

.ui-checkbox__box[data-disabled] {
  cursor: not-allowed;
  opacity: 0.55;
}

.ui-checkbox__icon {
  width: 12px;
  height: 12px;
  fill: none;
  stroke: var(--color-surface-default);
  stroke-width: 2.2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.ui-checkbox__label {
  font-size: var(--font-size-md);
  color: var(--color-ink-900);
  cursor: pointer;
}
</style>
