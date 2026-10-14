import {
  Bot,
  Users,
  BarChart3,
  Trash2,
  Pencil,
  Eye,
  EyeOff,
} from 'lucide-react'
import type { TraderInfo, Exchange } from '../../shared/types'
import type { Language } from '../../shared/i18n/translations'
import { t } from '../../shared/i18n/translations'
import { PunkAvatar, getTraderAvatar } from '../../shared/ui/PunkAvatar'
import { getModelDisplayName, getExchangeDisplayName } from './model-constants'

interface TradersListProps {
  traders: TraderInfo[] | undefined
  isLoading: boolean
  allExchanges: Exchange[]
  configuredModelsCount: number
  configuredExchangesCount: number
  language: Language
  onTraderSelect?: (traderId: string) => void
  onNavigate: (path: string) => void
  onEditTrader: (traderId: string) => void
  onToggleTrader: (traderId: string, running: boolean) => void
  onToggleCompetition: (
    traderId: string,
    currentShowInCompetition: boolean
  ) => void
  onDeleteTrader: (traderId: string) => void
}

export function TradersList({
  traders,
  isLoading,
  allExchanges,
  configuredModelsCount,
  configuredExchangesCount,
  language,
  onTraderSelect,
  onNavigate,
  onEditTrader,
  onToggleTrader,
  onToggleCompetition,
  onDeleteTrader,
}: TradersListProps) {
  return (
    <div className="binance-card p-4 md:p-6">
      <div className="flex items-center justify-between mb-4 md:mb-5">
        <h2
          className="text-lg md:text-xl font-bold flex items-center gap-2"
          style={{ color: '#1A1813' }}
        >
          <Users
            className="w-5 h-5 md:w-6 md:h-6"
            style={{ color: '#E0483B' }}
          />
          {t('currentTraders', language)}
        </h2>
      </div>

      {isLoading ? (
        <TradersLoadingSkeleton />
      ) : traders && traders.length > 0 ? (
        <div className="space-y-3 md:space-y-4">
          {traders.map((trader) => (
            <TraderRow
              key={trader.trader_id}
              trader={trader}
              allExchanges={allExchanges}
              language={language}
              onTraderSelect={onTraderSelect}
              onNavigate={onNavigate}
              onEditTrader={onEditTrader}
              onToggleTrader={onToggleTrader}
              onToggleCompetition={onToggleCompetition}
              onDeleteTrader={onDeleteTrader}
            />
          ))}
        </div>
      ) : (
        <TradersEmptyState
          configuredModelsCount={configuredModelsCount}
          configuredExchangesCount={configuredExchangesCount}
          language={language}
        />
      )}
    </div>
  )
}

function TradersLoadingSkeleton() {
  return (
    <div className="space-y-3 md:space-y-4">
      {[1, 2, 3].map((i) => (
        <div
          key={i}
          className="flex flex-col md:flex-row md:items-center justify-between p-3 md:p-4 rounded gap-3 md:gap-4 animate-pulse"
          style={{
            background: '#F7F4EC',
            border: '1px solid rgba(26,24,19,0.14)',
          }}
        >
          <div className="flex items-center gap-3 md:gap-4">
            <div className="w-10 h-10 md:w-12 md:h-12 rounded-full skeleton"></div>
            <div className="min-w-0 space-y-2">
              <div className="skeleton h-5 w-32"></div>
              <div className="skeleton h-3 w-24"></div>
            </div>
          </div>
          <div className="flex items-center gap-3 md:gap-4">
            <div className="skeleton h-6 w-16"></div>
            <div className="skeleton h-6 w-16"></div>
            <div className="skeleton h-8 w-20"></div>
          </div>
        </div>
      ))}
    </div>
  )
}

function TradersEmptyState({
  configuredModelsCount,
  configuredExchangesCount,
  language,
}: {
  configuredModelsCount: number
  configuredExchangesCount: number
  language: Language
}) {
  return (
    <div className="text-center py-12 md:py-16" style={{ color: '#8A8478' }}>
      <Bot className="w-16 h-16 md:w-24 md:h-24 mx-auto mb-3 md:mb-4 opacity-50" />
      <div className="text-base md:text-lg font-semibold mb-2">
        {t('noTraders', language)}
      </div>
      <div className="text-xs md:text-sm mb-3 md:mb-4">
        {t('createFirstTrader', language)}
      </div>
      {(configuredModelsCount === 0 || configuredExchangesCount === 0) && (
        <div className="text-xs md:text-sm text-nofx-gold">
          {configuredModelsCount === 0 && configuredExchangesCount === 0
            ? t('configureModelsAndExchangesFirst', language)
            : configuredModelsCount === 0
              ? t('configureModelsFirst', language)
              : t('configureExchangesFirst', language)}
        </div>
      )}
    </div>
  )
}

function TraderRow({
  trader,
  allExchanges,
  language,
  onTraderSelect,
  onNavigate,
  onEditTrader,
  onToggleTrader,
  onToggleCompetition,
  onDeleteTrader,
}: {
  trader: TraderInfo
  allExchanges: Exchange[]
  language: Language
  onTraderSelect?: (traderId: string) => void
  onNavigate: (path: string) => void
  onEditTrader: (traderId: string) => void
  onToggleTrader: (traderId: string, running: boolean) => void
  onToggleCompetition: (
    traderId: string,
    currentShowInCompetition: boolean
  ) => void
  onDeleteTrader: (traderId: string) => void
}) {
  return (
    <div
      className="flex flex-col md:flex-row md:items-center justify-between p-3 md:p-4 rounded transition-all hover:translate-y-[-1px] gap-3 md:gap-4"
      style={{ background: '#F7F4EC', border: '1px solid rgba(26,24,19,0.14)' }}
    >
      <div className="flex items-center gap-3 md:gap-4">
        <div className="flex-shrink-0">
          <PunkAvatar
            seed={getTraderAvatar(trader.trader_id, trader.trader_name)}
            size={48}
            className="rounded-lg hidden md:block"
          />
          <PunkAvatar
            seed={getTraderAvatar(trader.trader_id, trader.trader_name)}
            size={40}
            className="rounded-lg md:hidden"
          />
        </div>
        <div className="min-w-0">
          <div
            className="font-bold text-base md:text-lg truncate"
            style={{ color: '#1A1813' }}
          >
            {trader.trader_name}
          </div>
          <div
            className="text-xs md:text-sm truncate"
            style={{
              color: trader.ai_model.includes('deepseek')
                ? '#E0483B'
                : '#E0483B',
            }}
          >
            {getModelDisplayName(
              trader.ai_model.split('_').pop() || trader.ai_model
            )}{' '}
            Model • {getExchangeDisplayName(trader.exchange_id, allExchanges)}
          </div>
        </div>
      </div>

      <div className="flex items-center gap-3 md:gap-4 flex-wrap md:flex-nowrap">
        {/* Status */}
        <div className="text-center">
          <div
            className={`px-2 md:px-3 py-1 rounded text-xs font-bold ${
              trader.is_running
                ? 'bg-nofx-success/10 text-nofx-success'
                : 'bg-nofx-danger/10 text-nofx-danger'
            }`}
            style={
              trader.is_running
                ? {
                    background: 'rgba(46, 139, 87, 0.1)',
                    color: '#2E8B57',
                  }
                : {
                    background: 'rgba(214, 67, 58, 0.1)',
                    color: '#D6433A',
                  }
            }
          >
            {trader.is_running
              ? t('running', language)
              : t('stopped', language)}
          </div>
        </div>

        {/* Actions */}
        <div className="flex gap-1.5 md:gap-2 flex-nowrap overflow-x-auto items-center">
          <button
            onClick={() => {
              if (onTraderSelect) {
                onTraderSelect(trader.trader_id)
              } else {
                const slug = `${trader.trader_name}-${trader.trader_id.slice(0, 4)}`
                onNavigate(`/dashboard?trader=${encodeURIComponent(slug)}`)
              }
            }}
            className="px-2 md:px-3 py-1.5 md:py-2 rounded text-xs md:text-sm font-semibold transition-all hover:scale-105 flex items-center gap-1 whitespace-nowrap"
            style={{
              background: 'rgba(224, 72, 59, 0.1)',
              color: '#E0483B',
            }}
          >
            <BarChart3 className="w-3 h-3 md:w-4 md:h-4" />
            {t('view', language)}
          </button>

          <button
            onClick={() => onEditTrader(trader.trader_id)}
            disabled={trader.is_running}
            className="px-2 md:px-3 py-1.5 md:py-2 rounded text-xs md:text-sm font-semibold transition-all hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap flex items-center gap-1"
            style={{
              background: trader.is_running
                ? 'rgba(138, 132, 120, 0.1)'
                : 'rgba(224, 72, 59, 0.1)',
              color: trader.is_running ? '#8A8478' : '#E0483B',
            }}
          >
            <Pencil className="w-3 h-3 md:w-4 md:h-4" />
            {t('edit', language)}
          </button>

          <button
            onClick={() =>
              onToggleTrader(trader.trader_id, trader.is_running || false)
            }
            className="px-2 md:px-3 py-1.5 md:py-2 rounded text-xs md:text-sm font-semibold transition-all hover:scale-105 whitespace-nowrap"
            style={
              trader.is_running
                ? {
                    background: 'rgba(214, 67, 58, 0.1)',
                    color: '#D6433A',
                  }
                : {
                    background: 'rgba(46, 139, 87, 0.1)',
                    color: '#2E8B57',
                  }
            }
          >
            {trader.is_running ? t('stop', language) : t('start', language)}
          </button>

          <button
            onClick={() =>
              onToggleCompetition(
                trader.trader_id,
                trader.show_in_competition ?? true
              )
            }
            className="px-2 md:px-3 py-1.5 md:py-2 rounded text-xs md:text-sm font-semibold transition-all hover:scale-105 whitespace-nowrap flex items-center gap-1"
            style={
              trader.show_in_competition !== false
                ? {
                    background: 'rgba(46, 139, 87, 0.1)',
                    color: '#2E8B57',
                  }
                : {
                    background: 'rgba(138, 132, 120, 0.1)',
                    color: '#8A8478',
                  }
            }
            title={
              trader.show_in_competition !== false
                ? 'Shown in arena'
                : 'Hidden from arena'
            }
          >
            {trader.show_in_competition !== false ? (
              <Eye className="w-3 h-3 md:w-4 md:h-4" />
            ) : (
              <EyeOff className="w-3 h-3 md:w-4 md:h-4" />
            )}
          </button>

          <button
            onClick={() => onDeleteTrader(trader.trader_id)}
            className="px-2 md:px-3 py-1.5 md:py-2 rounded text-xs md:text-sm font-semibold transition-all hover:scale-105"
            style={{
              background: 'rgba(214, 67, 58, 0.1)',
              color: '#D6433A',
            }}
          >
            <Trash2 className="w-3 h-3 md:w-4 md:h-4" />
          </button>
        </div>
      </div>
    </div>
  )
}
