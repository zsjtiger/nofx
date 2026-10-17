import { useMemo, useState } from 'react'
import { HelpCircle } from 'lucide-react'
import { DeepVoidBackground } from '../../shared/ui/DeepVoidBackground'
import { FAQSearchBar } from './FAQSearchBar'
import { FAQSidebar } from './FAQSidebar'
import { FAQContent } from './FAQContent'
import { faqCategories, faqItemSearchText } from './faqData'
import type { FAQCategory } from './faqData'

export function FAQLayout() {
  const [searchTerm, setSearchTerm] = useState('')
  const [activeItemId, setActiveItemId] = useState<string | null>(null)

  const filteredCategories = useMemo(() => {
    if (!searchTerm.trim()) return faqCategories

    const term = searchTerm.toLowerCase()
    const filtered: FAQCategory[] = []
    faqCategories.forEach((category) => {
      const matchingItems = category.items.filter((item) =>
        faqItemSearchText(item).includes(term)
      )
      if (matchingItems.length > 0) {
        filtered.push({ ...category, items: matchingItems })
      }
    })
    return filtered
  }, [searchTerm])

  const totalItems = useMemo(
    () =>
      faqCategories.reduce((sum, category) => sum + category.items.length, 0),
    []
  )

  const handleItemClick = (_categoryId: string, itemId: string) => {
    const element = document.getElementById(itemId)
    if (!element) return
    const offset = 100
    const top =
      element.getBoundingClientRect().top + window.pageYOffset - offset
    window.scrollTo({ top, behavior: 'smooth' })
  }

  return (
    <DeepVoidBackground className="py-8 pt-24" disableAnimation>
      <div className="mx-auto w-full max-w-6xl px-4 md:px-8">
        {/* page header — same strip language as the other terminal pages */}
        <div className="mb-8 border-b border-crux-gold/20 pb-6">
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div className="flex items-center gap-4">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl border border-crux-gold/30 bg-crux-bg-lighter text-crux-gold md:h-14 md:w-14">
                <HelpCircle className="h-6 w-6 md:h-7 md:w-7" />
              </div>
              <div>
                <h1 className="font-mono text-2xl font-bold tracking-tight text-crux-text md:text-3xl">
                  FAQ
                </h1>
                <p className="mt-1 font-mono text-xs uppercase tracking-[0.14em] text-crux-text-muted">
                  {totalItems} answers · wallets · launch · trading ·
                  self-hosting
                </p>
              </div>
            </div>
            <div className="w-full md:w-80">
              <FAQSearchBar
                searchTerm={searchTerm}
                onSearchChange={setSearchTerm}
              />
            </div>
          </div>
        </div>

        {/* content */}
        <div className="flex gap-8">
          <aside className="hidden w-64 flex-shrink-0 lg:block">
            <FAQSidebar
              categories={filteredCategories}
              activeItemId={activeItemId}
              onItemClick={handleItemClick}
            />
          </aside>

          <main className="min-w-0 flex-1">
            {filteredCategories.length > 0 ? (
              <FAQContent
                categories={filteredCategories}
                onActiveItemChange={setActiveItemId}
              />
            ) : (
              <div className="rounded-xl border border-crux-gold/20 bg-crux-bg-lighter py-16 text-center">
                <p className="font-mono text-sm text-crux-text-muted">
                  No matching questions for “{searchTerm}”.
                </p>
                <button
                  onClick={() => setSearchTerm('')}
                  className="mt-4 rounded-lg border border-crux-gold/30 bg-crux-gold/10 px-5 py-2 font-mono text-xs font-bold uppercase tracking-[0.12em] text-crux-gold hover:bg-crux-gold/20"
                >
                  Clear search
                </button>
              </div>
            )}
          </main>
        </div>
      </div>
    </DeepVoidBackground>
  )
}
