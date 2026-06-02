<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { Alert } from '@/components/ui/alert'
import { RefreshCw, Search, X, Loader2 } from 'lucide-vue-next'
import { useConversations } from '@/composables/useConversations'
import { useLanguage } from '@/composables/useLanguage'
import { useStatus } from '@/composables/useStatus'
import ConversationCard from './ConversationCard.vue'
import { ref, computed } from 'vue'
import type { ChannelSequenceEntry, ConversationInfo } from '@/services/admin-api'

const { status } = useStatus()
const { tf } = useLanguage()
const {
  conversations,
  total,
  channelsByKind,
  overrides,
  loading,
  error,
  activeKind,
  fetchConversations,
  setOverride,
  removeOverride,
} = useConversations()

const searchQuery = ref('')
const overrideConversation = ref<ConversationInfo | null>(null)
const selectedOverrideIndexes = ref<number[]>([])
const overrideSaving = ref(false)
const overrideError = ref('')

const filteredConversations = computed(() => {
  if (!searchQuery.value.trim()) return conversations.value
  const q = searchQuery.value.toLowerCase()
  return conversations.value.filter(c =>
    c.title?.toLowerCase().includes(q) ||
    c.id.toLowerCase().includes(q) ||
    c.userId.toLowerCase().includes(q) ||
    c.lastModel.toLowerCase().includes(q)
  )
})

const kinds = computed(() => Object.keys(channelsByKind.value))

const availableOverrideChannels = computed<ChannelSequenceEntry[]>(() => {
  if (!overrideConversation.value) return []
  return channelsByKind.value[overrideConversation.value.kind] || []
})

const selectedOverrideSequence = computed(() => {
  const selected = new Set(selectedOverrideIndexes.value)
  return availableOverrideChannels.value.filter(channel => selected.has(channel.channelIndex))
})

function handleRefresh() {
  fetchConversations(activeKind.value || undefined)
}

function openOverrideDialog(conversationId: string) {
  const conversation = conversations.value.find(item => item.id === conversationId)
  if (!conversation) return

  overrideConversation.value = conversation
  overrideError.value = ''
  const existingSequence = overrides.value[conversation.id]?.sequence || []
  selectedOverrideIndexes.value = existingSequence.length
    ? existingSequence.map(item => item.channelIndex)
    : availableOverrideChannels.value
      .filter(item => item.channelName === conversation.channelName || String(item.channelIndex) === conversation.currentChannel)
      .map(item => item.channelIndex)
}

function closeOverrideDialog() {
  if (overrideSaving.value) return
  overrideConversation.value = null
  overrideError.value = ''
  selectedOverrideIndexes.value = []
}

function toggleOverrideChannel(channel: ChannelSequenceEntry) {
  const exists = selectedOverrideIndexes.value.includes(channel.channelIndex)
  selectedOverrideIndexes.value = exists
    ? selectedOverrideIndexes.value.filter(index => index !== channel.channelIndex)
    : [...selectedOverrideIndexes.value, channel.channelIndex]
}

async function saveOverrideSequence() {
  if (!overrideConversation.value) return
  const sequence = selectedOverrideSequence.value
  if (sequence.length === 0) {
    overrideError.value = tf('console.conversations.overrideRequired', '请至少选择一个渠道')
    return
  }

  overrideSaving.value = true
  overrideError.value = ''
  try {
    await setOverride(overrideConversation.value.id, sequence)
    closeOverrideDialog()
  } catch (e) {
    overrideError.value = e instanceof Error ? e.message : String(e)
  } finally {
    overrideSaving.value = false
  }
}

async function handleRemoveOverride(conversationId: string) {
  overrideError.value = ''
  try {
    await removeOverride(conversationId)
  } catch (e) {
    overrideError.value = e instanceof Error ? e.message : String(e)
  }
}

onMounted(() => {
  if (status.value.running) {
    fetchConversations()
  }
})

watch(() => status.value.running, (running) => {
  if (running) fetchConversations()
})
</script>

<template>
  <div class="space-y-4">
    <!-- 错误提示 -->
    <Alert v-if="error" variant="destructive">
      <p class="text-sm">{{ error }}</p>
    </Alert>

    <!-- 工具栏 -->
    <div class="flex items-center gap-3">
      <div class="relative flex-1 max-w-sm">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
        <Input
          v-model="searchQuery"
          :placeholder="tf('console.conversations.search', '搜索会话...')"
          class="pl-9"
        />
      </div>

      <Select
        v-if="kinds.length > 0"
        :model-value="activeKind || '__all__'"
        @update:model-value="(v: any) => { activeKind = (!v || v === '__all__') ? '' : String(v); fetchConversations(activeKind || undefined) }"
      >
        <SelectTrigger class="w-[140px]">
          <SelectValue :placeholder="tf('console.conversations.allKinds', '所有类型')" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="__all__">所有类型</SelectItem>
          <SelectItem v-for="kind in kinds" :key="kind" :value="kind">
            {{ kind }}
          </SelectItem>
        </SelectContent>
      </Select>

      <Button variant="outline" size="sm" :disabled="loading" @click="handleRefresh">
        <RefreshCw class="w-3.5 h-3.5 mr-1.5" :class="{ 'animate-spin': loading }" />
        {{ tf('console.conversations.refresh', '刷新') }}
      </Button>

      <span class="text-xs text-muted-foreground">
        {{ tf('console.conversations.total', '共 {count} 个会话', { count: String(total) }) }}
      </span>
    </div>

    <!-- 会话列表 -->
    <div v-if="loading && conversations.length === 0" class="space-y-3">
      <Skeleton v-for="i in 3" :key="i" class="h-24 w-full rounded-lg" />
    </div>

    <div v-else-if="filteredConversations.length === 0" class="text-center py-12">
      <p class="text-sm text-muted-foreground">
        {{ searchQuery
          ? tf('console.conversations.noSearchResults', '没有匹配的会话')
          : tf('console.conversations.empty', '暂无活跃会话')
        }}
      </p>
    </div>

    <div v-else class="space-y-2">
      <ConversationCard
        v-for="conv in filteredConversations"
        :key="conv.id"
        :conversation="conv"
        :override="overrides[conv.id]"
        :channels-by-kind="channelsByKind[conv.kind]"
        @set-override="openOverrideDialog"
        @remove-override="handleRemoveOverride"
      />
    </div>

    <!-- Override 渠道序列选择对话框 -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="overrideConversation" class="fixed inset-0 z-50 flex items-center justify-center" @keydown.escape="closeOverrideDialog">
          <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="closeOverrideDialog" />

          <div class="relative z-10 flex max-h-[70vh] w-[520px] max-w-[92vw] flex-col border border-border bg-card shadow-2xl">
            <div class="flex shrink-0 items-center justify-between border-b border-border p-4">
              <div>
                <div class="text-xs font-bold uppercase tracking-[0.18em] text-primary">
                  {{ tf('console.conversations.overrideTitle', '会话渠道覆盖') }}
                </div>
                <h3 class="text-sm font-semibold mt-0.5">
                  {{ overrideConversation.title || overrideConversation.id }}
                </h3>
                <p class="text-xs text-muted-foreground">
                  Kind: {{ overrideConversation.kind }} · Current: {{ overrideConversation.currentChannel }}
                </p>
              </div>
              <Button variant="ghost" size="icon-sm" :disabled="overrideSaving" @click="closeOverrideDialog">
                <X class="h-4 w-4" />
              </Button>
            </div>

            <div class="min-h-0 flex-1 overflow-y-auto p-4 space-y-2">
              <div v-if="overrideError" class="border border-destructive/30 bg-destructive/10 p-2 text-sm text-destructive">
                {{ overrideError }}
              </div>

              <p class="text-xs text-muted-foreground mb-2">
                {{ tf('console.conversations.overrideHint', '选择该会话使用的渠道序列（按优先级排列）：') }}
              </p>

              <div v-if="availableOverrideChannels.length === 0" class="text-center py-8 text-sm text-muted-foreground">
                {{ tf('console.conversations.noChannelsForKind', '该类型暂无可用渠道') }}
              </div>

              <label
                v-for="channel in availableOverrideChannels"
                :key="channel.channelIndex"
                class="flex items-center gap-3 border border-border bg-background/50 p-3 cursor-pointer transition-colors hover:bg-accent/30"
                :class="{ 'border-primary/40 bg-primary/5': selectedOverrideIndexes.includes(channel.channelIndex) }"
              >
                <input
                  type="checkbox"
                  :checked="selectedOverrideIndexes.includes(channel.channelIndex)"
                  class="h-4 w-4 rounded border-border text-primary focus:ring-primary"
                  @change="toggleOverrideChannel(channel)"
                />
                <div class="min-w-0 flex-1">
                  <div class="text-sm font-medium truncate">{{ channel.channelName }}</div>
                  <div class="text-xs text-muted-foreground">#{{ channel.channelIndex }}</div>
                </div>
              </label>
            </div>

            <div v-if="selectedOverrideSequence.length" class="shrink-0 border-t border-border bg-muted/30 px-4 py-2">
              <div class="text-[10px] font-bold uppercase tracking-[0.18em] text-muted-foreground mb-1">
                {{ tf('console.conversations.overrideSequence', '覆盖序列') }}
              </div>
              <div class="text-xs font-mono text-foreground">
                {{ selectedOverrideSequence.map(item => item.channelName).join(' → ') }}
              </div>
            </div>

            <div class="flex shrink-0 items-center justify-end gap-2 border-t border-border p-4">
              <Button variant="ghost" :disabled="overrideSaving" @click="closeOverrideDialog">
                {{ tf('common.cancel', '取消') }}
              </Button>
              <Button :disabled="selectedOverrideIndexes.length === 0 || overrideSaving" @click="saveOverrideSequence">
                <Loader2 v-if="overrideSaving" class="mr-2 h-4 w-4 animate-spin" />
                {{ tf('console.conversations.saveOverride', '保存覆盖') }}
              </Button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>
