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
  const lookup = (locale: Language): unknown =>
    key.split('.').reduce<any>((node, k) => node?.[k], translations[locale])

  // Fall back to English when a translation is missing, instead of showing the raw key.
  const value = lookup(lang) ?? lookup('en')

  let text = typeof value === 'string' ? value : key

  // Replace parameters like {count}, {gap}, etc.
  if (params) {
    Object.entries(params).forEach(([param, value]) => {
      text = text.replace(`{${param}}`, String(value))
    })
  }

  return text
}
