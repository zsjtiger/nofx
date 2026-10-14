import { render, screen } from '@testing-library/react'
import { beforeEach, expect, it } from 'vitest'
import { ExecutionLog } from './ExecutionLog'
import { LanguageProvider } from '../../shared/contexts/LanguageContext'
import type { DecisionRecord } from '../../shared/types'

beforeEach(() => localStorage.setItem('language', 'en'))

it('marks legacy success=true cycles as faults when their orders failed and shows the date', () => {
  const record = {
    cycle_number: 2581,
    timestamp: '2026-08-15T12:00:00Z',
    success: true,
    decisions: [
      {
        action: 'open_long',
        symbol: 'SOL',
        success: false,
        timestamp: '2026-08-15T12:00:00Z',
        error: 'denied',
      },
    ],
  } as DecisionRecord
  render(
    <LanguageProvider>
      <ExecutionLog decisions={[record]} />
    </LanguageProvider>
  )
  expect(
    screen.getByText('FAULT · 1/1 trades failed/blocked')
  ).toBeInTheDocument()
  expect(screen.getAllByText(/2026-08-15/).length).toBeGreaterThan(0)
})

const protocolError =
  'paid request outcome unknown; refusing automatic replay: malformed HTTP response "\\x00\\x00\\x12\\x04"'

it('keeps historical decision faults visible but formats their protocol error', () => {
  const record = {
    cycle_number: 2664,
    timestamp: '2026-10-14T03:41:23Z',
    success: false,
    decisions: [],
    error_message: protocolError,
    execution_log: [protocolError],
  } as DecisionRecord
  render(
    <LanguageProvider>
      <ExecutionLog decisions={[record]} />
    </LanguageProvider>
  )
  expect(screen.getByText('FAULT')).toBeInTheDocument()
  expect(screen.getAllByText(/Payment result is unconfirmed/)).toHaveLength(2)
  expect(document.body.textContent).not.toContain('\\x00')
  expect(record.error_message).toBe(protocolError)
})
