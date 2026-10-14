// id UI strings. Keys are read with t() in translations.ts.
const id = {
  leadPnL: 'PnL Teratas',
  loadingChartData: 'Memuat data grafik...',
  loadingError: 'Gagal memuat data',
  loggingIn: 'Sedang masuk...',
  // Header
  appTitle: 'Crux',
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
  footerTitle: 'Crux - Sistem Trading AI',
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
  dataWillAppear: 'Kurva ekuitas akan muncul setelah beberapa siklus berjalan',
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
  noStrategyHint: 'Belum ada strategi, buat di Strategy Studio terlebih dahulu',
  strategyDetails: 'Detail Strategi',
  activating: 'Mengaktifkan',
  coinSource: 'Sumber Koin',
  marginLimit: 'Batas Margin',
  tradingParams: 'Parameter Trading',
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
  passwordNotMeetRequirements: 'Kata sandi tidak memenuhi persyaratan keamanan',
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
    'Demi keamanan, pemulihan kata sandi tidak lagi tersedia dari browser. Jalankan perintah ini di server tempat Crux dipasang:',
  resetPasswordCliSecurityNote:
    'Ini memerlukan akses shell ke server, sehingga akun Anda tetap aman bahkan saat Crux terekspos ke internet.',
  resetAccountCliIntro:
    'Untuk menghapus semua data dan memulai dari awal, jalankan perintah ini di server tempat Crux dipasang:',
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
      'Buka Crux melalui HTTPS (atau http://localhost saat pengembangan).',
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
    scanInterval: '{minutes} menit',
    scanIntervalLabel: 'Interval Scan',
    strategyUsed: 'Strategi Digunakan',
    strategyName: 'Nama Strategi',
    close: 'Tutup',
    yes: 'Ya',
    no: 'Tidak',
  },
} as const

export default id
