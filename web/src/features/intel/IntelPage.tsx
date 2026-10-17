import { useCallback, useEffect, useState } from 'react'
import { Loader2 } from 'lucide-react'
import { useLanguage } from '../../shared/contexts/LanguageContext'
import { api } from '../../shared/lib/api'
import { formatRequestError } from '../../shared/lib/requestErrors'
import { t } from '../../shared/i18n/translations'
import type { AIModel, IntelItem, IntelRunResult } from '../../shared/types'
import { parseIntelUrls } from './parseUrls'

const RISK_COLOR: Record<string, string> = {
  high: 'var(--crux-red-text)',
  low: 'var(--crux-warning-text)',
  none: 'var(--crux-green-text)',
  unknown: 'var(--text-tertiary)',
}

const panelStyle = {
  background: 'var(--panel-bg)',
  border: '1px solid var(--panel-border)',
  borderRadius: 12,
} as const

const inputStyle = {
  background: 'var(--panel-bg-hover)',
  border: '1px solid var(--panel-border)',
  borderRadius: 8,
  color: 'var(--text-primary)',
  padding: '8px 10px',
  width: '100%',
  fontFamily: 'var(--tm-mono)',
  fontSize: 14,
} as const

export function IntelPage() {
  const { language } = useLanguage()
  const [symbol, setSymbol] = useState('')
  const [urlText, setUrlText] = useState('')
  const [models, setModels] = useState<AIModel[]>([])
  const [modelId, setModelId] = useState('')
  const [running, setRunning] = useState(false)
  const [results, setResults] = useState<IntelRunResult[] | null>(null)
  const [history, setHistory] = useState<IntelItem[]>([])
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    api
      .getModelConfigs()
      .then((list) => {
        const enabled = list.filter((m) => m.enabled)
        setModels(enabled)
        setModelId((current) => current || enabled[0]?.id || '')
      })
      .catch(() => setModels([]))
  }, [])

  const loadHistory = useCallback(
    async (rawSymbol: string) => {
      const sym = rawSymbol.trim().toUpperCase()
      if (!sym) return
      try {
        setHistory(await api.getIntel(sym))
      } catch (err) {
        setError(formatRequestError(err, language))
      }
    },
    [language]
  )

  const runAnalysis = async () => {
    setError(null)
    const sym = symbol.trim().toUpperCase()
    if (!sym) {
      setError(t('intel.symbol', language))
      return
    }
    const parsed = parseIntelUrls(urlText)
    if (parsed.error === 'empty') {
      setError(t('intel.invalidUrls', language))
      return
    }
    if (parsed.error === 'tooMany') {
      setError(t('intel.tooMany', language))
      return
    }
    if (!modelId) {
      setError(t('intel.noModel', language))
      return
    }

    setRunning(true)
    try {
      const res = await api.runIntel({
        symbol: sym,
        urls: parsed.urls,
        model_id: modelId,
      })
      setResults(res)
      await loadHistory(sym)
    } catch (err) {
      const notConfigured = (err as { statusCode?: number })?.statusCode === 503
      setError(
        notConfigured
          ? t('intel.notConfigured', language)
          : formatRequestError(err, language)
      )
    } finally {
      setRunning(false)
    }
  }

  return (
    <div
      className="mx-auto max-w-5xl px-4 py-6 space-y-6"
      style={{ fontFamily: 'var(--tm-mono)' }}
    >
      <header>
        <h1 className="text-xl" style={{ color: 'var(--text-primary)' }}>
          {t('intel.title', language)}
        </h1>
        <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>
          {t('intel.subtitle', language)}
        </p>
      </header>

      <section className="p-4 space-y-4" style={panelStyle}>
        <div className="grid gap-4 md:grid-cols-2">
          <label className="text-sm space-y-1">
            <span style={{ color: 'var(--text-secondary)' }}>
              {t('intel.symbol', language)}
            </span>
            <input
              style={inputStyle}
              value={symbol}
              placeholder={t('intel.symbolHint', language)}
              onChange={(e) => setSymbol(e.target.value)}
              onBlur={() => loadHistory(symbol)}
            />
          </label>
          <label className="text-sm space-y-1">
            <span style={{ color: 'var(--text-secondary)' }}>
              {t('intel.model', language)}
            </span>
            <select
              style={inputStyle}
              value={modelId}
              onChange={(e) => setModelId(e.target.value)}
              disabled={models.length === 0}
            >
              {models.length === 0 && (
                <option value="">{t('intel.noModel', language)}</option>
              )}
              {models.map((m) => (
                <option key={m.id} value={m.id}>
                  {m.name}
                </option>
              ))}
            </select>
          </label>
        </div>

        <label className="text-sm space-y-1 block">
          <span style={{ color: 'var(--text-secondary)' }}>
            {t('intel.urls', language)}
            <span
              className="ml-2 text-xs"
              style={{ color: 'var(--text-tertiary)' }}
            >
              {t('intel.urlsHint', language)}
            </span>
          </span>
          <textarea
            style={{ ...inputStyle, minHeight: 110 }}
            value={urlText}
            onChange={(e) => setUrlText(e.target.value)}
          />
        </label>

        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={runAnalysis}
            disabled={running}
            className="btn-crux px-4 py-2 rounded-lg text-sm inline-flex items-center gap-2"
          >
            {running && <Loader2 className="h-4 w-4 animate-spin" />}
            {running ? t('intel.running', language) : t('intel.run', language)}
          </button>
          {error && (
            <span
              role="alert"
              className="text-sm"
              style={{ color: 'var(--crux-red-text)' }}
            >
              {error}
            </span>
          )}
        </div>
      </section>

      {results && (
        <section className="space-y-3">
          <h2 className="text-base" style={{ color: 'var(--text-primary)' }}>
            {t('intel.results', language)}
          </h2>
          {results.map((r) => (
            <IntelRow
              key={r.url}
              url={r.url}
              item={r.item}
              error={r.error}
              language={language}
            />
          ))}
        </section>
      )}

      <section className="space-y-3">
        <h2 className="text-base" style={{ color: 'var(--text-primary)' }}>
          {t('intel.history', language)}
          {symbol.trim() && (
            <span
              className="ml-2 text-sm"
              style={{ color: 'var(--text-tertiary)' }}
            >
              {symbol.trim().toUpperCase()}
            </span>
          )}
        </h2>
        {history.length === 0 ? (
          <p className="text-sm" style={{ color: 'var(--text-tertiary)' }}>
            {t('intel.empty', language)}
          </p>
        ) : (
          history.map((item) => (
            <IntelRow
              key={item.id}
              url={item.source_url}
              item={item}
              language={language}
            />
          ))
        )}
      </section>
    </div>
  )
}

function IntelRow({
  url,
  item,
  error,
  language,
}: {
  url: string
  item?: IntelItem
  error?: string
  language: 'zh' | 'en' | 'id'
}) {
  return (
    <article className="p-4 space-y-2" style={panelStyle}>
      <a
        href={url}
        target="_blank"
        rel="noreferrer"
        className="text-sm break-all"
        style={{ color: 'var(--crux-accent-text)' }}
      >
        {url}
      </a>
      {error && (
        <p className="text-sm" style={{ color: 'var(--crux-red-text)' }}>
          {t('intel.failed', language)}: {error}
        </p>
      )}
      {item && (
        <>
          <div className="flex flex-wrap items-center gap-3 text-xs">
            <span
              className="px-2 py-0.5 rounded"
              style={{
                color: RISK_COLOR[item.risk] ?? RISK_COLOR.unknown,
                border: '1px solid currentColor',
              }}
            >
              {t(
                `intel.risk.${item.risk in RISK_COLOR ? item.risk : 'unknown'}`,
                language
              )}
            </span>
            <span
              style={{
                color:
                  item.sentiment >= 0
                    ? 'var(--crux-green-text)'
                    : 'var(--crux-red-text)',
              }}
            >
              {t('intel.sentiment', language)} {item.sentiment >= 0 ? '+' : ''}
              {item.sentiment.toFixed(2)}
            </span>
            {item.symbols && (
              <span style={{ color: 'var(--text-secondary)' }}>
                {t('intel.mentions', language)}: {item.symbols}
              </span>
            )}
          </div>
          <p className="text-sm" style={{ color: 'var(--text-primary)' }}>
            {item.summary}
          </p>
        </>
      )}
    </article>
  )
}
