<script module lang="ts">
  export {
    nextTourStep,
    prevTourStep,
    startTour,
    stopTour,
    type TourStep,
  } from '$lib/components/custom/tour-state.svelte'
</script>

<script lang="ts">
  import { browser } from '$app/environment'
  import {
    nextTourStep as advanceTourStep,
    prevTourStep as goToPrevTourStep,
    stopTour as closeTour,
    useTour,
    type TourStep as TourStepConfig,
  } from '$lib/components/custom/tour-state.svelte'
  import { Button } from '$lib/components/ui/button'
  import * as Card from '$lib/components/ui/card'
  import { onDestroy, tick } from 'svelte'
  import { fade, scale } from 'svelte/transition'

  type SpotlightRect = {
    bottom: number
    height: number
    left: number
    right: number
    top: number
    width: number
  }

  const viewportGap = 16
  const tooltipGap = 20
  const scrollMargin = 48
  const fallbackCardHeight = 220
  const fallbackCardWidth = 360

  const { tour } = useTour()

  let cardElement = $state<HTMLDivElement>()
  let activeRect = $state<SpotlightRect | null>(null)
  let currentPlacement = $state<'bottom' | 'center' | 'top'>('center')
  let cardPosition = $state({
    left: viewportGap,
    maxWidth: 420,
    top: viewportGap,
  })
  let rafId = 0

  const currentStep = $derived(tour.steps[tour.index])
  const hasStep = $derived(Boolean(currentStep))
  const isLastStep = $derived(tour.index === tour.steps.length - 1)
  const layoutKey = $derived(`${tour.open}:${tour.index}:${currentStep?.element ?? ''}`)

  function clamp(value: number, min: number, max: number) {
    return Math.min(Math.max(value, min), max)
  }

  function setCardElement(currentLayoutKey: string) {
    return (node: HTMLDivElement) => {
      cardElement = node
      void currentLayoutKey

      void tick().then(() => updateLayout(true))

      return () => {
        if (cardElement === node) {
          cardElement = undefined
        }
      }
    }
  }

  function getTargetElement(step: TourStepConfig | undefined) {
    if (!browser || !step?.element) return null

    try {
      const element = document.querySelector(step.element)
      return element instanceof HTMLElement ? element : null
    } catch {
      return null
    }
  }

  function scrollElementIntoView(element: HTMLElement) {
    const rect = element.getBoundingClientRect()
    const visibleTop = rect.top >= scrollMargin
    const visibleBottom = rect.bottom <= window.innerHeight - scrollMargin

    if (visibleTop && visibleBottom) return

    element.scrollIntoView({
      behavior: 'smooth',
      block: 'center',
      inline: 'nearest',
    })
  }

  function updateLayout(shouldScroll = false) {
    if (!browser || !tour.open || !hasStep) return

    cancelAnimationFrame(rafId)
    rafId = requestAnimationFrame(() => {
      const targetElement = getTargetElement(currentStep)

      if (shouldScroll && targetElement) {
        scrollElementIntoView(targetElement)
      }

      requestAnimationFrame(() => {
        if (!tour.open || !hasStep) return

        const viewportWidth = window.innerWidth
        const viewportHeight = window.innerHeight
        const cardWidth = cardElement?.offsetWidth ?? fallbackCardWidth
        const cardHeight = cardElement?.offsetHeight ?? fallbackCardHeight
        const maxLeft = Math.max(viewportGap, viewportWidth - cardWidth - viewportGap)
        const maxTop = Math.max(viewportGap, viewportHeight - cardHeight - viewportGap)

        cardPosition.maxWidth = Math.min(420, viewportWidth - viewportGap * 2)

        if (!targetElement) {
          activeRect = null
          currentPlacement = 'center'
          cardPosition.left = clamp((viewportWidth - cardWidth) / 2, viewportGap, maxLeft)
          cardPosition.top = clamp((viewportHeight - cardHeight) / 2, viewportGap, maxTop)
          return
        }

        const rect = targetElement.getBoundingClientRect()
        activeRect = {
          bottom: rect.bottom,
          height: rect.height,
          left: rect.left,
          right: rect.right,
          top: rect.top,
          width: rect.width,
        }

        const availableTop = rect.top - viewportGap
        const availableBottom = viewportHeight - rect.bottom - viewportGap
        const preferredPlacement = currentStep?.placement ?? 'auto'

        currentPlacement =
          preferredPlacement === 'top'
            ? 'top'
            : preferredPlacement === 'bottom'
              ? 'bottom'
              : availableBottom >= cardHeight + tooltipGap || availableBottom >= availableTop
                ? 'bottom'
                : 'top'

        const nextTop = currentPlacement === 'bottom' ? rect.bottom + tooltipGap : rect.top - cardHeight - tooltipGap
        const nextLeft = rect.left + rect.width / 2 - cardWidth / 2

        cardPosition.left = clamp(nextLeft, viewportGap, maxLeft)
        cardPosition.top = clamp(nextTop, viewportGap, maxTop)
      })
    })
  }

  function handleKeydown(event: KeyboardEvent) {
    if (!tour.open) return

    if (event.key === 'Escape') {
      closeTour()
      return
    }

    if (event.key === 'ArrowLeft') {
      goToPrevTourStep()
      return
    }

    if (event.key === 'ArrowRight' || event.key === 'Enter') {
      advanceTourStep()
    }
  }

  onDestroy(() => {
    cancelAnimationFrame(rafId)
  })
</script>

<svelte:window onkeydown={handleKeydown} onresize={() => updateLayout()} onscroll={() => updateLayout()} />

{#if tour.open && hasStep}
  <div class="pointer-events-none fixed inset-0 z-100 bg-background/38 backdrop-blur-[1px]" transition:fade>
    {#if activeRect}
      <div
        class="absolute rounded-3xl border-2 border-primary bg-background/6 shadow-[0_0_0_9999px_color-mix(in_oklab,var(--color-background)_52%,transparent),0_0_0_1px_color-mix(in_oklab,var(--color-primary)_28%,transparent),0_22px_44px_-18px_color-mix(in_oklab,var(--color-primary)_45%,transparent)] ring-4 ring-primary/18 transition-all duration-300"
        style={`left:${activeRect.left - 8}px;top:${activeRect.top - 8}px;width:${activeRect.width + 16}px;height:${activeRect.height + 16}px;`}
      ></div>
    {/if}
  </div>

  <div class="fixed inset-0 z-110">
    <div
      aria-atomic="true"
      aria-live="polite"
      class="pointer-events-auto absolute w-full"
      role="dialog"
      style={`left:${cardPosition.left}px;top:${cardPosition.top}px;max-width:${cardPosition.maxWidth}px;`}
      transition:scale={{ duration: 180, start: 0.96 }}
      {@attach setCardElement(layoutKey)}
    >
      <Card.Root
        class="rounded-4xl border-border/80 bg-card/95 text-card-foreground shadow-2xl ring-1 shadow-foreground/10 ring-foreground/5 backdrop-blur"
      >
        <Card.Header class="gap-4 px-5 pt-5 sm:px-6 sm:pt-6">
          <div class="flex items-start justify-between gap-4">
            <div class="space-y-3">
              <div class="flex items-center gap-3">
                <span
                  class="inline-flex rounded-full border border-primary/15 bg-primary/10 px-3 py-1 text-xs font-semibold tracking-[0.24em] text-primary uppercase"
                >
                  Tour step {tour.index + 1}
                </span>
                <span class="text-xs font-medium tracking-[0.24em] text-muted-foreground uppercase">
                  {currentPlacement === 'center' ? 'Overview' : currentPlacement}
                </span>
              </div>

              {#if currentStep?.title}
                <Card.Title class="text-lg tracking-tight text-balance">
                  {currentStep.title}
                </Card.Title>
              {/if}

              <Card.Description class="max-w-[34ch] text-sm leading-6 text-muted-foreground">
                {currentStep?.message}
              </Card.Description>
            </div>

            <Button type="button" variant="ghost" size="icon-lg" class="rounded-full" onclick={closeTour}>
              <span class="sr-only">Close tour</span>
              ×
            </Button>
          </div>
        </Card.Header>

        <Card.Content class="px-5 sm:px-6">
          <div class="h-2 overflow-hidden rounded-full bg-muted">
            <div
              class="h-full rounded-full bg-linear-to-r from-primary via-chart-1 to-chart-2 transition-[width] duration-300"
              style={`width:${((tour.index + 1) / tour.steps.length) * 100}%`}
            ></div>
          </div>
        </Card.Content>

        <Card.Footer
          class="mt-5 flex flex-wrap items-center justify-between gap-3 border-0 bg-transparent px-5 pt-0 pb-5 sm:px-6 sm:pb-6"
        >
          <p class="text-sm text-muted-foreground">{tour.index + 1} of {tour.steps.length}</p>

          <div class="flex items-center gap-2">
            <Button
              type="button"
              variant="outline"
              class="min-w-24 rounded-full"
              disabled={tour.index === 0}
              onclick={goToPrevTourStep}
            >
              Previous
            </Button>

            <Button type="button" class="min-w-24 rounded-full" onclick={advanceTourStep}>
              {isLastStep ? 'Finish' : 'Next'}
            </Button>
          </div>
        </Card.Footer>
      </Card.Root>
    </div>
  </div>
{/if}
