import { describe, expect, it } from 'vitest'
import { currencyForExchangeType, inferEquityMarket } from './equityMarket'

describe('inferEquityMarket', () => {
  it('treats A-share codes as cn_equity', () => {
    expect(inferEquityMarket('600519.SH')).toBe('cn_equity')
    expect(inferEquityMarket('000001')).toBe('cn_equity')
    expect(inferEquityMarket(' 300750.sz ')).toBe('cn_equity')
  })

  it('treats tickers as us_equity', () => {
    expect(inferEquityMarket('AAPL')).toBe('us_equity')
    expect(inferEquityMarket('brk.b')).toBe('us_equity')
  })
})

describe('currencyForExchangeType', () => {
  it('maps each broker to its account currency', () => {
    expect(currencyForExchangeType('alpaca_paper')).toBe('USD')
    expect(currencyForExchangeType('cn_paper')).toBe('CNY')
    expect(currencyForExchangeType('unknown')).toBe('')
    expect(currencyForExchangeType(undefined)).toBe('')
  })
})
