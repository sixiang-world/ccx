#!/usr/bin/env bun
/**
 * Batch 1: 将 desktop console.* key 映射到 web 命名空间。
 * 同时修改 JSON 文件和组件中的 t()/tf() 调用。
 * 用法：bun run scripts/align-i18n-batch1.ts
 */

import { readFileSync, writeFileSync } from 'fs'
import { resolve } from 'path'

const desktopDir = resolve(__dirname, '..')
const jsonEn = resolve(desktopDir, 'src/locales/en.json')
const jsonZh = resolve(desktopDir, 'src/locales/zh-CN.json')

// 映射表：old key → new key
const keyMap: Record<string, string> = {
  // console.conversations.* → cockpit.*
  'console.conversations.search': 'cockpit.searchPlaceholder',
  'console.conversations.active': 'cockpit.active',
  'console.conversations.override': 'cockpit.override',
  'console.conversations.empty': 'cockpit.empty',
  'console.conversations.noSearchResults': 'cockpit.noMatches',
  'console.conversations.duration.default': 'cockpit.durationDefault',
  'console.conversations.duration.never': 'cockpit.durationNever',
  'console.conversations.duration.15min': 'cockpit.duration15min',
  'console.conversations.duration.1hour': 'cockpit.duration1hour',
  'console.conversations.duration.2hours': 'cockpit.duration2hours',
  'console.conversations.duration.4hours': 'cockpit.duration4hours',
  'console.conversations.duration.8hours': 'cockpit.duration8hours',
  'console.conversations.duration.12hours': 'cockpit.duration12hours',
  'console.conversations.duration.24hours': 'cockpit.duration24hours',
  // 桌面端有的但 web 没有的 conversations 键 → 新 cockpit.* 键
  'console.conversations.allKinds': 'cockpit.allKinds',
  'console.conversations.refresh': 'cockpit.refresh',
  'console.conversations.total': 'cockpit.total',
  'console.conversations.overrideTitle': 'cockpit.overrideTitle',
  'console.conversations.overrideHint': 'cockpit.overrideHint',
  'console.conversations.overrideSequence': 'cockpit.overrideSequence',
  'console.conversations.overrideRequired': 'cockpit.overrideRequired',
  'console.conversations.saveOverride': 'cockpit.saveOverride',
  'console.conversations.noChannelsForKind': 'cockpit.noChannelsForKind',

  // console.logs.* → channelLogs.*
  'console.logs.title': 'channelLogs.title',
  'console.logs.empty': 'channelLogs.empty',
  'console.logs.time': 'channelLogs.duration.total',
  'console.logs.source': 'channelLogs.sourceCapabilityTest',
  'console.logs.retry': 'channelLogs.retry',
  'console.logs.statusCode': 'channelLogs.status.connecting',
  'console.logs.autoRefresh': 'channelLogs.autoRefresh',
  'console.logs.autoRefreshing': 'channelLogs.autoRefreshing',

  // console.capability.* → capability.*
  'console.capability.title': 'capability.title',
  'console.capability.start': 'capability.startTest',
  'console.capability.cancel': 'capability.cancel',
  'console.capability.supported': 'capability.supported',
  'console.capability.unsupported': 'capability.unsupported',
  'console.capability.partial': 'capability.partial',
  'console.capability.protocolRunning': 'capability.protocolRunning',
  'console.capability.models': 'capability.modelsLabel',
  'console.capability.duration': 'capability.duration',
  'console.capability.noResults': 'capability.noResults',

  // console.actions.* → orchestration.* / app.actions.*
  'console.actions.label': 'orchestration.edit',
  'console.actions.edit': 'orchestration.edit',
  'console.actions.capability': 'capability.startTest',
  'console.actions.logs': 'orchestration.logs',
  'console.actions.copy': 'orchestration.copyConfig',
  'console.actions.website': 'orchestration.openWebsite',
  'console.actions.ping': 'app.actions.ping',
  'console.actions.enable': 'orchestration.enable',
  'console.actions.suspend': 'orchestration.pause',
  'console.actions.resume': 'orchestration.resume',
  'console.actions.promote': 'orchestration.promotion',
  'console.actions.disable': 'orchestration.moveToPool',
  'console.actions.delete': 'orchestration.delete',
  'console.actions.resetCircuit': 'orchestration.resumeReset',
  'console.actions.refresh': 'app.actions.refresh',

  // console.mode.* → orchestration.*
  'console.mode.multi': 'orchestration.multiChannel',
  'console.mode.single': 'orchestration.singleChannel',

  // console.pool.* → orchestration.*
  'console.pool.active': 'orchestration.failoverSequence',
  'console.pool.current': 'orchestration.failoverSequence',
  'console.pool.inactive': 'orchestration.standbyPool',

  // console.keys.* → channelCard.*
  'console.keys.active': 'channelCard.configuredKeys',
  'console.keys.disabled': 'channelCard.disabledKeys',

  // console.channel.* → orchestration.*
  'console.channel.cacheWriteHigh': 'orchestration.cacheWriteHigh',
  'console.channel.cacheWriteHighHint': 'orchestration.cacheWriteHighHint',

  // console.channelStatus.* → orchestration.* / channelCard.*
  'console.channelStatus.active': 'orchestration.enable',
  'console.channelStatus.suspended': 'channelCard.suspended',
  'console.channelStatus.disabled': 'orchestration.moveToPool',

  // console.circuit.* → status.*
  'console.circuit.open': 'status.tripped',
  'console.circuit.halfOpen': 'status.tripped',

  // console.fuzzy* → toast.* / tooltip.*
  'console.fuzzyEnabled': 'tooltip.fuzzyEnabled',
  'console.fuzzyDisabled': 'tooltip.fuzzyDisabled',
  'console.fuzzyLoadFailed': 'toast.loadFuzzyFailed',

  // console.cbSettings → tooltip.*
  'console.cbSettings': 'tooltip.circuitBreakerSettings',

  // console.searchChannels → orchestration.*
  'console.searchChannels': 'orchestration.searchPlaceholder',
  'console.pingAll': 'app.actions.ping',
  'console.addChannel': 'app.actions.addChannel',
  'console.noChannels': 'orchestration.noActiveChannels',
  'console.noSearchResults': 'orchestration.searchPlaceholder',
  'console.channelsTab': 'app.tabs.messages',
  'console.conversationsTab': 'app.tabs.conversations',
}

// 读取 JSON
function readJSON(path: string): Record<string, string> {
  return JSON.parse(readFileSync(path, 'utf-8'))
}

function writeJSON(path: string, data: Record<string, string>) {
  writeFileSync(path, JSON.stringify(data, null, 2) + '\n', 'utf-8')
}

// 重命名 JSON keys
function renameKeys(data: Record<string, string>, map: Record<string, string>): Record<string, string> {
  const result: Record<string, string> = {}
  for (const [key, value] of Object.entries(data)) {
    const newKey = map[key] ?? key
    // 如果目标 key 已存在（被其他映射覆盖），保留已存在的值（web 对齐的值更权威）
    if (result[newKey] !== undefined) {
      console.log(`  ⚠ 冲突: ${key} → ${newKey}（保留已有值）`)
      continue
    }
    result[newKey] = value
  }
  return result
}

// 替换组件文件中的 key 引用
function replaceInFiles(map: Record<string, string>) {
  const glob = new Bun.Glob('src/**/*.vue')
  const tsGlob = new Bun.Glob('src/**/*.ts')
  const files: string[] = []
  for (const match of glob.scanSync({ cwd: desktopDir })) {
    files.push(resolve(desktopDir, match))
  }
  for (const match of tsGlob.scanSync({ cwd: desktopDir })) {
    const f = resolve(desktopDir, match)
    if (f.includes('locales/') || f.includes('node_modules/') || f.includes('.test.')) continue
    files.push(f)
  }

  let totalReplacements = 0
  for (const file of files) {
    let content = readFileSync(file, 'utf-8')
    let changed = false
    for (const [oldKey, newKey] of Object.entries(map)) {
      // 匹配 t('oldKey') 和 tf('oldKey' 两种模式
      const escaped = oldKey.replace(/\./g, '\\.')
      const regex = new RegExp(`(t|tf)\\('${escaped}'`, 'g')
      if (regex.test(content)) {
        content = content.replace(new RegExp(`(t|tf)\\('${escaped}'`, 'g'), `$1('${newKey}'`)
        changed = true
        totalReplacements++
      }
    }
    if (changed) {
      writeFileSync(file, content, 'utf-8')
      console.log(`  ✓ 更新组件: ${file.replace(desktopDir + '/', '')}`)
    }
  }
  console.log(`  共替换 ${totalReplacements} 处引用`)
}

// 执行
console.log('=== Batch 1: console.* → web 命名空间 ===\n')

console.log('1. 重命名 en.json keys...')
const enData = readJSON(jsonEn)
const newEn = renameKeys(enData, keyMap)
writeJSON(jsonEn, newEn)
console.log(`   en.json: ${Object.keys(enData).length} → ${Object.keys(newEn).length} keys\n`)

console.log('2. 重命名 zh-CN.json keys...')
const zhData = readJSON(jsonZh)
const newZh = renameKeys(zhData, keyMap)
writeJSON(jsonZh, newZh)
console.log(`   zh-CN.json: ${Object.keys(zhData).length} → ${Object.keys(newZh).length} keys\n`)

console.log('3. 替换组件中的 t()/tf() 调用...')
replaceInFiles(keyMap)

console.log('\n✅ Batch 1 完成')
