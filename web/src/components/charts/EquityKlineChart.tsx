import useSWR from 'swr'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import { api } from '../../lib/api'
import { useAuth } from '../../contexts/AuthContext'
import type { EquityMarket } from '../../lib/equityMarket'

interface EquityKlineChartProps {
  market: EquityMarket
  symbol: string
}

// Daily closing prices for one stock. The backend returns 503 until a trader for
// the market has been created, so that case is shown as a message, not an error.
export function EquityKlineChart({ market, symbol }: EquityKlineChartProps) {
  const { user, token } = useAuth()
  const { data, error } = useSWR(
    user && token && symbol ? `equity-klines-${market}-${symbol}` : null,
    () => api.getEquityKlines(market, symbol, 120, true),
    { revalidateOnFocus: false, dedupingInterval: 60000 }
  )

  if (error) {
    return <p className="text-sm text-nofx-text-muted">Daily bars unavailable for {symbol}.</p>
  }
  if (!data) {
    return <p className="text-sm text-nofx-text-muted">Loading {symbol}…</p>
  }

  const points = data.klines.map((k) => ({
    date: new Date(k.openTime).toISOString().slice(0, 10),
    close: k.close,
  }))

  return (
    <div className="h-56 w-full">
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={points}>
          <CartesianGrid strokeDasharray="3 3" opacity={0.2} />
          <XAxis dataKey="date" minTickGap={24} tick={{ fontSize: 11 }} />
          <YAxis domain={['auto', 'auto']} tick={{ fontSize: 11 }} width={60} />
          <Tooltip />
          <Line type="monotone" dataKey="close" dot={false} strokeWidth={2} />
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}
