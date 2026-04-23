<script lang="ts">
  import { initAnalyticsIfConsented, setConsentChoice } from '$lib/analytics/tracking'
  import { onMount } from 'svelte'

  let isVisible = $state(false)

  onMount(() => {
    const existingChoice = localStorage.getItem('cookie-consent')
    isVisible = existingChoice !== 'accepted' && existingChoice !== 'declined'
  })

  function acceptCookies() {
    setConsentChoice('accepted')
    initAnalyticsIfConsented()
    isVisible = false
  }

  function declineCookies() {
    setConsentChoice('declined')
    isVisible = false
  }
</script>

{#if isVisible}
  <section
    class="fixed inset-x-4 bottom-4 z-50 rounded-2xl border border-border/70 bg-card/95 p-4 text-card-foreground shadow-2xl backdrop-blur sm:inset-x-auto sm:right-6 sm:bottom-6 sm:max-w-md"
    aria-live="polite"
    aria-label="Cookie consent"
  >
    <h2 class="font-heading text-base font-semibold">Cookie preferences</h2>
    <p class="mt-2 text-sm text-muted-foreground">
      We use optional analytics cookies to understand traffic and improve the app. You can accept or decline anytime by
      clearing your browser storage.
    </p>
    <div class="mt-4 flex gap-2">
      <button
        class="rounded-md bg-primary px-3 py-2 text-sm font-medium text-primary-foreground hover:opacity-90"
        type="button"
        onclick={acceptCookies}
      >
        Accept
      </button>
      <button
        class="rounded-md border border-border bg-background px-3 py-2 text-sm font-medium text-foreground hover:bg-muted"
        type="button"
        onclick={declineCookies}
      >
        Decline
      </button>
    </div>
  </section>
{/if}
