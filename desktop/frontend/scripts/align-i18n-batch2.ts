#!/usr/bin/env bun
/**
 * Batch 2: 将 desktop console.form.* key 映射到 channelEditor.* 命名空间。
 * 用法：bun run scripts/align-i18n-batch2.ts
 */

import { readFileSync, writeFileSync } from 'fs'
import { resolve } from 'path'

const desktopDir = resolve(__dirname, '..')
const jsonEn = resolve(desktopDir, 'src/locales/en.json')
const jsonZh = resolve(desktopDir, 'src/locales/zh-CN.json')

const keyMap: Record<string, string> = {
  // ── Dialog Titles / Actions ──
  'console.form.addChannel': 'channelEditor.title.create',
  'console.form.editChannel': 'channelEditor.title.edit',
  'console.form.cancel': 'channelEditor.actions.cancel',
  'console.form.create': 'channelEditor.actions.create',
  'console.form.save': 'channelEditor.actions.save',

  // ── Basic Info ──
  'console.form.basicInfo': 'channelEditor.nav.basic',
  'console.form.basicInfoTitle': 'channelEditor.basic.info.title',
  'console.form.name': 'channelEditor.basic.name.label',
  'console.form.nameRequired': 'channelEditor.basic.name.required',
  'console.form.serviceType': 'channelEditor.basic.serviceType.label',
  'console.form.selectServiceType': 'channelEditor.basic.serviceType.placeholder',
  'console.form.serviceTypeRequired': 'channelEditor.basic.serviceType.required',
  'console.form.description': 'channelEditor.basic.description.label',
  'console.form.descriptionPlaceholder': 'channelEditor.basic.description.placeholder',
  'console.form.website': 'channelEditor.basic.website.label',
  'console.form.sectionBasic': 'channelEditor.nav.basic',

  // ── Connection ──
  'console.form.connection': 'channelEditor.nav.connection',
  'console.form.connectionTitle': 'channelEditor.connection.title',
  'console.form.baseUrl': 'channelEditor.basic.baseUrl.label',
  'console.form.baseUrlRequired': 'channelEditor.basic.baseUrl.required',
  'console.form.baseUrlPlaceholder': 'channelEditor.basic.baseUrl.placeholder',
  'console.form.additionalUrls': 'channelEditor.basic.additionalUrls.label',
  'console.form.multiLineFailover': 'channelEditor.basic.multiLineFailover',
  'console.form.expectedEndpoint': 'channelEditor.basic.expectedEndpoint',
  'console.form.proxyUrl': 'channelEditor.transport.proxyUrl.label',
  'console.form.proxyUrlLabel': 'channelEditor.transport.proxyUrl.label',
  'console.form.proxyUrlHint': 'channelEditor.transport.proxyUrl.hint',
  'console.form.routePrefix': 'channelEditor.transport.routePrefix.label',
  'console.form.routePrefixLabel': 'channelEditor.transport.routePrefix.label',
  'console.form.routePrefixHint': 'channelEditor.transport.routePrefix.hint',
  'console.form.insecureSkipVerify': 'channelEditor.transport.skipTls.label',
  'console.form.insecureSkipVerifyHint': 'channelEditor.transport.skipTls.hint',
  'console.form.transportTitle': 'channelEditor.transport.title',

  // ── Authentication ──
  'console.form.authentication': 'channelEditor.nav.auth',
  'console.form.sectionAuth': 'channelEditor.nav.auth',
  'console.form.apiKeys': 'channelEditor.auth.keys.label',
  'console.form.apiKeyRequired': 'channelEditor.auth.apiKeyRequired',
  'console.form.disabledKeys': 'channelEditor.auth.disabledKeys.label',
  'console.form.restoreKey': 'channelEditor.auth.restoreKey',
  'console.form.historicalKeys': 'channelEditor.auth.historicalKeys',

  // ── Model Scope ──
  'console.form.models': 'channelEditor.mapping.models',
  'console.form.modelScope': 'channelEditor.mapping.modelScope',
  'console.form.supportedModels': 'channelEditor.mapping.supportedModels',
  'console.form.supportedModelsLabel': 'channelEditor.mapping.supportedModels.label',
  'console.form.fetchModels': 'channelEditor.mapping.fetchModels',
  'console.form.fetchingModels': 'channelEditor.mapping.fetchingModels',
  'console.form.modelFetchNeedsConfig': 'channelEditor.mapping.modelFetchNeedsConfig',

  // ── Model Redirect ──
  'console.form.sectionRedirect': 'channelEditor.mapping.redirect.title',
  'console.form.modelRedirect': 'channelEditor.mapping.redirect.title',
  'console.form.modelMapping': 'channelEditor.mapping.configured.label',
  'console.form.mappingHint': 'channelEditor.mapping.hint',
  'console.form.modelMappingExisting': 'channelEditor.mapping.configured.label',
  'console.form.modelMappingAdd': 'channelEditor.mapping.addRedirect',
  'console.form.reasoningEffort': 'channelEditor.mapping.reasoningEffort.label',

  // ── Vision ──
  'console.form.visionTitle': 'channelEditor.compat.vision.title',
  'console.form.visionEnabled': 'channelEditor.compat.visionEnabled',
  'console.form.visionDisabled': 'channelEditor.compat.visionDisabled',
  'console.form.visionFallbackModel': 'channelEditor.compat.visionFallback.label',
  'console.form.visionFallbackHint': 'channelEditor.compat.visionFallback.hint',
  'console.form.noVision': 'channelEditor.compat.noVision.label',
  'console.form.noVisionHint': 'channelEditor.compat.noVision.hint',
  'console.form.noVisionModels': 'channelEditor.compat.noVisionModels',
  'console.form.historicalImageTurnLimit': 'channelEditor.compat.historicalImageLimit.label',
  'console.form.historicalImageTurnLimitHint': 'channelEditor.compat.historicalImageLimit.hint',

  // ── Compatibility ──
  'console.form.protocolOptions': 'channelEditor.compat.title',
  'console.form.sectionAdvanced': 'channelEditor.nav.advanced',
  'console.form.advancedFlags': 'channelEditor.compat.flags',
  'console.form.compatibilityTitle': 'channelEditor.compat.title',
  'console.form.generationParams': 'channelEditor.compat.generationParams',
  'console.form.reasoningParamStyle': 'channelEditor.compat.reasoningStyle.label',
  'console.form.reasoningParamStyleHint': 'channelEditor.compat.reasoningStyle.hint',
  'console.form.textVerbosity': 'channelEditor.compat.textVerbosity.label',
  'console.form.textVerbosityStyle': 'channelEditor.compat.textVerbosity.style',
  'console.form.textVerbosityPlaceholder': 'channelEditor.compat.textVerbosity.placeholder',

  // ── Advanced Flags ──
  'console.form.passbackReasoning': 'channelEditor.compat.passbackReasoning.label',
  'console.form.passbackReasoningHint': 'channelEditor.compat.passbackReasoning.hint',
  'console.form.passbackThinking': 'channelEditor.compat.passbackThinking.label',
  'console.form.passbackThinkingHint': 'channelEditor.compat.passbackThinking.hint',
  'console.form.fastMode': 'channelEditor.compat.fastMode.label',
  'console.form.fastModeHint': 'channelEditor.compat.fastMode.hint',
  'console.form.lowQuality': 'channelEditor.transport.lowQuality.label',
  'console.form.lowQualityHint': 'channelEditor.transport.lowQuality.hint',
  'console.form.injectDummySignature': 'channelEditor.compat.injectDummySignature.label',
  'console.form.injectDummySignatureHint': 'channelEditor.compat.injectDummySignature.hint',
  'console.form.stripThoughtSignature': 'channelEditor.compat.stripThoughtSignature.label',
  'console.form.stripThoughtSignatureHint': 'channelEditor.compat.stripThoughtSignature.hint',
  'console.form.stripEmptyBlocks': 'channelEditor.compat.stripEmptyBlocks.label',
  'console.form.stripEmptyBlocksHint': 'channelEditor.compat.stripEmptyBlocks.hint',
  'console.form.normalizeSystem': 'channelEditor.compat.normalizeSystem.label',
  'console.form.normalizeSystemHint': 'channelEditor.compat.normalizeSystem.hint',
  'console.form.normalizeUserId': 'channelEditor.compat.normalizeUserId.label',
  'console.form.normalizeUserIdHint': 'channelEditor.compat.normalizeUserId.hint',
  'console.form.stripBillingHeader': 'channelEditor.compat.stripBillingHeader.label',
  'console.form.stripBillingHeaderHint': 'channelEditor.compat.stripBillingHeader.hint',
  'console.form.normalizeChatRoles': 'channelEditor.compat.normalizeRoles.label',
  'console.form.normalizeChatRolesHint': 'channelEditor.compat.normalizeRoles.hint',
  'console.form.autoBlacklist': 'channelEditor.runtime.autoBlacklist.label',
  'console.form.autoBlacklistHint': 'channelEditor.runtime.autoBlacklist.hint',
  'console.form.autoBlacklistBalanceLabel': 'channelEditor.compat.autoBlacklistBalance.label',
  'console.form.autoBlacklistBalanceHint': 'channelEditor.compat.autoBlacklistBalance.hint',
  'console.form.codexNativeTools': 'channelEditor.compat.codexNativeTools.label',
  'console.form.codexNativeToolsHint': 'channelEditor.compat.codexNativeTools.hint',
  'console.form.codexCompat': 'channelEditor.compat.codexCompat.label',
  'console.form.codexCompatHint': 'channelEditor.compat.codexCompat.hint',
  'console.form.stripCodexTools': 'channelEditor.compat.stripCodexTools.label',
  'console.form.compactModel': 'channelEditor.compat.compactModel.label',
  'console.form.compactModelHint': 'channelEditor.compat.compactModel.hint',

  // ── Custom Headers ──
  'console.form.customHeaders': 'channelEditor.nav.custom',
  'console.form.sectionHeaders': 'channelEditor.nav.custom',

  // ── Request Timeout ──
  'console.form.requestTimeoutMs': 'channelEditor.transport.requestTimeout.label',
  'console.form.requestTimeoutLabel': 'channelEditor.transport.requestTimeout.label',
  'console.form.requestTimeoutMsHint': 'channelEditor.transport.requestTimeout.hint',
  'console.form.requestTimeoutInvalid': 'channelEditor.transport.requestTimeout.invalid',

  // ── Stream Timeout ──
  'console.form.streamTimeouts': 'channelEditor.streamTimeout.title',
  'console.form.streamTimeoutTitle': 'channelEditor.streamTimeout.title',
  'console.form.streamFirstContentTimeoutOverrideLabel': 'channelEditor.streamTimeout.firstContent.overrideLabel',
  'console.form.streamFirstContentTimeoutLabel': 'channelEditor.streamTimeout.firstContent.label',
  'console.form.streamInactivityTimeoutOverrideLabel': 'channelEditor.streamTimeout.inactivity.overrideLabel',
  'console.form.streamInactivityTimeoutLabel': 'channelEditor.streamTimeout.inactivity.label',
  'console.form.streamToolCallIdleTimeoutOverrideLabel': 'channelEditor.streamTimeout.toolCallIdle.overrideLabel',
  'console.form.streamToolCallIdleTimeoutLabel': 'channelEditor.streamTimeout.toolCallIdle.label',
  'console.form.streamTimeoutOverrideHint': 'channelEditor.streamTimeout.overrideHint',
  'console.form.streamTimeoutInheritHint': 'channelEditor.streamTimeout.inheritHint',
  'console.form.firstByteWait': 'channelEditor.streamTimeout.firstContent.label',
  'console.form.firstByteWaitHint': 'channelEditor.streamTimeout.firstContent.hint',
  'console.form.idleAfterFirstByte': 'channelEditor.streamTimeout.inactivity.label',
  'console.form.idleAfterFirstByteHint': 'channelEditor.streamTimeout.inactivity.hint',
  'console.form.toolCallIdle': 'channelEditor.streamTimeout.toolCallIdle.label',
  'console.form.toolCallIdleHint': 'channelEditor.streamTimeout.toolCallIdle.hint',
  'console.form.timeoutThreshold': 'channelEditor.streamTimeout.timeoutThreshold',
  'console.form.presetGentle': 'channelEditor.streamTimeout.preset.gentle',
  'console.form.presetBalanced': 'channelEditor.streamTimeout.preset.balanced',
  'console.form.presetAggressive': 'channelEditor.streamTimeout.preset.aggressive',

  // ── Rate Limit ──
  'console.form.rateLimitTitle': 'channelEditor.rateLimit.title',
  'console.form.rateLimitSectionLabel': 'channelEditor.rateLimit.section.label',
  'console.form.rateLimitSectionHint': 'channelEditor.rateLimit.section.hint',
  'console.form.rateLimitRpmLabel': 'channelEditor.rateLimit.rpm.label',
  'console.form.rateLimitRpmHint': 'channelEditor.rateLimit.rpm.hint',
  'console.form.rateLimitBurstLabel': 'channelEditor.rateLimit.burst.label',
  'console.form.rateLimitBurstHint': 'channelEditor.rateLimit.burst.hint',
  'console.form.rateLimitMaxConcurrentLabel': 'channelEditor.rateLimit.maxConcurrent.label',
  'console.form.rateLimitMaxConcurrentHint': 'channelEditor.rateLimit.maxConcurrent.hint',
  'console.form.rateLimitAutoFromHeadersLabel': 'channelEditor.runtime.autoLearnRateLimits.label',
  'console.form.rateLimitAutoFromHeadersHint': 'channelEditor.runtime.autoLearnRateLimits.hint',
  'console.form.rateLimitAutoLabel': 'channelEditor.runtime.autoLearnRateLimits.label',
  'console.form.rateLimitAutoHint': 'channelEditor.runtime.autoLearnRateLimits.hint',
  'console.form.rpmLabel': 'channelEditor.rateLimit.rpm.label',
  'console.form.rpmPlaceholder': 'channelEditor.rateLimit.rpm.placeholder',
  'console.form.windowLabel': 'channelEditor.rateLimit.window.label',
  'console.form.windowPlaceholder': 'channelEditor.rateLimit.window.placeholder',
  'console.form.maxConcurrentLabel': 'channelEditor.rateLimit.maxConcurrent.label',
  'console.form.maxConcurrentPlaceholder': 'channelEditor.rateLimit.maxConcurrent.placeholder',

  // ── Runtime ──
  'console.form.runtimeTitle': 'channelEditor.runtime.title',

  // ── Miscellaneous ──
  'console.form.outline': 'channelEditor.nav.outline',
  'console.form.selectDefault': 'channelEditor.compat.selectDefault',
  'console.form.selectDefaultLabel': 'channelEditor.compat.selectDefault',
}

function readJSON(path: string): Record<string, string> {
  return JSON.parse(readFileSync(path, 'utf-8'))
}
function writeJSON(path: string, data: Record<string, string>) {
  writeFileSync(path, JSON.stringify(data, null, 2) + '\n', 'utf-8')
}
function renameKeys(data: Record<string, string>, map: Record<string, string>): Record<string, string> {
  const result: Record<string, string> = {}
  for (const [key, value] of Object.entries(data)) {
    const newKey = map[key] ?? key
    if (result[newKey] !== undefined) {
      console.log(`  ⚠ 冲突: ${key} → ${newKey}（保留已有值）`)
      continue
    }
    result[newKey] = value
  }
  return result
}
function replaceInFiles(map: Record<string, string>) {
  const glob = new Bun.Glob('src/**/*.vue')
  const tsGlob = new Bun.Glob('src/**/*.ts')
  const files: string[] = []
  for (const match of glob.scanSync({ cwd: desktopDir })) files.push(resolve(desktopDir, match))
  for (const match of tsGlob.scanSync({ cwd: desktopDir })) {
    const f = resolve(desktopDir, match)
    if (f.includes('locales/') || f.includes('node_modules/') || f.includes('.test.')) continue
    files.push(f)
  }
  let total = 0
  for (const file of files) {
    let content = readFileSync(file, 'utf-8')
    let changed = false
    for (const [oldKey, newKey] of Object.entries(map)) {
      const escaped = oldKey.replace(/\./g, '\\.')
      const regex = new RegExp(`(t|tf)\\('${escaped}'`, 'g')
      if (regex.test(content)) {
        content = content.replace(new RegExp(`(t|tf)\\('${escaped}'`, 'g'), `$1('${newKey}'`)
        changed = true
        total++
      }
    }
    if (changed) {
      writeFileSync(file, content, 'utf-8')
      console.log(`  ✓ ${file.replace(desktopDir + '/', '')}`)
    }
  }
  console.log(`  共替换 ${total} 处引用`)
}

console.log('=== Batch 2: console.form.* → channelEditor.* ===\n')

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

console.log('\n✅ Batch 2 完成')
