<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { AlertTriangle, XCircle, HardDrive, Network, Clock, ShieldAlert, X } from 'lucide-vue-next'
import { useLanguage } from '@/composables/useLanguage'

const props = defineProps<{
  error: string
}>()

const emit = defineEmits<{
  dismiss: []
}>()

const { t } = useLanguage()

const dismissed = ref(false)

watch(
  () => props.error,
  (val) => {
    dismissed.value = false
    if (!val) dismissed.value = false
  },
)

const visible = computed(() => props.error && !dismissed.value)

type ErrorKind = 'binary' | 'port' | 'health' | 'permission' | 'generic'

interface DiagnosticInfo {
  kind: ErrorKind
  icon: typeof AlertTriangle
  title: string
  suggestions: string[]
  color: string
}

const patterns: { re: RegExp; kind: ErrorKind }[] = [
  { re: /未找到.*二进制|binary.*not.*found/i, kind: 'binary' },
  { re: /端口.*冲突|端口.*占用|no.*available.*port/i, kind: 'port' },
  { re: /connection.*refused|连接.*拒绝/i, kind: 'port' },
  { re: /health.*超时|等待.*health|health.*timeout/i, kind: 'health' },
  { re: /permission.*denied|权限|access.*denied|不允许/i, kind: 'permission' },
]

const kindDefaults = computed<Record<ErrorKind, Omit<DiagnosticInfo, 'kind'>>>(() => ({
  binary: {
    icon: HardDrive,
    title: t('diagnostic.binaryTitle'),
    color: 'text-amber-400',
    suggestions: [
      t('diagnostic.binarySuggestionBuild'),
      t('diagnostic.binarySuggestionCheckDataDir'),
      t('diagnostic.binarySuggestionDownload'),
    ],
  },
  port: {
    icon: Network,
    title: t('diagnostic.portTitle'),
    color: 'text-orange-400',
    suggestions: [
      t('diagnostic.portSuggestionInstance'),
      t('diagnostic.portSuggestionEnv'),
      t('diagnostic.portSuggestionInspect'),
    ],
  },
  health: {
    icon: Clock,
    title: t('diagnostic.healthTitle'),
    color: 'text-amber-400',
    suggestions: [
      t('diagnostic.healthSuggestionLogs'),
      t('diagnostic.healthSuggestionEnv'),
      t('diagnostic.healthSuggestionChannels'),
      t('diagnostic.healthSuggestionRestart'),
    ],
  },
  permission: {
    icon: ShieldAlert,
    title: t('diagnostic.permissionTitle'),
    color: 'text-rose-400',
    suggestions: [
      t('diagnostic.permissionSuggestionDataDir'),
      t('diagnostic.permissionSuggestionExecutable'),
      t('diagnostic.permissionSuggestionWindows'),
    ],
  },
  generic: {
    icon: XCircle,
    title: t('diagnostic.genericTitle'),
    color: 'text-rose-400',
    suggestions: [
      t('diagnostic.genericSuggestionLogs'),
      t('diagnostic.genericSuggestionRestart'),
    ],
  },
}))

const diagnostic = computed<DiagnosticInfo>(() => {
  const msg = props.error
  for (const { re, kind } of patterns) {
    if (re.test(msg)) {
      return { kind, ...kindDefaults.value[kind] }
    }
  }
  return { kind: 'generic', ...kindDefaults.value.generic }
})
</script>

<template>
  <div
    v-if="visible"
    class="rounded-lg border border-rose-500/20 bg-rose-500/5 backdrop-blur-sm px-4 py-3"
  >
    <div class="flex items-start gap-3">
      <component :is="diagnostic.icon" :class="['h-5 w-5 mt-0.5 shrink-0', diagnostic.color]" />
      <div class="flex-1 min-w-0 space-y-2">
        <div class="flex items-center justify-between gap-2">
          <h4 :class="['text-sm font-semibold', diagnostic.color]">{{ diagnostic.title }}</h4>
          <button
            class="text-slate-500 hover:text-slate-300 transition-colors shrink-0"
            @click="dismissed = true; emit('dismiss')"
          >
            <X class="h-4 w-4" />
          </button>
        </div>
        <p class="text-xs text-slate-400 font-mono break-all leading-relaxed">{{ error }}</p>
        <ul class="space-y-1 pt-1">
          <li
            v-for="(suggestion, i) in diagnostic.suggestions"
            :key="i"
            class="text-xs text-slate-400 flex items-start gap-1.5"
          >
            <span class="text-slate-600 mt-px">-</span>
            <span>{{ suggestion }}</span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>
