import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useRoutingStore } from '@/stores/routing'
import { Route, ListRoutingLogs } from '../../wailsjs/go/main/App'
import { db, api } from '../../wailsjs/go/models'

describe('routing store error handling', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.mocked(Route).mockReset()
    vi.mocked(ListRoutingLogs).mockReset()
  })

  it('captures a failed fetchLogs() as store.logsError and clears loading', async () => {
    vi.mocked(ListRoutingLogs).mockRejectedValueOnce(new Error('network down'))
    const store = useRoutingStore()

    await store.fetchLogs()

    expect(store.logsLoading).toBe(false)
    expect(store.logsError).toContain('network down')
  })

  it('clears the previous error on a successful retry', async () => {
    vi.mocked(ListRoutingLogs).mockRejectedValueOnce(new Error('network down'))
    const store = useRoutingStore()
    await store.fetchLogs()
    expect(store.logsError).not.toBeNull()

    vi.mocked(ListRoutingLogs).mockResolvedValueOnce([])
    await store.fetchLogs()

    expect(store.logsError).toBeNull()
  })

  it('route() stores the result and refreshes history', async () => {
    vi.mocked(Route).mockResolvedValueOnce(
      api.RouteResult.createFrom({
        log: db.RoutingLog.createFrom({ id: 'r1', task_type: 'summarize', chosen_model_id: 'anthropic:claude-haiku-4-5', reasoning: 'cheapest' }),
        chosen_model: db.ModelSpec.createFrom({ id: 'anthropic:claude-haiku-4-5', provider: 'anthropic', model_id: 'claude-haiku-4-5' }),
      }),
    )
    vi.mocked(ListRoutingLogs).mockResolvedValueOnce([
      db.RoutingLog.createFrom({ id: 'r1', task_type: 'summarize', chosen_model_id: 'anthropic:claude-haiku-4-5' }),
    ])

    const store = useRoutingStore()
    await store.route('manual', 'summarize', [], 'economy')

    expect(store.routeError).toBeNull()
    expect(store.lastResult?.chosen_model?.model_id).toBe('claude-haiku-4-5')
    expect(store.logs).toHaveLength(1)
  })
})
