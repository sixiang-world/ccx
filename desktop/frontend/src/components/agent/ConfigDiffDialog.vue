<script setup lang="ts">
import { computed } from 'vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import type { ConfigDiffResult, FileDiff } from '@/types'

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
  props.mode === 'apply' ? '应用配置预览' : '恢复配置预览'
)

const confirmLabel = computed(() =>
  props.mode === 'apply' ? '确认应用' : '确认恢复'
)

const platformLabel = computed(() =>
  props.platform === 'claude' ? 'Claude Code' : 'Codex'
)

const actionLabel = (action: string) => {
  switch (action) {
    case 'create': return '创建'
    case 'delete': return '删除'
    default: return '修改'
  }
}

const actionBadgeClass = (action: string) => {
  switch (action) {
    case 'create': return 'bg-emerald-500/20 text-emerald-400 border-0'
    case 'delete': return 'bg-red-500/20 text-red-400 border-0'
    default: return 'bg-blue-500/20 text-blue-400 border-0'
  }
}
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
              <div class="text-sm text-white/40">计算变更中...</div>
            </div>

            <!-- No changes -->
            <div v-else-if="!result || result.files.length === 0" class="flex items-center justify-center py-12">
              <div class="text-sm text-white/40">无变更</div>
            </div>

            <!-- Diff blocks -->
            <template v-else>
              <div
                v-for="file in result.files"
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
                      <tr
                        v-for="(line, idx) in file.lines"
                        :key="idx"
                        :class="{
                          'bg-emerald-500/[0.07]': line.type === 'added',
                          'bg-red-500/[0.07]': line.type === 'removed',
                        }"
                      >
                        <td class="w-8 text-right pr-2 py-0.5 select-none text-white/20 align-top">
                          {{ idx + 1 }}
                        </td>
                        <td class="w-4 text-center py-0.5 select-none align-top"
                          :class="{
                            'text-emerald-400': line.type === 'added',
                            'text-red-400': line.type === 'removed',
                            'text-white/20': line.type === 'context',
                          }"
                        >
                          {{ line.type === 'added' ? '+' : line.type === 'removed' ? '-' : ' ' }}
                        </td>
                        <td class="py-0.5 pr-4 whitespace-pre-wrap break-all"
                          :class="{
                            'text-emerald-300/80': line.type === 'added',
                            'text-red-300/80': line.type === 'removed',
                            'text-white/30': line.type === 'context',
                          }"
                        >
                          {{ line.content || ' ' }}
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </template>
          </div>

          <!-- Footer -->
          <div class="flex justify-end gap-2 border-t border-white/[0.05] px-6 py-4 shrink-0">
            <Button variant="ghost" size="sm" @click="emit('cancel')">
              取消
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
