import { useCallback, useEffect, useState } from 'react'
import { Copy, Loader2, Plus, Save, Trash2, Zap } from 'lucide-react'
import { useAuth } from '../../shared/contexts/AuthContext'
import { useLanguage } from '../../shared/contexts/LanguageContext'
import { DeepVoidBackground } from '../../shared/ui/DeepVoidBackground'
import { api } from '../../shared/lib/api'
import { confirmToast, notify } from '../../shared/lib/notify'
import type { Strategy, StrategyConfig } from '../../shared/types'

type Market = 'us_equity' | 'cn_equity'

const text = (language: string, zh: string, en: string) =>
  language === 'zh' ? zh : en

const SYMBOL_HINT: Record<Market, { zh: string; en: string; example: string }> =
  {
    us_equity: {
      zh: '美股代码，例如 AAPL、MSFT、BRK.B',
      en: 'US tickers, e.g. AAPL, MSFT, BRK.B',
      example: 'AAPL, MSFT, NVDA',
    },
    cn_equity: {
      zh: 'A 股代码，例如 600519.SH、000001.SZ，可省略交易所后缀',
      en: 'A-share codes, e.g. 600519.SH, 000001.SZ; the exchange suffix can be omitted',
      example: '600519.SH, 000001.SZ',
    },
  }

const INDICATOR_TOGGLES: Array<{
  key:
    | 'enable_ema'
    | 'enable_macd'
    | 'enable_rsi'
    | 'enable_atr'
    | 'enable_boll'
    | 'enable_volume'
  zh: string
  en: string
}> = [
  { key: 'enable_ema', zh: 'EMA 均线', en: 'EMA' },
  { key: 'enable_macd', zh: 'MACD', en: 'MACD' },
  { key: 'enable_rsi', zh: 'RSI', en: 'RSI' },
  { key: 'enable_atr', zh: 'ATR 波动', en: 'ATR' },
  { key: 'enable_boll', zh: '布林带', en: 'Bollinger Bands' },
  { key: 'enable_volume', zh: '成交量', en: 'Volume' },
]

// Symbols are stored upper-cased and de-duplicated; the server canonicalises them per market.
function parseSymbols(raw: string): string[] {
  const seen = new Set<string>()
  const out: string[] = []
  for (const part of raw.split(/[\s,，]+/)) {
    const symbol = part.trim().toUpperCase()
    if (symbol && !seen.has(symbol)) {
      seen.add(symbol)
      out.push(symbol)
    }
  }
  return out
}

function clampNumber(value: number, min: number, max: number) {
  if (Number.isNaN(value)) return min
  return Math.min(max, Math.max(min, value))
}

export function StrategyStudioPage() {
  const { token } = useAuth()
  const { language } = useLanguage()

  const [strategies, setStrategies] = useState<Strategy[]>([])
  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [draft, setDraft] = useState<StrategyConfig | null>(null)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [symbolsText, setSymbolsText] = useState('')
  const [market, setMarket] = useState<Market>('us_equity')
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [previewText, setPreviewText] = useState<string | null>(null)

  const selected = strategies.find((s) => s.id === selectedId) ?? null

  const loadStrategies = useCallback(async () => {
    const list = await api.getStrategies()
    setStrategies(list)
    return list
  }, [])

  useEffect(() => {
    if (!token) return
    let cancelled = false
    loadStrategies()
      .then((list) => {
        if (cancelled || list.length === 0) return
        setSelectedId(
          (current) =>
            current ?? list.find((s) => s.is_active)?.id ?? list[0].id
        )
      })
      .catch((err) => notify.error(String(err)))
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [token, loadStrategies])

  // Load the selected strategy into the editor.
  useEffect(() => {
    if (!selected) return
    setDraft(selected.config)
    setName(selected.name)
    setDescription(selected.description ?? '')
    setSymbolsText(
      (selected.config.ai_config.coin_source.static_coins ?? []).join(', ')
    )
    setPreviewText(null)
  }, [selectedId, selected?.updated_at])

  const updateAI = (patch: Partial<StrategyConfig['ai_config']>) => {
    setDraft((prev) =>
      prev ? { ...prev, ai_config: { ...prev.ai_config, ...patch } } : prev
    )
  }

  const updateRisk = (
    patch: Partial<StrategyConfig['ai_config']['risk_control']>
  ) => {
    if (!draft) return
    updateAI({ risk_control: { ...draft.ai_config.risk_control, ...patch } })
  }

  const updateIndicators = (
    patch: Partial<StrategyConfig['ai_config']['indicators']>
  ) => {
    if (!draft) return
    updateAI({ indicators: { ...draft.ai_config.indicators, ...patch } })
  }

  const updatePrompt = (
    patch: Partial<NonNullable<StrategyConfig['ai_config']['prompt_sections']>>
  ) => {
    if (!draft) return
    updateAI({
      prompt_sections: { ...(draft.ai_config.prompt_sections ?? {}), ...patch },
    })
  }

  const handleNew = async () => {
    try {
      const config = await api.getDefaultStrategyConfig()
      const created = await api.createStrategy({
        name: text(language, '新策略', 'New strategy'),
        description: '',
        config,
      })
      await loadStrategies()
      setSelectedId(created.id)
    } catch (err) {
      notify.error(String(err))
    }
  }

  const handleDuplicate = async () => {
    if (!selected) return
    try {
      const copy = await api.duplicateStrategy(selected.id)
      await loadStrategies()
      setSelectedId(copy.id)
    } catch (err) {
      notify.error(String(err))
    }
  }

  const handleDelete = async () => {
    if (!selected) return
    const ok = await confirmToast(
      text(
        language,
        `删除策略「${selected.name}」？`,
        `Delete strategy "${selected.name}"?`
      )
    )
    if (!ok) return
    try {
      await api.deleteStrategy(selected.id)
      const list = await loadStrategies()
      setSelectedId(list[0]?.id ?? null)
    } catch (err) {
      notify.error(String(err))
    }
  }

  const handleActivate = async () => {
    if (!selected) return
    try {
      await api.activateStrategy(selected.id)
      await loadStrategies()
      notify.success(
        text(language, '已设为当前策略', 'Set as the active strategy')
      )
    } catch (err) {
      notify.error(String(err))
    }
  }

  const handleSave = async () => {
    if (!selected || !draft) return
    const symbols = parseSymbols(symbolsText)
    if (symbols.length === 0) {
      notify.error(
        text(language, '请至少填写一个股票代码', 'Enter at least one symbol')
      )
      return
    }
    const config: StrategyConfig = {
      ...draft,
      ai_config: {
        ...draft.ai_config,
        coin_source: {
          ...draft.ai_config.coin_source,
          source_type: 'static',
          static_coins: symbols,
        },
      },
    }
    setSaving(true)
    try {
      await api.updateStrategy(selected.id, {
        name: name.trim() || selected.name,
        description,
        config,
      })
      await loadStrategies()
      notify.success(text(language, '已保存', 'Saved'))
    } catch (err) {
      notify.error(String(err))
    } finally {
      setSaving(false)
    }
  }

  const handlePreview = async () => {
    if (!draft) return
    try {
      const result = await api.previewPrompt(draft, market)
      setPreviewText(result.system_prompt)
    } catch (err) {
      notify.error(String(err))
    }
  }

  if (loading) {
    return (
      <div className="flex min-h-[70vh] items-center justify-center">
        <Loader2 className="h-7 w-7 animate-spin text-nofx-gold" />
      </div>
    )
  }

  const inputClass =
    'w-full rounded-lg border border-[rgba(26,24,19,0.2)] bg-nofx-bg px-3 py-2 text-sm text-nofx-text'
  const labelClass = 'mb-1 block text-sm text-nofx-text-muted'
  const sectionClass =
    'rounded-xl border border-[rgba(26,24,19,0.14)] bg-nofx-bg-deeper p-5 space-y-4'

  return (
    <DeepVoidBackground className="min-h-[calc(100vh-64px)] bg-nofx-bg">
      <div className="border-b border-[rgba(26,24,19,0.14)] bg-nofx-bg/75 px-5 py-4 backdrop-blur">
        <h1 className="text-xl font-semibold text-nofx-text">
          {text(language, '策略工作室', 'Strategy Studio')}
        </h1>
        <p className="mt-1 text-sm text-nofx-text-muted">
          {text(
            language,
            '基于日线数据的 AI 股票策略。风控规则由代码强制执行。',
            'AI stock strategies on daily bars. Risk limits are enforced in code.'
          )}
        </p>
      </div>

      <div className="grid grid-cols-1 gap-6 p-5 lg:grid-cols-[260px_1fr]">
        <aside className="space-y-3">
          <div className="flex gap-2">
            <button
              type="button"
              onClick={handleNew}
              className="inline-flex flex-1 items-center justify-center gap-2 rounded-lg bg-nofx-gold px-3 py-2 text-sm font-semibold text-nofx-bg"
            >
              <Plus className="h-4 w-4" /> {text(language, '新建', 'New')}
            </button>
          </div>
          <ul className="space-y-2">
            {strategies.map((s) => (
              <li key={s.id}>
                <button
                  type="button"
                  onClick={() => setSelectedId(s.id)}
                  className={`w-full rounded-lg border px-3 py-2 text-left text-sm ${
                    s.id === selectedId
                      ? 'border-nofx-gold bg-nofx-gold/10 text-nofx-text'
                      : 'border-[rgba(26,24,19,0.14)] text-nofx-text-muted hover:text-nofx-text'
                  }`}
                >
                  <span className="block font-medium">{s.name}</span>
                  {s.is_active && (
                    <span className="mt-1 inline-block rounded bg-nofx-gold/20 px-1.5 py-0.5 text-xs text-nofx-gold">
                      {text(language, '当前', 'Active')}
                    </span>
                  )}
                </button>
              </li>
            ))}
          </ul>
        </aside>

        {selected && draft ? (
          <div className="space-y-6">
            <div className="flex flex-wrap items-center gap-2">
              <button
                type="button"
                onClick={handleSave}
                disabled={saving}
                className="inline-flex items-center gap-2 rounded-lg bg-nofx-gold px-4 py-2 text-sm font-semibold text-nofx-bg disabled:opacity-50"
              >
                {saving ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <Save className="h-4 w-4" />
                )}
                {text(language, '保存', 'Save')}
              </button>
              <button
                type="button"
                onClick={handleActivate}
                disabled={selected.is_active}
                className="inline-flex items-center gap-2 rounded-lg border border-nofx-gold px-4 py-2 text-sm text-nofx-gold disabled:opacity-40"
              >
                <Zap className="h-4 w-4" />{' '}
                {text(language, '设为当前', 'Set active')}
              </button>
              <button
                type="button"
                onClick={handlePreview}
                className="inline-flex items-center gap-2 rounded-lg border border-[rgba(26,24,19,0.2)] px-4 py-2 text-sm text-nofx-text"
              >
                {text(language, '预览提示词', 'Preview prompt')}
              </button>
              <button
                type="button"
                onClick={handleDuplicate}
                className="inline-flex items-center gap-2 rounded-lg border border-[rgba(26,24,19,0.2)] px-4 py-2 text-sm text-nofx-text"
              >
                <Copy className="h-4 w-4" />{' '}
                {text(language, '复制', 'Duplicate')}
              </button>
              <button
                type="button"
                onClick={handleDelete}
                className="ml-auto inline-flex items-center gap-2 rounded-lg px-4 py-2 text-sm text-red-600"
              >
                <Trash2 className="h-4 w-4" />{' '}
                {text(language, '删除', 'Delete')}
              </button>
            </div>

            <section className={sectionClass}>
              <h2 className="text-base font-semibold text-nofx-text">
                {text(language, '基本信息', 'Basics')}
              </h2>
              <div>
                <label className={labelClass}>
                  {text(language, '名称', 'Name')}
                </label>
                <input
                  className={inputClass}
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                />
              </div>
              <div>
                <label className={labelClass}>
                  {text(language, '描述', 'Description')}
                </label>
                <input
                  className={inputClass}
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                />
              </div>
              <div>
                <label className={labelClass}>
                  {text(language, '语言', 'Language')}
                </label>
                <select
                  className={inputClass}
                  value={draft.language ?? 'en'}
                  onChange={(e) =>
                    setDraft({ ...draft, language: e.target.value })
                  }
                >
                  <option value="en">English</option>
                  <option value="zh">中文</option>
                </select>
              </div>
            </section>

            <section className={sectionClass}>
              <h2 className="text-base font-semibold text-nofx-text">
                {text(language, '股票池', 'Universe')}
              </h2>
              <div>
                <label className={labelClass}>
                  {text(language, '股票代码', 'Symbols')}
                </label>
                <textarea
                  className={`${inputClass} min-h-[72px]`}
                  value={symbolsText}
                  onChange={(e) => setSymbolsText(e.target.value)}
                  placeholder={SYMBOL_HINT[market].example}
                />
                <p className="mt-1 text-xs text-nofx-text-muted">
                  {text(
                    language,
                    SYMBOL_HINT[market].zh,
                    SYMBOL_HINT[market].en
                  )}
                </p>
              </div>
              <div>
                <label className={labelClass}>
                  {text(language, '预览市场', 'Preview market')}
                </label>
                <select
                  className={inputClass}
                  value={market}
                  onChange={(e) => setMarket(e.target.value as Market)}
                >
                  <option value="us_equity">
                    {text(language, '美股', 'US equities')}
                  </option>
                  <option value="cn_equity">
                    {text(language, 'A 股', 'A-shares')}
                  </option>
                </select>
              </div>
            </section>

            <section className={sectionClass}>
              <h2 className="text-base font-semibold text-nofx-text">
                {text(language, '数据', 'Data')}
              </h2>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className={labelClass}>
                    {text(language, '日线根数', 'Daily bars')}
                  </label>
                  <input
                    type="number"
                    className={inputClass}
                    min={10}
                    max={120}
                    value={draft.ai_config.indicators.klines.primary_count}
                    onChange={(e) =>
                      updateIndicators({
                        klines: {
                          ...draft.ai_config.indicators.klines,
                          primary_count: clampNumber(
                            Number(e.target.value),
                            10,
                            120
                          ),
                        },
                      })
                    }
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-2 sm:grid-cols-3">
                {INDICATOR_TOGGLES.map((item) => (
                  <label
                    key={item.key}
                    className="flex items-center gap-2 text-sm text-nofx-text"
                  >
                    <input
                      type="checkbox"
                      checked={Boolean(draft.ai_config.indicators[item.key])}
                      onChange={(e) =>
                        updateIndicators({ [item.key]: e.target.checked })
                      }
                    />
                    {text(language, item.zh, item.en)}
                  </label>
                ))}
              </div>
            </section>

            <section className={sectionClass}>
              <h2 className="text-base font-semibold text-nofx-text">
                {text(language, '风控', 'Risk')}
              </h2>
              <div className="grid grid-cols-2 gap-4 sm:grid-cols-3">
                <div>
                  <label className={labelClass}>
                    {text(language, '最多持仓数', 'Max positions')}
                  </label>
                  <input
                    type="number"
                    className={inputClass}
                    min={1}
                    max={8}
                    value={draft.ai_config.risk_control.max_positions}
                    onChange={(e) =>
                      updateRisk({
                        max_positions: clampNumber(
                          Number(e.target.value),
                          1,
                          8
                        ),
                      })
                    }
                  />
                </div>
                <div>
                  <label className={labelClass}>
                    {text(
                      language,
                      '单仓上限（倍数权益）',
                      'Position cap (x equity)'
                    )}
                  </label>
                  <input
                    type="number"
                    step="0.1"
                    className={inputClass}
                    min={0.1}
                    max={10}
                    value={
                      draft.ai_config.risk_control
                        .altcoin_max_position_value_ratio
                    }
                    onChange={(e) =>
                      updateRisk({
                        altcoin_max_position_value_ratio: clampNumber(
                          Number(e.target.value),
                          0.1,
                          10
                        ),
                      })
                    }
                  />
                </div>
                <div>
                  <label className={labelClass}>
                    {text(language, '最低置信度', 'Min confidence')}
                  </label>
                  <input
                    type="number"
                    className={inputClass}
                    min={50}
                    max={100}
                    value={draft.ai_config.risk_control.min_confidence}
                    onChange={(e) =>
                      updateRisk({
                        min_confidence: clampNumber(
                          Number(e.target.value),
                          50,
                          100
                        ),
                      })
                    }
                  />
                </div>
                <div>
                  <label className={labelClass}>
                    {text(language, '最低盈亏比', 'Min reward / risk')}
                  </label>
                  <input
                    type="number"
                    step="0.5"
                    className={inputClass}
                    min={1}
                    max={10}
                    value={draft.ai_config.risk_control.min_risk_reward_ratio}
                    onChange={(e) =>
                      updateRisk({
                        min_risk_reward_ratio: clampNumber(
                          Number(e.target.value),
                          1,
                          10
                        ),
                      })
                    }
                  />
                </div>
                <div className="flex items-start gap-2 pt-6">
                  <input
                    id="profit-protection"
                    type="checkbox"
                    className="mt-1"
                    checked={
                      draft.ai_config.risk_control.profit_protection ?? false
                    }
                    onChange={(e) =>
                      updateRisk({ profit_protection: e.target.checked })
                    }
                  />
                  <label
                    htmlFor="profit-protection"
                    className="text-sm text-nofx-text-muted"
                  >
                    {text(
                      language,
                      '利润保护：盈利超过 5% 后回吐 40% 自动平仓',
                      'Profit protection: close after a 5% gain gives back 40%'
                    )}
                  </label>
                </div>
              </div>
              <p className="text-xs text-nofx-text-muted">
                {text(
                  language,
                  '股票不使用杠杆。A 股不做空，按整手（100 股）交易，遵守 T+1。',
                  'Stocks use no leverage. A-shares cannot be shorted, trade in whole lots of 100, and follow T+1.'
                )}
              </p>
            </section>

            <section className={sectionClass}>
              <h2 className="text-base font-semibold text-nofx-text">
                {text(language, '提示词', 'Prompt')}
              </h2>
              {(
                [
                  ['role_definition', text(language, '角色', 'Role')],
                  [
                    'trading_frequency',
                    text(language, '交易频率', 'Trading frequency'),
                  ],
                  [
                    'entry_standards',
                    text(language, '入场标准', 'Entry standards'),
                  ],
                  [
                    'decision_process',
                    text(language, '决策流程', 'Decision process'),
                  ],
                ] as const
              ).map(([key, label]) => (
                <div key={key}>
                  <label className={labelClass}>{label}</label>
                  <textarea
                    className={`${inputClass} min-h-[72px]`}
                    value={draft.ai_config.prompt_sections?.[key] ?? ''}
                    onChange={(e) => updatePrompt({ [key]: e.target.value })}
                  />
                </div>
              ))}
              <div>
                <label className={labelClass}>
                  {text(language, '自定义指令', 'Custom instructions')}
                </label>
                <textarea
                  className={`${inputClass} min-h-[72px]`}
                  value={draft.ai_config.custom_prompt ?? ''}
                  onChange={(e) => updateAI({ custom_prompt: e.target.value })}
                />
              </div>
            </section>

            {previewText !== null && (
              <section className={sectionClass}>
                <h2 className="text-base font-semibold text-nofx-text">
                  {text(language, '提示词预览', 'Prompt preview')}
                </h2>
                <pre className="max-h-[480px] overflow-auto whitespace-pre-wrap rounded-lg bg-nofx-bg p-4 text-xs text-nofx-text">
                  {previewText}
                </pre>
              </section>
            )}
          </div>
        ) : (
          <div className="text-sm text-nofx-text-muted">
            {text(
              language,
              '还没有策略，点击「新建」开始。',
              'No strategy yet. Click New to start.'
            )}
          </div>
        )}
      </div>
    </DeepVoidBackground>
  )
}

export default StrategyStudioPage
