import { act, render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it } from 'vitest'
import { LanguageProvider, useLanguage } from './LanguageContext'

function LanguageProbe() {
  const { language, setLanguage } = useLanguage()
  return (
    <div>
      <span data-testid="lang">{language}</span>
      <button type="button" onClick={() => setLanguage('en')}>
        en
      </button>
      <button type="button" onClick={() => setLanguage('zh')}>
        zh
      </button>
      <button type="button" onClick={() => setLanguage('id' as never)}>
        id
      </button>
    </div>
  )
}

function renderProbe() {
  return render(
    <LanguageProvider>
      <LanguageProbe />
    </LanguageProvider>
  )
}

describe('language selection', () => {
  beforeEach(() => localStorage.clear())

  it('defaults to Chinese when nothing is stored', () => {
    renderProbe()
    expect(screen.getByTestId('lang').textContent).toBe('zh')
  })

  it('uses a stored English choice', () => {
    localStorage.setItem('language', 'en')
    renderProbe()
    expect(screen.getByTestId('lang').textContent).toBe('en')
  })

  it('falls back to Chinese for unsupported stored values', () => {
    localStorage.setItem('language', 'xx')
    renderProbe()
    expect(screen.getByTestId('lang').textContent).toBe('zh')
  })

  it('switches language and persists the choice', () => {
    renderProbe()
    act(() => screen.getByRole('button', { name: 'en' }).click())
    expect(screen.getByTestId('lang').textContent).toBe('en')
    expect(localStorage.getItem('language')).toBe('en')

    act(() => screen.getByRole('button', { name: 'zh' }).click())
    expect(screen.getByTestId('lang').textContent).toBe('zh')
    expect(localStorage.getItem('language')).toBe('zh')
  })

  it('ignores switches to languages that are not offered', () => {
    renderProbe()
    act(() => screen.getByRole('button', { name: 'id' }).click())
    expect(screen.getByTestId('lang').textContent).toBe('zh')
    expect(localStorage.getItem('language')).toBeNull()
  })
})
