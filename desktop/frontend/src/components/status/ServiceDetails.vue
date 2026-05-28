<script setup lang="ts">
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { FolderOpen } from 'lucide-vue-next'
import { OpenDirectory } from '@bindings/github.com/BenedictKing/ccx/desktop/desktopservice'
import type { DesktopStatus } from '@/types'
import { useLanguage } from '@/composables/useLanguage'

defineProps<{
  status: DesktopStatus
}>()

const { t } = useLanguage()

const openDir = (path: string) => {
  OpenDirectory(path).catch(() => {})
}
</script>

<template>
  <Card>
    <CardHeader class="pb-3">
      <CardTitle class="text-sm font-medium text-muted-foreground">{{ t('details.title') }}</CardTitle>
    </CardHeader>
    <CardContent class="space-y-3">
      <div v-for="item in [
        { label: t('details.binary'), value: status.binaryPath || t('details.binaryMissing'), action: status.binaryPath ? 'reveal' : null, actionPath: status.binaryPath },
        { label: t('details.dataDir'), value: status.dataDir || t('details.dataDirMissing'), action: status.dataDir ? 'open' : null, actionPath: status.dataDir },
        { label: 'PID', value: String(status.pid || '-'), action: null },
        { label: t('details.healthStatus'), value: status.health?.status || 'unknown', action: null },
      ]" :key="item.label" class="grid grid-cols-[5rem_minmax(0,1fr)] items-center gap-3 text-sm">
        <span class="text-muted-foreground">{{ item.label }}</span>
        <div class="flex min-w-0 items-center justify-end gap-2">
          <code
            class="inline-block min-w-0 max-w-full rounded-md bg-secondary px-2 py-1 text-right text-xs"
            :class="item.action ? 'break-all' : 'whitespace-nowrap'"
          >{{ item.value }}</code>
          <Button
            v-if="item.action"
            variant="ghost"
            size="icon-sm"
            :title="item.action === 'reveal' ? t('details.revealDir') : t('details.openDir')"
            class="shrink-0"
            @click="openDir(item.actionPath!)"
          >
            <FolderOpen class="w-3.5 h-3.5" />
          </Button>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
