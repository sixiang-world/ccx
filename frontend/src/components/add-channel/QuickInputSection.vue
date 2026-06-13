<template>
  <div class="quick-input-section">
    <v-alert
      type="info"
      variant="tonal"
      class="mb-4"
    >
      <div class="text-body-2">
        {{ t('addChannel.quickInputHint') }}
      </div>
    </v-alert>

    <v-textarea
      v-model="quickInput"
      :label="t('addChannel.quickInputLabel')"
      :placeholder="t('addChannel.quickInputPlaceholder')"
      prepend-inner-icon="mdi-lightning-bolt"
      variant="outlined"
      rows="8"
      no-resize
      auto-grow
      hide-details="auto"
      @paste="handlePaste"
    />

    <!-- 检测结果 -->
    <div v-if="detectedType" class="detection-result mt-4">
      <v-chip
        color="success"
        prepend-icon="mdi-check-circle"
        class="mb-2"
      >
        {{ t('addChannel.detectedType', { type: detectedType }) }}
      </v-chip>

      <div v-if="detectedName" class="text-caption text-medium-emphasis">
        {{ t('addChannel.detectedName') }}: {{ detectedName }}
      </div>

      <div v-if="detectedBaseUrl" class="text-caption text-medium-emphasis">
        {{ t('addChannel.detectedBaseUrl') }}: {{ detectedBaseUrl }}
      </div>

      <div v-if="detectedApiKey" class="text-caption text-medium-emphasis">
        {{ t('addChannel.detectedApiKey') }}: {{ maskApiKey(detectedApiKey) }}
      </div>
    </div>

    <v-btn
      color="primary"
      size="large"
      block
      class="mt-4"
      prepend-icon="mdi-check"
      :disabled="!quickInput.trim()"
      @click="handleQuickSubmit"
    >
      {{ t('addChannel.quickSubmitButton') }}
    </v-btn>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from '../../i18n'

interface QuickSubmitData {
  name?: string
  baseUrl?: string
  apiKey?: string
  serviceType?: string
}

const emit = defineEmits<{
  'quick-submit': [QuickSubmitData]
}>()

const { t } = useI18n()

const quickInput = ref('')
const detectedType = ref('')
const detectedName = ref('')
const detectedBaseUrl = ref('')
const detectedApiKey = ref('')

const handlePaste = (event: ClipboardEvent) => {
  const text = event.clipboardData?.getData('text') || ''
  parseQuickInput(text)
}

const parseQuickInput = (text: string) => {
  // 简单的启发式检测
  detectedType.value = ''
  detectedName.value = ''
  detectedBaseUrl.value = ''
  detectedApiKey.value = ''

  // 检测 URL
  const urlMatch = text.match(/https?:\/\/[^\s]+/i)
  if (urlMatch) {
    detectedBaseUrl.value = urlMatch[0]
  }

  // 检测 API Key (sk- 开头或其他常见格式)
  const keyMatch = text.match(/(?:sk-|api[_-]?key[:\s]+)([a-zA-Z0-9_-]+)/i)
  if (keyMatch) {
    detectedApiKey.value = keyMatch[1]
  }

  // 检测服务类型
  if (text.includes('openai') || text.includes('gpt')) {
    detectedType.value = 'OpenAI'
  } else if (text.includes('claude') || text.includes('anthropic')) {
    detectedType.value = 'Claude'
  } else if (text.includes('gemini') || text.includes('google')) {
    detectedType.value = 'Gemini'
  }
}

const maskApiKey = (key: string): string => {
  if (key.length <= 8) return '***'
  return key.slice(0, 4) + '...' + key.slice(-4)
}

const handleQuickSubmit = () => {
  parseQuickInput(quickInput.value)

  emit('quick-submit', {
    name: detectedName.value,
    baseUrl: detectedBaseUrl.value,
    apiKey: detectedApiKey.value,
    serviceType: detectedType.value,
  })
}
</script>

<style scoped>
.detection-result {
  padding: 12px;
  background: rgba(var(--v-theme-surface-variant), 0.3);
  border-radius: 8px;
}
</style>
