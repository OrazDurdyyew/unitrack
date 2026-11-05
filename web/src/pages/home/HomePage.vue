<script setup lang="ts">
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>

<template>
  <main class="home-page">
    <h1 class="home-page__title">{{ t('app.name') }}</h1>
  </main>
</template>

<style scoped>
.home-page {
  padding: var(--space-6);
}

.home-page__title {
  font-size: var(--font-size-2xl);
}
</style>
