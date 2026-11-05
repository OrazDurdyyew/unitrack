<script setup lang="ts">
import { useId } from 'vue';

const props = withDefaults(
  defineProps<{
    label: string;
    error?: string;
    hint?: string;
    type?: 'text' | 'password' | 'email' | 'tel' | 'search';
    placeholder?: string;
    disabled?: boolean;
    autocomplete?: string;
  }>(),
  { type: 'text', disabled: false },
);

const model = defineModel<string>({ required: true });

const id = useId();
const errorId = `${id}-error`;
const hintId = `${id}-hint`;

const describedBy = (): string | undefined => {
  if (props.error) return errorId;
  if (props.hint) return hintId;
  return undefined;
};
</script>

<template>
  <div class="ui-text-field" :class="{ 'ui-text-field--invalid': Boolean(error) }">
    <label class="ui-text-field__label" :for="id">{{ label }}</label>
    <input
      :id="id"
      v-model="model"
      class="ui-text-field__input"
      :type="type"
      :placeholder="placeholder"
      :disabled="disabled"
      :autocomplete="autocomplete"
      :aria-invalid="error ? 'true' : undefined"
      :aria-describedby="describedBy()"
    />
    <p v-if="error" :id="errorId" class="ui-text-field__error">{{ error }}</p>
    <p v-else-if="hint" :id="hintId" class="ui-text-field__hint">{{ hint }}</p>
  </div>
</template>

<style scoped>
.ui-text-field {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.ui-text-field__label {
  font-size: var(--font-size-sm);
  font-weight: 600;
  color: var(--color-ink-700);
}

.ui-text-field__input {
  height: 36px;
  padding: 0 var(--space-3);
  font-size: var(--font-size-md);
  color: var(--color-ink-900);
  background: var(--color-surface-default);
  border: 1px solid var(--color-line-strong);
  border-radius: var(--radius-sm);
  outline: none;
  transition: border-color 0.12s;
}

.ui-text-field__input::placeholder {
  color: var(--color-ink-400);
}

.ui-text-field__input:focus {
  border-color: var(--color-green-600);
  box-shadow: var(--shadow-focus);
}

.ui-text-field__input:disabled {
  color: var(--color-ink-400);
  background: var(--color-bg-default);
}

.ui-text-field--invalid .ui-text-field__input {
  border-color: var(--color-red-700);
}

.ui-text-field__error {
  font-size: var(--font-size-sm);
  color: var(--color-red-700);
}

.ui-text-field__hint {
  font-size: var(--font-size-sm);
  color: var(--color-ink-500);
}
</style>
