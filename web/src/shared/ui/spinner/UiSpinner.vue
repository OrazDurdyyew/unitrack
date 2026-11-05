<script setup lang="ts">
import { useI18n } from 'vue-i18n';

withDefaults(defineProps<{ size?: 'sm' | 'md' }>(), { size: 'md' });

const { t } = useI18n();
</script>

<template>
  <span
    class="ui-spinner"
    :class="`ui-spinner--${size}`"
    role="status"
    :aria-label="t('ui.loading')"
  />
</template>

<style scoped>
.ui-spinner {
  display: inline-block;
  width: 20px;
  height: 20px;
  border: 2px solid var(--color-line-strong);
  border-top-color: var(--color-green-600);
  border-radius: var(--radius-full);
  animation: ui-spinner-rotate 0.8s linear infinite;
}

.ui-spinner--sm {
  width: 14px;
  height: 14px;
}

@keyframes ui-spinner-rotate {
  to {
    transform: rotate(360deg);
  }
}
</style>
