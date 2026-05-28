import { defaultLocale, messages } from './messages'
import type { MessageKey, SupportedLocale } from './messages'

const supportedLocales: SupportedLocale[] = ['en', 'zh-CN']

export const isSupportedLocale = (value?: string | null): value is SupportedLocale => {
  if (!value) return false
  return supportedLocales.some((locale) => locale === value)
}

export const normalizeLocale = (value?: string | null): SupportedLocale => {
  if (!value) return defaultLocale
  const normalized = value.replace(/_/g, '-')
  const lower = normalized.toLowerCase()
  if (lower === 'zh' || lower.startsWith('zh-')) {
    return 'zh-CN'
  }
  if (lower === 'en' || lower.startsWith('en-')) {
    return 'en'
  }
  return defaultLocale
}

export const resolveInitialLocale = (persistedLocale?: string | null, systemLocale?: string | null): SupportedLocale => {
  if (persistedLocale && isSupportedLocale(persistedLocale)) {
    return persistedLocale
  }
  return normalizeLocale(systemLocale)
}

export const translate = (locale: SupportedLocale, key: MessageKey, params?: Record<string, string>): string => {
  const value = messages[locale]?.[key] ?? messages[defaultLocale][key] ?? key
  if (!params) return value
  return Object.entries(params).reduce((acc, [k, v]) => acc.replaceAll(`{${k}}`, v), value)
}

export const applyDocumentLanguage = (locale: SupportedLocale) => {
  try {
    document.documentElement.lang = locale
  } catch {
    // SSR 或测试环境可忽略
  }
}
