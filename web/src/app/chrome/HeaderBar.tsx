import { useState, useEffect, useRef } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { motion, AnimatePresence } from 'framer-motion'
import { Menu, X, ChevronDown, Settings } from 'lucide-react'
import { t, type Language } from '../../shared/i18n/translations'
import {
  DEFAULT_LANGUAGE,
  SWITCHABLE_LANGUAGES,
} from '../../shared/contexts/LanguageContext'
import {
  getCurrentPageForPath,
  ROUTES,
  type Page,
} from '../../shared/routes/paths'

interface HeaderBarProps {
  onLoginClick?: () => void
  isLoggedIn?: boolean
  isHomePage?: boolean
  currentPage?: Page
  language?: Language
  onLanguageChange?: (lang: Language) => void
  user?: { email: string } | null
  onLogout?: () => void
  onPageChange?: (page: Page) => void
  onLoginRequired?: (featureName: string) => void
}

export default function HeaderBar({
  isLoggedIn = false,
  isHomePage = false,
  currentPage,
  language = DEFAULT_LANGUAGE,
  onLanguageChange,
  user,
  onLogout,
  onPageChange,
  onLoginRequired,
}: HeaderBarProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const [userDropdownOpen, setUserDropdownOpen] = useState(false)

  const userDropdownRef = useRef<HTMLDivElement>(null)
  const resolvedCurrentPage =
    currentPage ?? getCurrentPageForPath(location.pathname)

  const navigateInApp = (path: string) => {
    navigate(path)
  }

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        userDropdownRef.current &&
        !userDropdownRef.current.contains(event.target as Node)
      ) {
        setUserDropdownOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [])

  return (
    <nav className="fixed top-0 w-full z-50 header-bar">
      <div className="flex items-center justify-between h-16 px-4 sm:px-6 max-w-[1920px] mx-auto">
        {/* Logo - Always go to home page */}
        <div
          onClick={() => {
            navigateInApp(ROUTES.home)
          }}
          className="flex items-center gap-2 hover:opacity-80 transition-opacity cursor-pointer"
        >
          <span
            className="flex items-center justify-center w-8 h-8 rounded-md overflow-hidden shrink-0"
            style={{
              background: '#fff',
              border: '1px solid rgba(26,24,19,0.12)',
            }}
          >
            <img
              src="/icons/crux-logo.png"
              alt="Crux Logo"
              className="w-8 h-8"
            />
          </span>
          <span className="text-lg font-bold text-crux-gold tracking-wide">
            Crux
          </span>
        </div>

        {/* Desktop Menu */}
        <div className="hidden min-w-0 flex-1 items-center justify-between gap-3 ml-6 md:flex">
          {/* Left Side - Navigation Tabs - Always show all tabs */}
          <div className="flex min-w-0 flex-1 items-center gap-2">
            <div className="hidden shrink-0 items-center gap-0.5 lg:flex xl:gap-1 2xl:gap-2">
              {/* Navigation tabs configuration */}
              {(() => {
                // Define all navigation tabs
                const navTabs: {
                  page: Page
                  path: string
                  label: string
                  requiresAuth: boolean
                  badge?: string
                  hidden?: boolean
                }[] = [
                  {
                    page: 'traders',
                    path: ROUTES.traders,
                    label: t('configNav', language),
                    requiresAuth: true,
                  },
                  {
                    page: 'trader',
                    path: ROUTES.dashboard,
                    label: t('dashboardNav', language),
                    requiresAuth: true,
                  },
                  {
                    page: 'strategy',
                    path: ROUTES.strategy,
                    label: t('strategyNav', language),
                    requiresAuth: true,
                  },
                  {
                    page: 'competition',
                    path: ROUTES.competition,
                    label: t('realtimeNav', language),
                    requiresAuth: true,
                  },
                  {
                    page: 'intel',
                    path: ROUTES.intel,
                    label: t('intelNav', language),
                    requiresAuth: true,
                  },
                  {
                    page: 'faq',
                    path: ROUTES.faq,
                    label: t('faqNav', language),
                    requiresAuth: false,
                  },
                ]

                const handleNavClick = (tab: (typeof navTabs)[0]) => {
                  // If requires auth and not logged in, show login prompt
                  if (tab.requiresAuth && !isLoggedIn) {
                    onLoginRequired?.(tab.label)
                    return
                  }
                  // Navigate normally
                  if (onPageChange) {
                    onPageChange(tab.page)
                  }
                  navigateInApp(tab.path)
                }

                return navTabs
                  .filter((tab) => !tab.hidden)
                  .map((tab) => (
                    <button
                      key={tab.page}
                      onClick={() => handleNavClick(tab)}
                      className={`shrink-0 whitespace-nowrap text-xs 2xl:text-sm font-bold transition-all duration-300 relative focus:outline-2 focus:outline-yellow-500 px-1.5 xl:px-2 2xl:px-3 py-2 rounded-lg
                    ${resolvedCurrentPage === tab.page ? 'text-crux-gold' : 'text-crux-text-muted hover:text-crux-gold'}`}
                    >
                      {resolvedCurrentPage === tab.page && (
                        <span className="absolute inset-0 rounded-lg bg-crux-gold/15 -z-10" />
                      )}
                      {tab.label}
                      {tab.badge && (
                        <span className="ml-1 text-[10px] px-1.5 py-0.5 rounded-full bg-crux-gold/20 text-crux-gold font-semibold uppercase align-top relative -top-1">
                          {tab.badge}
                        </span>
                      )}
                    </button>
                  ))
              })()}
            </div>
            {/* Dashboard context slot — terminal selector + status portals in here */}
            <div
              id="dash-header-slot"
              className="hidden min-w-0 flex-1 items-center justify-center overflow-hidden lg:flex"
            />
          </div>

          {/* Right Side - Social Links and User Actions */}
          <div className="flex shrink-0 items-center gap-2">
            {/* Language switch */}
            <div className="flex items-center gap-0.5 rounded-lg border border-[rgba(26,24,19,0.14)] p-0.5">
              {SWITCHABLE_LANGUAGES.map((code) => (
                <button
                  key={code}
                  type="button"
                  onClick={() => onLanguageChange?.(code)}
                  aria-pressed={language === code}
                  className={`px-2 py-1 rounded text-xs font-semibold transition-all ${
                    language === code
                      ? 'bg-crux-gold/15 text-crux-gold'
                      : 'text-crux-text-muted hover:text-crux-text'
                  }`}
                >
                  {code === 'zh' ? '中文' : 'EN'}
                </button>
              ))}
            </div>
            {/* User Info and Actions */}
            {isLoggedIn && user ? (
              <div className="flex items-center gap-3">
                {/* User Info with Dropdown */}
                <div className="relative" ref={userDropdownRef}>
                  <button
                    onClick={() => setUserDropdownOpen(!userDropdownOpen)}
                    className="flex items-center gap-2 px-3 py-2 rounded transition-colors bg-crux-bg-lighter border border-crux-gold/20 hover:bg-[rgba(26,24,19,0.06)]"
                  >
                    <div className="w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold bg-crux-gold text-white">
                      {user.email[0].toUpperCase()}
                    </div>
                    <span className="hidden text-sm text-crux-text-muted 2xl:inline">
                      {user.email}
                    </span>
                    <ChevronDown className="w-4 h-4 text-crux-text-muted" />
                  </button>

                  {userDropdownOpen && (
                    <div className="absolute right-0 top-full mt-2 w-48 rounded-lg shadow-lg overflow-hidden z-50 bg-crux-bg-lighter border border-crux-gold/20">
                      <div className="px-3 py-2 border-b border-crux-gold/20">
                        <div className="text-xs text-crux-text-muted">
                          {t('loggedInAs', language)}
                        </div>
                        <div className="text-sm font-medium text-crux-text-muted">
                          {user.email}
                        </div>
                      </div>
                      <button
                        onClick={() => {
                          navigateInApp(ROUTES.settings)
                          setUserDropdownOpen(false)
                        }}
                        className="w-full flex items-center gap-2 px-3 py-2 text-sm transition-colors hover:bg-[rgba(26,24,19,0.06)] text-crux-text-muted hover:text-crux-text"
                      >
                        <Settings className="w-3.5 h-3.5" />
                        Settings
                      </button>

                      {onLogout && (
                        <button
                          onClick={() => {
                            onLogout()
                            setUserDropdownOpen(false)
                          }}
                          className="w-full px-3 py-2 text-sm font-semibold transition-colors hover:opacity-80 text-center bg-crux-danger/20 text-crux-danger"
                        >
                          {t('exitLogin', language)}
                        </button>
                      )}
                    </div>
                  )}
                </div>
              </div>
            ) : (
              /* Show login/register buttons when not logged in and not on login/register pages */
              resolvedCurrentPage !== 'login' &&
              resolvedCurrentPage !== 'register' && (
                <div className="flex items-center gap-3">
                  <button
                    type="button"
                    onClick={() => navigateInApp(ROUTES.login)}
                    className="px-3 py-2 text-sm font-medium transition-colors rounded text-crux-text-muted hover:text-crux-text"
                  >
                    {t('signIn', language)}
                  </button>
                </div>
              )
            )}

            {/* Language switcher removed — the product UI is English-only. */}
          </div>
        </div>

        {/* Mobile Menu Button */}
        <motion.button
          onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
          className="shrink-0 text-crux-text-muted hover:text-crux-text lg:hidden"
          aria-label={mobileMenuOpen ? 'Close navigation' : 'Open navigation'}
          whileTap={{ scale: 0.9 }}
        >
          {mobileMenuOpen ? (
            <X className="w-6 h-6" />
          ) : (
            <Menu className="w-6 h-6" />
          )}
        </motion.button>
      </div>

      {/* Mobile Menu Overlay */}
      <AnimatePresence>
        {mobileMenuOpen && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="fixed inset-0 z-40 bg-black/90 backdrop-blur-xl lg:hidden"
            style={{ top: '64px' }} // Below header
          >
            <motion.div
              initial={{ y: -20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.1, duration: 0.3 }}
              className="flex flex-col h-[calc(100vh-64px)] overflow-y-auto px-6 py-8"
            >
              {/* Navigation Links */}
              <div className="flex flex-col gap-6 mb-12">
                {(() => {
                  const navTabs: {
                    page: Page
                    path: string
                    label: string
                    requiresAuth: boolean
                    badge?: string
                    hidden?: boolean
                  }[] = [
                    {
                      page: 'data',
                      path: ROUTES.data,
                      label:
                        language === 'zh'
                          ? 'Data'
                          : language === 'id'
                            ? 'Data'
                            : 'Data',
                      requiresAuth: false,
                    },
                    {
                      page: 'traders',
                      path: ROUTES.traders,
                      label: t('configNav', language),
                      requiresAuth: true,
                    },
                    {
                      page: 'trader',
                      path: ROUTES.dashboard,
                      label: t('dashboardNav', language),
                      requiresAuth: true,
                    },
                    {
                      page: 'strategy',
                      path: ROUTES.strategy,
                      label: t('strategyNav', language),
                      requiresAuth: true,
                    },
                    {
                      page: 'competition',
                      path: ROUTES.competition,
                      label: t('realtimeNav', language),
                      requiresAuth: true,
                    },
                    {
                      page: 'intel',
                      path: ROUTES.intel,
                      label: t('intelNav', language),
                      requiresAuth: true,
                    },
                    {
                      page: 'faq',
                      path: ROUTES.faq,
                      label: t('faqNav', language),
                      requiresAuth: false,
                    },
                  ]

                  const handleMobileNavClick = (tab: (typeof navTabs)[0]) => {
                    if (tab.requiresAuth && !isLoggedIn) {
                      onLoginRequired?.(tab.label)
                      setMobileMenuOpen(false)
                      return
                    }
                    if (onPageChange) {
                      onPageChange(tab.page)
                    }
                    navigateInApp(tab.path)
                    setMobileMenuOpen(false)
                  }

                  return navTabs
                    .filter((tab) => !tab.hidden)
                    .map((tab, i) => (
                      <motion.button
                        key={tab.page}
                        initial={{ x: -20, opacity: 0 }}
                        animate={{ x: 0, opacity: 1 }}
                        transition={{ delay: 0.1 + i * 0.05 }}
                        onClick={() => handleMobileNavClick(tab)}
                        className={`text-2xl font-black tracking-tight text-left flex items-center gap-3
                        ${resolvedCurrentPage === tab.page ? 'text-crux-gold' : 'text-zinc-500'}`}
                      >
                        {resolvedCurrentPage === tab.page && (
                          <motion.div
                            layoutId="active-indicator"
                            className="w-1.5 h-1.5 rounded-full bg-crux-gold"
                          />
                        )}
                        {tab.label}
                        {tab.badge && (
                          <span className="text-[10px] px-1.5 py-0.5 rounded-full bg-crux-gold/20 text-crux-gold font-semibold uppercase align-middle relative -top-1">
                            {tab.badge}
                          </span>
                        )}
                        {tab.requiresAuth && !isLoggedIn && (
                          <span className="text-[10px] px-1.5 py-0.5 rounded border border-zinc-800 text-zinc-500 font-normal tracking-wide uppercase align-middle relative -top-1">
                            LOGIN_REQ
                          </span>
                        )}
                      </motion.button>
                    ))
                })()}

                {/* Original Page Links */}
                {isHomePage && (
                  <div className="pt-6 border-t border-white/5 space-y-4">
                    {[
                      { key: 'features', label: t('features', language) },
                      { key: 'howItWorks', label: t('howItWorks', language) },
                    ].map((item, i) => (
                      <motion.a
                        key={item.key}
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        transition={{ delay: 0.5 + i * 0.1 }}
                        href={`#${item.key === 'features' ? 'features' : 'how-it-works'}`}
                        className="block text-lg font-mono text-zinc-600 hover:text-white"
                        onClick={() => setMobileMenuOpen(false)}
                      >
                        {'>'} {item.label}
                      </motion.a>
                    ))}
                  </div>
                )}
              </div>

              {/* Bottom Actions */}
              <div className="mt-auto space-y-8">
                {/* Account (language switcher removed — English-only UI) */}
                <div className="grid grid-cols-1 gap-4">
                  {/* Auth Actions */}
                  {isLoggedIn && user ? (
                    <button
                      onClick={() => {
                        onLogout?.()
                        setMobileMenuOpen(false)
                      }}
                      className="bg-red-500/10 border border-red-500/20 text-red-500 rounded-lg font-bold text-sm hover:bg-red-500/20 transition-colors"
                    >
                      {t('exitLogin', language)}
                    </button>
                  ) : (
                    resolvedCurrentPage !== 'login' &&
                    resolvedCurrentPage !== 'register' && (
                      <button
                        type="button"
                        onClick={() => {
                          navigateInApp(ROUTES.login)
                          setMobileMenuOpen(false)
                        }}
                        className="flex items-center justify-center bg-crux-gold text-white rounded-lg font-bold text-sm hover:opacity-90 transition-colors"
                      >
                        {t('signIn', language)}
                      </button>
                    )
                  )}
                </div>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </nav>
  )
}
