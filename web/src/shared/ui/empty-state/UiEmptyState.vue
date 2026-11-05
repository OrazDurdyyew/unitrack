<script setup lang="ts">
defineProps<{ title: string; description?: string }>();
</script>

<template>
  <div class="ui-empty-state">
    <h3 class="ui-empty-state__title">{{ title }}</h3>
    <p v-if="description" class="ui-empty-state__description">{{ description }}</p>
    <div v-if="$slots.default" class="ui-empty-state__actions"><slot /></div>
  </div>
</template>

<style scoped>
.ui-empty-state {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  align-items: center;
  padding: var(--space-8) var(--space-6);
  text-align: center;
  border: 1px dashed var(--color-line-strong);
  border-radius: var(--radius-md);
}

.ui-empty-state__title {
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--color-ink-700);
}

.ui-empty-state__description {
  max-width: 420px;
  font-size: var(--font-size-md);
  color: var(--color-ink-500);
}

.ui-empty-state__actions {
  margin-top: var(--space-2);
}
</style>
