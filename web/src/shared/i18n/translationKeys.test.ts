import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join, resolve } from 'node:path'
import { describe, expect, it } from 'vitest'
import { translations } from './translations'

// Every key passed as a string literal to t() must exist in English. Keys built
// at runtime (loginRequired.*, comparisonChart.*) are checked by their own tests.
function walk(dir: string): string[] {
  return readdirSync(dir).flatMap((name) => {
    const full = join(dir, name)
    if (statSync(full).isDirectory()) return walk(full)
    return /\.(ts|tsx)$/.test(name) &&
      !name.endsWith('.test.ts') &&
      !name.endsWith('.test.tsx') &&
      name !== 'translations.ts'
      ? [full]
      : []
  })
}

function hasKey(root: unknown, dotted: string): boolean {
  let value: any = root
  for (const part of dotted.split('.')) {
    if (value == null || typeof value !== 'object' || !(part in value))
      return false
    value = value[part]
  }
  return typeof value === 'string'
}

describe('translation keys', () => {
  it('every literal t() key used in src exists in English', () => {
    const srcDir = resolve(__dirname, '../..')
    const pattern = /\bt\(\s*'([A-Za-z0-9_.]+)'/g
    const missing = new Set<string>()
    for (const file of walk(srcDir)) {
      const text = readFileSync(file, 'utf8')
      for (const match of text.matchAll(pattern)) {
        const key = match[1]
        if (!hasKey(translations.en, key))
          missing.add(`${key} (${file.replace(srcDir, 'src')})`)
      }
    }
    expect([...missing]).toEqual([])
  })

  it('the four keys that were referenced without a translation now resolve', () => {
    for (const key of [
      'leadPnL',
      'loadingChartData',
      'loadingError',
      'loggingIn',
    ]) {
      expect(hasKey(translations.en, key)).toBe(true)
      expect(hasKey(translations.zh, key)).toBe(true)
      expect(hasKey(translations.id, key)).toBe(true)
    }
  })
})
