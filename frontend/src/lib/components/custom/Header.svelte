<script lang="ts">
  import { resolve } from '$app/paths'
  import ThemeToggler from '$lib/components/custom/ThemeToggler.svelte'
  import TaskifyLogo from '$lib/components/svgs/TaskifyLogo.svelte'

  const navItems = [
    { href: '/' as const, label: 'Home' },
    { href: '/contact' as const, label: 'Contact' },
    { href: '/login' as const, label: 'Login' },
    { href: '/register' as const, label: 'Register' },
  ]
</script>

<header class="sticky top-0 z-50 border-b border-border bg-background shadow-sm shadow-secondary backdrop-blur">
  <div class="mx-auto flex max-w-6xl items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
    <a
      href={resolve('/')}
      draggable="false"
      class="group flex items-center gap-2 text-lg font-semibold text-foreground"
    >
      <TaskifyLogo width={40} height={40} />
      <span class="text-2xl">Taskify</span>
    </a>

    <div class="flex items-center gap-3">
      <nav class="flex items-center gap-2 text-sm text-muted-foreground">
        {#each navItems as item (item.href)}
          <a href={resolve(item.href)} class="rounded-md px-3 py-2 transition hover:bg-muted hover:text-foreground">
            {item.label}
          </a>
        {/each}
      </nav>

      <ThemeToggler />
    </div>
  </div>
</header>
