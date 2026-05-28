<script setup lang="ts">
import { computed, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import type { ConfigDiffResult, DiffLine, FileDiff } from '@/types'
import { useLanguage } from '@/composables/useLanguage'

const CONTEXT_THRESHOLD = 4
const CONTEXT_KEEP = 2

const { t } = useLanguage()

type DisplayLine =
  | { kind: 'line'; line: DiffLine; origIndex: number }
  | { kind: 'collapsed'; id: string; lines: DiffLine[]; startOrigIndex: number }

const props = defineProps<{
  open: boolean
  mode: 'apply' | 'restore'
  platform: string
  result: ConfigDiffResult | null
  loading: boolean
}>()

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()

const title = computed(() =>
  props.mode === 'apply' ? t('agent.diffPreviewApply') : t('agent.diffPreviewRestore')
)

const confirmLabel = computed(() =>
  props.mode === 'apply' ? t('agent.diffConfirmApply') : t('agent.diffConfirmRestore')
)

const platformLabel = computed(() =>
  props.platform === 'claude' ? 'Claude Code' : 'Codex'
)

const actionLabel = (action: string) => {
  switch (action) {
    case 'create': return t('agent.diffActionCreate')
    case 'delete': return t('agent.diffActionDelete')
    default: return t('agent.diffActionModify')
  }
}

const actionBadgeClass = (action: string) => {
  switch (action) {
    case 'create': return 'bg-emerald-500/20 text-emerald-400 border-0'
    case 'delete': return 'bg-red-500/20 text-red-400 border-0'
    default: return 'bg-blue-500/20 text-blue-400 border-0'
  }
}

// --- Context folding ---

const collapsedSections = ref(new Set<string>())

function toggleCollapse(id: string) {
  if (collapsedSections.value.has(id)) {
    collapsedSections.value.delete(id)
  } else {
    collapsedSections.value.add(id)
  }
}

function isExpanded(id: string): boolean {
  return collapsedSections.value.has(id)
}

function collapseContextLines(file: FileDiff, fileIndex: number): DisplayLine[] {
  const result: DisplayLine[] = []
  let runStart = -1

  const flushRun = (end: number) => {
    const run = file.lines.slice(runStart, end)
    if (run.length <= CONTEXT_THRESHOLD) {
      run.forEach((line, i) => result.push({ kind: 'line', line, origIndex: runStart + i }))
    } else {
      const hidden = run.slice(CONTEXT_KEEP, run.length - CONTEXT_KEEP)
      const id = `${fileIndex}-c-${runStart}`
      // head
      for (let i = 0; i < CONTEXT_KEEP; i++) {
        result.push({ kind: 'line', line: run[i], origIndex: runStart + i })
      }
      // collapsed marker
      result.push({ kind: 'collapsed', id, lines: hidden, startOrigIndex: runStart + CONTEXT_KEEP })
      // tail
      for (let i = run.length - CONTEXT_KEEP; i < run.length; i++) {
        result.push({ kind: 'line', line: run[i], origIndex: runStart + i })
      }
    }
  }

  for (let i = 0; i < file.lines.length; i++) {
    if (file.lines[i].type === 'context') {
      if (runStart === -1) runStart = i
    } else {
      if (runStart !== -1) {
        flushRun(i)
        runStart = -1
      }
      result.push({ kind: 'line', line: file.lines[i], origIndex: i })
    }
  }
  if (runStart !== -1) flushRun(file.lines.length)

  return result
}

const processedFiles = computed(() => {
  if (!props.result) return []
  return props.result.files.map((file, fi) => ({
    file,
    displayLines: collapseContextLines(file, fi),
  }))
})
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="open"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
        @click.self="emit('cancel')"
      >
        <div
          class="w-[min(720px,90vw)] max-h-[85vh] overflow-hidden rounded-2xl border border-white/[0.06] bg-[#0f1421] shadow-2xl flex flex-col"
        >
          <!-- Header -->
          <div class="border-b border-white/[0.05] px-6 py-4 shrink-0">
            <div class="flex items-baseline justify-between">
              <h2 class="text-lg font-semibold text-white">{{ title }}</h2>
              <div class="text-xs text-white/40">{{ platformLabel }}</div>
            </div>
          </div>

          <!-- Content -->
          <div class="flex-1 overflow-y-auto px-6 py-4 space-y-4 min-h-0">
            <!-- Loading state -->
            <div v-if="loading" class="flex items-center justify-center py-12">
              <div class="text-sm text-white/40">{{ t('agent.diffComputing') }}</div>
            </div>

            <!-- No changes -->
            <div v-else-if="!result || result.files.length === 0" class="flex items-center justify-center py-12">
              <div class="text-sm text-white/40">{{ t('agent.diffNoChanges') }}</div>
            </div>

            <!-- Diff blocks -->
            <template v-else>
              <div
                v-for="{ file, displayLines } in processedFiles"
                :key="file.path"
                class="rounded-lg border border-white/[0.06] overflow-hidden"
              >
                <!-- File header -->
                <div class="flex items-center justify-between px-4 py-2 bg-white/[0.03] border-b border-white/[0.05]">
                  <code class="text-xs text-white/70 break-all">{{ file.path }}</code>
                  <Badge :class="actionBadgeClass(file.action)" class="text-[10px] px-1.5 py-0">
                    {{ actionLabel(file.action) }}
                  </Badge>
                </div>

                <!-- Diff lines -->
                <div class="overflow-x-auto">
                  <table class="w-full text-xs font-mono">
                    <tbody>
                      <template v-for="item in displayLines" :key="item.kind === 'collapsed' ? item.id : item.origIndex">
                        <!-- Normal diff line -->
                        <tr
                          v-if="item.kind === 'line'"
                          :class="{
                            'bg-emerald-500/[0.07]': item.line.type === 'added',
                            'bg-red-500/[0.07]': item.line.type === 'removed',
                          }"
                        >
                          <td class="w-8 text-right pr-2 py-0.5 select-none text-white/20 align-top">
                            {{ item.origIndex + 1 }}
                          </td>
                          <td class="w-4 text-center py-0.5 select-none align-top"
                            :class="{
                              'text-emerald-400': item.line.type === 'added',
                              'text-red-400': item.line.type === 'removed',
                              'text-white/20': item.line.type === 'context',
                            }"
                          >
                            {{ item.line.type === 'added' ? '+' : item.line.type === 'removed' ? '-' : ' ' }}
                          </td>
                          <td class="py-0.5 pr-4 whitespace-pre-wrap break-all"
                            :class="{
                              'text-emerald-300/80': item.line.type === 'added',
                              'text-red-300/80': item.line.type === 'removed',
                              'text-white/30': item.line.type === 'context',
                            }"
                          >
                            {{ item.line.content || ' ' }}
                          </td>
                        </tr>

                        <!-- Collapsed marker (not yet expanded) -->
                        <tr v-else-if="!isExpanded(item.id)" class="group">
                          <td colspan="3" class="py-1 px-4">
                            <button
                              class="w-full flex items-center justify-center gap-1.5 py-1 rounded text-white/50 hover:text-white/80 hover:bg-white/[0.04] transition-colors cursor-pointer"
                              @click="toggleCollapse(item.id)"
                            >
                              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M4 8h16M4 16h16" />
                              </svg>
                              <span>{{ t('agent.diffExpandContext', { count: String(item.lines.length) }) }}</span>
                            </button>
                          </td>
                        </tr>

                        <!-- Expanded hidden lines -->
                        <template v-else>
                          <tr class="group">
                            <td colspan="3" class="py-1 px-4">
                              <button
                                class="w-full flex items-center justify-center gap-1.5 py-1 rounded text-white/50 hover:text-white/80 hover:bg-white/[0.04] transition-colors cursor-pointer"
                                @click="toggleCollapse(item.id)"
                              >
                                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 15l7-7 7 7" />
                                </svg>
                                <span>{{ t('agent.diffCollapseContext', { count: String(item.lines.length) }) }}</span>
                              </button>
                            </td>
                          </tr>
                          <tr
                            v-for="(hidden, hi) in item.lines"
                            :key="`${item.id}-${hi}`"
                          >
                            <td class="w-8 text-right pr-2 py-0.5 select-none text-white/20 align-top">
                              {{ item.startOrigIndex + hi + 1 }}
                            </td>
                            <td class="w-4 text-center py-0.5 select-none text-white/20 align-top"> </td>
                            <td class="py-0.5 pr-4 whitespace-pre-wrap break-all text-white/30">
                              {{ hidden.content || ' ' }}
                            </td>
                          </tr>
                        </template>
                      </template>
                    </tbody>
                  </table>
                </div>
              </div>
            </template>
          </div>

          <!-- Footer -->
          <div class="flex justify-end gap-2 border-t border-white/[0.05] px-6 py-4 shrink-0">
            <Button variant="ghost" size="sm" @click="emit('cancel')">
              {{ t('agent.diffCancel') }}
            </Button>
            <Button size="sm" :disabled="loading" @click="emit('confirm')">
              {{ confirmLabel }}
            </Button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.18s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
