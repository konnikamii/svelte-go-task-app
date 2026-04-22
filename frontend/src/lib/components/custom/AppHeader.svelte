<script lang="ts">
  import { page } from '$app/state'
  import { navItems } from '$lib/components/custom/Navbar.svelte'
  import ThemeToggler from '$lib/components/custom/ThemeToggler.svelte'
</script>

<header
  class="sticky top-0 z-50 h-[70px] border-b border-border bg-background shadow-sm shadow-secondary backdrop-blur"
>
  <div class="flex h-full items-center justify-between px-12">
    <div>{navItems.find((e) => e.href === page.url.pathname)?.message}</div>
    <div class="flex items-center gap-3">
      <ThemeToggler />
    </div>
  </div>
</header>
