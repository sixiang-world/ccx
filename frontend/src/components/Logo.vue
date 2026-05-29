<script setup lang="ts">
import { computed } from 'vue'
import { useTheme } from 'vuetify'

const props = withDefaults(defineProps<{
  size?: number | string
  animated?: boolean
}>(), {
  size: 32,
  animated: true
})

const theme = useTheme()
const isDark = computed(() => theme.global.current.value.dark)

const sizeStyle = computed(() => {
  const s = typeof props.size === 'number' ? `${props.size}px` : props.size
  return { width: s, height: s }
})

// 根据主题选择对应的 gradient ID 前缀
const gradientPrefix = computed(() => isDark.value ? 'web' : 'light')
</script>

<template>
  <div class="ccx-logo-container" :style="sizeStyle">
    <svg
      viewBox="0 0 100 100"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      class="ccx-logo-svg"
      aria-hidden="true"
    >
      <defs>
        <!-- 与 App Icon 同源的终端网关流光渐变 -->
        <linearGradient :id="`${gradientPrefix}-gateway-flow`" x1="18%" y1="28%" x2="82%" y2="72%">
          <template v-if="isDark">
            <stop offset="0%" stop-color="#38bdf8" />
            <stop offset="48%" stop-color="#6366f1" />
            <stop offset="100%" stop-color="#10b981" />
          </template>
          <template v-else>
            <stop offset="0%" stop-color="#0284c7" />
            <stop offset="48%" stop-color="#4f46e5" />
            <stop offset="100%" stop-color="#059669" />
          </template>
        </linearGradient>

        <!-- 玻璃面板渐变 -->
        <linearGradient :id="`${gradientPrefix}-gateway-panel`" x1="15%" y1="20%" x2="85%" y2="82%">
          <template v-if="isDark">
            <stop offset="0%" stop-color="#102a56" stop-opacity="0.95" />
            <stop offset="52%" stop-color="#06142a" stop-opacity="0.92" />
            <stop offset="100%" stop-color="#042f2e" stop-opacity="0.95" />
          </template>
          <template v-else>
            <stop offset="0%" stop-color="#f0f9ff" stop-opacity="0.98" />
            <stop offset="52%" stop-color="#eff6ff" stop-opacity="0.95" />
            <stop offset="100%" stop-color="#ecfdf5" stop-opacity="0.98" />
          </template>
        </linearGradient>

        <radialGradient :id="`${gradientPrefix}-gateway-bg`" cx="70%" cy="70%" r="86%">
          <template v-if="isDark">
            <stop offset="0%" stop-color="#064e3b" />
            <stop offset="40%" stop-color="#082f49" />
            <stop offset="100%" stop-color="#020617" />
          </template>
          <template v-else>
            <stop offset="0%" stop-color="#d1fae5" />
            <stop offset="40%" stop-color="#dbeafe" />
            <stop offset="100%" stop-color="#f8fafc" />
          </template>
        </radialGradient>

        <filter :id="`${gradientPrefix}-gateway-glow`" x="-28%" y="-28%" width="156%" height="156%">
          <feGaussianBlur :stdDeviation="isDark ? 2.2 : 1.8" result="blur" />
          <feMerge>
            <feMergeNode in="blur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
      </defs>

      <!-- 1. App 图标同源深色圆角底 -->
      <rect x="5" y="5" width="90" height="90" rx="22" :fill="`url(#${gradientPrefix}-gateway-bg)`" />
      <rect
        x="7.5" y="7.5" width="85" height="85" rx="20"
        fill="none"
        :stroke="isDark ? '#93c5fd' : '#3b82f6'"
        stroke-width="0.9"
        :opacity="isDark ? 0.32 : 0.3"
      />

      <!-- 2. 玻璃终端窗口 -->
      <rect
        x="15" y="20" width="70" height="62" rx="14"
        :fill="`url(#${gradientPrefix}-gateway-panel)`"
        :stroke="isDark ? '#93c5fd' : '#3b82f6'"
        stroke-width="1.4"
        opacity="0.98"
      />
      <path d="M 19 32 H 81" :stroke="isDark ? '#bae6fd' : '#94a3b8'" stroke-width="0.8" :opacity="isDark ? 0.18 : 0.25" />
      <circle cx="25" cy="26.5" r="2.3" :fill="isDark ? '#10b981' : '#059669'" />
      <circle cx="32" cy="26.5" r="2.3" :fill="isDark ? '#38bdf8' : '#0284c7'" :opacity="isDark ? 0.78 : 0.85" />
      <circle cx="39" cy="26.5" r="2.3" :fill="isDark ? '#6366f1' : '#4f46e5'" :opacity="isDark ? 0.78 : 0.85" />

      <!-- 3. 终端网关提示符与 X 路由束 -->
      <g :filter="`url(#${gradientPrefix}-gateway-glow)`" stroke-linecap="round" stroke-linejoin="round">
        <path d="M 28 39 L 42 51 L 28 63" :stroke="`url(#${gradientPrefix}-gateway-flow)`" stroke-width="8" />
        <path d="M 52 38 L 73 64" :stroke="`url(#${gradientPrefix}-gateway-flow)`" stroke-width="8" />
        <path d="M 73 38 L 52 64" :stroke="`url(#${gradientPrefix}-gateway-flow)`" stroke-width="8" />
      </g>

      <!-- 4. 底部网关状态线与在线节点 -->
      <path d="M 22 74 H 50" :stroke="isDark ? '#10b981' : '#059669'" stroke-width="2.6" stroke-linecap="round" :opacity="isDark ? 0.46 : 0.55" />
      <path d="M 55 74 H 68" :stroke="isDark ? '#38bdf8' : '#0284c7'" stroke-width="2.6" stroke-linecap="round" :opacity="isDark ? 0.34 : 0.42" />
      <g :class="{ 'animate-gateway-pulse': animated }">
        <circle cx="76" cy="74" r="2.4" :fill="isDark ? '#5eead4' : '#14b8a6'" />
        <circle cx="76" cy="74" r="5.5" :stroke="isDark ? '#5eead4' : '#14b8a6'" stroke-width="1.1" :opacity="isDark ? 0.24 : 0.3" />
      </g>
    </svg>
  </div>
</template>

<style scoped>
.ccx-logo-container {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.ccx-logo-svg {
  width: 100%;
  height: 100%;
}

/* 在线网关节点呼吸脉冲 */
@keyframes gateway-pulse {
  0%, 100% {
    transform: scale(0.92);
    transform-origin: 76px 74px;
    opacity: 0.82;
  }
  50% {
    transform: scale(1.12);
    transform-origin: 76px 74px;
    opacity: 1;
  }
}

.animate-gateway-pulse {
  animation: gateway-pulse 2.4s infinite ease-in-out;
}
</style>
