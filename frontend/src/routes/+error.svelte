<script lang="ts">
  import { resolve } from '$app/paths'
  import { page } from '$app/state'
  import * as Button from '$lib/components/ui/button'
  import * as Card from '$lib/components/ui/card'
  import { Home, RotateCcw, SearchX, TriangleAlert } from '@lucide/svelte'

  const currentStatus = $derived(page.status ?? 500)
  const currentError = $derived(page.error)
  const isNotFound = $derived(currentStatus === 404)
  const title = $derived(isNotFound ? 'This page slipped off the board' : 'Something interrupted the flow')
  const description = $derived(
    isNotFound
      ? 'The link is stale, the route moved, or the page never existed in this workspace.'
      : (currentError?.message ?? 'An unexpected error occurred while loading this page.'),
  )

  function goBack() {
    if (typeof window !== 'undefined' && window.history.length > 1) {
      window.history.back()
    }
  }
</script>

<section class="relative flex min-h-[70dvh] items-center justify-center py-10">
  <div class="pointer-events-none absolute inset-0 -z-10 overflow-hidden">
    <div class="absolute top-0 left-1/2 h-56 w-56 -translate-x-1/2 rounded-full bg-primary/12 blur-3xl"></div>
    <div class="absolute right-0 bottom-0 h-64 w-64 rounded-full bg-chart-2/12 blur-3xl"></div>
    <div class="absolute top-12 left-0 h-40 w-40 rounded-full bg-chart-3/10 blur-3xl"></div>
  </div>

  <Card.Root class="relative w-full max-w-4xl border-border/70 bg-card/95 shadow-2xl backdrop-blur">
    <div class="absolute inset-x-0 top-0 h-24 bg-linear-to-b from-muted/60 via-muted/20 to-transparent"></div>

    <div class="grid gap-0 md:grid-cols-[1.15fr_0.85fr]">
      <div class="relative border-b border-border/60 p-8 md:border-r md:border-b-0 md:p-10">
        <div class="mb-6 flex items-start justify-between gap-4">
          <div>
            <div
              class="mb-3 inline-flex items-center rounded-full border border-border bg-background/80 px-3 py-1 text-xs font-medium tracking-[0.24em] text-muted-foreground uppercase"
            >
              Error {currentStatus}
            </div>
            <Card.Title class="text-3xl leading-tight sm:text-4xl">{title}</Card.Title>
            <Card.Description class="mt-3 max-w-xl text-base leading-7">
              {description}
            </Card.Description>
          </div>

          <div
            class="flex size-14 shrink-0 items-center justify-center rounded-2xl bg-muted text-foreground ring-1 ring-border/60"
          >
            {#if isNotFound}
              <SearchX class="size-7" />
            {:else}
              <TriangleAlert class="size-7" />
            {/if}
          </div>
        </div>

        <Card.Content class="px-0">
          <div class="grid gap-3 sm:grid-cols-3">
            <div class="rounded-2xl border border-border/60 bg-muted/60 p-4">
              <div class="text-xs tracking-[0.2em] text-muted-foreground uppercase">Status</div>
              <div class="mt-2 text-2xl font-semibold">{currentStatus}</div>
            </div>
            <div class="rounded-2xl border border-border/60 bg-muted/60 p-4">
              <div class="text-xs tracking-[0.2em] text-muted-foreground uppercase">State</div>
              <div class="mt-2 text-sm font-medium">{isNotFound ? 'Missing route' : 'Unexpected failure'}</div>
            </div>
            <div class="rounded-2xl border border-border/60 bg-muted/60 p-4">
              <div class="text-xs tracking-[0.2em] text-muted-foreground uppercase">Hint</div>
              <div class="mt-2 text-sm font-medium">
                {isNotFound ? 'Check the URL or head back home.' : 'Retry or return to a stable page.'}
              </div>
            </div>
          </div>
        </Card.Content>

        <Card.Footer class="mt-8 justify-start gap-3 border-0 bg-transparent p-0">
          <Button.Root href={resolve('/')} class="min-w-36">
            <Home class="size-4" />
            Go home
          </Button.Root>

          <Button.Root variant="outline" class="min-w-36" onclick={goBack}>
            <RotateCcw class="size-4" />
            Go back
          </Button.Root>
        </Card.Footer>
      </div>

      <div class="relative flex min-h-72 items-center justify-center p-8 md:p-10">
        <div class="absolute inset-6 rounded-[2rem] bg-linear-to-br from-primary/14 via-background to-chart-2/10"></div>
        <div
          class="relative flex aspect-square w-full max-w-xs flex-col items-center justify-center rounded-[2rem] bg-background/85 p-8 text-center shadow-lg ring-1 ring-border/60 backdrop-blur"
        >
          <div class="mb-3 text-xs tracking-[0.32em] text-muted-foreground uppercase">Taskify</div>
          <div
            class="bg-linear-to-b from-foreground to-muted-foreground bg-clip-text text-7xl leading-none font-semibold text-transparent sm:text-8xl"
          >
            {currentStatus}
          </div>
          <p class="mt-4 max-w-[16rem] text-sm leading-6 text-muted-foreground">
            {isNotFound
              ? 'The destination is not on the current route map.'
              : 'The page failed before it could finish rendering cleanly.'}
          </p>
        </div>
      </div>
    </div>
  </Card.Root>
</section>
