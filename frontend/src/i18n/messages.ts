import { enMessages } from './messages-en'
import { idMessages } from './messages-id'
import { zhCNMessages } from './messages-zh-cn'

export type SupportedLocale = 'en' | 'id' | 'zh-CN'

// 自动从实际的消息对象推导类型，避免手动维护
export type MessageKey = keyof typeof zhCNMessages

export const messages: Record<SupportedLocale, Record<MessageKey, string>> = {
  en: enMessages,
  id: idMessages,
  'zh-CN': zhCNMessages
}
