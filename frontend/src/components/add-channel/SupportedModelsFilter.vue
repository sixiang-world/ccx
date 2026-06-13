<template>
  <div class="supported-models-filter">
    <v-combobox
      :model-value="modelValue"
      :label="t('addChannel.supportedModelsLabel')"
      :placeholder="t('addChannel.supportedModelsPlaceholder')"
      prepend-inner-icon="mdi-brain"
      :hint="t('addChannel.supportedModelsHint')"
      :error-messages="error ? [error] : []"
      persistent-hint
      clearable
      multiple
      chips
      closable-chips
      variant="outlined"
      density="comfortable"
      eager
      @update:model-value="$emit('update:modelValue', $event)"
      @update:menu="$emit('menu-update', $event)"
    />
    <div class="d-flex align-center flex-wrap ga-2 mt-2">
      <div class="text-caption text-primary">{{ t('addChannel.commonFilters') }}</div>
      <v-chip
        v-for="filter in commonFilters"
        :key="filter"
        size="small"
        :color="selectedFilters.includes(filter) ? 'primary' : 'default'"
        :variant="selectedFilters.includes(filter) ? 'flat' : 'tonal'"
        @click="$emit('append-filter', filter)"
      >
        {{ filter }}
      </v-chip>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '../../i18n'

interface Props {
  modelValue: string[]
  error?: string
  commonFilters: string[]
  selectedFilters: string[]
}

defineProps<Props>()

defineEmits<{
  'update:modelValue': [string[]]
  'append-filter': [string]
  'menu-update': [boolean]
}>()

const { t } = useI18n()
</script>
