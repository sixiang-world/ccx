<script setup lang="ts">
import { onMounted } from 'vue'
import AgentCard from '@/components/agent/AgentCard.vue'
import ConfigDiffDialog from '@/components/agent/ConfigDiffDialog.vue'
import { useStatus } from '@/composables/useStatus'
import { useAgentConfig } from '@/composables/useAgentConfig'
import type { AgentPlatform } from '@/types'

const { status, actionError } = useStatus()
const {
  agentStatuses,
  configLoading,
  selectedClaudeProvider,
  claudeProviderKeys,
  savedProviderKeys,
  codexOpenAIKey,
  claudeMimoBaseUrl,
  selectedMimoPlan,
  agentLabels,
  agentPlatforms,
  claudeProviderLabel,
  claudeTargetBaseUrl,
  agentStatusText,
  agentStatusClass,
  loadAgentStatuses,
  canApplyAgent,
  selectedCodexProvider,
  codexProviderLabels,
  codexProviderLabel,
  codexTargetBaseUrl,
  // Diff preview
  diffDialogOpen,
  diffResult,
  diffMode,
  diffLoading,
  diffPendingPlatform,
  showApplyPreview,
  showRestorePreview,
  confirmApply,
  confirmRestore,
  closeDiffDialog,
} = useAgentConfig()

onMounted(() => {
  loadAgentStatuses()
})

const handleApply = async (platform: AgentPlatform) => {
  actionError.value = ''
  try {
    await showApplyPreview(platform)
  } catch (error) {
    actionError.value = error instanceof Error ? error.message : String(error)
  }
}

const handleRestore = async (platform: AgentPlatform) => {
  actionError.value = ''
  try {
    await showRestorePreview(platform)
  } catch (error) {
    actionError.value = error instanceof Error ? error.message : String(error)
  }
}

const handleConfirm = async () => {
  actionError.value = ''
  try {
    if (diffMode.value === 'apply') {
      await confirmApply()
    } else {
      await confirmRestore()
    }
    await loadAgentStatuses()
  } catch (error) {
    actionError.value = error instanceof Error ? error.message : String(error)
  }
}
</script>

<template>
  <div class="space-y-4">
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <AgentCard
        v-for="platform in agentPlatforms"
        :key="platform"
        :platform="platform"
        :agent-status="agentStatuses[platform]"
        :config-loading="configLoading"
        :service-running="status.running"
        :agent-label="agentLabels[platform]"
        :agent-status-text="agentStatusText(agentStatuses[platform])"
        :agent-status-class="agentStatusClass(agentStatuses[platform])"
        :can-apply="canApplyAgent(platform, status.running)"
        :selected-claude-provider="selectedClaudeProvider"
        :claude-provider-keys="claudeProviderKeys"
        :saved-provider-keys="savedProviderKeys"
        :claude-mimo-base-url="claudeMimoBaseUrl"
        :selected-mimo-plan="selectedMimoPlan"
        :claude-provider-label="claudeProviderLabel"
        :claude-target-base-url="claudeTargetBaseUrl"
        :selected-codex-provider="selectedCodexProvider"
        :codex-open-a-i-key="codexOpenAIKey"
        :codex-provider-labels="codexProviderLabels"
        :codex-provider-label="codexProviderLabel"
        :codex-target-base-url="codexTargetBaseUrl"
        @apply="handleApply(platform)"
        @restore="handleRestore(platform)"
        @update:selected-claude-provider="selectedClaudeProvider = $event"
        @update:claude-provider-keys="claudeProviderKeys = $event"
        @update:mimo-base-url="claudeMimoBaseUrl = $event"
        @update:selected-mimo-plan="selectedMimoPlan = $event"
        @update:selected-codex-provider="selectedCodexProvider = $event"
        @update:codex-open-a-i-key="codexOpenAIKey = $event"
      />
    </div>
    <p v-if="actionError" class="text-sm text-destructive-foreground">{{ actionError }}</p>

    <ConfigDiffDialog
      :open="diffDialogOpen"
      :mode="diffMode"
      :platform="diffPendingPlatform"
      :result="diffResult"
      :loading="diffLoading"
      @confirm="handleConfirm"
      @cancel="closeDiffDialog"
    />
  </div>
</template>
