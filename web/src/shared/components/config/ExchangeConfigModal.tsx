import React, { useEffect, useState } from 'react'
import type { Exchange } from '../../types'
import type { Language } from '../../i18n/translations'
import { getExchangeIcon } from '../../ui/ExchangeIcons'
import { Trash2, Check } from 'lucide-react'
import { getShortName } from '../../../features/traders/utils'

// Paper accounts are the only exchange types. The market is fixed by the type.
export const SUPPORTED_EXCHANGE_TEMPLATES = [
  {
    exchange_type: 'alpaca_paper',
    name: 'Alpaca Paper (US Equities)',
    market: 'us_equity',
    keyLabel: 'API Key ID',
    secretLabel: 'Secret Key',
    needsSecret: true,
  },
  {
    exchange_type: 'cn_paper',
    name: 'A-Share Paper (CN)',
    market: 'cn_equity',
    keyLabel: 'Tushare Token',
    secretLabel: '',
    needsSecret: false,
  },
] as const

type Template = (typeof SUPPORTED_EXCHANGE_TEMPLATES)[number]

interface ExchangeConfigModalProps {
  allExchanges: Exchange[]
  editingExchangeId: string | null
  initialExchangeType?: string | null
  onSave: (
    exchangeId: string | null,
    exchangeType: string,
    accountName: string,
    apiKey: string,
    secretKey?: string
  ) => Promise<void>
  onDelete: (exchangeId: string) => void
  onClose: () => void
  language: Language
}

function templateFor(
  exchangeType: string | null | undefined
): Template | undefined {
  return SUPPORTED_EXCHANGE_TEMPLATES.find(
    (t) => t.exchange_type === exchangeType
  )
}

export function ExchangeConfigModal({
  allExchanges,
  editingExchangeId,
  initialExchangeType,
  onSave,
  onDelete,
  onClose,
  language,
}: ExchangeConfigModalProps) {
  const zh = language === 'zh'
  const existing = editingExchangeId
    ? allExchanges.find((e) => e.id === editingExchangeId)
    : undefined

  const [template, setTemplate] = useState<Template | undefined>(
    templateFor(existing?.exchange_type ?? initialExchangeType)
  )
  const [accountName, setAccountName] = useState(existing?.account_name ?? '')
  const [apiKey, setApiKey] = useState('')
  const [secretKey, setSecretKey] = useState('')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setTemplate(templateFor(existing?.exchange_type ?? initialExchangeType))
  }, [existing?.exchange_type, initialExchangeType])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!template) return
    // Editing keeps the stored secret when the field is left blank.
    if (!existing && !apiKey.trim()) {
      setError(zh ? '请填写密钥' : 'Please enter the key')
      return
    }
    if (!existing && template.needsSecret && !secretKey.trim()) {
      setError(zh ? '请填写 Secret Key' : 'Please enter the secret key')
      return
    }
    setSaving(true)
    setError(null)
    try {
      await onSave(
        existing?.id ?? null,
        template.exchange_type,
        accountName.trim(),
        apiKey.trim(),
        template.needsSecret ? secretKey.trim() : undefined
      )
      onClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err))
    } finally {
      setSaving(false)
    }
  }

  const inputStyle: React.CSSProperties = {
    width: '100%',
    padding: '10px 12px',
    borderRadius: 8,
    border: '1px solid rgba(26,24,19,0.2)',
    background: '#FFFDF7',
    color: '#1A1813',
  }

  return (
    <div
      className="w-full max-w-lg rounded-2xl p-6"
      style={{ background: '#F7F4EC', color: '#1A1813' }}
    >
      <div className="flex items-center justify-between mb-5">
        <h2 className="text-lg font-bold">
          {existing
            ? zh
              ? '编辑纸上账户'
              : 'Edit paper account'
            : zh
              ? '添加纸上账户'
              : 'Add paper account'}
        </h2>
        <button
          type="button"
          onClick={onClose}
          className="text-sm opacity-70 hover:opacity-100"
        >
          {zh ? '关闭' : 'Close'}
        </button>
      </div>

      {!existing && (
        <div className="grid grid-cols-2 gap-3 mb-5">
          {SUPPORTED_EXCHANGE_TEMPLATES.map((t) => (
            <button
              key={t.exchange_type}
              type="button"
              onClick={() => setTemplate(t)}
              className="flex flex-col items-center gap-2 p-4 rounded-xl"
              style={{
                background:
                  template?.exchange_type === t.exchange_type
                    ? 'rgba(224,72,59,0.15)'
                    : '#FFFDF7',
                border:
                  template?.exchange_type === t.exchange_type
                    ? '2px solid #E0483B'
                    : '2px solid rgba(26,24,19,0.14)',
              }}
            >
              <div className="relative">
                {getExchangeIcon(t.exchange_type, { width: 40, height: 40 })}
                {template?.exchange_type === t.exchange_type && (
                  <div
                    className="absolute -top-1 -right-1 w-5 h-5 rounded-full flex items-center justify-center"
                    style={{ background: '#2E8B57' }}
                  >
                    <Check className="w-3 h-3 text-white" />
                  </div>
                )}
              </div>
              <span className="text-sm font-semibold">
                {getShortName(t.name)}
              </span>
            </button>
          ))}
        </div>
      )}

      {template && (
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm mb-1">
              {zh ? '账户名称' : 'Account name'}
            </label>
            <input
              value={accountName}
              onChange={(e) => setAccountName(e.target.value)}
              placeholder={zh ? '例如：主账户' : 'e.g. main'}
              style={inputStyle}
            />
          </div>
          <div>
            <label className="block text-sm mb-1">{template.keyLabel}</label>
            <input
              type="password"
              autoComplete="off"
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
              placeholder={
                existing?.has_api_key
                  ? zh
                    ? '已保存，留空则不修改'
                    : 'Saved — leave blank to keep'
                  : ''
              }
              style={inputStyle}
            />
          </div>
          {template.needsSecret && (
            <div>
              <label className="block text-sm mb-1">
                {template.secretLabel}
              </label>
              <input
                type="password"
                autoComplete="off"
                value={secretKey}
                onChange={(e) => setSecretKey(e.target.value)}
                placeholder={
                  existing?.has_secret_key
                    ? zh
                      ? '已保存，留空则不修改'
                      : 'Saved — leave blank to keep'
                    : ''
                }
                style={inputStyle}
              />
            </div>
          )}
          <p className="text-xs opacity-70">
            {zh
              ? '纸上账户只模拟交易，不会下真实订单。A 股行情来自 Tushare。'
              : 'Paper accounts simulate trades and never place real orders. A-share prices come from Tushare.'}
          </p>
          {error && <p className="text-sm text-red-600">{error}</p>}
          <div className="flex items-center justify-between pt-2">
            {existing ? (
              <button
                type="button"
                onClick={() => onDelete(existing.id)}
                className="flex items-center gap-2 text-sm text-red-600"
              >
                <Trash2 className="w-4 h-4" /> {zh ? '删除' : 'Delete'}
              </button>
            ) : (
              <span />
            )}
            <button
              type="submit"
              disabled={saving}
              className="px-5 py-2 rounded-lg font-semibold disabled:opacity-50"
              style={{ background: '#E0483B', color: '#fff' }}
            >
              {saving ? (zh ? '保存中…' : 'Saving…') : zh ? '保存' : 'Save'}
            </button>
          </div>
        </form>
      )}
    </div>
  )
}

export default ExchangeConfigModal
