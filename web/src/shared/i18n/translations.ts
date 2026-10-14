import en from './locales/en'
import zh from './locales/zh'
import id from './locales/id'

export type Language = 'en' | 'zh' | 'id'

export const translations = { en, zh, id }
export function t(
  key: string,
  lang: Language,
  params?: Record<string, string | number>
): string {
  // Handle nested keys like 'twoStageKey.title'
  const keys = key.split('.')
  let value: any = translations[lang]

  for (const k of keys) {
    value = value?.[k]
  }

  let text = typeof value === 'string' ? value : key

  // Replace parameters like {count}, {gap}, etc.
  if (params) {
    Object.entries(params).forEach(([param, value]) => {
      text = text.replace(`{${param}}`, String(value))
    })
  }

  return text
}
