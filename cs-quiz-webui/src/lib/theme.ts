const THEME_KEY = 'csquiz-theme'

export type Theme = 'light' | 'dark' | 'system'

function systemPrefersDark(): boolean {
  return window.matchMedia('(prefers-color-scheme: dark)').matches
}

export function getTheme(): Theme {
  const stored = localStorage.getItem(THEME_KEY)
  if (stored === 'light' || stored === 'dark' || stored === 'system') return stored
  return 'system'
}

export function resolveTheme(theme: Theme = getTheme()): 'light' | 'dark' {
  if (theme === 'system') return systemPrefersDark() ? 'dark' : 'light'
  return theme
}

export function applyTheme(theme: Theme): void {
  const root = document.documentElement
  const resolved = resolveTheme(theme)
  if (resolved === 'dark') root.classList.add('dark')
  else root.classList.remove('dark')
  localStorage.setItem(THEME_KEY, theme)
}

export function setTheme(theme: Theme): Theme {
  applyTheme(theme)
  return theme
}

export function toggleTheme(): Theme {
  const next: Theme = resolveTheme() === 'dark' ? 'light' : 'dark'
  applyTheme(next)
  return next
}

let mediaListener: ((e: MediaQueryListEvent) => void) | null = null

export function initTheme(): void {
  applyTheme(getTheme())
  const mq = window.matchMedia('(prefers-color-scheme: dark)')
  if (mediaListener) mq.removeEventListener('change', mediaListener)
  mediaListener = () => {
    if (getTheme() === 'system') applyTheme('system')
  }
  mq.addEventListener('change', mediaListener)
}
