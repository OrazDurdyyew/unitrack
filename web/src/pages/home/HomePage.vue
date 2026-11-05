<script setup lang="ts">
import { useI18n } from 'vue-i18n';

import { useHealthQuery } from '@/entities/system';
import { UiPageHeader, UiSpinner, UiStatus } from '@/shared/ui';

const { t } = useI18n();
const health = useHealthQuery();
</script>

<template>
  <main class="home-page">
    <UiPageHeader :title="t('pages.home.title')" :subtitle="t('pages.home.subtitle')" />
    <section class="home-page__server">
      <UiSpinner v-if="health.isPending.value" />
      <UiStatus
        v-else-if="health.isError.value"
        tone="danger"
        :label="t('pages.home.serverDown')"
      />
      <UiStatus
        v-else-if="health.data.value"
        tone="ok"
        :label="t('pages.home.serverOk', { version: health.data.value.version })"
      />
    </section>
  </main>
</template>

<style scoped>
.home-page {
  max-width: 960px;
  padding: var(--space-6);
  margin: 0 auto;
}

.home-page__server {
  display: flex;
  align-items: center;
  min-height: 32px;
}
</style>
