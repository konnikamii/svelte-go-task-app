<script lang="ts">
  import { useGetMe } from '$lib/api/users/users'
  import { Button } from '$lib/components/ui/button'
  import * as Card from '$lib/components/ui/card'

  type HomeView = 'dashboard' | 'tasks' | 'settings'
  type Density = 'comfortable' | 'compact'

  const userQuery = useGetMe()

  let homeView = $state<HomeView>('dashboard')
  let density = $state<Density>('comfortable')
  let digestEnabled = $state(true)
  let deadlineWarnings = $state(true)
  let focusMode = $state(false)

  const user = $derived(userQuery.data)
  const enabledCount = $derived.by(() => [digestEnabled, deadlineWarnings, focusMode].filter(Boolean).length)
  const profileInitials = $derived.by(() => {
    const source = user?.username ?? user?.email ?? 'TU'
    return source.slice(0, 2).toUpperCase()
  })

  function togglePreference(key: 'digestEnabled' | 'deadlineWarnings' | 'focusMode') {
    if (key === 'digestEnabled') digestEnabled = !digestEnabled
    if (key === 'deadlineWarnings') deadlineWarnings = !deadlineWarnings
    if (key === 'focusMode') focusMode = !focusMode
  }

  function resetPreferences() {
    homeView = 'dashboard'
    density = 'comfortable'
    digestEnabled = true
    deadlineWarnings = true
    focusMode = false
  }
</script>

<svelte:head>
  <title>Taskify | Settings</title>
</svelte:head>

<div class="mx-auto flex w-full max-w-7xl flex-col gap-6 px-4 py-6 sm:px-6 lg:px-8">
  <section class="grid gap-6 lg:grid-cols-[1.1fr_0.9fr]">
    <Card.Root
      class="overflow-hidden rounded-4xl border-primary/15 bg-linear-to-br from-card via-card to-primary/6 shadow-xl shadow-primary/5"
    >
      <Card.Header class="gap-5 px-6 py-8 sm:px-8">
        <div
          class="inline-flex w-fit rounded-full border border-primary/15 bg-primary/10 px-3 py-1 text-xs font-semibold tracking-[0.28em] text-primary uppercase"
        >
          Settings
        </div>
        <div class="space-y-3">
          <Card.Title class="max-w-3xl text-3xl font-semibold tracking-tight sm:text-4xl">
            Tune the workspace around how you actually work.
          </Card.Title>
          <Card.Description class="max-w-2xl text-sm leading-6 sm:text-base">
            Set the defaults you want to keep close and shape the workspace around your preferred flow.
          </Card.Description>
        </div>
      </Card.Header>

      <Card.Content class="grid gap-4 px-6 pb-8 sm:grid-cols-3 sm:px-8">
        <Card.Root size="sm" class="rounded-3xl border-border/70 bg-muted/45 py-0 shadow-none">
          <Card.Header>
            <Card.Description>Home route</Card.Description>
            <Card.Title class="text-2xl capitalize">{homeView}</Card.Title>
          </Card.Header>
        </Card.Root>

        <Card.Root size="sm" class="rounded-3xl border-border/70 bg-muted/45 py-0 shadow-none">
          <Card.Header>
            <Card.Description>Density</Card.Description>
            <Card.Title class="text-2xl capitalize">{density}</Card.Title>
          </Card.Header>
        </Card.Root>

        <Card.Root size="sm" class="rounded-3xl border-border/70 bg-muted/45 py-0 shadow-none">
          <Card.Header>
            <Card.Description>Active assists</Card.Description>
            <Card.Title class="text-2xl">{enabledCount}/3</Card.Title>
          </Card.Header>
        </Card.Root>
      </Card.Content>
    </Card.Root>

    <Card.Root class="rounded-4xl border-border/80 shadow-sm">
      <Card.Header class="px-6 py-6 sm:px-7">
        <Card.Description class="tracking-[0.28em] text-chart-2 uppercase">Account</Card.Description>
        <Card.Title class="text-2xl">Profile snapshot</Card.Title>
      </Card.Header>
      <Card.Content class="space-y-4 px-6 pb-6 sm:px-7">
        {#if userQuery.isLoading}
          <div class="rounded-3xl bg-muted/45 p-4 text-sm text-muted-foreground">Loading your profile...</div>
        {:else if userQuery.isError}
          <div class="rounded-3xl border border-destructive/20 bg-destructive/10 p-4 text-sm text-destructive">
            Profile details could not be loaded right now.
          </div>
        {:else}
          <div class="flex items-center gap-4 rounded-3xl border border-border/70 bg-muted/35 p-4">
            <div
              class="inline-flex size-14 items-center justify-center rounded-full bg-primary/12 text-lg font-semibold text-primary"
            >
              {profileInitials}
            </div>
            <div class="min-w-0">
              <p class="truncate text-lg font-semibold">{user?.username}</p>
              <p class="truncate text-sm text-muted-foreground">{user?.email}</p>
            </div>
          </div>
        {/if}
      </Card.Content>
    </Card.Root>
  </section>

  <section class="grid gap-6 lg:grid-cols-[1fr_1fr]">
    <Card.Root class="rounded-4xl border-border/80 shadow-sm">
      <Card.Header class="px-6 py-6 sm:px-7">
        <div class="flex items-start justify-between gap-4">
          <div>
            <Card.Description class="tracking-[0.28em] text-chart-5 uppercase">Workspace defaults</Card.Description>
            <Card.Title class="text-2xl">Interactive preferences</Card.Title>
          </div>
          <Button variant="ghost" size="sm" class="rounded-full" onclick={resetPreferences}>Reset</Button>
        </div>
      </Card.Header>
      <Card.Content class="space-y-5 px-6 pb-6 sm:px-7">
        <div class="space-y-3">
          <p class="text-sm font-medium text-foreground">Default landing route</p>
          <div class="flex flex-wrap gap-2">
            <Button
              variant={homeView === 'dashboard' ? 'default' : 'outline'}
              size="sm"
              class="rounded-full"
              onclick={() => (homeView = 'dashboard')}>Dashboard</Button
            >
            <Button
              variant={homeView === 'tasks' ? 'default' : 'outline'}
              size="sm"
              class="rounded-full"
              onclick={() => (homeView = 'tasks')}>Tasks</Button
            >
            <Button
              variant={homeView === 'settings' ? 'default' : 'outline'}
              size="sm"
              class="rounded-full"
              onclick={() => (homeView = 'settings')}>Settings</Button
            >
          </div>
        </div>

        <div class="space-y-3">
          <p class="text-sm font-medium text-foreground">Workspace density</p>
          <div class="flex flex-wrap gap-2">
            <Button
              variant={density === 'comfortable' ? 'default' : 'outline'}
              size="sm"
              class="rounded-full"
              onclick={() => (density = 'comfortable')}>Comfortable</Button
            >
            <Button
              variant={density === 'compact' ? 'default' : 'outline'}
              size="sm"
              class="rounded-full"
              onclick={() => (density = 'compact')}>Compact</Button
            >
          </div>
        </div>

        <div class="space-y-3">
          <p class="text-sm font-medium text-foreground">Assistive controls</p>
          <div class="grid gap-3">
            <button
              type="button"
              class="flex items-center justify-between rounded-3xl border border-border/70 bg-muted/35 p-4 text-left transition hover:border-primary/20 hover:bg-muted/50"
              onclick={() => togglePreference('digestEnabled')}
            >
              <div>
                <p class="font-medium">Daily digest</p>
                <p class="text-sm text-muted-foreground">Start with a short summary of due work and completions.</p>
              </div>
              <div
                class={`rounded-full px-3 py-1 text-xs font-medium ${digestEnabled ? 'bg-primary/10 text-primary' : 'bg-secondary text-secondary-foreground'}`}
              >
                {digestEnabled ? 'On' : 'Off'}
              </div>
            </button>

            <button
              type="button"
              class="flex items-center justify-between rounded-3xl border border-border/70 bg-muted/35 p-4 text-left transition hover:border-primary/20 hover:bg-muted/50"
              onclick={() => togglePreference('deadlineWarnings')}
            >
              <div>
                <p class="font-medium">Deadline warnings</p>
                <p class="text-sm text-muted-foreground">
                  Highlight tasks that are due soon before they become urgent.
                </p>
              </div>
              <div
                class={`rounded-full px-3 py-1 text-xs font-medium ${deadlineWarnings ? 'bg-chart-4/10 text-chart-4' : 'bg-secondary text-secondary-foreground'}`}
              >
                {deadlineWarnings ? 'On' : 'Off'}
              </div>
            </button>

            <button
              type="button"
              class="flex items-center justify-between rounded-3xl border border-border/70 bg-muted/35 p-4 text-left transition hover:border-primary/20 hover:bg-muted/50"
              onclick={() => togglePreference('focusMode')}
            >
              <div>
                <p class="font-medium">Focus mode</p>
                <p class="text-sm text-muted-foreground">
                  Reduce visual noise and emphasize the single most important task.
                </p>
              </div>
              <div
                class={`rounded-full px-3 py-1 text-xs font-medium ${focusMode ? 'bg-chart-2/10 text-chart-2' : 'bg-secondary text-secondary-foreground'}`}
              >
                {focusMode ? 'On' : 'Off'}
              </div>
            </button>
          </div>
        </div>
      </Card.Content>
    </Card.Root>

    <Card.Root class="rounded-4xl border-border/80 shadow-sm">
      <Card.Header class="px-6 py-6 sm:px-7">
        <Card.Description class="tracking-[0.28em] text-chart-1 uppercase">Preview</Card.Description>
        <Card.Title class="text-2xl">Workspace summary</Card.Title>
      </Card.Header>
      <Card.Content class="space-y-5 px-6 pb-6 sm:px-7">
        <div class="rounded-3xl border border-border/70 bg-muted/40 p-5">
          <p class="text-sm font-medium tracking-[0.24em] text-muted-foreground uppercase">Workspace summary</p>
          <p class="mt-3 text-xl font-semibold">{homeView.charAt(0).toUpperCase() + homeView.slice(1)} opens first</p>
          <p class="mt-2 text-sm leading-6 text-muted-foreground">
            The workspace is set to a {density} layout with {enabledCount} assistive feature{enabledCount === 1
              ? ''
              : 's'} enabled.
          </p>
        </div>

        <div class="grid gap-3 sm:grid-cols-3">
          <Card.Root size="sm" class="rounded-3xl border-border/70 bg-muted/40 py-0 shadow-none">
            <Card.Header>
              <Card.Description>Digest</Card.Description>
              <Card.Title class="text-xl">{digestEnabled ? 'Enabled' : 'Muted'}</Card.Title>
            </Card.Header>
          </Card.Root>
          <Card.Root size="sm" class="rounded-3xl border-border/70 bg-muted/40 py-0 shadow-none">
            <Card.Header>
              <Card.Description>Warnings</Card.Description>
              <Card.Title class="text-xl">{deadlineWarnings ? 'Visible' : 'Hidden'}</Card.Title>
            </Card.Header>
          </Card.Root>
          <Card.Root size="sm" class="rounded-3xl border-border/70 bg-muted/40 py-0 shadow-none">
            <Card.Header>
              <Card.Description>Focus mode</Card.Description>
              <Card.Title class="text-xl">{focusMode ? 'Active' : 'Inactive'}</Card.Title>
            </Card.Header>
          </Card.Root>
        </div>
      </Card.Content>
    </Card.Root>
  </section>
</div>
