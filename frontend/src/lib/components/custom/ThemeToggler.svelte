<script lang="ts">
  import ThemeToggleDark from '$lib/components/svgs/ThemeToggleDark.svelte'
  import ThemeToggleLight from '$lib/components/svgs/ThemeToggleLight.svelte'
  import { useTheme } from '$lib/utils/theme.svelte'

  const { theme, toggleTheme } = useTheme()
</script>

<button
  type="button"
  aria-label="Toggle theme"
  aria-pressed={theme.isDark}
  class="relative h-6 w-12 cursor-pointer rounded-full border border-gray-400 bg-gray-400 px-[6px] py-[3px] shadow-inner shadow-gray-700 transition-all"
  onclick={toggleTheme}
>
  <ThemeToggleLight
    class="absolute left-1 size-4 translate-x-0 -translate-y-1/2 rotate-[-360deg] rounded-full bg-white opacity-100 transition-all duration-300 dark:translate-x-6 dark:rotate-0 dark:opacity-0"
    color="oklch(79.5% 0.184 86.047)"
  />
  <ThemeToggleDark
    class="absolute left-1 size-4 translate-x-0 -translate-y-1/2 rotate-[-360deg] rounded-full bg-white fill-black opacity-0 transition-all duration-300 dark:translate-x-6 dark:rotate-0 dark:opacity-100"
  />
</button>
