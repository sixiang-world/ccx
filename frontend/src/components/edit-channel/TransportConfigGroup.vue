<template>
  <v-card variant="outlined" class="pa-4">
    <div class="text-caption font-weight-bold text-uppercase text-medium-emphasis mb-3">
      <v-icon size="small" color="primary" class="mr-1">mdi-network</v-icon>
      {{ t('addChannel.transportTitle') }}
    </div>

    <v-row dense>
      <!-- 代理 URL -->
      <v-col cols="12">
        <v-text-field
          :model-value="form.proxyUrl"
          :label="t('addChannel.proxyUrlLabel')"
          :placeholder="t('addChannel.proxyUrlPlaceholder')"
          prepend-inner-icon="mdi-shield-lock-outline"
          :hint="t('addChannel.proxyUrlHint')"
          persistent-hint
          clearable
          variant="outlined"
          density="comfortable"
          @update:model-value="updateField('proxyUrl', $event)"
        />
      </v-col>

      <!-- 请求超时 -->
      <v-col cols="12">
        <v-text-field
          :model-value="form.requestTimeoutMs"
          :label="t('addChannel.requestTimeoutMsLabel')"
          :placeholder="t('addChannel.requestTimeoutMsPlaceholder')"
          prepend-inner-icon="mdi-timer-sand"
          :hint="t('addChannel.requestTimeoutMsHint')"
          :rules="[rules.requestTimeoutMs]"
          persistent-hint
          clearable
          variant="outlined"
          density="comfortable"
          type="number"
          min="1"
          step="1000"
          @update:model-value="updateField('requestTimeoutMs', $event)"
        />
      </v-col>

      <slot name="stream-timeout" />

      <!-- 路由前缀 -->
      <v-col cols="12">
        <v-text-field
          :model-value="form.routePrefix"
          :label="t('addChannel.routePrefixLabel')"
          :placeholder="t('addChannel.routePrefixPlaceholder')"
          prepend-inner-icon="mdi-routes"
          :hint="t('addChannel.routePrefixHint')"
          persistent-hint
          clearable
          variant="outlined"
          density="comfortable"
          @update:model-value="updateField('routePrefix', $event)"
        />
      </v-col>
    </v-row>
  </v-card>
</template>

<script setup lang="ts">
import { useI18n } from '../../i18n'

interface FormData {
  proxyUrl: string
  requestTimeoutMs: string | number | null
  routePrefix?: string
}

interface Props {
  form: FormData
  rules: Record<string, (value: any) => boolean | string>
}

defineProps<Props>()

const emit = defineEmits<{
  'update:field': [field: keyof FormData, value: unknown]
}>()

const { t } = useI18n()

const updateField = (field: keyof FormData, value: unknown) => {
  emit('update:field', field, value)
}
</script>
