<script lang="ts" module>
  import { LayoutDashboardIcon, Settings2, PackagePlus } from '@lucide/svelte'

  export const navItems = [
    {
      href: '/app/dashboard' as const,
      label: 'Dashboard',
      message: 'Welcome to your dashboard! Here you can see a quick overview of completed, upcoming and urgent tasks.',
      icon: LayoutDashboardIcon,
    },
    {
      href: '/app/tasks' as const,
      label: 'Tasks',
      message: 'Here you can see all your tasks, filter them and create new ones.',
      icon: PackagePlus,
    },
    {
      href: '/app/settings' as const,
      label: 'Settings',
      message: 'Manage your account settings, preferences and more.',
      icon: Settings2,
    },
  ]
</script>

<script lang="ts">
  import { resolve } from '$app/paths'
  import { page } from '$app/state'
  import { useLogout } from '$lib/api/auth/auth'
  import TaskifyLogo from '$lib/components/svgs/TaskifyLogo.svelte'
  import { LogOut, ChevronLeft } from '@lucide/svelte'

  let navbarOpen = $state(true)

  const logoutMutation = useLogout()
</script>

{#snippet link(label: string)}
  <span
    class={`overflow-hidden whitespace-nowrap transition-[max-width,opacity] duration-150 ${navbarOpen ? 'max-w-28 opacity-100' : 'max-w-0 opacity-0 group-hover:max-w-28 group-hover:opacity-100'}`}
  >
    {label}
  </span>
{/snippet}

<aside
  class={`flex h-screen flex-col border-r border-border transition-all duration-150 ${navbarOpen ? 'w-[230px]' : 'w-[80px]'}`}
>
  <div class="relative flex h-[70px] items-center justify-center overflow-visible border-b border-border p-4">
    <a
      href={resolve('/app/dashboard')}
      draggable="false"
      class={`group flex items-center text-lg font-semibold text-foreground transition-[gap] duration-150 ${navbarOpen ? 'gap-2' : 'gap-0'}`}
    >
      <TaskifyLogo width={40} height={40} />
      <span
        class={`overflow-hidden text-2xl whitespace-nowrap transition-[max-width,opacity] duration-150 ${navbarOpen ? 'max-w-28 opacity-100' : 'max-w-0 opacity-0'}`}
      >
        Taskify
      </span>
    </a>
    <div class="absolute top-1/2 right-0 z-[60] translate-x-1/2 -translate-y-1/2">
      <ChevronLeft
        onclick={() => (navbarOpen = !navbarOpen)}
        class={`size-5 cursor-pointer rounded-full bg-accent-foreground text-black transition duration-150 hover:opacity-80 ${navbarOpen ? '' : 'rotate-180'}`}
      />
    </div>
  </div>
  <div class="flex grow flex-col gap-2 p-4">
    {#each navItems as item (item.href)}
      <a
        href={resolve(item.href)}
        aria-label={item.label}
        class={`group flex items-center rounded-md px-3 py-2 transition-all duration-150 ${navbarOpen ? 'w-full gap-2 hover:opacity-80' : 'w-[44px] gap-0 hover:w-fit hover:gap-2'} ${page.url.pathname.startsWith(item.href) ? 'bg-accent' : 'bg-secondary'}`}
      >
        <item.icon class="size-5" />
        {@render link(item.label)}
      </a>
    {/each}
  </div>
  <div class="flex h-[70px] flex-col border-t border-border p-4">
    <button
      type="button"
      class={`group flex cursor-pointer items-center rounded-md bg-secondary px-3 py-2 transition-all duration-150 ${navbarOpen ? 'w-full gap-2 hover:opacity-80' : 'w-[44px] gap-0 hover:w-fit hover:gap-2'}`}
      onclick={() => logoutMutation.mutate()}
    >
      <LogOut class="size-5" />
      {@render link('Logout')}
    </button>
  </div>
</aside>
