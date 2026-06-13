<template>
  <div class="custom-headers-section">
    <v-card variant="outlined">
      <v-card-title class="section-card-title d-flex align-center ga-2">
        <v-icon size="small">mdi-web</v-icon>
        {{ t('addChannel.customHeadersLabel') }}
      </v-card-title>
      <v-card-text>
        <div class="text-caption text-medium-emphasis mb-3">
          {{ t('addChannel.customHeadersHint') }}
        </div>

        <!-- 已添加的请求头列表 -->
        <v-list v-if="headers.length > 0" density="compact" class="mb-3">
          <v-list-item
            v-for="(header, index) in headers"
            :key="`${header.key}-${index}`"
            class="px-2"
          >
            <template #prepend>
              <v-icon size="small" color="primary">mdi-tag</v-icon>
            </template>
            <v-list-item-title class="text-body-2">
              <code>{{ header.key }}</code>: <span class="text-medium-emphasis">{{ header.value }}</span>
            </v-list-item-title>
            <template #append>
              <v-btn
                icon="mdi-delete"
                size="x-small"
                variant="text"
                color="error"
                @click="removeHeader(index)"
              />
            </template>
          </v-list-item>
        </v-list>

        <!-- 添加新请求头 -->
        <div class="d-flex ga-2 align-center">
          <v-text-field
            v-model="newHeaderKey"
            :label="t('addChannel.headerNameLabel')"
            placeholder="X-Custom-Header"
            variant="outlined"
            density="compact"
            hide-details
            style="flex: 1"
          />
          <v-text-field
            v-model="newHeaderValue"
            :label="t('addChannel.headerValueLabel')"
            placeholder="value"
            variant="outlined"
            density="compact"
            hide-details
            style="flex: 2"
          />
          <v-btn
            icon="mdi-plus"
            size="small"
            color="primary"
            variant="tonal"
            :disabled="!newHeaderKey.trim() || !newHeaderValue.trim()"
            @click="addHeader"
          />
        </div>
      </v-card-text>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from '../../i18n'

interface CustomHeader {
  key: string
  value: string
}

interface Props {
  headers: CustomHeader[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:headers': [CustomHeader[]]
}>()

const { t } = useI18n()
const newHeaderKey = ref('')
const newHeaderValue = ref('')

const addHeader = () => {
  const key = newHeaderKey.value.trim()
  const value = newHeaderValue.value.trim()
  if (!key || !value) return
  emit('update:headers', [...props.headers, { key, value }])
  newHeaderKey.value = ''
  newHeaderValue.value = ''
}

const removeHeader = (index: number) => {
  const updated = props.headers.filter((_, i) => i !== index)
  emit('update:headers', updated)
}
</script>

<style scoped>
.section-card-title {
  font-size: 0.875rem !important;
  line-height: 1.4;
  font-weight: 600;
}
</style>
