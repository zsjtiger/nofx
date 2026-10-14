export type Language = 'en' | 'zh' | 'id'

export const translations = {
  en: {
    leadPnL: 'Lead PnL',
    loadingChartData: 'Loading chart data...',
    loadingError: 'Failed to load data',
    loggingIn: 'Signing in...',
    // Header
    appTitle: 'NOFX',
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

    // Footer
    footerTitle: 'NOFX - AI Trading System',
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
    leverage: 'Leverage',
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
    marginMode: 'Margin Mode',
    crossMargin: 'Cross Margin',
    isolatedMargin: 'Isolated Margin',
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
    // Hyperliquid Agent Wallet (New Security Model)
    // Aster API Pro Configuration

    // LIGHTER Configuration

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
    pleaseDeleteTradersFirst:
      'Please delete or reconfigure these traders first',
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
      'For security, password recovery is no longer available from the browser. Run this command on the server where NOFX is installed:',
    resetPasswordCliSecurityNote:
      'This requires shell access to the server, which keeps your account safe even when NOFX is exposed to the internet.',
    resetAccountCliIntro:
      'To wipe everything and start over, run this command on the server where NOFX is installed:',
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
        'Open NOFX over HTTPS (or http://localhost during development) and avoid insecure iframes/reverse proxies so the browser can enable Web Crypto.',
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
      lev: 'Lev',
      pnl: 'P&L',
      duration: 'Duration',
      closedAt: 'Closed At',
    },

    // Data Page

    // Strategy Market Page
    strategyMarket: {
      title: 'STRATEGY MARKET',
      subtitle: 'GLOBAL STRATEGY DATABASE',
      description:
        'Discover, analyze, and clone high-performance trading algorithms',
      search: 'SEARCH PARAMETERS...',
      all: 'ALL PROTOCOLS',
      recent: 'LATEST',
      createdAt: 'TIMESTAMP',
      copied: 'COPIED',
      indicators: 'INDICATORS',
      maxPositions: 'POS_LIMIT',
      loading: 'INITIALIZING...',
    },

    // Strategy Studio Page

    // Metric Tooltip
    metricTooltip: {
      formula: 'Formula',
    },

    // Login Required Overlay
    loginRequired: {
      title: 'SYSTEM ACCESS DENIED',
      accessDenied: 'ACCESS DENIED',
      subtitleWithFeature:
        'Module "{featureName}" requires elevated privileges',
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
      confirmCloseAllPositions:
        'Market-close ALL {count} open positions?',
      allPositionsClosed: 'All positions closed',
      closeAllPartial: '{failed} of {count} positions failed to close',
      action: 'Action',
      entry: 'Entry',
      mark: 'Mark',
      qty: 'Qty',
      value: 'Value',
      lev: 'Lev.',
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

    // TelegramConfigModal
    telegram: {
      botSetup: 'Telegram Bot Setup',
      createBot: 'Create Bot',
      bindAccount: 'Bind Account',
      done: 'Done',
      invalidTokenFormat:
        'Invalid Bot Token format. Expected "numbers:alphanumeric"',
      tokenSaved: 'Bot Token saved, waiting for binding',
      saveFailed: 'Save failed, please verify the token',
      unbound: 'Telegram account unbound',
      unbindFailed: 'Unbind failed',
      step1Title: 'Step 1: Create your Bot in Telegram',
      step1Desc1: 'Open Telegram, search for',
      step1Desc2: 'Send',
      step1Desc2Suffix: 'command',
      step1Desc3: 'Follow prompts to set bot name and username',
      step1Desc4: 'BotFather will return a Token, copy it',
      openBotFather: 'Open @BotFather',
      pasteToken: 'Paste Bot Token',
      tokenFormat: 'Format: numbers:alphanumeric, e.g. 123456789:ABCdef...',
      selectAiModel: 'Select AI Model (optional)',
      noEnabledModels: 'No enabled models. Configure one in AI Models first.',
      autoSelect: '— Auto-select (recommended)',
      autoUseEnabled: 'Leave blank to auto-use any enabled model',
      savingToken: 'Saving...',
      saveAndContinue: 'Save & Continue',
      step2Title: 'Step 2: Send /start to your Bot',
      step2Desc1: 'Search for your newly created Bot in Telegram',
      step2Desc2: 'Click Start or send',
      step2Desc3: 'Bot will automatically bind to your account',
      currentToken: 'Current Token',
      waitingForStart:
        'Waiting for you to send /start... Refresh page after sending',
      reconfigureToken: 'Reconfigure Token',
      bindSuccess: 'Bound successfully!',
      noStartReceived:
        'No /start received yet. Please send /start to your Bot first',
      checkFailed: 'Check failed',
      checkStatus: 'Check Status',
      botActive: 'Telegram Bot is Active!',
      botActiveDesc:
        'You can now control the trading system via natural language in Telegram',
      supportedCommands: 'Supported Commands',
      cmdHelp: 'Show all commands',
      cmdStatus: 'Show trader status',
      cmdNaturalLang: 'Natural language',
      cmdStartStop: 'Start/stop trader',
      cmdControl: 'Natural language control',
      cmdPositions: 'View positions',
      cmdPositionsDesc: 'Real-time position query',
      cmdStrategy: 'Configure strategy',
      cmdStrategyDesc: 'Modify trading strategy',
      unbinding: 'Unbinding...',
      unbindAccount: 'Unbind Account',
      aiModelLabel: 'AI Model (for natural language)',
      aiModelAutoSelect: '— Auto-select',
      modelUpdated: 'AI model updated',
      modelUpdateFailed: 'Update failed',
      save: 'Save',
      loading: 'Loading...',
    },

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
      marginMode: 'Margin Mode',
      crossMargin: 'Cross',
      isolatedMargin: 'Isolated',
      scanInterval: '{minutes} minutes',
      scanIntervalLabel: 'Scan Interval',
      strategyUsed: 'Strategy Used',
      strategyName: 'Strategy Name',
      close: 'Close',
      yes: 'Yes',
      no: 'No',
    },
  },
  zh: {
    leadPnL: '领先收益',
    loadingChartData: '正在加载图表数据...',
    loadingError: '加载数据失败',
    loggingIn: '登录中...',
    // Header
    appTitle: 'NOFX',
    subtitle: '多AI模型交易平台',
    aiTraders: 'AI交易员',
    details: '详情',
    competition: '竞赛',
    running: '运行中',
    stopped: '已停止',
    logout: '退出',
    view: '查看',

    // Navigation
    realtimeNav: '排行榜',
    configNav: '配置',
    dashboardNav: '看板',
    strategyNav: '策略',
    faqNav: '常见问题',

    // Footer
    footerTitle: 'NOFX - AI交易系统',
    footerWarning: '⚠️ 交易有风险，请谨慎使用。',

    // Stats Cards
    totalEquity: '总净值',
    availableBalance: '可用余额',
    totalPnL: '总盈亏',
    positions: '持仓',
    margin: '保证金',
    free: '空闲',

    // Positions Table
    currentPositions: '当前持仓',
    active: '活跃',
    symbol: '币种',
    side: '方向',
    entryPrice: '入场价',
    stopLoss: '止损',
    takeProfit: '止盈',
    riskReward: '风险回报比',
    markPrice: '标记价',
    quantity: '数量',
    positionValue: '仓位价值',
    leverage: '杠杆',
    unrealizedPnL: '未实现盈亏',
    liqPrice: '强平价',
    long: '多头',
    short: '空头',
    noPositions: '无持仓',
    noActivePositions: '当前没有活跃的交易持仓',

    // Recent Decisions
    recentDecisions: '最近决策',
    lastCycles: '最近 {count} 个交易周期',
    noDecisionsYet: '暂无决策',
    aiDecisionsWillAppear: 'AI交易决策将显示在这里',
    cycle: '周期',
    success: '成功',
    failed: '失败',
    aiThinking: '💭 AI思维链分析',
    collapse: '▼ 收起',
    expand: '▶ 展开',

    // Equity Chart
    accountEquityCurve: '账户净值曲线',
    noHistoricalData: '暂无历史数据',
    dataWillAppear: '运行几个周期后将显示收益率曲线',
    initialBalance: '初始余额',
    currentEquity: '当前净值',
    historicalCycles: '历史周期',
    displayRange: '显示范围',
    recent: '最近',
    allData: '全部数据',
    cycles: '个',

    // Comparison Chart
    dataPoints: '数据点数',
    currentGap: '当前差距',
    count: '{count} 个',

    // TradingView Chart
    fullscreen: '全屏',

    // Competition Page
    aiCompetition: 'AI竞赛',
    traders: '交易员',
    liveBattle: '实时对战',
    leader: '领先者',
    leaderboard: '排行榜',
    live: '实时',
    performanceComparison: '表现对比',
    realTimePnL: '实时收益率',
    headToHead: '正面对决',
    leadingBy: '领先 {gap}%',
    behindBy: '落后 {gap}%',
    equity: '权益',
    pnl: '收益',
    pos: '持仓',

    // AI Traders Management
    aiModels: 'AI模型',
    exchanges: '交易所',
    createTrader: '创建交易员',
    configured: '已配置',
    currentTraders: '当前交易员',
    noTraders: '暂无AI交易员',
    createFirstTrader: '创建您的第一个AI交易员开始使用',
    dashboardEmptyTitle: '开始使用吧！',
    dashboardEmptyDescription:
      '创建您的第一个 AI 交易员，自动化您的交易策略。连接交易所、选择 AI 模型，几分钟内即可开始交易！',
    goToTradersPage: '创建您的第一个交易员',
    configureModelsFirst: '请先配置AI模型',
    configureExchangesFirst: '请先配置交易所',
    configureModelsAndExchangesFirst: '请先配置AI模型和交易所',
    modelNotConfigured: '所选模型未配置',
    exchangeNotConfigured: '所选交易所未配置',
    confirmDeleteTrader: '确定要删除这个交易员吗？',
    status: '状态',
    start: '启动',
    stop: '停止',
    traderName: '交易员名称',
    cancel: '取消',
    create: '创建',
    aiScanInterval: 'AI 扫描决策间隔 (分钟)',
    scanIntervalRecommend: '建议: 15-30分钟',
    enabled: '启用',
    save: '保存',

    // TraderConfigModal - New keys for hardcoded Chinese strings
    saving: '正在保存…',
    saveFailed: '保存失败',
    editTraderConfig: '修改交易员配置',
    selectStrategyAndConfigParams: '选择策略并配置基础参数',
    basicConfig: '基础配置',
    traderNameRequired: '交易员名称 *',
    enterTraderNamePlaceholder: '请输入交易员名称',
    aiModelRequired: 'AI模型 *',
    exchangeRequired: '交易所 *',
    selectTradingStrategy: '选择交易策略',
    useStrategy: '使用策略',
    noStrategyManual: '-- 不使用策略（手动配置） --',
    strategyActive: ' (当前激活)',
    strategyDefault: ' [默认]',
    noStrategyHint: '暂无策略，请先在策略工作室创建策略',
    strategyDetails: '策略详情',
    activating: '激活中',
    coinSource: '币种来源',
    marginLimit: '保证金上限',
    tradingParams: '交易参数',
    marginMode: '保证金模式',
    crossMargin: '全仓',
    isolatedMargin: '逐仓',
    competitionDisplay: '竞技场显示',
    show: '显示',
    hide: '隐藏',
    hiddenInCompetition: '隐藏后将不在竞技场页面显示此交易员',
    fetching: '获取中...',
    autoFetchBalanceInfo: '系统将自动获取您的账户净值作为初始余额',
    editTrader: '保存修改',
    createTraderButton: '创建交易员',

    // AI Model Configuration
    apiKey: 'API密钥',
    enterAPIKey: '请输入API密钥',

    // Exchange Configuration
    secretKey: '密钥',
    walletAddress: '钱包地址',
    user: '用户名',
    // Hyperliquid 代理钱包 (新安全模型)
    // Aster API Pro 配置

    // LIGHTER 配置

    // Exchange names

    // Secure input

    // Two Stage Key Modal

    // Trader Configuration

    // System Prompt Templates

    // Loading & Error
    loading: '加载中...',

    // AI Traders Page - Additional
    noModelsConfigured: '暂无已配置的AI模型',
    noExchangesConfigured: '暂无已配置的交易所',
    information: '说明',
    editAIModel: '编辑AI模型',
    addAIModel: '添加AI模型',
    confirmDeleteModel: '确定要删除此AI模型配置吗？',
    cannotDeleteModelInUse: '无法删除此AI模型，因为有交易员正在使用',
    tradersUsing: '正在使用此配置的交易员',
    pleaseDeleteTradersFirst: '请先删除或重新配置这些交易员',
    selectModel: '选择AI模型',
    customBaseURL: 'Base URL (可选)',
    customBaseURLPlaceholder: '自定义API基础URL，如: https://api.openai.com/v1',
    leaveBlankForDefault: '留空则使用默认API地址',
    modelConfigInfo1: '• 使用官方 API 时，只需填写 API Key，其他字段留空即可',
    modelConfigInfo2:
      '• 自定义 Base URL 和 Model Name 仅在使用第三方代理时需要填写',
    modelConfigInfo3: '• API Key 加密存储，不会明文展示',
    defaultModel: '默认模型',
    kimiApiNote:
      'Kimi 需要从国际站申请 API Key (moonshot.ai)，中国区 Key 不通用',
    leaveBlankForDefaultModel: '留空使用默认模型名称',
    customModelName: 'Model Name (可选)',
    customModelNamePlaceholder: '例如: deepseek-chat, qwen3-max, gpt-4o',
    saveConfig: '保存配置',
    confirmDeleteExchange: '确定要删除此交易所配置吗？',
    cannotDeleteExchangeInUse: '无法删除此交易所，因为有交易员正在使用',
    edit: '编辑',

    // Error Messages
    createTraderFailed: '创建交易员失败',
    getTraderConfigFailed: '获取交易员配置失败',
    modelConfigNotExist: 'AI模型配置不存在或未启用',
    exchangeConfigNotExist: '交易所配置不存在或未启用',
    updateTraderFailed: '更新交易员失败',
    deleteTraderFailed: '删除交易员失败',
    operationFailed: '操作失败',
    deleteConfigFailed: '删除配置失败',
    modelNotExist: '模型不存在',
    saveConfigFailed: '保存配置失败',
    exchangeNotExist: '交易所不存在',
    deleteExchangeConfigFailed: '删除交易所配置失败',

    // Login & Register
    login: '登录',
    register: '注册',
    email: '邮箱',
    password: '密码',
    confirmPassword: '确认密码',
    emailPlaceholder: '请输入邮箱地址',
    passwordPlaceholder: '请输入密码（至少6位）',
    passwordRuleMinLength: '至少 8 位',
    passwordRuleUppercase: '至少 1 个大写字母',
    passwordRuleLowercase: '至少 1 个小写字母',
    passwordRuleNumber: '至少 1 个数字',
    passwordRuleSpecial: '至少 1 个特殊字符（@#$%!&*?）',
    passwordRuleMatch: '两次密码一致',
    passwordNotMeetRequirements: '密码不符合安全要求',
    loginButton: '登录',
    registerButton: '注册',
    back: '返回',
    forgotPassword: '忘记密码？',
    forgotAccount: '忘记账户？',
    resetPassword: '重置密码',
    resetPasswordTitle: '重置您的密码',
    newPassword: '新密码',
    backToLogin: '返回登录',
    resetPasswordCliIntro:
      '出于安全考虑，密码找回不再通过浏览器进行。请在部署 NOFX 的服务器上运行以下命令：',
    resetPasswordCliSecurityNote:
      '该操作需要服务器的 shell 访问权限，因此即使 NOFX 暴露在公网上，你的账户依然安全。',
    resetAccountCliIntro:
      '如需清空所有数据并重新开始，请在部署 NOFX 的服务器上运行以下命令：',
    copy: '复制',
    loginFailed: '登录失败，请检查您的邮箱和密码。',
    registrationFailed: '注册失败，请重试。',
    sessionExpired: '登录已过期，请重新登录',
    medium: '中',
    strong: '强',

    // Landing Page
    features: '功能',
    howItWorks: '如何运作',
    community: '社区',
    language: '语言',
    loggedInAs: '已登录为',
    exitLogin: '退出登录',
    signIn: '登录',
    registrationClosed: '注册已关闭',
    registrationClosedMessage:
      '平台当前不开放新用户注册，如需访问请联系管理员获取账号。',

    // Hero Section

    // Landing Page CTA

    // Features Section
    secureDesc: '企业级安全保障，完全掌控你的资金和交易策略。',

    // About Section

    // How It Works Section
    step1Title: '拉取 GitHub 仓库',
    step2Title: '配置环境',

    // Community Section (testimonials are kept as-is since they are quotes)

    // Footer Section
    links: '链接',

    // Login Modal

    // Candidate Coins Warnings
    candidateCoins: '候选币种',

    // FAQ Page

    // FAQ Categories

    // ===== 入门指南 =====






    // ===== 安装部署 =====






    // ===== 配置设置 =====






    // ===== 交易相关 =====








    // ===== 技术问题 =====








    // ===== 安全相关 =====




    // ===== 功能介绍 =====



    // ===== AI 模型 =====




    // ===== 参与贡献 =====




    // Web Crypto Environment Check
    environmentCheck: {
      button: '一键检测环境',
      checking: '正在检测...',
      description: '系统将自动检测当前浏览器是否允许使用 Web Crypto。',
      secureTitle: '环境安全，已启用 Web Crypto',
      secureDesc: '页面处于安全上下文，可继续输入敏感信息并使用加密传输。',
      insecureTitle: '检测到非安全环境',
      insecureDesc:
        '当前访问未通过 HTTPS 或可信 localhost，浏览器会阻止 Web Crypto 调用。',
      tipsTitle: '修改建议：',
      tipHTTPS:
        '通过 HTTPS 访问（即使是 IP 也需证书），或部署到支持 TLS 的域名。',
      tipLocalhost: '开发阶段请使用 http://localhost 或 127.0.0.1。',
      tipIframe:
        '避免把应用嵌入在不安全的 HTTP iframe 或会降级协议的反向代理中。',
      unsupportedTitle: '浏览器未提供 Web Crypto',
      unsupportedDesc:
        '请通过 HTTPS 或本机 localhost 访问 NOFX，并避免嵌入不安全 iframe/反向代理，以符合浏览器的 Web Crypto 规则。',
      summary: '当前来源：{origin} · 协议：{protocol}',
      disabledTitle: '传输加密已禁用',
      disabledDesc:
        '服务端传输加密已关闭，API 密钥将以明文传输。如需增强安全性，请设置 TRANSPORT_ENCRYPTION=true。',
    },


    // Two-Stage Key Modal
    twoStageKey: {
      title: '两阶段私钥输入',
      stage1Description: '请输入私钥的前 {length} 位字符',
      stage2Description: '请输入私钥的后 {length} 位字符',
      stage1InputLabel: '第一部分',
      stage2InputLabel: '第二部分',
      characters: '位字符',
      processing: '处理中...',
      nextButton: '下一步',
      cancelButton: '取消',
      backButton: '返回',
      encryptButton: '加密并提交',
      obfuscationCopied: '混淆数据已复制到剪贴板',
      obfuscationInstruction: '请粘贴其他内容清空剪贴板，然后继续',
      obfuscationManual: '需要手动混淆',
    },

    // Error Messages
    errors: {
      privatekeyIncomplete: '请输入至少 {expected} 位字符',
      privatekeyInvalidFormat: '私钥格式无效（应为64位十六进制字符）',
      privatekeyObfuscationFailed: '剪贴板混淆失败',
    },

    // Position History
    positionHistory: {
      title: '历史仓位',
      loading: '加载历史仓位...',
      noHistory: '暂无历史仓位',
      noHistoryDesc: '平仓后的仓位记录将显示在此处',
      showingPositions: '显示 {count} / {total} 条记录',
      totalPnL: '总盈亏',
      // Stats
      totalTrades: '总交易次数',
      winLoss: '盈利: {win} / 亏损: {loss}',
      winRate: '胜率',
      profitFactor: '盈利因子',
      profitFactorDesc: '总盈利 / 总亏损',
      plRatio: '盈亏比',
      plRatioDesc: '平均盈利 / 平均亏损',
      sharpeRatio: '夏普比率',
      sharpeRatioDesc: '风险调整收益',
      maxDrawdown: '最大回撤',
      avgWin: '平均盈利',
      avgLoss: '平均亏损',
      netPnL: '净盈亏',
      netPnLDesc: '扣除手续费后',
      fee: '手续费',
      // Direction Stats
      trades: '交易次数',
      avgPnL: '平均盈亏',
      // Symbol Performance
      symbolPerformance: '品种表现',
      // Filters
      symbol: '交易对',
      allSymbols: '全部交易对',
      side: '方向',
      all: '全部',
      sort: '排序',
      latestFirst: '最新优先',
      oldestFirst: '最早优先',
      highestPnL: '盈利最高',
      lowestPnL: '亏损最多',
      // Table Headers
      entry: '开仓价',
      exit: '平仓价',
      qty: '数量',
      value: '仓位价值',
      lev: '杠杆',
      pnl: '盈亏',
      duration: '持仓时长',
      closedAt: '平仓时间',
    },

    // Data Page

    // Strategy Market Page
    strategyMarket: {
      title: '策略市场',
      subtitle: 'STRATEGY MARKETPLACE',
      description: '发现、学习并复用社区精英交易员的策略配置',
      search: '搜索参数...',
      all: '全部协议',
      recent: '最新提交',
      createdAt: 'TIMESTAMP',
      copied: 'COPIED',
      indicators: 'INDICATORS',
      maxPositions: 'POS_LIMIT',
      loading: 'INITIALIZING...',
    },

    // Strategy Studio Page

    // Metric Tooltip
    metricTooltip: {
      formula: '计算公式',
    },

    // Login Required Overlay
    loginRequired: {
      title: '系统访问受限',
      accessDenied: '访问被拒绝',
      subtitleWithFeature: '访问「{featureName}」需要更高权限',
      subtitleDefault: '此模块需要授权访问',
      description:
        '初始化身份验证协议以解锁完整系统功能：AI 交易员配置、策略市场数据流。',
      benefit1: 'AI 交易员控制权',
      benefit2: '高频策略核心市场',
      benefit4: '全系统数据可视化',
      loginButton: '执行登录指令',
      registerButton: '注册新用户 ID',
      abort: '中止操作',
    },

    // Advanced Chart

    // Chart With Orders

    // Comparison Chart
    comparisonChart: {
      '1d': '1天',
      '3d': '3天',
      '7d': '7天',
      '30d': '30天',
      all: '全部',
    },

    traderDashboard: {
      connectionFailed: '无法连接到服务器',
      connectionFailedDesc: '请确认后端服务已启动。',
      retry: '重试',
      confirmClosePosition: '确定要平仓 {symbol} {side} 吗？',
      confirmClose: '确认平仓',
      confirm: '确认',
      cancel: '取消',
      positionClosed: '平仓成功',
      closeFailed: '平仓失败',
      closeAll: '一键全平',
      confirmCloseAllPositions: '确定要市价平掉全部 {count} 个持仓吗？',
      allPositionsClosed: '全部持仓已平',
      closeAllPartial: '{count} 个持仓中有 {failed} 个平仓失败',
      action: '操作',
      entry: '入场价',
      mark: '标记价',
      qty: '数量',
      value: '价值',
      lev: '杠杆',
      uPnL: '未实现盈亏',
      liq: '强平价',
      closePosition: '平仓',
      close: '平仓',
      showingPositions: '显示 {shown} / {total} 个持仓',
      perPage: '每页',
      accountFetchFailed: 'DATA_FETCH::FAILED — 账户数据请求失败，请检查连接',
      positionsFetchFailed: '持仓数据请求失败',
      decisionsFetchFailed: '决策记录请求失败',
    },

    aiTradersToast: {
      creating: '正在创建…',
      created: '创建成功',
      saving: '正在保存…',
      saved: '保存成功',
      saveFailed: '保存失败',
      deleted: '删除成功',
      stopping: '正在停止…',
      stopped: '已停止',
      starting: '正在启动…',
      started: '已启动',
      configUpdated: '配置已更新',
      showInCompetition: '已在竞技场显示',
      hideInCompetition: '已在竞技场隐藏',
      modelConfigUpdated: '模型配置已更新',
      exchangeDeleted: '交易所账户已删除',
      exchangeConfigUpdated: '交易所配置已更新',
      exchangeCreated: '交易所账户已创建',
    },

    modelConfig: {
      selectModel: '选择模型',
      configure: '配置',
      configureApi: '配置 API',
      chooseProvider: '选择 AI 模型提供商',
      selectAiModel: '① 选择 AI 模型',
      exportKey: '可以用 MetaMask、Rabby 等钱包导出私钥',
      back: '返回',
      modelsConfigured: '带金色标记的模型已配置',
      getApiKey: '获取 API Key',
      validating: '验证中...',
      walletAddress: '钱包地址',
    },


    telegram: {
      botSetup: 'Telegram Bot 配置',
      createBot: '创建 Bot',
      bindAccount: '绑定账号',
      done: '完成',
      invalidTokenFormat: 'Bot Token 格式不正确，应为 "数字:字母数字串"',
      tokenSaved: 'Bot Token 已保存，等待绑定',
      saveFailed: '保存失败，请检查 Token 是否正确',
      unbound: '已解绑 Telegram 账号',
      unbindFailed: '解绑失败',
      step1Title: '第一步：在 Telegram 创建你的 Bot',
      step1Desc1: '打开 Telegram，搜索',
      step1Desc2: '发送',
      step1Desc2Suffix: '命令',
      step1Desc3: '按提示输入 Bot 名称和用户名',
      step1Desc4: 'BotFather 会返回一个 Token，复制它',
      openBotFather: '打开 @BotFather',
      pasteToken: '粘贴 Bot Token',
      tokenFormat: 'Token 格式：数字:字母数字串，如 123456789:ABCdef...',
      selectAiModel: '选择 AI 模型（可选）',
      noEnabledModels: '暂无启用的模型，请先在「AI 模型」中配置',
      autoSelect: '— 自动选择（推荐）',
      autoUseEnabled: '不选则自动使用已启用的模型',
      savingToken: '保存中...',
      saveAndContinue: '保存并继续',
      step2Title: '第二步：向你的 Bot 发送 /start',
      step2Desc1: '在 Telegram 中搜索你刚创建的 Bot',
      step2Desc2: '点击 Start 或发送',
      step2Desc3: 'Bot 会自动绑定到你的账号',
      currentToken: '当前 Token',
      waitingForStart: '⏳ 等待你发送 /start... 发送后刷新页面查看状态',
      reconfigureToken: '重新配置 Token',
      bindSuccess: '绑定成功！',
      noStartReceived: '尚未收到 /start，请先向 Bot 发送 /start',
      checkFailed: '检查失败',
      checkStatus: '检查绑定状态',
      botActive: 'Telegram Bot 已绑定！',
      botActiveDesc: '你现在可以通过 Telegram 用自然语言控制交易系统',
      supportedCommands: '支持的命令',
      cmdHelp: '查看所有命令',
      cmdStatus: '查看交易员状态',
      cmdNaturalLang: '自然语言查询',
      cmdStartStop: '启动/停止交易员',
      cmdControl: '自然语言控制',
      cmdPositions: '查看持仓',
      cmdPositionsDesc: '实时持仓查询',
      cmdStrategy: '配置策略',
      cmdStrategyDesc: '修改交易策略',
      unbinding: '解绑中...',
      unbindAccount: '解绑账号',
      aiModelLabel: 'AI 模型（用于自然语言解析）',
      aiModelAutoSelect: '— 自动选择',
      modelUpdated: 'AI 模型已更新',
      modelUpdateFailed: '更新失败',
      save: '保存',
      loading: '加载中...',
    },

    traderConfigView: {
      traderConfig: '交易员配置',
      configInfo: '{name} 的配置信息',
      running: '运行中',
      stopped: '已停止',
      basicInfo: '基础信息',
      traderName: '交易员名称',
      aiModel: 'AI模型',
      exchange: '交易所',
      initialBalance: '初始余额',
      marginMode: '保证金模式',
      crossMargin: '全仓',
      isolatedMargin: '逐仓',
      scanInterval: '{minutes} 分钟',
      scanIntervalLabel: '扫描间隔',
      strategyUsed: '使用策略',
      strategyName: '策略名称',
      close: '关闭',
      yes: '是',
      no: '否',
    },
  },
  id: {
    leadPnL: 'PnL Teratas',
    loadingChartData: 'Memuat data grafik...',
    loadingError: 'Gagal memuat data',
    loggingIn: 'Sedang masuk...',
    // Header
    appTitle: 'NOFX',
    subtitle: 'Platform Trading Multi-AI',
    aiTraders: 'Trader AI',
    details: 'Detail',
    competition: 'Kompetisi',
    running: 'BERJALAN',
    stopped: 'BERHENTI',
    logout: 'Keluar',
    view: 'Lihat',

    // Navigation
    realtimeNav: 'Papan Peringkat',
    configNav: 'Konfigurasi',
    dashboardNav: 'Dasbor',
    strategyNav: 'Strategi',
    faqNav: 'FAQ',

    // Footer
    footerTitle: 'NOFX - Sistem Trading AI',
    footerWarning: '⚠️ Trading memiliki risiko. Gunakan dengan bijak.',

    // Stats Cards
    totalEquity: 'Total Ekuitas',
    availableBalance: 'Saldo Tersedia',
    totalPnL: 'Total L/R',
    positions: 'Posisi',
    margin: 'Margin',
    free: 'Bebas',

    // Positions Table
    currentPositions: 'Posisi Saat Ini',
    active: 'Aktif',
    symbol: 'Simbol',
    side: 'Arah',
    entryPrice: 'Harga Masuk',
    stopLoss: 'Stop Loss',
    takeProfit: 'Take Profit',
    riskReward: 'Risiko/Imbalan',
    markPrice: 'Harga Tanda',
    quantity: 'Jumlah',
    positionValue: 'Nilai Posisi',
    leverage: 'Leverage',
    unrealizedPnL: 'L/R Belum Terealisasi',
    liqPrice: 'Harga Likuidasi',
    long: 'LONG',
    short: 'SHORT',
    noPositions: 'Tidak Ada Posisi',
    noActivePositions: 'Tidak ada posisi trading yang aktif',

    // Recent Decisions
    recentDecisions: 'Keputusan Terbaru',
    lastCycles: '{count} siklus trading terakhir',
    noDecisionsYet: 'Belum Ada Keputusan',
    aiDecisionsWillAppear: 'Keputusan trading AI akan muncul di sini',
    cycle: 'Siklus',
    success: 'Berhasil',
    failed: 'Gagal',
    aiThinking: 'Rantai Pemikiran AI',
    collapse: 'Tutup',
    expand: 'Buka',

    // Equity Chart
    accountEquityCurve: 'Kurva Ekuitas Akun',
    noHistoricalData: 'Tidak Ada Data Historis',
    dataWillAppear:
      'Kurva ekuitas akan muncul setelah beberapa siklus berjalan',
    initialBalance: 'Saldo Awal',
    currentEquity: 'Ekuitas Saat Ini',
    historicalCycles: 'Siklus Historis',
    displayRange: 'Rentang Tampilan',
    recent: 'Terbaru',
    allData: 'Semua Data',
    cycles: 'Siklus',

    // Comparison Chart
    dataPoints: 'Titik Data',
    currentGap: 'Selisih Saat Ini',
    count: '{count} poin',

    // TradingView Chart
    fullscreen: 'Layar Penuh',

    // Competition Page
    aiCompetition: 'Kompetisi AI',
    traders: 'trader',
    liveBattle: 'Pertarungan Langsung',
    leader: 'Pemimpin',
    leaderboard: 'Papan Peringkat',
    live: 'LIVE',
    performanceComparison: 'Perbandingan Performa',
    realTimePnL: 'L/R Realtime %',
    headToHead: 'Pertarungan Langsung',
    leadingBy: 'Unggul {gap}%',
    behindBy: 'Tertinggal {gap}%',
    equity: 'Ekuitas',
    pnl: 'L/R',
    pos: 'Pos',

    // AI Traders Management
    aiModels: 'Model AI',
    exchanges: 'Bursa',
    createTrader: 'Buat Trader',
    configured: 'Terkonfigurasi',
    currentTraders: 'Trader Saat Ini',
    noTraders: 'Tidak Ada Trader AI',
    createFirstTrader: 'Buat trader AI pertama Anda untuk memulai',
    dashboardEmptyTitle: 'Mari Mulai!',
    dashboardEmptyDescription:
      'Buat trader AI pertama Anda untuk mengotomatisasi strategi trading. Hubungkan bursa, pilih model AI, dan mulai trading dalam hitungan menit!',
    goToTradersPage: 'Buat Trader Pertama Anda',
    configureModelsFirst: 'Silakan konfigurasi model AI terlebih dahulu',
    configureExchangesFirst: 'Silakan konfigurasi bursa terlebih dahulu',
    configureModelsAndExchangesFirst:
      'Silakan konfigurasi model AI dan bursa terlebih dahulu',
    modelNotConfigured: 'Model yang dipilih belum dikonfigurasi',
    exchangeNotConfigured: 'Bursa yang dipilih belum dikonfigurasi',
    confirmDeleteTrader: 'Apakah Anda yakin ingin menghapus trader ini?',
    status: 'Status',
    start: 'Mulai',
    stop: 'Berhenti',
    traderName: 'Nama Trader',
    cancel: 'Batal',
    create: 'Buat',
    aiScanInterval: 'Interval Keputusan AI (menit)',
    scanIntervalRecommend: 'Disarankan: 15-30 menit',
    enabled: 'Aktif',
    save: 'Simpan',

    // TraderConfigModal
    saving: 'Menyimpan...',
    saveFailed: 'Gagal menyimpan',
    editTraderConfig: 'Edit Konfigurasi Trader',
    selectStrategyAndConfigParams:
      'Pilih Strategi dan Konfigurasi Parameter Dasar',
    basicConfig: 'Konfigurasi Dasar',
    traderNameRequired: 'Nama Trader *',
    enterTraderNamePlaceholder: 'Masukkan nama trader',
    aiModelRequired: 'Model AI *',
    exchangeRequired: 'Bursa *',
    selectTradingStrategy: 'Pilih Strategi Trading',
    useStrategy: 'Gunakan Strategi',
    noStrategyManual: '-- Tanpa Strategi (Konfigurasi Manual) --',
    strategyActive: ' (Aktif)',
    strategyDefault: ' [Default]',
    noStrategyHint:
      'Belum ada strategi, buat di Strategy Studio terlebih dahulu',
    strategyDetails: 'Detail Strategi',
    activating: 'Mengaktifkan',
    coinSource: 'Sumber Koin',
    marginLimit: 'Batas Margin',
    tradingParams: 'Parameter Trading',
    marginMode: 'Mode Margin',
    crossMargin: 'Cross Margin',
    isolatedMargin: 'Isolated Margin',
    competitionDisplay: 'Tampilkan di Kompetisi',
    show: 'Tampilkan',
    hide: 'Sembunyikan',
    hiddenInCompetition:
      'Trader ini tidak akan ditampilkan di halaman kompetisi saat disembunyikan',
    fetching: 'Mengambil...',
    autoFetchBalanceInfo:
      'Sistem akan otomatis mengambil ekuitas akun Anda sebagai saldo awal',
    editTrader: 'Simpan Perubahan',
    createTraderButton: 'Buat Trader',

    // AI Model Configuration
    apiKey: 'API Key',
    enterAPIKey: 'Masukkan API Key',

    // Exchange Configuration
    secretKey: 'Secret Key',
    walletAddress: 'Alamat Wallet',
    user: 'Pengguna',

    // Trader Configuration
    loading: 'Memuat...',

    // AI Traders Page - Additional
    noModelsConfigured: 'Belum ada model AI yang dikonfigurasi',
    noExchangesConfigured: 'Belum ada bursa yang dikonfigurasi',
    information: 'Informasi',
    editAIModel: 'Edit Model AI',
    addAIModel: 'Tambah Model AI',
    confirmDeleteModel:
      'Apakah Anda yakin ingin menghapus konfigurasi model AI ini?',
    cannotDeleteModelInUse:
      'Tidak dapat menghapus model AI ini karena sedang digunakan oleh trader',
    tradersUsing: 'Trader yang menggunakan konfigurasi ini',
    pleaseDeleteTradersFirst:
      'Silakan hapus atau konfigurasi ulang trader ini terlebih dahulu',
    selectModel: 'Pilih Model AI',
    customBaseURL: 'Base URL (Opsional)',
    customBaseURLPlaceholder:
      'URL base API kustom, misal: https://api.openai.com/v1',
    leaveBlankForDefault: 'Kosongkan untuk menggunakan alamat API default',
    modelConfigInfo1:
      '• Untuk API resmi, hanya API Key yang diperlukan, biarkan kolom lain kosong',
    modelConfigInfo2:
      '• Base URL dan Nama Model kustom hanya diperlukan untuk proxy pihak ketiga',
    modelConfigInfo3: '• API Key dienkripsi dan disimpan dengan aman',
    defaultModel: 'Model default',
    kimiApiNote:
      'Kimi memerlukan API Key dari situs internasional (moonshot.ai), key region China tidak kompatibel',
    leaveBlankForDefaultModel: 'Kosongkan untuk menggunakan model default',
    customModelName: 'Nama Model (Opsional)',
    customModelNamePlaceholder: 'misal: deepseek-chat, qwen3-max, gpt-4o',
    saveConfig: 'Simpan Konfigurasi',
    confirmDeleteExchange:
      'Apakah Anda yakin ingin menghapus konfigurasi bursa ini?',
    cannotDeleteExchangeInUse:
      'Tidak dapat menghapus bursa ini karena sedang digunakan oleh trader',
    edit: 'Edit',

    // Error Messages
    createTraderFailed: 'Gagal membuat trader',
    getTraderConfigFailed: 'Gagal mendapatkan konfigurasi trader',
    modelConfigNotExist: 'Konfigurasi model tidak ada atau tidak diaktifkan',
    exchangeConfigNotExist: 'Konfigurasi bursa tidak ada atau tidak diaktifkan',
    updateTraderFailed: 'Gagal memperbarui trader',
    deleteTraderFailed: 'Gagal menghapus trader',
    operationFailed: 'Operasi gagal',
    deleteConfigFailed: 'Gagal menghapus konfigurasi',
    modelNotExist: 'Model tidak ada',
    saveConfigFailed: 'Gagal menyimpan konfigurasi',
    exchangeNotExist: 'Bursa tidak ada',
    deleteExchangeConfigFailed: 'Gagal menghapus konfigurasi bursa',

    // Login & Register
    login: 'Masuk',
    register: 'Daftar',
    email: 'Email',
    password: 'Kata Sandi',
    confirmPassword: 'Konfirmasi Kata Sandi',
    emailPlaceholder: 'email@anda.com',
    passwordPlaceholder: 'Masukkan kata sandi',
    passwordRuleMinLength: 'Minimal 8 karakter',
    passwordRuleUppercase: 'Minimal 1 huruf besar',
    passwordRuleLowercase: 'Minimal 1 huruf kecil',
    passwordRuleNumber: 'Minimal 1 angka',
    passwordRuleSpecial: 'Minimal 1 karakter khusus (@#$%!&*?)',
    passwordRuleMatch: 'Kata sandi cocok',
    passwordNotMeetRequirements:
      'Kata sandi tidak memenuhi persyaratan keamanan',
    loginButton: 'Masuk',
    registerButton: 'Daftar',
    back: 'Kembali',
    forgotPassword: 'Lupa kata sandi?',
    forgotAccount: 'Lupa akun?',
    resetPassword: 'Reset Kata Sandi',
    resetPasswordTitle: 'Reset kata sandi Anda',
    newPassword: 'Kata Sandi Baru',
    backToLogin: 'Kembali ke Login',
    resetPasswordCliIntro:
      'Demi keamanan, pemulihan kata sandi tidak lagi tersedia dari browser. Jalankan perintah ini di server tempat NOFX dipasang:',
    resetPasswordCliSecurityNote:
      'Ini memerlukan akses shell ke server, sehingga akun Anda tetap aman bahkan saat NOFX terekspos ke internet.',
    resetAccountCliIntro:
      'Untuk menghapus semua data dan memulai dari awal, jalankan perintah ini di server tempat NOFX dipasang:',
    copy: 'Salin',
    loginFailed: 'Gagal masuk. Periksa email dan kata sandi Anda.',
    registrationFailed: 'Gagal mendaftar. Silakan coba lagi.',
    sessionExpired: 'Sesi berakhir, silakan masuk kembali',
    medium: 'Sedang',
    strong: 'Kuat',

    // Landing Page
    features: 'Fitur',
    howItWorks: 'Cara Kerja',
    community: 'Komunitas',
    language: 'Bahasa',
    loggedInAs: 'Masuk sebagai',
    exitLogin: 'Keluar',
    signIn: 'Masuk',
    registrationClosed: 'Pendaftaran Ditutup',
    registrationClosedMessage:
      'Pendaftaran pengguna saat ini dinonaktifkan. Silakan hubungi administrator untuk akses.',
    secureDesc:
      'Keamanan tingkat enterprise, kontrol penuh atas dana dan strategi trading Anda.',
    step1Title: 'Clone Repository GitHub',
    step2Title: 'Konfigurasi Lingkungan',
    links: 'Tautan',
    candidateCoins: 'Koin Kandidat',

    // FAQ Page

    // Web Crypto Environment Check
    environmentCheck: {
      button: 'Periksa Lingkungan Aman',
      checking: 'Memeriksa...',
      description:
        'Memverifikasi otomatis apakah konteks browser ini memungkinkan Web Crypto sebelum memasukkan kunci sensitif.',
      secureTitle: 'Konteks aman terdeteksi',
      secureDesc:
        'API Web Crypto tersedia. Anda dapat melanjutkan memasukkan rahasia dengan enkripsi diaktifkan.',
      insecureTitle: 'Konteks tidak aman terdeteksi',
      insecureDesc:
        'Halaman ini tidak berjalan melalui HTTPS atau origin localhost tepercaya.',
      tipsTitle: 'Cara memperbaiki:',
      tipHTTPS: 'Sajikan dasbor melalui HTTPS dengan sertifikat valid.',
      tipLocalhost:
        'Selama pengembangan, buka aplikasi via http://localhost atau 127.0.0.1.',
      tipIframe:
        'Hindari menyematkan aplikasi dalam iframe HTTP yang tidak aman.',
      unsupportedTitle: 'Browser tidak mengekspos Web Crypto',
      unsupportedDesc:
        'Buka NOFX melalui HTTPS (atau http://localhost saat pengembangan).',
      summary: 'Origin saat ini: {origin} · Protokol: {protocol}',
      disabledTitle: 'Enkripsi transport dinonaktifkan',
      disabledDesc:
        'Enkripsi transport sisi server dinonaktifkan. API key akan ditransmisikan dalam plaintext. Aktifkan TRANSPORT_ENCRYPTION=true untuk keamanan yang lebih baik.',
    },
    twoStageKey: {
      title: 'Input Private Key Dua Tahap',
      stage1Description: 'Masukkan {length} karakter pertama private key Anda',
      stage2Description: 'Masukkan {length} karakter sisa private key Anda',
      stage1InputLabel: 'Bagian Pertama',
      stage2InputLabel: 'Bagian Kedua',
      characters: 'karakter',
      processing: 'Memproses...',
      nextButton: 'Lanjut',
      cancelButton: 'Batal',
      backButton: 'Kembali',
      encryptButton: 'Enkripsi & Kirim',
      obfuscationCopied: 'Data pengacak disalin ke clipboard',
      obfuscationInstruction:
        'Tempelkan sesuatu yang lain untuk membersihkan clipboard, lalu lanjutkan',
      obfuscationManual: 'Diperlukan pengacakan manual',
    },
    errors: {
      privatekeyIncomplete: 'Masukkan minimal {expected} karakter',
      privatekeyInvalidFormat:
        'Format private key tidak valid (harus 64 karakter heksadesimal)',
      privatekeyObfuscationFailed: 'Pengacakan clipboard gagal',
    },
    positionHistory: {
      title: 'Riwayat Posisi',
      loading: 'Memuat riwayat posisi...',
      noHistory: 'Tidak Ada Riwayat Posisi',
      noHistoryDesc: 'Posisi yang ditutup akan muncul di sini setelah trading.',
      showingPositions: 'Menampilkan {count} dari {total} posisi',
      totalPnL: 'Total L/R',
      totalTrades: 'Total Trading',
      winLoss: 'Menang: {win} / Kalah: {loss}',
      winRate: 'Win Rate',
      profitFactor: 'Profit Factor',
      profitFactorDesc: 'Total Profit / Total Loss',
      plRatio: 'Rasio L/R',
      plRatioDesc: 'Rata-rata Menang / Rata-rata Kalah',
      sharpeRatio: 'Rasio Sharpe',
      sharpeRatioDesc: 'Return yang Disesuaikan Risiko',
      maxDrawdown: 'Drawdown Maksimum',
      avgWin: 'Rata-rata Menang',
      avgLoss: 'Rata-rata Kalah',
      netPnL: 'L/R Bersih',
      netPnLDesc: 'Setelah Biaya',
      fee: 'Biaya',
      trades: 'Trading',
      avgPnL: 'Rata-rata L/R',
      symbolPerformance: 'Performa Simbol',
      symbol: 'Simbol',
      allSymbols: 'Semua Simbol',
      side: 'Arah',
      all: 'Semua',
      sort: 'Urutkan',
      latestFirst: 'Terbaru Dulu',
      oldestFirst: 'Terlama Dulu',
      highestPnL: 'L/R Tertinggi',
      lowestPnL: 'L/R Terendah',
      entry: 'Masuk',
      exit: 'Keluar',
      qty: 'Jml',
      value: 'Nilai',
      lev: 'Lev',
      pnl: 'L/R',
      duration: 'Durasi',
      closedAt: 'Ditutup Pada',
    },

    // Data Page

    // Strategy Market Page
    strategyMarket: {
      title: 'PASAR STRATEGI',
      subtitle: 'DATABASE STRATEGI GLOBAL',
      description:
        'Temukan, analisis, dan kloning algoritma trading berperforma tinggi',
      search: 'CARI PARAMETER...',
      all: 'SEMUA PROTOKOL',
      recent: 'TERBARU',
      createdAt: 'TIMESTAMP',
      copied: 'DISALIN',
      indicators: 'INDIKATOR',
      maxPositions: 'BATAS_POS',
      loading: 'MENGINISIALISASI...',
    },

    // Strategy Studio Page

    // Metric Tooltip
    metricTooltip: {
      formula: 'Formula',
    },

    // Login Required Overlay
    loginRequired: {
      title: 'AKSES SISTEM DITOLAK',
      accessDenied: 'AKSES DITOLAK',
      subtitleWithFeature:
        'Modul "{featureName}" memerlukan hak akses lebih tinggi',
      subtitleDefault: 'Otorisasi diperlukan untuk modul ini',
      description:
        'Inisialisasi protokol autentikasi untuk membuka kemampuan sistem penuh: konfigurasi Trader AI dan aliran data Pasar Strategi.',
      benefit1: 'Kontrol Trader AI',
      benefit2: 'Pasar Strategi HFT',
      benefit4: 'Visualisasi Sistem Penuh',
      loginButton: 'JALANKAN LOGIN',
      registerButton: 'DAFTAR ID BARU',
      abort: 'BATALKAN',
    },

    // Advanced Chart

    // Chart With Orders

    // Comparison Chart
    comparisonChart: {
      '1d': '1H',
      '3d': '3H',
      '7d': '7H',
      '30d': '30H',
      all: 'Semua',
    },

    traderDashboard: {
      connectionFailed: 'Koneksi Gagal',
      connectionFailedDesc: 'Silakan periksa apakah layanan backend berjalan.',
      retry: 'Coba Lagi',
      confirmClosePosition: 'Yakin ingin menutup posisi {symbol} {side}?',
      confirmClose: 'Konfirmasi Tutup',
      confirm: 'Konfirmasi',
      cancel: 'Batal',
      positionClosed: 'Posisi berhasil ditutup',
      closeFailed: 'Gagal menutup posisi',
      closeAll: 'Tutup Semua',
      confirmCloseAllPositions:
        'Tutup SEMUA {count} posisi terbuka dengan harga pasar?',
      allPositionsClosed: 'Semua posisi ditutup',
      closeAllPartial: '{failed} dari {count} posisi gagal ditutup',
      action: 'Aksi',
      entry: 'Entry',
      mark: 'Mark',
      qty: 'Qty',
      value: 'Nilai',
      lev: 'Lev.',
      uPnL: 'uPnL',
      liq: 'Liq.',
      closePosition: 'Tutup Posisi',
      close: 'Tutup',
      showingPositions: 'Menampilkan {shown} dari {total} posisi',
      perPage: 'Per halaman',
      accountFetchFailed:
        'DATA_FETCH::FAILED — Data akun tidak tersedia, periksa koneksi',
      positionsFetchFailed: 'Data posisi tidak tersedia',
      decisionsFetchFailed: 'Data keputusan tidak tersedia',
    },

    aiTradersToast: {
      creating: 'Membuat...',
      created: 'Berhasil dibuat',
      saving: 'Menyimpan...',
      saved: 'Berhasil disimpan',
      saveFailed: 'Gagal menyimpan',
      deleted: 'Berhasil dihapus',
      stopping: 'Menghentikan...',
      stopped: 'Dihentikan',
      starting: 'Memulai...',
      started: 'Dimulai',
      configUpdated: 'Konfigurasi diperbarui',
      showInCompetition: 'Ditampilkan di kompetisi',
      hideInCompetition: 'Disembunyikan dari kompetisi',
      modelConfigUpdated: 'Konfigurasi model diperbarui',
      exchangeDeleted: 'Akun exchange dihapus',
      exchangeConfigUpdated: 'Konfigurasi exchange diperbarui',
      exchangeCreated: 'Akun exchange dibuat',
    },

    modelConfig: {
      selectModel: 'Pilih Model',
      configure: 'Konfigurasi',
      configureApi: 'Konfigurasi API',
      chooseProvider: 'Pilih Penyedia AI Anda',
      selectAiModel: 'Pilih Model AI',
      exportKey: 'Ekspor private key dari MetaMask, Rabby, dll.',
      back: 'Kembali',
      modelsConfigured: 'Model dengan lencana emas sudah dikonfigurasi',
      getApiKey: 'Dapatkan API Key',
      validating: 'Memvalidasi...',
      walletAddress: 'Alamat Wallet',
    },


    telegram: {
      botSetup: 'Setup Telegram Bot',
      createBot: 'Buat Bot',
      bindAccount: 'Hubungkan Akun',
      done: 'Selesai',
      invalidTokenFormat:
        'Format Bot Token tidak valid. Seharusnya "angka:alfanumerik"',
      tokenSaved: 'Bot Token tersimpan, menunggu binding',
      saveFailed: 'Gagal menyimpan, silakan periksa token',
      unbound: 'Akun Telegram terputus',
      unbindFailed: 'Gagal memutuskan',
      step1Title: 'Langkah 1: Buat Bot di Telegram',
      step1Desc1: 'Buka Telegram, cari',
      step1Desc2: 'Kirim',
      step1Desc2Suffix: 'perintah',
      step1Desc3: 'Ikuti petunjuk untuk mengatur nama dan username bot',
      step1Desc4: 'BotFather akan mengembalikan Token, salin itu',
      openBotFather: 'Buka @BotFather',
      pasteToken: 'Tempel Bot Token',
      tokenFormat: 'Format: angka:alfanumerik, mis. 123456789:ABCdef...',
      selectAiModel: 'Pilih Model AI (opsional)',
      noEnabledModels:
        'Belum ada model aktif. Konfigurasi di AI Models terlebih dahulu.',
      autoSelect: '— Pilih otomatis (disarankan)',
      autoUseEnabled: 'Kosongkan untuk otomatis menggunakan model aktif',
      savingToken: 'Menyimpan...',
      saveAndContinue: 'Simpan & Lanjut',
      step2Title: 'Langkah 2: Kirim /start ke Bot Anda',
      step2Desc1: 'Cari Bot yang baru dibuat di Telegram',
      step2Desc2: 'Klik Start atau kirim',
      step2Desc3: 'Bot akan otomatis terhubung ke akun Anda',
      currentToken: 'Token Saat Ini',
      waitingForStart:
        'Menunggu Anda mengirim /start... Refresh halaman setelah mengirim',
      reconfigureToken: 'Konfigurasi Ulang Token',
      bindSuccess: 'Berhasil terhubung!',
      noStartReceived:
        'Belum menerima /start. Silakan kirim /start ke Bot Anda terlebih dahulu',
      checkFailed: 'Pemeriksaan gagal',
      checkStatus: 'Periksa Status',
      botActive: 'Telegram Bot Aktif!',
      botActiveDesc:
        'Anda sekarang dapat mengontrol sistem trading melalui bahasa alami di Telegram',
      supportedCommands: 'Perintah yang Didukung',
      cmdHelp: 'Tampilkan semua perintah',
      cmdStatus: 'Tampilkan status trader',
      cmdNaturalLang: 'Bahasa alami',
      cmdStartStop: 'Mulai/hentikan trader',
      cmdControl: 'Kontrol bahasa alami',
      cmdPositions: 'Lihat posisi',
      cmdPositionsDesc: 'Kueri posisi real-time',
      cmdStrategy: 'Konfigurasi strategi',
      cmdStrategyDesc: 'Ubah strategi trading',
      unbinding: 'Memutuskan...',
      unbindAccount: 'Putuskan Akun',
      aiModelLabel: 'Model AI (untuk bahasa alami)',
      aiModelAutoSelect: '— Pilih otomatis',
      modelUpdated: 'Model AI diperbarui',
      modelUpdateFailed: 'Gagal memperbarui',
      save: 'Simpan',
      loading: 'Memuat...',
    },

    traderConfigView: {
      traderConfig: 'Konfigurasi Trader',
      configInfo: 'Detail konfigurasi {name}',
      running: 'Berjalan',
      stopped: 'Berhenti',
      basicInfo: 'Informasi Dasar',
      traderName: 'Nama Trader',
      aiModel: 'Model AI',
      exchange: 'Exchange',
      initialBalance: 'Saldo Awal',
      marginMode: 'Mode Margin',
      crossMargin: 'Cross',
      isolatedMargin: 'Isolated',
      scanInterval: '{minutes} menit',
      scanIntervalLabel: 'Interval Scan',
      strategyUsed: 'Strategi Digunakan',
      strategyName: 'Nama Strategi',
      close: 'Tutup',
      yes: 'Ya',
      no: 'Tidak',
    },
  },
}

export function t(
  key: string,
  lang: Language,
  params?: Record<string, string | number>
): string {
  // Handle nested keys like 'twoStageKey.title'
  const keys = key.split('.')
  let value: any = translations[lang]

  for (const k of keys) {
    value = value?.[k]
  }

  let text = typeof value === 'string' ? value : key

  // Replace parameters like {count}, {gap}, etc.
  if (params) {
    Object.entries(params).forEach(([param, value]) => {
      text = text.replace(`{${param}}`, String(value))
    })
  }

  return text
}
