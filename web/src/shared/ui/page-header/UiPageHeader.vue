<script setup lang="ts">
defineProps<{ title: string; subtitle?: string }>();
</script>

<template>
  <header class="ui-page-header">
    <div class="ui-page-header__text">
      <h1 class="ui-page-header__title">{{ title }}</h1>
      <p v-if="subtitle" class="ui-page-header__subtitle">{{ subtitle }}</p>
    </div>
    <div v-if="$slots.actions" class="ui-page-header__actions"><slot name="actions" /></div>
  </header>
</template>

<style scoped>
.ui-page-header {
  display: flex;
  gap: var(--space-4);
  align-items: flex-end;
  justify-content: space-between;
  margin-bottom: var(--space-5);
}

.ui-page-header__title {
  font-size: var(--font-size-2xl);
  font-weight: 600;
  letter-spacing: -0.015em;
  color: var(--color-ink-900);
}

.ui-page-header__subtitle {
  margin-top: var(--space-1);
  font-size: var(--font-size-md);
  color: var(--color-ink-500);
}

.ui-page-header__actions {
  display: flex;
  gap: var(--space-2);
}
</style>
