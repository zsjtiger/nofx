// zh UI strings. Keys are read with t() in translations.ts.
const zh = {
  leadPnL: '领先收益',
  loadingChartData: '正在加载图表数据...',
  loadingError: '加载数据失败',
  loggingIn: '登录中...',
  // Header
  appTitle: 'Crux',
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
  footerTitle: 'Crux - AI交易系统',
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
  kimiApiNote: 'Kimi 需要从国际站申请 API Key (moonshot.ai)，中国区 Key 不通用',
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
    '出于安全考虑，密码找回不再通过浏览器进行。请在部署 Crux 的服务器上运行以下命令：',
  resetPasswordCliSecurityNote:
    '该操作需要服务器的 shell 访问权限，因此即使 Crux 暴露在公网上，你的账户依然安全。',
  resetAccountCliIntro:
    '如需清空所有数据并重新开始，请在部署 Crux 的服务器上运行以下命令：',
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
      '请通过 HTTPS 或本机 localhost 访问 Crux，并避免嵌入不安全 iframe/反向代理，以符合浏览器的 Web Crypto 规则。',
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
} as const

export default zh
