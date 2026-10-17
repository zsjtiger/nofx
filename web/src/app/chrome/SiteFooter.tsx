import { t, type Language } from '../../shared/i18n/translations'

interface SiteFooterProps {
  language: Language
}

export function SiteFooter({ language }: SiteFooterProps) {
  return (
    <footer
      className="mt-16"
      style={{
        borderTop: '1px solid rgba(26,24,19,0.14)',
        background: '#F7F4EC',
      }}
    >
      <div
        className="max-w-[1920px] mx-auto px-6 py-6 text-center text-sm"
        style={{ color: '#8A8478' }}
      >
        <p>{t('footerTitle', language)}</p>
        <p className="mt-1">{t('footerWarning', language)}</p>
      </div>
    </footer>
  )
}
