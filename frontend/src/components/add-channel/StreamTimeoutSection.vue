<template>
  <div class="stream-timeout-section">
    <div class="section-header">
      <h3 class="text-h6">{{ t('addChannel.streamTimeoutTitle') }}</h3>
      <p class="text-caption text-medium-emphasis">
        {{ t('addChannel.streamTimeoutDescription') }}
      </p>
    </div>

    <!-- 预设按钮组 -->
    <div class="preset-buttons mb-6">
      <v-btn-toggle
        :model-value="currentPreset"
        mandatory
        variant="outlined"
        color="primary"
        @update:model-value="applyPreset"
      >
        <v-btn value="default">
          {{ t('addChannel.timeoutPresetDefault') }}
        </v-btn>
        <v-btn value="relaxed">
          {{ t('addChannel.timeoutPresetRelaxed') }}
        </v-btn>
        <v-btn value="strict">
          {{ t('addChannel.timeoutPresetStrict') }}
        </v-btn>
        <v-btn value="custom">
          {{ t('addChannel.timeoutPresetCustom') }}
        </v-btn>
      </v-btn-toggle>
    </div>

    <!-- 首次响应超时 -->
    <div class="timeout-control mb-6">
      <label class="control-label">
        {{ t('addChannel.streamFirstByteTimeoutLabel') }}
        <span class="value-display">{{ firstByteTimeout }}s</span>
      </label>
      <div class="slider-container neo-brutalism">
        <input
          type="range"
          :value="firstByteTimeout"
          min="5"
          max="120"
          step="5"
          class="timeout-slider"
          @input="updateTimeout('firstByte', $event)"
        />
      </div>
      <div class="hint-text">
        <span class="text-caption text-medium-emphasis">
          {{ t('addChannel.streamFirstByteTimeoutHint') }}
        </span>
      </div>
    </div>

    <!-- 持续响应超时 -->
    <div class="timeout-control mb-6">
      <label class="control-label">
        {{ t('addChannel.streamChunkIntervalTimeoutLabel') }}
        <span class="value-display">{{ chunkIntervalTimeout }}s</span>
      </label>
      <div class="slider-container neo-brutalism">
        <input
          type="range"
          :value="chunkIntervalTimeout"
          min="10"
          max="180"
          step="10"
          class="timeout-slider"
          @input="updateTimeout('chunkInterval', $event)"
        />
      </div>
      <div class="hint-text">
        <span class="text-caption text-medium-emphasis">
          {{ t('addChannel.streamChunkIntervalTimeoutHint') }}
        </span>
      </div>
    </div>

    <!-- 整体超时 -->
    <div class="timeout-control">
      <label class="control-label">
        {{ t('addChannel.streamOverallTimeoutLabel') }}
        <span class="value-display">{{ overallTimeout }}s</span>
      </label>
      <div class="slider-container neo-brutalism">
        <input
          type="range"
          :value="overallTimeout"
          min="60"
          max="600"
          step="30"
          class="timeout-slider"
          @input="updateTimeout('overall', $event)"
        />
      </div>
      <div class="hint-text">
        <span class="text-caption text-medium-emphasis">
          {{ t('addChannel.streamOverallTimeoutHint') }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from '../../i18n'

interface Props {
  firstByteTimeout: number
  chunkIntervalTimeout: number
  overallTimeout: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:firstByteTimeout': [number]
  'update:chunkIntervalTimeout': [number]
  'update:overallTimeout': [number]
}>()

const { t } = useI18n()

// 预设配置
const PRESETS = {
  default: { firstByte: 30, chunkInterval: 60, overall: 300 },
  relaxed: { firstByte: 60, chunkInterval: 120, overall: 600 },
  strict: { firstByte: 15, chunkInterval: 30, overall: 180 },
}

// 检测当前预设
const currentPreset = computed(() => {
  for (const [key, preset] of Object.entries(PRESETS)) {
    if (
      props.firstByteTimeout === preset.firstByte &&
      props.chunkIntervalTimeout === preset.chunkInterval &&
      props.overallTimeout === preset.overall
    ) {
      return key
    }
  }
  return 'custom'
})

const applyPreset = (preset: string) => {
  if (preset === 'custom') return

  const config = PRESETS[preset as keyof typeof PRESETS]
  if (config) {
    emit('update:firstByteTimeout', config.firstByte)
    emit('update:chunkIntervalTimeout', config.chunkInterval)
    emit('update:overallTimeout', config.overall)
  }
}

const updateTimeout = (type: 'firstByte' | 'chunkInterval' | 'overall', event: Event) => {
  const value = parseInt((event.target as HTMLInputElement).value, 10)

  if (type === 'firstByte') {
    emit('update:firstByteTimeout', value)
  } else if (type === 'chunkInterval') {
    emit('update:chunkIntervalTimeout', value)
  } else {
    emit('update:overallTimeout', value)
  }
}
</script>

<style scoped>
.section-header {
  margin-bottom: 24px;
}

.preset-buttons {
  display: flex;
  justify-content: center;
}

.timeout-control {
  margin-bottom: 24px;
}

.control-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-size: 0.875rem;
  font-weight: 500;
}

.value-display {
  font-weight: 700;
  color: rgb(var(--v-theme-primary));
}

/* Neo-Brutalism 滑块样式 */
.slider-container.neo-brutalism {
  padding: 16px;
  background: rgba(var(--v-theme-surface), 1);
  border: 3px solid rgb(var(--v-theme-on-surface));
  border-radius: 0;
  box-shadow: 6px 6px 0 rgb(var(--v-theme-on-surface));
}

.timeout-slider {
  width: 100%;
  height: 8px;
  -webkit-appearance: none;
  appearance: none;
  background: rgba(var(--v-theme-primary), 0.2);
  outline: none;
  border: 2px solid rgb(var(--v-theme-on-surface));
}

.timeout-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 24px;
  height: 24px;
  background: rgb(var(--v-theme-primary));
  border: 3px solid rgb(var(--v-theme-on-surface));
  cursor: pointer;
  box-shadow: 3px 3px 0 rgb(var(--v-theme-on-surface));
}

.timeout-slider::-moz-range-thumb {
  width: 24px;
  height: 24px;
  background: rgb(var(--v-theme-primary));
  border: 3px solid rgb(var(--v-theme-on-surface));
  cursor: pointer;
  box-shadow: 3px 3px 0 rgb(var(--v-theme-on-surface));
  border-radius: 0;
}

.hint-text {
  margin-top: 8px;
  padding-left: 4px;
}
</style>
