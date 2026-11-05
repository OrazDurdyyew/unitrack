<script setup lang="ts">
import {
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogOverlay,
  DialogPortal,
  DialogRoot,
  DialogTitle,
} from 'reka-ui';
import { useI18n } from 'vue-i18n';

defineProps<{ title: string; description?: string }>();

const open = defineModel<boolean>('open', { required: true });
const { t } = useI18n();
</script>

<template>
  <DialogRoot v-model:open="open">
    <DialogPortal>
      <DialogOverlay class="ui-dialog__overlay" />
      <DialogContent class="ui-dialog">
        <header class="ui-dialog__header">
          <DialogTitle class="ui-dialog__title">{{ title }}</DialogTitle>
          <DialogClose class="ui-dialog__close" :aria-label="t('ui.close')">
            <svg viewBox="0 0 16 16" aria-hidden="true"><path d="M4 4l8 8M12 4l-8 8" /></svg>
          </DialogClose>
        </header>
        <DialogDescription v-if="description" class="ui-dialog__description">{{
          description
        }}</DialogDescription>
        <div class="ui-dialog__body"><slot /></div>
        <footer v-if="$slots.footer" class="ui-dialog__footer"><slot name="footer" /></footer>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>

<style scoped>
.ui-dialog__overlay {
  position: fixed;
  inset: 0;
  z-index: 40;
  background: var(--color-overlay-default);
}

.ui-dialog {
  position: fixed;
  top: 50%;
  left: 50%;
  z-index: 41;
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  width: min(520px, calc(100vw - var(--space-8)));
  max-height: calc(100vh - var(--space-8));
  padding: var(--space-5) var(--space-6);
  overflow: auto;
  background: var(--color-surface-default);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-md);
  transform: translate(-50%, -50%);
}

.ui-dialog__header {
  display: flex;
  gap: var(--space-3);
  align-items: flex-start;
  justify-content: space-between;
}

.ui-dialog__title {
  font-size: var(--font-size-xl);
  font-weight: 600;
  color: var(--color-ink-900);
}

.ui-dialog__close {
  display: grid;
  flex: none;
  place-items: center;
  width: 28px;
  height: 28px;
  padding: 0;
  background: transparent;
  border: 0;
  border-radius: var(--radius-sm);
}

.ui-dialog__close:hover {
  background: var(--color-bg-default);
}

.ui-dialog__close svg {
  width: 14px;
  height: 14px;
  fill: none;
  stroke: var(--color-ink-500);
  stroke-width: 1.8;
  stroke-linecap: round;
}

.ui-dialog__description {
  font-size: var(--font-size-md);
  color: var(--color-ink-500);
}

.ui-dialog__footer {
  display: flex;
  gap: var(--space-2);
  justify-content: flex-end;
}
</style>
