// Infers the stock market from a symbol as the backend stores it:
// A-shares are six digits with an optional .SH/.SZ/.BJ suffix, US tickers are letters.
export type EquityMarket = 'us_equity' | 'cn_equity'

export function inferEquityMarket(symbol: string): EquityMarket {
  const s = symbol.trim().toUpperCase()
  return /^\d{6}(\.(SH|SZ|BJ))?$/.test(s) ? 'cn_equity' : 'us_equity'
}

// Account currency for a broker type, e.g. alpaca_paper -> USD, cn_paper -> CNY.
export function currencyForExchangeType(exchangeType?: string): string {
  switch (exchangeType) {
    case 'alpaca_paper':
      return 'USD'
    case 'cn_paper':
      return 'CNY'
    default:
      return ''
  }
}
