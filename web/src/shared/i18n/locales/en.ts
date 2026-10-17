// en UI strings. Keys are read with t() in translations.ts.
const en = {
  leadPnL: 'Lead PnL',
  loadingChartData: 'Loading chart data...',
  loadingError: 'Failed to load data',
  loggingIn: 'Signing in...',
  // Header
  appTitle: 'Crux',
  subtitle: 'Multi-AI Model Trading Platform',
  aiTraders: 'AI Traders',
  details: 'Details',
  competition: 'Competition',
  running: 'RUNNING',
  stopped: 'STOPPED',
  logout: 'Logout',
  view: 'View',

  // Navigation
  realtimeNav: 'Leaderboard',
  configNav: 'Config',
  dashboardNav: 'Dashboard',
  strategyNav: 'Strategy',
  faqNav: 'FAQ',
  intelNav: 'Intel',

  // Footer
  footerTitle: 'Crux - AI Trading System',
  footerWarning: '⚠️ Trading involves risk. Use at your own discretion.',

  // Stats Cards
  totalEquity: 'Total Equity',
  availableBalance: 'Available Balance',
  totalPnL: 'Total P&L',
  positions: 'Positions',
  margin: 'Margin',
  free: 'Free',

  // Positions Table
  currentPositions: 'Current Positions',
  active: 'Active',
  symbol: 'Symbol',
  side: 'Side',
  entryPrice: 'Entry Price',
  stopLoss: 'Stop Loss',
  takeProfit: 'Take Profit',
  riskReward: 'Risk/Reward',
  markPrice: 'Mark Price',
  quantity: 'Quantity',
  positionValue: 'Position Value',
  unrealizedPnL: 'Unrealized P&L',
  liqPrice: 'Liq. Price',
  long: 'LONG',
  short: 'SHORT',
  noPositions: 'No Positions',
  noActivePositions: 'No active trading positions',

  // Recent Decisions
  recentDecisions: 'Recent Decisions',
  lastCycles: 'Last {count} trading cycles',
  noDecisionsYet: 'No Decisions Yet',
  aiDecisionsWillAppear: 'AI trading decisions will appear here',
  cycle: 'Cycle',
  success: 'Success',
  failed: 'Failed',
  aiThinking: 'AI Chain of Thought',
  collapse: 'Collapse',
  expand: 'Expand',

  // Equity Chart
  accountEquityCurve: 'Account Equity Curve',
  noHistoricalData: 'No Historical Data',
  dataWillAppear: 'Equity curve will appear after running a few cycles',
  initialBalance: 'Initial Balance',
  currentEquity: 'Current Equity',
  historicalCycles: 'Historical Cycles',
  displayRange: 'Display Range',
  recent: 'Recent',
  allData: 'All Data',
  cycles: 'Cycles',

  // Comparison Chart
  dataPoints: 'Data Points',
  currentGap: 'Current Gap',
  count: '{count} pts',

  // TradingView Chart
  fullscreen: 'Fullscreen',

  // Competition Page
  aiCompetition: 'AI Competition',
  traders: 'traders',
  liveBattle: 'Live Battle',
  leader: 'Leader',
  leaderboard: 'Leaderboard',
  live: 'LIVE',
  performanceComparison: 'Performance Comparison',
  realTimePnL: 'Real-time PnL %',
  headToHead: 'Head-to-Head Battle',
  leadingBy: 'Leading by {gap}%',
  behindBy: 'Behind by {gap}%',
  equity: 'Equity',
  pnl: 'P&L',
  pos: 'Pos',

  // AI Traders Management
  aiModels: 'AI Models',
  exchanges: 'Exchanges',
  createTrader: 'Create Trader',
  configured: 'Configured',
  currentTraders: 'Current Traders',
  noTraders: 'No AI Traders',
  createFirstTrader: 'Create your first AI trader to get started',
  dashboardEmptyTitle: "Let's Get Started!",
  dashboardEmptyDescription:
    'Create your first AI trader to automate your trading strategy. Connect an exchange, choose an AI model, and start trading in minutes!',
  goToTradersPage: 'Create Your First Trader',
  configureModelsFirst: 'Please configure AI models first',
  configureExchangesFirst: 'Please configure exchanges first',
  configureModelsAndExchangesFirst:
    'Please configure AI models and exchanges first',
  modelNotConfigured: 'Selected model is not configured',
  exchangeNotConfigured: 'Selected exchange is not configured',
  confirmDeleteTrader: 'Are you sure you want to delete this trader?',
  status: 'Status',
  start: 'Start',
  stop: 'Stop',
  traderName: 'Trader Name',
  cancel: 'Cancel',
  create: 'Create',
  aiScanInterval: 'AI Scan Decision Interval (minutes)',
  scanIntervalRecommend: 'Recommended: 15-30 minutes',
  enabled: 'Enabled',
  save: 'Save',

  // TraderConfigModal - New keys for hardcoded Chinese strings
  saving: 'Saving...',
  saveFailed: 'Save failed',
  editTraderConfig: 'Edit Trader Configuration',
  selectStrategyAndConfigParams:
    'Select Strategy and Configure Basic Parameters',
  basicConfig: 'Basic Configuration',
  traderNameRequired: 'Trader Name *',
  enterTraderNamePlaceholder: 'Enter trader name',
  aiModelRequired: 'AI Model *',
  exchangeRequired: 'Exchange *',
  selectTradingStrategy: 'Select Trading Strategy',
  useStrategy: 'Use Strategy',
  noStrategyManual: '-- No Strategy (Manual Configuration) --',
  strategyActive: ' (Active)',
  strategyDefault: ' [Default]',
  noStrategyHint: 'No strategies yet, please create in Strategy Studio first',
  strategyDetails: 'Strategy Details',
  activating: 'Activating',
  coinSource: 'Coin Source',
  marginLimit: 'Margin Limit',
  tradingParams: 'Trading Parameters',
  competitionDisplay: 'Show in Competition',
  show: 'Show',
  hide: 'Hide',
  hiddenInCompetition:
    'This trader will not be shown in the competition page when hidden',
  fetching: 'Fetching...',
  autoFetchBalanceInfo:
    'The system will automatically fetch your account equity as the initial balance',
  editTrader: 'Save Changes',
  createTraderButton: 'Create Trader',

  // AI Model Configuration
  apiKey: 'API Key',
  enterAPIKey: 'Enter API Key',

  // Exchange Configuration
  secretKey: 'Secret Key',
  walletAddress: 'Wallet Address',
  user: 'User',

  // Exchange names

  // Secure input

  // Two Stage Key Modal

  // Trader Configuration

  // System Prompt Templates

  // Loading & Error
  loading: 'Loading...',

  // AI Traders Page - Additional
  noModelsConfigured: 'No configured AI models',
  noExchangesConfigured: 'No configured exchanges',
  information: 'Information',
  editAIModel: 'Edit AI Model',
  addAIModel: 'Add AI Model',
  confirmDeleteModel:
    'Are you sure you want to delete this AI model configuration?',
  cannotDeleteModelInUse:
    'Cannot delete this AI model because it is being used by traders',
  tradersUsing: 'Traders using this configuration',
  pleaseDeleteTradersFirst: 'Please delete or reconfigure these traders first',
  selectModel: 'Select AI Model',
  customBaseURL: 'Base URL (Optional)',
  customBaseURLPlaceholder:
    'Custom API base URL, e.g.: https://api.openai.com/v1',
  leaveBlankForDefault: 'Leave blank to use default API address',
  modelConfigInfo1:
    '• For official API, only API Key is required, leave other fields blank',
  modelConfigInfo2:
    '• Custom Base URL and Model Name only needed for third-party proxies',
  modelConfigInfo3: '• API Key is encrypted and stored securely',
  defaultModel: 'Default model',
  kimiApiNote:
    'Kimi requires API Key from international site (moonshot.ai), China region keys are not compatible',
  leaveBlankForDefaultModel: 'Leave blank to use default model',
  customModelName: 'Model Name (Optional)',
  customModelNamePlaceholder: 'e.g.: deepseek-chat, qwen3-max, gpt-4o',
  saveConfig: 'Save Configuration',
  confirmDeleteExchange:
    'Are you sure you want to delete this exchange configuration?',
  cannotDeleteExchangeInUse:
    'Cannot delete this exchange because it is being used by traders',
  edit: 'Edit',

  // Error Messages
  createTraderFailed: 'Failed to create trader',
  getTraderConfigFailed: 'Failed to get trader configuration',
  modelConfigNotExist: 'Model configuration does not exist or is not enabled',
  exchangeConfigNotExist:
    'Exchange configuration does not exist or is not enabled',
  updateTraderFailed: 'Failed to update trader',
  deleteTraderFailed: 'Failed to delete trader',
  operationFailed: 'Operation failed',
  deleteConfigFailed: 'Failed to delete configuration',
  modelNotExist: 'Model does not exist',
  saveConfigFailed: 'Failed to save configuration',
  exchangeNotExist: 'Exchange does not exist',
  deleteExchangeConfigFailed: 'Failed to delete exchange configuration',

  // Login & Register
  login: 'Sign In',
  register: 'Sign Up',
  email: 'Email',
  password: 'Password',
  confirmPassword: 'Confirm Password',
  emailPlaceholder: 'your@email.com',
  passwordPlaceholder: 'Enter your password',
  passwordRuleMinLength: 'Minimum 8 characters',
  passwordRuleUppercase: 'At least 1 uppercase letter',
  passwordRuleLowercase: 'At least 1 lowercase letter',
  passwordRuleNumber: 'At least 1 number',
  passwordRuleSpecial: 'At least 1 special character (@#$%!&*?)',
  passwordRuleMatch: 'Passwords match',
  passwordNotMeetRequirements:
    'Password does not meet the security requirements',
  loginButton: 'Sign In',
  registerButton: 'Sign Up',
  back: 'Back',
  forgotPassword: 'Forgot password?',
  forgotAccount: 'Forgot account?',
  resetPassword: 'Reset Password',
  resetPasswordTitle: 'Reset your password',
  newPassword: 'New Password',
  backToLogin: 'Back to Login',
  resetPasswordCliIntro:
    'For security, password recovery is no longer available from the browser. Run this command on the server where Crux is installed:',
  resetPasswordCliSecurityNote:
    'This requires shell access to the server, which keeps your account safe even when Crux is exposed to the internet.',
  resetAccountCliIntro:
    'To wipe everything and start over, run this command on the server where Crux is installed:',
  copy: 'Copy',
  loginFailed: 'Login failed. Please check your email and password.',
  registrationFailed: 'Registration failed. Please try again.',
  sessionExpired: 'Session expired, please login again',
  medium: 'Medium',
  strong: 'Strong',

  // Landing Page
  features: 'Features',
  howItWorks: 'How it Works',
  community: 'Community',
  language: 'Language',
  loggedInAs: 'Logged in as',
  exitLogin: 'Sign Out',
  signIn: 'Sign In',
  registrationClosed: 'Registration Closed',
  registrationClosedMessage:
    'User registration is currently disabled. Please contact the administrator for access.',

  // Hero Section

  // Landing Page CTA

  // Features Section
  secureDesc:
    'Enterprise-grade security, complete control over your funds and trading strategies.',

  // About Section

  // How It Works Section
  step1Title: 'Clone GitHub Repository',
  step2Title: 'Configure Environment',

  // Community Section (testimonials are kept as-is since they are quotes)

  // Footer Section
  links: 'Links',

  // Login Modal

  // Candidate Coins Warnings
  candidateCoins: 'Candidate Coins',

  // FAQ Page

  // FAQ Categories

  // ===== GETTING STARTED =====

  // ===== INSTALLATION =====

  // ===== CONFIGURATION =====

  // ===== TRADING =====

  // ===== TECHNICAL ISSUES =====

  // ===== SECURITY =====

  // ===== FEATURES =====

  // ===== AI MODELS =====

  // ===== CONTRIBUTING =====

  // Web Crypto Environment Check
  environmentCheck: {
    button: 'Check Secure Environment',
    checking: 'Checking...',
    description:
      'Automatically verifying whether this browser context allows Web Crypto before entering sensitive keys.',
    secureTitle: 'Secure context detected',
    secureDesc:
      'Web Crypto API is available. You can continue entering secrets with encryption enabled.',
    insecureTitle: 'Insecure context detected',
    insecureDesc:
      'This page is not running over HTTPS or a trusted localhost origin, so browsers block Web Crypto calls.',
    tipsTitle: 'How to fix:',
    tipHTTPS:
      'Serve the dashboard over HTTPS with a valid certificate (IP origins also need TLS).',
    tipLocalhost:
      'During development, open the app via http://localhost or 127.0.0.1.',
    tipIframe:
      'Avoid embedding the app in insecure HTTP iframes or reverse proxies that strip HTTPS.',
    unsupportedTitle: 'Browser does not expose Web Crypto',
    unsupportedDesc:
      'Open Crux over HTTPS (or http://localhost during development) and avoid insecure iframes/reverse proxies so the browser can enable Web Crypto.',
    summary: 'Current origin: {origin} • Protocol: {protocol}',
    disabledTitle: 'Transport encryption disabled',
    disabledDesc:
      'Server-side transport encryption is disabled. API keys will be transmitted in plaintext. Enable TRANSPORT_ENCRYPTION=true for enhanced security.',
  },

  // Two-Stage Key Modal
  twoStageKey: {
    title: 'Two-Stage Private Key Input',
    stage1Description:
      'Enter the first {length} characters of your private key',
    stage2Description:
      'Enter the remaining {length} characters of your private key',
    stage1InputLabel: 'First Part',
    stage2InputLabel: 'Second Part',
    characters: 'characters',
    processing: 'Processing...',
    nextButton: 'Next',
    cancelButton: 'Cancel',
    backButton: 'Back',
    encryptButton: 'Encrypt & Submit',
    obfuscationCopied: 'Obfuscation data copied to clipboard',
    obfuscationInstruction:
      'Paste something else to clear clipboard, then continue',
    obfuscationManual: 'Manual obfuscation required',
  },

  // Error Messages
  errors: {
    privatekeyIncomplete: 'Please enter at least {expected} characters',
    privatekeyInvalidFormat:
      'Invalid private key format (should be 64 hex characters)',
    privatekeyObfuscationFailed: 'Clipboard obfuscation failed',
  },

  // Position History
  positionHistory: {
    title: 'Position History',
    loading: 'Loading position history...',
    noHistory: 'No Position History',
    noHistoryDesc: 'Closed positions will appear here after trading.',
    showingPositions: 'Showing {count} of {total} positions',
    totalPnL: 'Total P&L',
    // Stats
    totalTrades: 'Total Trades',
    winLoss: 'Win: {win} / Loss: {loss}',
    winRate: 'Win Rate',
    profitFactor: 'Profit Factor',
    profitFactorDesc: 'Total Profit / Total Loss',
    plRatio: 'P/L Ratio',
    plRatioDesc: 'Avg Win / Avg Loss',
    sharpeRatio: 'Sharpe Ratio',
    sharpeRatioDesc: 'Risk-adjusted Return',
    maxDrawdown: 'Max Drawdown',
    avgWin: 'Avg Win',
    avgLoss: 'Avg Loss',
    netPnL: 'Net P&L',
    netPnLDesc: 'After Fees',
    fee: 'Fee',
    // Direction Stats
    trades: 'Trades',
    avgPnL: 'Avg P&L',
    // Symbol Performance
    symbolPerformance: 'Symbol Performance',
    // Filters
    symbol: 'Symbol',
    allSymbols: 'All Symbols',
    side: 'Side',
    all: 'All',
    sort: 'Sort',
    latestFirst: 'Latest First',
    oldestFirst: 'Oldest First',
    highestPnL: 'Highest P&L',
    lowestPnL: 'Lowest P&L',
    // Table Headers
    entry: 'Entry',
    exit: 'Exit',
    qty: 'Qty',
    value: 'Value',
    pnl: 'P&L',
    duration: 'Duration',
    closedAt: 'Closed At',
  },

  // Data Page

  // Strategy Studio Page

  // Metric Tooltip
  metricTooltip: {
    formula: 'Formula',
  },

  // Login Required Overlay
  loginRequired: {
    title: 'SYSTEM ACCESS DENIED',
    accessDenied: 'ACCESS DENIED',
    subtitleWithFeature: 'Module "{featureName}" requires elevated privileges',
    subtitleDefault: 'Authorization required for this module',
    description:
      'Initialize authentication protocol to unlock full system capabilities: AI Trader configuration and Strategy Market data streams.',
    benefit1: 'AI Trader Control',
    benefit2: 'HFT Strategy Market',
    benefit4: 'Full System Visualization',
    loginButton: 'EXECUTE LOGIN',
    registerButton: 'REGISTER NEW ID',
    abort: 'ABORT',
  },

  // Advanced Chart

  // Chart With Orders

  // Comparison Chart
  comparisonChart: {
    '1d': '1D',
    '3d': '3D',
    '7d': '7D',
    '30d': '30D',
    all: 'All',
  },

  // TraderDashboardPage
  traderDashboard: {
    connectionFailed: 'Connection Failed',
    connectionFailedDesc: 'Please check if the backend service is running.',
    retry: 'Retry',
    confirmClosePosition:
      'Are you sure you want to close {symbol} {side} position?',
    confirmClose: 'Confirm Close',
    confirm: 'Confirm',
    cancel: 'Cancel',
    positionClosed: 'Position closed successfully',
    closeFailed: 'Failed to close position',
    closeAll: 'Close All',
    confirmCloseAllPositions: 'Market-close ALL {count} open positions?',
    allPositionsClosed: 'All positions closed',
    closeAllPartial: '{failed} of {count} positions failed to close',
    action: 'Action',
    entry: 'Entry',
    mark: 'Mark',
    qty: 'Qty',
    value: 'Value',
    uPnL: 'uPnL',
    liq: 'Liq.',
    closePosition: 'Close Position',
    close: 'Close',
    showingPositions: 'Showing {shown} of {total} positions',
    perPage: 'Per page',
    accountFetchFailed:
      'DATA_FETCH::FAILED — Account data unavailable, check connection',
    positionsFetchFailed: 'Position data unavailable',
    decisionsFetchFailed: 'Decision data unavailable',
  },

  // AITradersPage toast messages
  aiTradersToast: {
    creating: 'Creating...',
    created: 'Created successfully',
    saving: 'Saving...',
    saved: 'Saved successfully',
    saveFailed: 'Save failed',
    deleted: 'Deleted successfully',
    stopping: 'Stopping...',
    stopped: 'Stopped',
    starting: 'Starting...',
    started: 'Started',
    configUpdated: 'Config updated',
    showInCompetition: 'Shown in competition',
    hideInCompetition: 'Hidden from competition',
    modelConfigUpdated: 'Model config updated',
    exchangeDeleted: 'Exchange account deleted',
    exchangeConfigUpdated: 'Exchange config updated',
    exchangeCreated: 'Exchange account created',
  },

  // ModelConfigModal
  modelConfig: {
    selectModel: 'Select Model',
    configure: 'Configure',
    configureApi: 'Configure API',
    chooseProvider: 'Choose Your AI Provider',
    selectAiModel: 'Choose AI Model',
    exportKey: 'Export private key from MetaMask, Rabby, etc.',
    back: 'Back',
    modelsConfigured: 'Models with gold badge are already configured',
    getApiKey: 'Get API Key',
    validating: 'Validating...',
    walletAddress: 'Wallet Address',
  },

  // ExchangeConfigModal

  // TraderConfigViewModal
  traderConfigView: {
    traderConfig: 'Trader Configuration',
    configInfo: '{name} configuration details',
    running: 'Running',
    stopped: 'Stopped',
    basicInfo: 'Basic Information',
    traderName: 'Trader Name',
    aiModel: 'AI Model',
    exchange: 'Exchange',
    initialBalance: 'Initial Balance',
    scanInterval: '{minutes} minutes',
    scanIntervalLabel: 'Scan Interval',
    strategyUsed: 'Strategy Used',
    strategyName: 'Strategy Name',
    close: 'Close',
    yes: 'Yes',
    no: 'No',
  },
  terminal: {
    closePosition: 'Close position',
    close: 'Close',
    cancel: 'Cancel',
    closeAll: 'Close all',
    flattenBook: 'Flatten book',
    closeFailed: 'Close failed',
    allPositionsClosed: 'All positions closed',
    tradingBlocked: 'Trading blocked',
    model: 'model',
    strategy: 'strategy',
    scan: 'scan',
    universe: 'universe',
    positions: 'positions',
    nextCycle: 'next cycle',
    equity: 'Equity',
    profitFactor: 'Profit factor',
    maxDrawdown: 'Max drawdown',
    currentPositionsLive: 'Current positions · live',
    symbol: 'symbol',
    hold: 'hold',
    closed: 'closed',
    side: 'side',
    bySymbolHistory: 'By-symbol history · trades/win/pnl',
    noSymbolHistory: 'No symbol history.',
  },
  intel: {
    title: 'Intelligence',
    subtitle:
      'Analyse public web pages about a stock with an AI model. Results feed the strategy screen and the open-position gate.',
    symbol: 'Stock code',
    symbolHint: 'e.g. 600519.SH or AAPL',
    urls: 'Page URLs',
    urlsHint: 'One per line, up to 5',
    model: 'AI model',
    noModel: 'Enable an AI model in Config first',
    run: 'Analyse',
    running: 'Analysing…',
    results: 'Results',
    history: 'Recent analysis',
    empty: 'No analysed pages yet for this stock.',
    notConfigured:
      'The intelligence crawler is not configured. Set CRAWL4AI_URL on the server.',
    invalidUrls: 'Enter at least one http(s) URL.',
    tooMany: 'At most 5 URLs per run.',
    risk: {
      none: 'No risk',
      low: 'Low risk',
      high: 'High risk',
      unknown: 'Unknown',
    },
    sentiment: 'Sentiment',
    mentions: 'Also mentions',
    failed: 'Failed',
    load: 'Load',
  },
} as const

export default en
