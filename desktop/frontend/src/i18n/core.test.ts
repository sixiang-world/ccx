import { describe, expect, it } from 'vitest'
import { applyDocumentLanguage, normalizeLocale, resolveInitialLocale, translate } from './core'

describe('normalizeLocale', () => {
  it('returns zh-CN for Chinese variants', () => {
    expect(normalizeLocale('zh')).toBe('zh-CN')
    expect(normalizeLocale('zh-CN')).toBe('zh-CN')
    expect(normalizeLocale('zh_CN')).toBe('zh-CN')
    expect(normalizeLocale('zh-Hans')).toBe('zh-CN')
    expect(normalizeLocale('zh-Hans-CN')).toBe('zh-CN')
    expect(normalizeLocale('zh_CN.UTF-8')).toBe('zh-CN')
  })

  it('returns en for English and unknown locales', () => {
    expect(normalizeLocale('en')).toBe('en')
    expect(normalizeLocale('en-US')).toBe('en')
    expect(normalizeLocale('fr-FR')).toBe('en')
    expect(normalizeLocale(undefined)).toBe('en')
    expect(normalizeLocale(null)).toBe('en')
  })
})

describe('resolveInitialLocale', () => {
  it('prefers supported persisted locale', () => {
    expect(resolveInitialLocale('zh-CN', 'en-US')).toBe('zh-CN')
    expect(resolveInitialLocale('en', 'zh_CN')).toBe('en')
  })

  it('falls back to normalized system locale when persisted unsupported', () => {
    expect(resolveInitialLocale(null, 'zh_CN.UTF-8')).toBe('zh-CN')
    expect(resolveInitialLocale(undefined, 'fr-FR')).toBe('en')
  })

  it('defaults to en when no data', () => {
    expect(resolveInitialLocale(null, null)).toBe('en')
  })
})

describe('translate', () => {
  it('returns locale message when available', () => {
    expect(translate('zh-CN', 'nav.status')).toBe('网关监控')
  })

  it('falls back to English message', () => {
    expect(translate('en', 'nav.status')).toBe('Status')
  })

  it('substitutes parameters in messages', () => {
    expect(translate('en', 'env.fieldMin', { field: 'Port', min: '1000' })).toBe('Port must be at least 1000')
  })

  it('falls back to key when message missing', () => {
    expect(translate('en', 'common.missing' as never)).toBe('common.missing')
  })
})

describe('applyDocumentLanguage', () => {
  it('sets document.documentElement.lang', () => {
    document.documentElement.lang = ''
    applyDocumentLanguage('zh-CN')
    expect(document.documentElement.lang).toBe('zh-CN')
  })
})
