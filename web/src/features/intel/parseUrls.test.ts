import { describe, expect, it } from 'vitest'
import { MAX_INTEL_URLS, parseIntelUrls } from './parseUrls'

describe('parseIntelUrls', () => {
  it('keeps http(s) lines, trims them and removes duplicates', () => {
    const text =
      ' https://a.example/1 \n\nhttp://b.example\nhttps://a.example/1\nnot a url\n'
    expect(parseIntelUrls(text)).toEqual({
      urls: ['https://a.example/1', 'http://b.example'],
    })
  })

  it('reports empty input when no valid URL is present', () => {
    expect(parseIntelUrls('  \nftp://x\n')).toEqual({
      urls: [],
      error: 'empty',
    })
  })

  it('reports too many URLs above the per-run cap', () => {
    const text = Array.from(
      { length: MAX_INTEL_URLS + 1 },
      (_, i) => `https://e.example/${i}`
    ).join('\n')
    expect(parseIntelUrls(text).error).toBe('tooMany')
  })
})
