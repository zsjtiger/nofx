export const MAX_INTEL_URLS = 5

export interface ParsedIntelUrls {
  urls: string[]
  /** 'empty' when no http(s) URL was entered, 'tooMany' above MAX_INTEL_URLS. */
  error?: 'empty' | 'tooMany'
}

/**
 * Reads the URL box: one URL per line. Blank lines and non-http(s) lines are
 * ignored, duplicates are removed, and the result is checked against the per-run cap.
 */
export function parseIntelUrls(text: string): ParsedIntelUrls {
  const seen = new Set<string>()
  const urls: string[] = []
  for (const raw of text.split(/\r?\n/)) {
    const url = raw.trim()
    if (!/^https?:\/\//i.test(url) || seen.has(url)) continue
    seen.add(url)
    urls.push(url)
  }
  if (urls.length === 0) return { urls, error: 'empty' }
  if (urls.length > MAX_INTEL_URLS) return { urls, error: 'tooMany' }
  return { urls }
}
