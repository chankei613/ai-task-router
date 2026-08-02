<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useModelsStore } from '@/stores/models'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const store = useModelsStore()

const provider = ref('')
const modelId = ref('')
const tier = ref('standard')
const inputPrice = ref(0)
const outputPrice = ref(0)
const capabilities = ref('')
const saving = ref(false)

onMounted(() => store.fetchModels())

async function save() {
  if (!provider.value.trim() || !modelId.value.trim()) return
  saving.value = true
  const caps = capabilities.value
    .split(',')
    .map((c) => c.trim())
    .filter(Boolean)
  await store.createModel(provider.value.trim(), modelId.value.trim(), tier.value, inputPrice.value, outputPrice.value, caps)
  saving.value = false
  provider.value = ''
  modelId.value = ''
  inputPrice.value = 0
  outputPrice.value = 0
  capabilities.value = ''
}

async function remove(id: string) {
  if (!confirm(t('models.delete.confirm'))) return
  await store.deleteModel(id)
}
</script>

<template>
  <div class="p-6 space-y-6 overflow-y-auto h-full">
    <h2 class="text-sm font-semibold">{{ t('models.title') }}</h2>

    <div v-if="store.error" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
      {{ t('error.prefix') }}{{ store.error }}
      <button class="ml-2 underline" @click="store.fetchModels">{{ t('error.retry') }}</button>
    </div>

    <div class="border border-border rounded-lg p-4 space-y-3 max-w-2xl">
      <h3 class="text-xs font-semibold text-muted-foreground">{{ t('models.new') }}</h3>
      <div class="grid grid-cols-2 gap-2">
        <input v-model="provider" :placeholder="t('models.new.provider')" class="text-sm border border-border rounded px-2 py-1.5" />
        <input v-model="modelId" :placeholder="t('models.new.modelId')" class="text-sm border border-border rounded px-2 py-1.5" />
      </div>
      <div class="grid grid-cols-3 gap-2">
        <select v-model="tier" class="text-sm border border-border rounded px-2 py-1.5">
          <option value="economy">{{ t('tier.economy') }}</option>
          <option value="standard">{{ t('tier.standard') }}</option>
          <option value="premium">{{ t('tier.premium') }}</option>
        </select>
        <input v-model.number="inputPrice" type="number" min="0" step="0.01" :placeholder="t('models.new.inputPrice')" class="text-sm border border-border rounded px-2 py-1.5" />
        <input v-model.number="outputPrice" type="number" min="0" step="0.01" :placeholder="t('models.new.outputPrice')" class="text-sm border border-border rounded px-2 py-1.5" />
      </div>
      <input v-model="capabilities" :placeholder="t('models.new.capabilities')" class="w-full text-sm border border-border rounded px-2 py-1.5" />
      <button
        :disabled="saving || !provider.trim() || !modelId.trim()"
        class="text-sm px-3 py-1.5 rounded bg-gray-900 text-white disabled:opacity-40"
        @click="save"
      >
        {{ t('models.new.save') }}
      </button>
    </div>

    <div v-if="store.loading" class="text-sm text-muted-foreground">{{ t('loading') }}</div>
    <div v-else-if="store.models.length === 0" class="text-sm text-muted-foreground">{{ t('models.empty') }}</div>

    <div v-else class="space-y-1.5">
      <div
        v-for="m in store.models"
        :key="m.id"
        class="flex items-center gap-3 border border-border rounded-lg px-4 py-2.5 border-l-4"
        :style="{ borderLeftColor: m.enabled ? '#5b6ee1' : '#c9c9c9' }"
      >
        <div class="flex-1">
          <div class="text-sm font-medium">{{ m.provider }}/{{ m.model_id }}</div>
          <div class="text-xs text-muted-foreground">
            {{ t('tier.' + m.quality_tier) }} · ${{ m.input_price_per_1m }}/${{ m.output_price_per_1m }} per 1M
            <span v-if="m.capabilities?.length"> · {{ m.capabilities.join(', ') }}</span>
          </div>
        </div>
        <button
          class="text-xs px-2 py-1 border border-border rounded hover:bg-gray-50"
          @click="store.setEnabled(m.id, !m.enabled)"
        >
          {{ m.enabled ? t('models.enabled') : t('models.disabled') }}
        </button>
        <button class="text-xs text-red-600 hover:underline" @click="remove(m.id)">{{ t('models.delete') }}</button>
      </div>
    </div>
  </div>
</template>
