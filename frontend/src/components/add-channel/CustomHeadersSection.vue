<template>
  <div class="custom-headers-section">
    <div class="section-header">
      <h3 class="text-h6">{{ t('addChannel.customHeadersTitle') }}</h3>
      <p class="text-caption text-medium-emphasis">
        {{ t('addChannel.customHeadersHint') }}
      </p>
    </div>

    <v-list class="headers-list">
      <v-list-item
        v-for="(header, index) in headers"
        :key="index"
        class="header-item"
      >
        <v-row dense>
          <v-col cols="12" md="5">
            <v-text-field
              :model-value="header.key"
              :label="t('addChannel.headerKeyLabel')"
              :placeholder="t('addChannel.headerKeyPlaceholder')"
              variant="outlined"
              density="compact"
              hide-details
              @update:model-value="updateHeader(index, 'key', $event)"
            />
          </v-col>
          <v-col cols="12" md="5">
            <v-text-field
              :model-value="header.value"
              :label="t('addChannel.headerValueLabel')"
              :placeholder="t('addChannel.headerValuePlaceholder')"
              variant="outlined"
              density="compact"
              hide-details
              @update:model-value="updateHeader(index, 'value', $event)"
            />
          </v-col>
          <v-col cols="12" md="2" class="d-flex align-center">
            <v-btn
              icon="mdi-delete"
              size="small"
              variant="text"
              color="error"
              @click="removeHeader(index)"
            />
          </v-col>
        </v-row>
      </v-list-item>
    </v-list>

    <v-btn
      prepend-icon="mdi-plus"
      variant="outlined"
      color="primary"
      @click="addHeader"
    >
      {{ t('addChannel.addHeaderButton') }}
    </v-btn>
  </div>
</template>

<script setup lang="ts">
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

const updateHeader = (index: number, field: 'key' | 'value', value: string) => {
  const updated = [...props.headers]
  updated[index] = { ...updated[index], [field]: value }
  emit('update:headers', updated)
}

const addHeader = () => {
  emit('update:headers', [...props.headers, { key: '', value: '' }])
}

const removeHeader = (index: number) => {
  const updated = props.headers.filter((_, i) => i !== index)
  emit('update:headers', updated)
}
</script>

<style scoped>
.section-header {
  margin-bottom: 16px;
}

.headers-list {
  margin-bottom: 16px;
  background: transparent;
}

.header-item {
  padding: 8px 0;
}
</style>
