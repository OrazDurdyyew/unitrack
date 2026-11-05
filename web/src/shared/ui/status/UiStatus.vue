<script setup lang="ts">
defineProps<{
  tone: 'ok' | 'wait' | 'danger' | 'neutral';
  label: string;
}>();
</script>

<template>
  <span class="ui-status" :class="`ui-status--${tone}`">
    <span class="ui-status__dot" aria-hidden="true" />
    <span class="ui-status__label">{{ label }}</span>
  </span>
</template>

<style scoped>
.ui-status {
  display: inline-flex;
  gap: var(--space-2);
  align-items: center;
  padding: var(--space-1) var(--space-2);
  font-size: var(--font-size-sm);
  font-weight: 600;
  line-height: 1.2;
  border-radius: var(--radius-full);
}

.ui-status__dot {
  width: 6px;
  height: 6px;
  background: currentcolor;
  border-radius: var(--radius-full);
}

.ui-status--ok {
  color: var(--color-green-700);
  background: var(--color-green-50);
}

.ui-status--wait {
  color: var(--color-amber-700);
  background: var(--color-amber-50);
}

.ui-status--danger {
  color: var(--color-red-700);
  background: var(--color-red-50);
}

.ui-status--neutral {
  color: var(--color-ink-500);
  background: var(--color-bg-default);
}
</style>
