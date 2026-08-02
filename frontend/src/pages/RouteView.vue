<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useRoutingStore } from '@/stores/routing'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const store = useRoutingStore()
const router = useRouter()

const taskType = ref('')
const capabilities = ref('')
const tier = ref('economy')
const source = ref('manual')

onMounted(() => store.fetchLogs())

async function submit() {
  const caps = capabilities.value
    .split(',')
    .map((c) => c.trim())
    .filter(Boolean)
  await store.route(source.value.trim(), taskType.value.trim(), caps, tier.value)
}

function openLog(id: string) {
  router.push(`/routes/${id}`)
}

function fmt(v: any): string {
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleString()
}
</script>

<template>
  <div class="p-6 space-y-6 overflow-y-auto h-full">
    <h2 class="text-sm font-semibold">{{ t('route.title') }}</h2>

    <section class="border border-border rounded-lg p-4 space-y-3 max-w-2xl">
      <input v-model="taskType" :placeholder="t('route.form.task')" class="w-full text-sm border border-border rounded px-2 py-1.5" />
      <div class="grid grid-cols-2 gap-2">
        <input
          v-model="capabilities"
          :placeholder="t('route.form.capabilities')"
          class="text-sm border border-border rounded px-2 py-1.5"
        />
        <select v-model="tier" class="text-sm border border-border rounded px-2 py-1.5">
          <option value="economy">{{ t('tier.economy') }}</option>
          <option value="standard">{{ t('tier.standard') }}</option>
          <option value="premium">{{ t('tier.premium') }}</option>
        </select>
      </div>
      <input v-model="source" :placeholder="t('route.form.source')" class="w-full text-sm border border-border rounded px-2 py-1.5" />
      <button :disabled="store.routing" class="text-sm px-3 py-1.5 rounded text-white disabled:opacity-40" style="background: #5b6ee1" @click="submit">
        {{ t('route.form.submit') }}
      </button>

      <div v-if="store.routeError" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
        {{ t('error.prefix') }}{{ store.routeError }}
      </div>

      <div v-if="store.lastResult" class="border border-border rounded-lg p-3 space-y-1 bg-gray-50">
        <div class="text-sm font-medium">
          {{ t('route.result.chosen') }}:
          <span v-if="store.lastResult.chosen_model" style="color: #5b6ee1">
            {{ store.lastResult.chosen_model.provider }}/{{ store.lastResult.chosen_model.model_id }}
          </span>
          <span v-else class="text-red-600">{{ t('route.result.none') }}</span>
        </div>
        <div class="text-xs text-muted-foreground">{{ t('route.result.reasoning') }}: {{ store.lastResult.log.reasoning }}</div>
      </div>
    </section>

    <section class="space-y-2">
      <h3 class="text-xs font-semibold text-muted-foreground">{{ t('history.title') }}</h3>
      <div v-if="store.logsError" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
        {{ t('error.prefix') }}{{ store.logsError }}
        <button class="ml-2 underline" @click="store.fetchLogs">{{ t('error.retry') }}</button>
      </div>
      <div v-if="store.logsLoading" class="text-xs text-muted-foreground">{{ t('loading') }}</div>
      <div v-else-if="store.logs.length === 0" class="text-xs text-muted-foreground">{{ t('history.empty') }}</div>
      <div v-else class="space-y-1.5">
        <div
          v-for="log in store.logs"
          :key="log.id"
          class="flex items-center gap-3 border border-border rounded-lg px-4 py-2.5 border-l-4"
          :style="{ borderLeftColor: log.chosen_model_id ? '#5b6ee1' : '#d03b3b' }"
        >
          <div class="flex-1">
            <div class="text-sm font-medium">
              {{ log.task_type || '(no task type)' }} —
              <span v-if="log.chosen_model_id">{{ log.chosen_model_id }}</span>
              <span v-else class="text-red-600">{{ t('history.none') }}</span>
            </div>
            <div class="text-xs text-muted-foreground">
              <span v-if="log.source">{{ log.source }} · </span>{{ fmt(log.requested_at) }}
            </div>
          </div>
          <button class="text-xs px-2 py-1 border border-border rounded hover:bg-gray-50" @click="openLog(log.id)">
            {{ t('history.open') }}
          </button>
        </div>
      </div>
    </section>
  </div>
</template>
