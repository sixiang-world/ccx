import { ref, computed } from 'vue'
import { useAdminApi } from '@/composables/useAdminApi'
import { useConsoleChannels } from '@/composables/useConsoleChannels'
import type {
  CapabilitySnapshot,
  CapabilityTestJob,
  CapabilityTestJobStartResponse,
  CapabilityProtocolJobResult,
  CapabilityLifecycle,
  CapabilityOutcome,
} from '@/services/admin-api'
import { getChannelTypeApi } from '@/utils/channel-type-api'
import type { ManagedChannelType } from '@/utils/channel-type-api'

// Module-level singletons
const activeJob = ref<CapabilityTestJob | null>(null)
const snapshot = ref<CapabilitySnapshot | null>(null)
const loading = ref(false)
const polling = ref(false)
const cancelling = ref(false)
const error = ref('')
const pollers = new Map<string, ReturnType<typeof setInterval>>()
const POLL_INTERVAL = 1000

export function useCapabilityTests() {
  const api = useAdminApi()
  const { refreshChannels } = useConsoleChannels()

  function clearError() {
    error.value = ''
  }

  // ── 基础 CRUD ──

  async function startTest(
    channelType: string,
    channelId: number,
    options?: { targetProtocols?: string[]; models?: string[]; rpm?: number; previousJobId?: string; sourceTab?: string },
  ) {
    loading.value = true
    clearError()
    try {
      const resp = await api.post<CapabilityTestJobStartResponse>(
        `/api/${channelType}/channels/${channelId}/capability-test`,
        options,
      )
      activeJob.value = resp.job || null
      if (resp.jobId) {
        startPolling(channelType, channelId, resp.jobId)
      }
      // 也为协议子 job 启动轮询
      const job = resp.job
      if (job?.protocolJobRefs) {
        for (const [, ref] of Object.entries(job.protocolJobRefs)) {
          if (ref.jobId && ref.jobId !== resp.jobId) {
            startPolling(channelType, channelId, ref.jobId)
          }
        }
      }
      return resp
    } catch (e) {
      error.value = e instanceof Error ? e.message : String(e)
      throw e
    } finally {
      loading.value = false
    }
  }

  /** 按协议启动单独测试（WebUI 的 testProtocol） */
  async function startProtocolTest(
    channelType: string,
    channelId: number,
    protocol: string,
    models?: string[],
    rpm?: number,
  ) {
    return startTest(channelType, channelId, {
      targetProtocols: [protocol],
      models,
      rpm,
      sourceTab: channelType,
    })
  }

  async function fetchSnapshot(channelType: string, channelId: number, sourceTab?: string) {
    try {
      const url = sourceTab
        ? `/api/${channelType}/channels/${channelId}/capability-snapshot?sourceTab=${sourceTab}`
        : `/api/${channelType}/channels/${channelId}/capability-snapshot`
      snapshot.value = await api.get<CapabilitySnapshot>(url)
    } catch {
      // snapshot 可能不存在，静默
    }
  }

  async function fetchJobStatus(channelType: string, channelId: number, jobId: string) {
    const job = await api.get<CapabilityTestJob>(
      `/api/${channelType}/channels/${channelId}/capability-test/${jobId}`,
    )
    activeJob.value = job
    if (job.status === 'completed' || job.status === 'failed' || job.status === 'cancelled') {
      stopPoller(jobId)
    }
    return job
  }

  // ── 轮询（支持多 job 并发） ──

  function startPolling(channelType: string, channelId: number, jobId: string) {
    if (pollers.has(jobId)) return
    polling.value = true
    const timer = setInterval(async () => {
      try {
        await fetchJobStatus(channelType, channelId, jobId)
      } catch (e) {
        stopPoller(jobId)
        error.value = e instanceof Error ? e.message : String(e)
      }
    }, POLL_INTERVAL)
    pollers.set(jobId, timer)
  }

  function stopPoller(jobId: string) {
    const timer = pollers.get(jobId)
    if (timer) {
      clearInterval(timer)
      pollers.delete(jobId)
    }
    if (pollers.size === 0) polling.value = false
  }

  function stopAllPolling() {
    for (const timer of pollers.values()) clearInterval(timer)
    pollers.clear()
    polling.value = false
  }

  // ── 取消 ──

  async function cancelTest(channelType: string, channelId: number, jobId: string) {
    cancelling.value = true
    try {
      await api.del(`/api/${channelType}/channels/${channelId}/capability-test/${jobId}`)
      stopAllPolling()
      // 重取快照
      await fetchSnapshot(channelType, channelId, channelType)
      // 检查是否有其他活跃 job 需要继续轮询
      if (snapshot.value?.protocolJobRefs) {
        for (const [, ref] of Object.entries(snapshot.value.protocolJobRefs)) {
          if (ref.jobId && ref.jobId !== jobId) {
            const job = await api.get<CapabilityTestJob>(
              `/api/${channelType}/channels/${channelId}/capability-test/${ref.jobId}`,
            )
            if (job.status === 'running' || job.status === 'queued') {
              startPolling(channelType, channelId, ref.jobId)
            }
          }
        }
      }
    } finally {
      cancelling.value = false
    }
  }

  // ── Retry ──

  async function retryModel(channelType: string, channelId: number, jobId: string) {
    await api.post(`/api/${channelType}/channels/${channelId}/capability-test/${jobId}/retry`)
    startPolling(channelType, channelId, jobId)
  }

  /** 指定协议+模型 retry（WebUI 的 retryCapabilityModel） */
  async function retryModelForProtocol(
    channelType: string,
    channelId: number,
    protocol: string,
    model: string,
  ) {
    const jobId = activeJob.value?.protocolJobRefs?.[protocol]?.jobId
    if (!jobId) {
      // 没有 jobId 则启动一个只测该模型的协议测试
      return startProtocolTest(channelType, channelId, protocol, [model])
    }
    await api.post(
      `/api/${channelType}/channels/${channelId}/capability-test/${jobId}/retry`,
      { protocol, model },
    )
    startPolling(channelType, channelId, jobId)
    // 也为子 job 启动轮询
    if (activeJob.value?.protocolJobRefs?.[protocol]?.jobId) {
      const subJobId = activeJob.value.protocolJobRefs[protocol].jobId
      if (subJobId !== jobId) startPolling(channelType, channelId, subJobId)
    }
  }

  // ── Copy to Tab ──

  /** 复制渠道配置到目标协议 tab（WebUI 的 copyToTab） */
  async function copyToTab(
    sourceChannelType: string,
    channelId: number,
    targetProtocol: string,
  ) {
    // 先获取当前渠道完整数据
    const typeApi = getChannelTypeApi(sourceChannelType as ManagedChannelType)
    const channels = await typeApi.getChannels()
    const sourceChannel = channels.find(ch => ch.index === channelId)
    if (!sourceChannel) {
      error.value = '找不到源渠道数据'
      return
    }

    // 构建 payload，去掉 index/status/latency 等 runtime 字段
    const payload: Record<string, unknown> = {}
    const skipKeys = new Set(['index', 'latency', 'status', 'disabledApiKeys', 'historicalApiKeys'])
    for (const [k, v] of Object.entries(sourceChannel)) {
      if (!skipKeys.has(k) && v !== undefined && v !== null) {
        payload[k] = v
      }
    }
    // 调整 serviceType 以适配目标协议
    payload.serviceType = getNativeServiceType(targetProtocol)

    const targetApi = getChannelTypeApi(targetProtocol as ManagedChannelType)
    await targetApi.addChannel(payload as any)
    await refreshChannels()
  }

  // ── Reset ──

  function reset() {
    stopAllPolling()
    activeJob.value = null
    snapshot.value = null
    error.value = ''
    loading.value = false
    cancelling.value = false
  }

  // ── Computed helpers ──

  /** 从 snapshot 或 activeJob 中获取协议结果 */
  const protocolResults = computed<CapabilityProtocolJobResult[]>(() => {
    if (activeJob.value?.tests?.length) return activeJob.value.tests
    if (snapshot.value?.tests?.length) return snapshot.value.tests
    return []
  })

  /** 从 job/snapshot 中获取兼容协议 */
  const compatibleProtocols = computed<string[]>(() => {
    return activeJob.value?.compatibleProtocols ?? snapshot.value?.compatibleProtocols ?? []
  })

  /** 聚合 job/snapshot 的 lifecycle */
  const lifecycle = computed<CapabilityLifecycle>(() => {
    return activeJob.value?.lifecycle ?? snapshot.value?.lifecycle ?? 'pending'
  })

  /** 聚合 job/snapshot 的 outcome */
  const outcome = computed<CapabilityOutcome>(() => {
    return activeJob.value?.outcome ?? snapshot.value?.outcome ?? 'unknown'
  })

  /** 是否处于活动状态 */
  const isActive = computed(() => {
    const s = activeJob.value?.status
    return s === 'running' || s === 'queued'
  })

  /** 对话框整体状态 */
  const state = computed<'initializing' | 'error' | 'idle' | 'pending' | 'running' | 'completed' | 'cancelled'>(() => {
    if (error.value) return 'error'
    if (loading.value && !activeJob.value && !snapshot.value) return 'initializing'
    const l = lifecycle.value
    if (l === 'pending') return 'pending'
    if (l === 'active') return 'running'
    if (l === 'cancelled') return 'cancelled'
    if (l === 'done') {
      const o = outcome.value
      if (o === 'success') return 'completed'
      if (o === 'partial') return 'completed'
      if (o === 'failed') return 'completed'
      return 'completed'
    }
    return 'idle'
  })

  return {
    activeJob,
    snapshot,
    loading,
    polling,
    cancelling,
    error,
    clearError,
    startTest,
    startProtocolTest,
    fetchSnapshot,
    fetchJobStatus,
    cancelTest,
    retryModel,
    retryModelForProtocol,
    copyToTab,
    reset,
    // computed helpers
    protocolResults,
    compatibleProtocols,
    lifecycle,
    outcome,
    isActive,
    state,
  }
}

/** 协议对应的原生 service type */
function getNativeServiceType(protocol: string): string {
  if (protocol === 'messages') return 'claude'
  if (protocol === 'chat') return 'openai'
  if (protocol === 'responses') return 'responses'
  if (protocol === 'gemini') return 'gemini'
  return 'openai'
}
