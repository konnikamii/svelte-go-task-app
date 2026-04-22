import { browser } from '$app/environment'
import { onMount } from 'svelte'

const THEME_STORAGE_KEY = 'theme'

const theme = $state({
  isDark: true,
  initialized: false,
})

function applyTheme() {
  if (!browser) return

  document.documentElement.classList.toggle('dark', theme.isDark)
  document.documentElement.style.colorScheme = theme.isDark ? 'dark' : 'light'
  localStorage.setItem(THEME_STORAGE_KEY, theme.isDark ? 'dark' : 'light')
}

function initializeTheme() {
  if (!browser || theme.initialized) return

  const savedTheme = localStorage.getItem(THEME_STORAGE_KEY)
  theme.isDark = savedTheme ? savedTheme === 'dark' : true
  theme.initialized = true
  applyTheme()
}

function toggleTheme() {
  theme.isDark = !theme.isDark
  applyTheme()
}

export function useTheme() {
  onMount(() => {
    initializeTheme()
  })

  return {
    theme,
    toggleTheme,
  }
}
