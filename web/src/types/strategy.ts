// Strategy types for stock markets (US equities and A-shares).

export interface Strategy {
  id: string;
  name: string;
  description: string;
  is_active: boolean;
  is_default: boolean;
  is_public: boolean;           // Whether published in the strategy marketplace
  config_visible: boolean;      // Whether the config parameters are publicly visible
  config: StrategyConfig;
  created_at: string;
  updated_at: string;
}

// Strategy usage statistics
export interface StrategyStats {
  clone_count: number;          // Number of times cloned
  active_users: number;         // Current number of users
}

export interface PromptSectionsConfig {
  role_definition?: string;
  trading_frequency?: string;
  entry_standards?: string;
  decision_process?: string;
}

// The stored and API shape: the decision settings live under ai_config.
export interface StrategyConfig {
  // Language setting: "zh" for Chinese, "en" for English
  language?: 'zh' | 'en' | string;
  ai_config: AIStrategyConfig;
  publish_config?: PublishStrategyConfig;
}

export interface AIStrategyConfig {
  coin_source: CoinSourceConfig;
  indicators: IndicatorConfig;
  custom_prompt?: string;
  risk_control: RiskControlConfig;
  prompt_sections?: PromptSectionsConfig;
}

export interface PublishStrategyConfig {
  is_public: boolean;
  config_visible: boolean;
}

// Symbols the AI may trade. Only a static list is supported.
export interface CoinSourceConfig {
  source_type: 'static';
  static_coins?: string[];
  excluded_coins?: string[];
}

export interface IndicatorConfig {
  klines: KlineConfig;
  enable_raw_klines: boolean;
  enable_ema: boolean;
  enable_macd: boolean;
  enable_rsi: boolean;
  enable_atr: boolean;
  enable_boll: boolean;
  enable_volume: boolean;
  ema_periods?: number[];
  rsi_periods?: number[];
  atr_periods?: number[];
  boll_periods?: number[];
  external_data_sources?: ExternalDataSource[];
}

export interface KlineConfig {
  primary_timeframe: string;    // stock data is daily: "1d"
  primary_count: number;
  longer_timeframe?: string;
  longer_count?: number;
  enable_multi_timeframe: boolean;
  selected_timeframes?: string[];
}

export interface ExternalDataSource {
  name: string;
  type: string;
  url: string;
  method: string;
  headers?: Record<string, string>;
  data_path?: string;
  refresh_secs?: number;
}

export interface RiskControlConfig {
  max_positions: number;
  altcoin_max_position_value_ratio: number; // single position value cap, multiple of equity
  max_margin_usage: number;
  min_position_size: number;
  min_risk_reward_ratio: number;
  min_confidence: number;
}
