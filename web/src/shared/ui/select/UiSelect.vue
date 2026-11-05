<script setup lang="ts" generic="T extends string | number">
import {
  SelectContent,
  SelectItem,
  SelectItemIndicator,
  SelectItemText,
  SelectPortal,
  SelectRoot,
  SelectTrigger,
  SelectValue,
  SelectViewport,
} from 'reka-ui';
import { computed, useId } from 'vue';

const props = defineProps<{
  label: string;
  options: { value: T; label: string }[];
  placeholder?: string;
  disabled?: boolean;
}>();

const model = defineModel<T>({ required: true });
const id = useId();

// Reka UI работает со строковыми значениями; храним соответствие строка → исходное значение.
const asString = (v: T): string => String(v);
const byString = computed(() => new Map(props.options.map((o) => [asString(o.value), o.value])));
const stringModel = computed({
  get: () => asString(model.value),
  set: (s: string) => {
    const v = byString.value.get(s);
    if (v !== undefined) model.value = v;
  },
});
</script>

<template>
  <div class="ui-select">
    <label class="ui-select__label" :for="id">{{ label }}</label>
    <SelectRoot v-model="stringModel" :disabled="disabled">
      <SelectTrigger :id="id" class="ui-select__trigger">
        <SelectValue :placeholder="placeholder ?? ''" />
        <svg viewBox="0 0 16 16" aria-hidden="true" class="ui-select__chevron">
          <path d="M4 6l4 4 4-4" />
        </svg>
      </SelectTrigger>
      <SelectPortal>
        <SelectContent class="ui-select__content" position="popper" :side-offset="4">
          <SelectViewport class="ui-select__viewport">
            <SelectItem
              v-for="o in options"
              :key="asString(o.value)"
              :value="asString(o.value)"
              class="ui-select__item"
            >
              <SelectItemText>{{ o.label }}</SelectItemText>
              <SelectItemIndicator class="ui-select__check">
                <svg viewBox="0 0 16 16" aria-hidden="true">
                  <path d="M3.5 8.5l2.8 2.8L12.5 5" />
                </svg>
              </SelectItemIndicator>
            </SelectItem>
          </SelectViewport>
        </SelectContent>
      </SelectPortal>
    </SelectRoot>
  </div>
</template>

<style scoped>
.ui-select {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.ui-select__label {
  font-size: var(--font-size-sm);
  font-weight: 600;
  color: var(--color-ink-700);
}

.ui-select__trigger {
  display: inline-flex;
  gap: var(--space-2);
  align-items: center;
  justify-content: space-between;
  height: 36px;
  padding: 0 var(--space-3);
  font-size: var(--font-size-md);
  color: var(--color-ink-900);
  background: var(--color-surface-default);
  border-width: 1px;
  border-style: solid;
  border-color: var(--color-line-strong);
  border-radius: var(--radius-sm);
}

.ui-select__trigger[data-placeholder] {
  color: var(--color-ink-400);
}

.ui-select__trigger:focus-visible {
  border-color: var(--color-green-600);
  outline: none;
  box-shadow: var(--shadow-focus);
}

.ui-select__chevron {
  width: 14px;
  height: 14px;
  fill: none;
  stroke: var(--color-ink-400);
  stroke-width: 1.8;
}

.ui-select__content {
  z-index: 50;
  min-width: var(--reka-select-trigger-width);
  background: var(--color-surface-default);
  border-width: 1px;
  border-style: solid;
  border-color: var(--color-line-default);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-md);
}

.ui-select__viewport {
  padding: var(--space-1);
}

.ui-select__item {
  display: flex;
  gap: var(--space-2);
  align-items: center;
  justify-content: space-between;
  padding: var(--space-2) var(--space-3);
  font-size: var(--font-size-md);
  color: var(--color-ink-900);
  border-radius: var(--radius-sm);
  outline: none;
}

.ui-select__item[data-highlighted] {
  background: var(--color-green-50);
}

.ui-select__check svg {
  width: 12px;
  height: 12px;
  fill: none;
  stroke: var(--color-green-700);
  stroke-width: 2.2;
}
</style>
