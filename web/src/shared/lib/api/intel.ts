import type { IntelItem, IntelRunRequest, IntelRunResult } from '../../types'
import { API_BASE, httpClient } from './helpers'

// Keeps the HTTP status so the page can tell "not configured" (503) from other failures.
function failure(
  message: string | undefined,
  statusCode: number | undefined
): Error {
  return Object.assign(new Error(message || 'Request failed'), { statusCode })
}

export const intelApi = {
  async getIntel(symbol: string, limit = 20): Promise<IntelItem[]> {
    const result = await httpClient.get<IntelItem[]>(`${API_BASE}/intel`, {
      symbol,
      limit,
    })
    if (!result.success) throw failure(result.message, result.statusCode)
    return Array.isArray(result.data) ? result.data : []
  },

  async runIntel(request: IntelRunRequest): Promise<IntelRunResult[]> {
    const result = await httpClient.post<{ results: IntelRunResult[] }>(
      `${API_BASE}/intel/run`,
      request
    )
    if (!result.success) throw failure(result.message, result.statusCode)
    return result.data?.results ?? []
  },
}
