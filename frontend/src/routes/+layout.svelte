<script lang="ts">
  import { page } from '$app/state'
  import { queryClient } from '$lib/api/client'
  import favicon from '$lib/assets/logo.svg'
  import AppHeader from '$lib/components/custom/AppHeader.svelte'
  import Header from '$lib/components/custom/Header.svelte'
  import Navbar from '$lib/components/custom/Navbar.svelte'
  import { QueryClientProvider } from '@tanstack/svelte-query'
  import { SvelteQueryDevtools } from '@tanstack/svelte-query-devtools'
  import { Toaster } from 'svelte-sonner'
  import './layout.css'

  let { children } = $props()
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
  <link rel="icon" type="image/svg+xml" href="/logo.svg" />
</svelte:head>

{#snippet main()}
  <main class="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
    {@render children()}
  </main>
{/snippet}

<QueryClientProvider client={queryClient}>
  <div class="min-h-screen bg-background text-foreground">
    {#if !page.url.pathname.startsWith('/app')}
      <Header />
      {@render main()}

      <Toaster richColors />
    {:else}
      <div class="flex">
        <Navbar />

        <div class="grow">
          <AppHeader />
          <div class="h-[calc(100vh-70px)] overflow-y-auto pb-10" style="scrollbar-width: thin;">
            {@render main()}
          </div>
        </div>
      </div>
    {/if}
  </div>
  <SvelteQueryDevtools />
</QueryClientProvider>
