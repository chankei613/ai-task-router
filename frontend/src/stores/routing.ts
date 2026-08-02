import { defineStore } from 'pinia'
import { Route, ListRoutingLogs, GetRoutingLog } from '../../wailsjs/go/main/App'
import { db, api } from '../../wailsjs/go/models'

export const useRoutingStore = defineStore('routing', {
  state: () => ({
    logs: [] as db.RoutingLog[],
    logsLoading: false,
    logsError: null as string | null,

    lastResult: null as api.RouteResult | null,
    routing: false,
    routeError: null as string | null,

    detail: null as api.RouteResult | null,
    detailLoading: false,
    detailError: null as string | null,
  }),
  actions: {
    async fetchLogs() {
      this.logsLoading = true
      this.logsError = null
      try {
        this.logs = ((await ListRoutingLogs()) ?? []).slice().reverse()
      } catch (e) {
        this.logsError = String(e)
      } finally {
        this.logsLoading = false
      }
    },
    async route(source: string, taskType: string, requiredCapabilities: string[], minQualityTier: string) {
      this.routing = true
      this.routeError = null
      try {
        this.lastResult = await Route(source, taskType, requiredCapabilities, minQualityTier)
        await this.fetchLogs()
      } catch (e) {
        this.routeError = String(e)
      } finally {
        this.routing = false
      }
    },
    async fetchDetail(id: string) {
      this.detailLoading = true
      this.detailError = null
      try {
        this.detail = await GetRoutingLog(id)
      } catch (e) {
        this.detailError = String(e)
      } finally {
        this.detailLoading = false
      }
    },
  },
})
