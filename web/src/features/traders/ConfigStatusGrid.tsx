import { Brain, Landmark } from 'lucide-react'
import type {
  AIModel,
  Exchange,
  ExchangeAccountState,
} from '../../shared/types'
import type { Language } from '../../shared/i18n/translations'
import { t } from '../../shared/i18n/translations'
import { getModelIcon } from '../../shared/ui/ModelIcons'
import { getExchangeIcon } from '../../shared/ui/ExchangeIcons'
import { getShortName, AI_PROVIDER_CONFIG } from '../../shared/domain/models'

interface UsageInfo {
  runningCount: number
  totalCount: number
}

interface ConfigStatusGridProps {
  configuredModels: AIModel[]
  configuredExchanges: Exchange[]
  exchangeAccountStates?: Record<string, ExchangeAccountState>
  isExchangeAccountStatesLoading?: boolean
  language: Language
  getModelUsageInfo: (modelId: string) => UsageInfo
  getExchangeUsageInfo: (exchangeId: string) => UsageInfo
  onModelClick: (modelId: string) => void
  onExchangeClick: (exchangeId: string) => void
}

export function ConfigStatusGrid({
  configuredModels,
  configuredExchanges,
  exchangeAccountStates,
  isExchangeAccountStatesLoading,
  language,
  getModelUsageInfo,
  getExchangeUsageInfo,
  onModelClick,
  onExchangeClick,
}: ConfigStatusGridProps) {
  const getExchangeStateMeta = (state: ExchangeAccountState | undefined) => {
    if (!state) {
      return {
        label: language === 'zh' ? 'NOT CHECKED' : 'NOT CHECKED',
        className:
          'text-crux-text-muted border-[rgba(26,24,19,0.14)] bg-crux-bg-deeper',
      }
    }

    switch (state.status) {
      case 'ok':
        return {
          label: state.display_balance || '0',
          className:
            'text-crux-success border-crux-success/20 bg-crux-success/10',
        }
      case 'disabled':
        return {
          label: language === 'zh' ? 'DISABLED' : 'DISABLED',
          className:
            'text-crux-text-muted border-[rgba(26,24,19,0.14)] bg-crux-bg-deeper',
        }
      case 'missing_credentials':
        return {
          label: language === 'zh' ? 'INCOMPLETE' : 'INCOMPLETE',
          className: 'text-crux-gold border-crux-gold/20 bg-crux-gold/10',
        }
      case 'invalid_credentials':
        return {
          label: language === 'zh' ? 'INVALID KEYS' : 'INVALID KEYS',
          className: 'text-crux-danger border-crux-danger/20 bg-crux-danger/10',
        }
      case 'permission_denied':
        return {
          label: language === 'zh' ? 'NO PERMISSION' : 'NO PERMISSION',
          className: 'text-crux-gold border-crux-gold/20 bg-crux-gold/10',
        }
      default:
        return {
          label: language === 'zh' ? 'UNAVAILABLE' : 'UNAVAILABLE',
          className:
            'text-crux-text-muted border-[rgba(26,24,19,0.14)] bg-crux-bg-deeper',
        }
    }
  }

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      {/* AI Models Card */}
      <div className="bg-crux-bg-lighter rounded-lg border border-[rgba(26,24,19,0.14)] overflow-hidden">
        <div className="px-4 py-3 border-b border-[rgba(26,24,19,0.14)] bg-crux-bg-deeper flex items-center gap-2">
          <Brain className="w-4 h-4 text-crux-gold" />
          <h3 className="text-sm font-mono tracking-widest text-crux-text uppercase">
            {t('aiModels', language)}
          </h3>
        </div>

        <div className="p-4 space-y-3">
          {configuredModels.map((model) => {
            const usageInfo = getModelUsageInfo(model.id)
            return (
              <div
                key={model.id}
                role="button"
                tabIndex={0}
                className="group relative flex cursor-pointer items-center justify-between rounded-md border border-transparent bg-crux-bg-deeper p-3 transition-all hover:border-crux-gold/20 hover:bg-crux-bg"
                onClick={() => onModelClick(model.id)}
                onKeyDown={(event) => {
                  if (event.key === 'Enter' || event.key === ' ') {
                    event.preventDefault()
                    onModelClick(model.id)
                  }
                }}
              >
                <div className="flex items-center gap-4">
                  <div className="relative">
                    <div className="w-10 h-10 rounded-full flex items-center justify-center bg-crux-bg border border-[rgba(26,24,19,0.14)] relative z-10">
                      {getModelIcon(model.provider || model.id, {
                        width: 20,
                        height: 20,
                      }) || (
                        <span className="text-xs font-bold text-crux-accent">
                          {getShortName(model.name)[0]}
                        </span>
                      )}
                    </div>
                  </div>

                  <div className="min-w-0">
                    <div className="font-mono text-sm text-crux-text group-hover:text-crux-gold transition-colors">
                      {getShortName(model.name)}
                    </div>
                    <div className="text-[10px] text-crux-text-muted font-mono flex items-center gap-2">
                      {model.customModelName ||
                        AI_PROVIDER_CONFIG[model.provider]?.defaultModel ||
                        ''}
                    </div>
                  </div>
                </div>

                <div className="text-right">
                  {usageInfo.totalCount > 0 ? (
                    <span
                      className={`text-[10px] font-mono px-2 py-1 rounded border ${
                        usageInfo.runningCount > 0
                          ? 'bg-crux-success/10 border-crux-success/30 text-crux-success'
                          : 'bg-crux-gold/10 border-crux-gold/30 text-crux-gold'
                      }`}
                    >
                      {usageInfo.runningCount}/{usageInfo.totalCount} ACTIVE
                    </span>
                  ) : (
                    <span className="text-[10px] font-mono text-crux-text-muted uppercase tracking-wider">
                      {language === 'zh' ? 'STANDBY' : 'STANDBY'}
                    </span>
                  )}
                </div>
              </div>
            )
          })}

          {configuredModels.length === 0 && (
            <div className="text-center py-10 border border-dashed border-[rgba(26,24,19,0.14)] rounded-lg bg-crux-bg-deeper">
              <Brain className="w-8 h-8 mx-auto mb-3 text-crux-text-muted" />
              <div className="text-xs font-mono text-crux-text-muted uppercase tracking-widest">
                {t('noModelsConfigured', language)}
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Exchanges Card */}
      <div className="bg-crux-bg-lighter rounded-lg border border-[rgba(26,24,19,0.14)] overflow-hidden">
        <div className="px-4 py-3 border-b border-[rgba(26,24,19,0.14)] bg-crux-bg-deeper flex items-center gap-2">
          <Landmark className="w-4 h-4 text-crux-gold" />
          <h3 className="text-sm font-mono tracking-widest text-crux-text uppercase">
            {t('exchanges', language)}
          </h3>
        </div>

        <div className="p-4 space-y-3">
          {configuredExchanges.map((exchange) => {
            const usageInfo = getExchangeUsageInfo(exchange.id)
            const state = exchangeAccountStates?.[exchange.id]
            const stateMeta = getExchangeStateMeta(state)
            return (
              <div
                key={exchange.id}
                role="button"
                tabIndex={0}
                className="group relative flex cursor-pointer items-center justify-between rounded-md border border-transparent bg-crux-bg-deeper p-3 transition-all hover:border-crux-gold/20 hover:bg-crux-bg"
                onClick={() => onExchangeClick(exchange.id)}
                onKeyDown={(event) => {
                  if (event.key === 'Enter' || event.key === ' ') {
                    event.preventDefault()
                    onExchangeClick(exchange.id)
                  }
                }}
              >
                <div className="flex items-center gap-4 min-w-0">
                  <div className="relative">
                    <div className="w-10 h-10 rounded-full flex items-center justify-center bg-crux-bg border border-[rgba(26,24,19,0.14)] relative z-10">
                      {getExchangeIcon(exchange.exchange_type || exchange.id, {
                        width: 20,
                        height: 20,
                      })}
                    </div>
                  </div>

                  <div className="min-w-0">
                    <div className="font-mono text-sm text-crux-text group-hover:text-crux-gold transition-colors truncate">
                      {exchange.exchange_type?.toUpperCase() ||
                        getShortName(exchange.name)}
                      <span className="text-[10px] text-crux-text-muted ml-2 border border-[rgba(26,24,19,0.14)] px-1 rounded">
                        {exchange.account_name || 'DEFAULT'}
                      </span>
                    </div>
                    <div className="text-[10px] text-crux-text-muted font-mono flex items-center gap-2">
                      {exchange.type?.toUpperCase() || 'CEX'}
                    </div>
                    <div className="mt-1 flex flex-wrap items-center gap-2 text-[10px] font-mono">
                      <span
                        className={`rounded border px-1.5 py-0.5 ${stateMeta.className}`}
                      >
                        {isExchangeAccountStatesLoading && !state
                          ? language === 'zh'
                            ? 'CHECKING...'
                            : 'CHECKING...'
                          : stateMeta.label}
                      </span>
                      {state?.status !== 'ok' && state?.error_message ? (
                        <span className="text-crux-text-muted truncate max-w-[220px]">
                          {state.error_message}
                        </span>
                      ) : null}
                    </div>
                  </div>
                </div>

                <div className="flex flex-col items-end gap-1">
                  {usageInfo.totalCount > 0 ? (
                    <span
                      className={`text-[10px] font-mono px-2 py-1 rounded border ${
                        usageInfo.runningCount > 0
                          ? 'bg-crux-success/10 border-crux-success/30 text-crux-success'
                          : 'bg-crux-gold/10 border-crux-gold/30 text-crux-gold'
                      }`}
                    >
                      {usageInfo.runningCount}/{usageInfo.totalCount} ACTIVE
                    </span>
                  ) : (
                    <span className="text-[10px] font-mono text-crux-text-muted uppercase tracking-wider">
                      {language === 'zh' ? 'STANDBY' : 'STANDBY'}
                    </span>
                  )}
                </div>
              </div>
            )
          })}
          {configuredExchanges.length === 0 && (
            <div className="text-center py-10 border border-dashed border-[rgba(26,24,19,0.14)] rounded-lg bg-crux-bg-deeper">
              <Landmark className="w-8 h-8 mx-auto mb-3 text-crux-text-muted" />
              <div className="text-xs font-mono text-crux-text-muted uppercase tracking-widest">
                {t('noExchangesConfigured', language)}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
