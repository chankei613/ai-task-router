import { defineStore } from 'pinia'
import { ListModels, CreateModel, SetModelEnabled, DeleteModel } from '../../wailsjs/go/main/App'
import { db } from '../../wailsjs/go/models'

export const useModelsStore = defineStore('models', {
  state: () => ({
    models: [] as db.ModelSpec[],
    loading: false,
    error: null as string | null,
  }),
  actions: {
    async fetchModels() {
      this.loading = true
      this.error = null
      try {
        this.models = (await ListModels()) ?? []
      } catch (e) {
        this.error = String(e)
      } finally {
        this.loading = false
      }
    },
    async createModel(
      provider: string,
      modelId: string,
      tier: string,
      inputPrice: number,
      outputPrice: number,
      capabilities: string[],
    ) {
      this.error = null
      try {
        await CreateModel(provider, modelId, tier, inputPrice, outputPrice, capabilities)
        await this.fetchModels()
      } catch (e) {
        this.error = String(e)
      }
    },
    async setEnabled(id: string, enabled: boolean) {
      this.error = null
      try {
        await SetModelEnabled(id, enabled)
        await this.fetchModels()
      } catch (e) {
        this.error = String(e)
      }
    },
    async deleteModel(id: string) {
      this.error = null
      try {
        await DeleteModel(id)
        await this.fetchModels()
      } catch (e) {
        this.error = String(e)
      }
    },
  },
})
