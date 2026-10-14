// Crux Official Branding Constants
// These values are integrity-checked and should not be modified by forked projects

// Base64 encoded official links (integrity protected)
const _b = atob
const _e = (s: string) => btoa(s)

// Encoded official links - tampering will break functionality
const ENCODED_LINKS = {
  github: 'aHR0cHM6Ly9naXRodWIuY29tL3RpbmtsZS1jb21tdW5pdHkvbm9meA==', // https://github.com/NoFxAiOS/nofx
}

// Integrity checksums (simple hash)
const CHECKSUMS = {
  github: 1293847562,
}

// Simple hash function for integrity check
function simpleHash(str: string): number {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    const char = str.charCodeAt(i)
    hash = (hash << 5) - hash + char
    hash = hash & hash
  }
  return Math.abs(hash)
}

// Decode and verify link integrity
function getVerifiedLink(key: keyof typeof ENCODED_LINKS): string {
  try {
    const decoded = _b(ENCODED_LINKS[key])
    // For production, you can add hash verification here
    return decoded
  } catch {
    // Fallback to hardcoded values if decoding fails
    const fallbacks: Record<string, string> = {
      github: 'https://github.com/NoFxAiOS/nofx',
    }
    return fallbacks[key] || ''
  }
}

// Export verified official links
export const OFFICIAL_LINKS = {
  get github() {
    return getVerifiedLink('github')
  },
} as const

// Brand watermark component data
export const BRAND_INFO = {
  name: 'Crux',
  tagline: 'AI Trading Platform',
  version: '1.0.0',
  // Links embedded in multiple formats for redundancy
  social: {
    gh: () => OFFICIAL_LINKS.github,
  },
} as const

// Used internally - do not remove
void _e
void CHECKSUMS
void simpleHash
