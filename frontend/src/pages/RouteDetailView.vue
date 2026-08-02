<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useRoutingStore } from '@/stores/routing'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const store = useRoutingStore()
const route = useRoute()
const router = useRouter()

const logId = computed(() => route.params.id as string)

function load() {
  store.fetchDetail(logId.value)
}

onMounted(load)
watch(logId, load)

function fmt(v: any): string {
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleString()
}
</script>

<template>
  <div class="p-6 space-y-6 overflow-y-auto h-full">
    <button class="text-xs text-muted-foreground hover:underline" @click="router.push('/route')">
      &larr; {{ t('detail.back') }}
    </button>

    <div v-if="store.detailError" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
      {{ t('error.prefix') }}{{ store.detailError }}
      <button class="ml-2 underline" @click="load">{{ t('error.retry') }}</button>
    </div>

    <div v-if="store.detailLoading" class="text-sm text-muted-foreground">{{ t('loading') }}</div>

    <template v-else-if="store.detail">
      <div>
        <h2 class="text-sm font-semibold">{{ t('detail.title') }}</h2>
        <p class="text-xs text-muted-foreground mt-0.5">
          {{ store.detail.log.task_type || '(no task type)' }}
          <span v-if="store.detail.log.source"> · {{ store.detail.log.source }}</span>
          · {{ fmt(store.detail.log.requested_at) }}
        </p>
      </div>

      <div class="border border-border rounded-lg p-4 space-y-2 max-w-2xl">
        <div class="text-sm font-medium">
          {{ t('route.result.chosen') }}:
          <span v-if="store.detail.chosen_model" style="color: #5b6ee1">
            {{ store.detail.chosen_model.provider }}/{{ store.detail.chosen_model.model_id }}
          </span>
          <span v-else class="text-red-600">{{ t('route.result.none') }}</span>
        </div>
        <div class="text-xs text-muted-foreground">{{ t('route.result.reasoning') }}: {{ store.detail.log.reasoning }}</div>
        <div v-if="store.detail.log.required_capabilities?.length" class="text-xs text-muted-foreground">
          {{ t('route.form.capabilities') }}: {{ store.detail.log.required_capabilities.join(', ') }}
        </div>
        <div class="text-xs text-muted-foreground">
          {{ t('route.form.tier') }}: {{ t('tier.' + store.detail.log.min_quality_tier) }}
        </div>
        <div v-if="store.detail.chosen_model" class="text-xs text-muted-foreground pt-2 border-t border-border">
          ${{ store.detail.chosen_model.input_price_per_1m }}/${{ store.detail.chosen_model.output_price_per_1m }} per 1M ·
          {{ store.detail.chosen_model.capabilities?.join(', ') }}
        </div>
      </div>
    </template>
  </div>
</template>
