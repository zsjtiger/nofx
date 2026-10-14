import {
  createContext,
  useCallback,
  useContext,
  useState,
  ReactNode,
} from 'react'
import type { Language } from '../i18n/translations'

export const DEFAULT_LANGUAGE: Language = 'zh'
export const SWITCHABLE_LANGUAGES: Language[] = ['zh', 'en']

const STORAGE_KEY = 'language'

interface LanguageContextType {
  language: Language
  setLanguage: (lang: Language) => void
}

const LanguageContext = createContext<LanguageContextType | undefined>(
  undefined
)

function readStoredLanguage(): Language {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored && SWITCHABLE_LANGUAGES.includes(stored as Language)) {
      return stored as Language
    }
  } catch {
    // Storage can be unavailable (private mode, blocked site data); use the default.
  }
  return DEFAULT_LANGUAGE
}

export function LanguageProvider({ children }: { children: ReactNode }) {
  const [language, setLanguageState] = useState<Language>(readStoredLanguage)

  const setLanguage = useCallback((lang: Language) => {
    if (!SWITCHABLE_LANGUAGES.includes(lang)) return
    setLanguageState(lang)
    try {
      localStorage.setItem(STORAGE_KEY, lang)
    } catch {
      // The choice still applies for this page view.
    }
  }, [])

  return (
    <LanguageContext.Provider value={{ language, setLanguage }}>
      {children}
    </LanguageContext.Provider>
  )
}

export function useLanguage() {
  const context = useContext(LanguageContext)
  if (!context) {
    throw new Error('useLanguage must be used within LanguageProvider')
  }
  return context
}
