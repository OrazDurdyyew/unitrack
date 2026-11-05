<script setup lang="ts">
import { TabsIndicator, TabsList, TabsRoot, TabsTrigger } from 'reka-ui';

defineProps<{ tabs: { value: string; label: string }[]; ariaLabel?: string }>();

const model = defineModel<string>({ required: true });
</script>

<template>
  <TabsRoot v-model="model" class="ui-tabs">
    <TabsList class="ui-tabs__list" :aria-label="ariaLabel">
      <TabsIndicator class="ui-tabs__indicator" />
      <TabsTrigger v-for="tab in tabs" :key="tab.value" :value="tab.value" class="ui-tabs__trigger">
        {{ tab.label }}
      </TabsTrigger>
    </TabsList>
    <slot />
  </TabsRoot>
</template>

<style scoped>
.ui-tabs__list {
  position: relative;
  display: flex;
  gap: var(--space-4);
  border-bottom: 1px solid var(--color-line-default);
}

.ui-tabs__trigger {
  padding: var(--space-2) 0 var(--space-3);
  font-size: var(--font-size-md);
  font-weight: 500;
  color: var(--color-ink-500);
  background: transparent;
  border: 0;
}

.ui-tabs__trigger[data-state='active'] {
  font-weight: 600;
  color: var(--color-green-700);
}

.ui-tabs__indicator {
  position: absolute;
  bottom: -1px;
  left: 0;
  width: var(--reka-tabs-indicator-size);
  height: 2px;
  background: var(--color-green-600);
  transition:
    width 0.15s,
    transform 0.15s;
  transform: translateX(var(--reka-tabs-indicator-position));
}
</style>
