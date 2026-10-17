import { BookOpen, Monitor, Shield, TrendingUp, Wrench } from 'lucide-react'
import type { LucideIcon } from 'lucide-react'

/**
 * FAQ content model. Answers are composed from typed blocks so the renderer
 * stays generic — no per-question JSX special cases. Inline `code` spans are
 * written with backticks and parsed by the renderer.
 */
export type FAQBlock =
  | { type: 'p'; text: string }
  | { type: 'list'; items: string[] }
  | { type: 'steps'; items: string[] }
  | { type: 'note'; text: string }
  | { type: 'links'; links: { label: string; href: string }[] }

export interface FAQItem {
  id: string
  question: string
  blocks: FAQBlock[]
}

export interface FAQCategory {
  id: string
  title: string
  icon: LucideIcon
  items: FAQItem[]
}

/** Plain text of an item, used by the search filter. */
export function faqItemSearchText(item: FAQItem): string {
  const parts: string[] = [item.question]
  for (const block of item.blocks) {
    if (block.type === 'p' || block.type === 'note') parts.push(block.text)
    else if (block.type === 'list' || block.type === 'steps')
      parts.push(block.items.join(' '))
    else if (block.type === 'links')
      parts.push(block.links.map((l) => l.label).join(' '))
  }
  return parts.join(' ').toLowerCase()
}

export const faqCategories: FAQCategory[] = [
  // ───────────────────────── Getting started ─────────────────────────
  {
    id: 'getting-started',
    title: 'Getting Started',
    icon: BookOpen,
    items: [
      {
        id: 'what-is-crux',
        question: 'What is Crux?',
        blocks: [
          {
            type: 'p',
            text: 'Crux is a self-hosted platform for AI-driven trading of US stocks and China A-shares. An AI model reads market data each cycle and proposes buy, sell or hold decisions; the platform checks each decision against market rules before it can reach a broker.',
          },
        ],
      },
      {
        id: 'which-markets',
        question: 'Which markets are supported?',
        blocks: [
          {
            type: 'list',
            items: [
              'US equities (`us_equity`): trades through an Alpaca paper account, priced in USD on the US Eastern clock.',
              'China A-shares (`cn_equity`): trades in a built-in paper account, priced in CNY on the Beijing clock. Market data comes from Tushare Pro.',
            ],
          },
          {
            type: 'note',
            text: 'Both markets currently run in paper trading only. No live broker is connected for A-shares yet.',
          },
        ],
      },
      {
        id: 'what-do-i-need',
        question: 'What do I need to get started?',
        blocks: [
          {
            type: 'list',
            items: [
              'An AI model with an OpenAI-compatible API key (for example OpenAI or DeepSeek, or a self-hosted gateway).',
              'For US stocks: an Alpaca paper trading API key pair.',
              'For A-shares: a Tushare Pro token.',
            ],
          },
        ],
      },
      {
        id: 'ai-models',
        question: 'Which AI models can I use?',
        blocks: [
          {
            type: 'p',
            text: 'Any model that exposes an OpenAI-compatible chat completions endpoint. You enter the base URL, model name and API key in the AI model settings. Keys are stored encrypted on your server.',
          },
        ],
      },
    ],
  },

  // ───────────────────────── Trading ─────────────────────────
  {
    id: 'trading',
    title: 'Trading',
    icon: TrendingUp,
    items: [
      {
        id: 'trading-rules',
        question: 'What trading rules does the platform enforce?',
        blocks: [
          {
            type: 'list',
            items: [
              'A-shares trade in whole lots of 100 shares. Odd-lot buys are rejected; a sell of an odd lot is allowed only when it closes the whole position.',
              'A-shares follow T+1: shares bought today cannot be sold until the next trading day.',
              'A-shares cannot be shorted or leveraged.',
              'A-share orders respect daily price limits: 10% for the main board, 20% for ChiNext and STAR, 30% for Beijing, and 5% for ST stocks.',
              'US stocks may be shorted and accept fractional quantities. Leverage is not used.',
            ],
          },
          {
            type: 'p',
            text: 'Decisions that break these rules are rejected before execution. The AI is also told about them in its prompt.',
          },
        ],
      },
      {
        id: 'data-sources',
        question: 'Where does market data come from?',
        blocks: [
          {
            type: 'list',
            items: [
              'US stocks: daily bars from Alpaca market data.',
              'A-shares: forward-adjusted (前复权) daily bars from Tushare Pro.',
            ],
          },
          {
            type: 'note',
            text: 'The platform currently uses daily bars for both markets. Intraday data is not used.',
          },
        ],
      },
      {
        id: 'decision-cycle',
        question: 'How often does the AI decide?',
        blocks: [
          {
            type: 'p',
            text: 'Each trader runs on its own cycle. On every cycle it loads the latest daily bars for the symbols in its strategy, asks the AI for decisions, validates them, and then places orders with the paper broker.',
          },
        ],
      },
      {
        id: 'paper-trading',
        question: 'Is this paper trading or live trading?',
        blocks: [
          {
            type: 'p',
            text: 'Paper trading only. US orders go to an Alpaca paper account. A-share orders go to a built-in paper account that applies commissions and stamp duty. No real money is moved.',
          },
        ],
      },
    ],
  },

  // ───────────────────────── Dashboard ─────────────────────────
  {
    id: 'dashboard',
    title: 'Dashboard',
    icon: Monitor,
    items: [
      {
        id: 'daily-chart',
        question: 'Why does the price chart show daily bars only?',
        blocks: [
          {
            type: 'p',
            text: 'The chart shows the daily close for the first open position. It stays empty until a trader for that market has been created, because the data source is registered when the trader is created.',
          },
        ],
      },
      {
        id: 'equity-curve',
        question: 'What does the equity curve show?',
        blocks: [
          {
            type: 'p',
            text: 'Total account equity recorded at each decision cycle, along with the change from the starting cash. It is taken from the paper account, so it reflects fills, commissions and stamp duty.',
          },
        ],
      },
    ],
  },

  // ───────────────────────── Security ─────────────────────────
  {
    id: 'security',
    title: 'Security',
    icon: Shield,
    items: [
      {
        id: 'key-storage',
        question: 'How are my API keys stored?',
        blocks: [
          {
            type: 'p',
            text: 'Broker, data and AI model keys are encrypted at rest with the `DATA_ENCRYPTION_KEY` you set on your server.',
          },
        ],
      },
      {
        id: 'who-controls-funds',
        question: 'Can Crux access my brokerage account?',
        blocks: [
          {
            type: 'p',
            text: 'Crux runs on your own server and calls your broker and data provider with the keys you configure. Anonymous usage statistics can be turned off with `EXPERIENCE_IMPROVEMENT=false`.',
          },
        ],
      },
    ],
  },

  // ───────────────────────── Self-hosting ─────────────────────────
  {
    id: 'self-hosting',
    title: 'Self-hosting',
    icon: Wrench,
    items: [
      {
        id: 'how-to-install',
        question: 'How do I run Crux?',
        blocks: [
          {
            type: 'steps',
            items: [
              'Copy `.env.example` to `.env` and set `JWT_SECRET`, `DATA_ENCRYPTION_KEY` and `RSA_PRIVATE_KEY`.',
              'Start the services with `docker compose up -d`.',
              'Open the web app, create the admin account, then add an AI model and an exchange account.',
              'Create a strategy in Strategy Studio and a trader that uses it.',
            ],
          },
        ],
      },
      {
        id: 'where-are-logs',
        question: 'Where are the logs?',
        blocks: [
          {
            type: 'p',
            text: 'The backend writes logs to standard output, so `docker compose logs -f` shows them. Paper account state is saved under `data/paper`.',
          },
        ],
      },
    ],
  },

  // ───────────────────────── Contributing ─────────────────────────
]
