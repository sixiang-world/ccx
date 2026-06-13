<template>
  <div class="supported-models-filter">
    <v-combobox
      :model-value="modelValue"
      :label="t('addChannel.supportedModelsLabel')"
      :placeholder="t('addChannel.supportedModelsPlaceholder')"
      prepend-inner-icon="mdi-filter"
      variant="outlined"
      density="comfortable"
      multiple
      chips
      closable-chips
      hide-details="auto"
      @update:model-value="$emit('update:modelValue', $event)"
    >
      <template #chip="{ item, props: chipProps }">
        <v-chip
          v-bind="chipProps"
          closable
          size="small"
        >
          {{ item.title }}
        </v-chip>
      </template>
    </v-combobox>
    <div class="hint-text">
      <span class="text-caption text-medium-emphasis">
        {{ t('addChannel.supportedModelsHint') }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '../../i18n'

interface Props {
  modelValue: string[]
}

defineProps<Props>()

defineEmits<{
  'update:modelValue': [string[]]
}>()

const { t } = useI18n()
</script>

<style scoped>
.hint-text {
  margin-top: 8px;
  padding-left: 12px;
}
</style>
