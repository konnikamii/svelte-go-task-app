<script lang="ts">
  import { resolve } from '$app/paths'
  import { Button } from '$lib/components/ui/button'
  import * as Card from '$lib/components/ui/card/index.js'

  type DemoTask = {
    id: number
    title: string
    lane: 'deep' | 'quick' | 'review'
    minutes: number
    note: string
    done: boolean
  }

  type DemoLane = 'all' | DemoTask['lane']

  const initialTasks: DemoTask[] = [
    {
      id: 1,
      title: 'Draft the landing copy',
      lane: 'deep',
      minutes: 45,
      note: 'Shape the value prop before the rest of the UI follows it.',
      done: false,
    },
    {
      id: 2,
      title: 'Reply to client questions',
      lane: 'quick',
      minutes: 15,
      note: 'Short, visible work that clears the queue fast.',
      done: true,
    },
    {
      id: 3,
      title: 'Review overdue tasks',
      lane: 'review',
      minutes: 20,
      note: 'Spot the items that need escalation before they slip further.',
      done: false,
    },
    {
      id: 4,
      title: "Plan tomorrow's focus block",
      lane: 'deep',
      minutes: 30,
      note: 'Use a quiet block to protect the next important deliverable.',
      done: false,
    },
  ]

  let demoLane = $state<DemoLane>('all')
  let demoTasks = $state(initialTasks.map((task) => ({ ...task })))
  let selectedTaskId = $state(initialTasks[0].id)

  const visibleTasks = $derived.by(() => demoTasks.filter((task) => demoLane === 'all' || task.lane === demoLane))
  const completedCount = $derived.by(() => demoTasks.filter((task) => task.done).length)
  const remainingMinutes = $derived.by(() =>
    demoTasks.filter((task) => !task.done).reduce((total, task) => total + task.minutes, 0),
  )
  const completionRate = $derived.by(() => Math.round((completedCount / demoTasks.length) * 100))
  const selectedTask = $derived.by(
    () => demoTasks.find((task) => task.id === selectedTaskId) ?? visibleTasks[0] ?? demoTasks[0],
  )
  const statusLabel = $derived.by(() => {
    if (completionRate >= 75) return 'Clear runway'
    if (completionRate >= 40) return 'Good momentum'
    return 'Needs focus'
  })

  function setLane(nextLane: DemoLane) {
    demoLane = nextLane

    const nextVisible = demoTasks.filter((task) => nextLane === 'all' || task.lane === nextLane)
    if (nextVisible.length > 0) {
      selectedTaskId = nextVisible[0].id
    }
  }

  function toggleTask(taskId: number) {
    const task = demoTasks.find((entry) => entry.id === taskId)
    if (!task) return

    task.done = !task.done
    selectedTaskId = taskId
  }

  function resetDemo() {
    demoTasks = initialTasks.map((task) => ({ ...task }))
    demoLane = 'all'
    selectedTaskId = initialTasks[0].id
  }

  function laneLabel(lane: DemoTask['lane']) {
    if (lane === 'deep') return 'Deep work'
    if (lane === 'quick') return 'Quick wins'
    return 'Review'
  }
</script>

<svelte:head>
  <title>Taskify | Focused task planning</title>
  <meta
    name="description"
    content="Taskify helps teams shape focused workdays with a clean Svelte frontend and Go backend. Try the interactive preview on the home page and step into the app."
  />
</svelte:head>

<div class="mx-auto flex w-full max-w-7xl flex-col gap-6 px-4 py-6 sm:px-6 lg:px-8">
  <section class="grid gap-6 lg:grid-cols-[1.15fr_0.95fr]">
    <Card.Root
      class="overflow-hidden rounded-4xl border-primary/15 bg-linear-to-br from-card via-card to-primary/6 shadow-xl shadow-primary/5"
    >
      <Card.Header class="gap-5 px-6 py-8 sm:px-8 sm:py-10">
        <div
          class="inline-flex w-fit rounded-full border border-primary/15 bg-primary/10 px-3 py-1 text-xs font-semibold tracking-[0.28em] text-primary uppercase"
        >
          Taskify
        </div>
        <div class="space-y-3">
          <Card.Title class="max-w-3xl text-4xl leading-tight font-semibold tracking-tight sm:text-5xl">
            Plan the day, test the flow, and move into the real workspace fast.
          </Card.Title>
          <Card.Description class="max-w-2xl text-base leading-7 text-muted-foreground sm:text-lg">
            Explore the workflow, move through a focused queue, and step into the workspace when you are ready.
          </Card.Description>
        </div>
      </Card.Header>

      <Card.Content class="px-6 pb-8 sm:px-8">
        <div class="flex flex-wrap gap-3">
          <Button href={resolve('/register')} size="lg" class="rounded-full px-5">Create account</Button>
          <Button href={resolve('/login')} variant="outline" size="lg" class="rounded-full px-5">Log in</Button>
          <Button href="/contact" variant="ghost" size="lg" class="rounded-full px-5">Ask a question</Button>
        </div>

        <div class="mt-8 grid gap-4 sm:grid-cols-3">
          <Card.Root size="sm" class="rounded-3xl border-border/70 bg-muted/45 py-0 shadow-none">
            <Card.Header>
              <Card.Description>Live preview</Card.Description>
              <Card.Title class="text-2xl">{completionRate}%</Card.Title>
            </Card.Header>
          </Card.Root>

          <Card.Root size="sm" class="rounded-3xl border-border/70 bg-muted/45 py-0 shadow-none">
            <Card.Header>
              <Card.Description>Minutes left</Card.Description>
              <Card.Title class="text-2xl">{remainingMinutes}</Card.Title>
            </Card.Header>
          </Card.Root>

          <Card.Root size="sm" class="rounded-3xl border-border/70 bg-muted/45 py-0 shadow-none">
            <Card.Header>
              <Card.Description>Momentum</Card.Description>
              <Card.Title class="text-2xl">{statusLabel}</Card.Title>
            </Card.Header>
          </Card.Root>
        </div>
      </Card.Content>
    </Card.Root>

    <Card.Root class="rounded-4xl border-border/80 shadow-lg shadow-foreground/5">
      <Card.Header class="px-6 py-6 sm:px-7">
        <div class="flex items-start justify-between gap-4">
          <div class="space-y-2">
            <Card.Description class="tracking-[0.28em] text-chart-2 uppercase">Interactive preview</Card.Description>
            <Card.Title class="text-2xl">Play with a sample workday</Card.Title>
          </div>
          <Button variant="ghost" size="sm" class="rounded-full" onclick={resetDemo}>Reset</Button>
        </div>
      </Card.Header>

      <Card.Content class="space-y-5 px-6 pb-6 sm:px-7">
        <div class="flex flex-wrap gap-2">
          <Button
            variant={demoLane === 'all' ? 'default' : 'outline'}
            size="sm"
            class="rounded-full"
            onclick={() => setLane('all')}
          >
            All
          </Button>
          <Button
            variant={demoLane === 'deep' ? 'default' : 'outline'}
            size="sm"
            class="rounded-full"
            onclick={() => setLane('deep')}
          >
            Deep work
          </Button>
          <Button
            variant={demoLane === 'quick' ? 'default' : 'outline'}
            size="sm"
            class="rounded-full"
            onclick={() => setLane('quick')}
          >
            Quick wins
          </Button>
          <Button
            variant={demoLane === 'review' ? 'default' : 'outline'}
            size="sm"
            class="rounded-full"
            onclick={() => setLane('review')}
          >
            Review
          </Button>
        </div>

        <div class="space-y-3">
          {#each visibleTasks as task (task.id)}
            <button
              type="button"
              class={`group w-full rounded-3xl border p-4 text-left transition ${task.id === selectedTaskId ? 'border-primary/40 bg-primary/8 shadow-lg shadow-primary/8' : 'border-border/70 bg-muted/35 hover:border-primary/20 hover:bg-muted/55'} ${task.done ? 'opacity-75' : ''}`}
              onclick={() => toggleTask(task.id)}
            >
              <div class="flex items-start justify-between gap-4">
                <div class="space-y-2">
                  <div class="flex flex-wrap items-center gap-2">
                    <span
                      class={`inline-flex rounded-full px-2.5 py-1 text-[11px] font-semibold tracking-[0.24em] uppercase ${task.lane === 'deep' ? 'bg-primary/10 text-primary' : task.lane === 'quick' ? 'bg-chart-2/10 text-chart-2' : 'bg-chart-4/10 text-chart-4'}`}
                    >
                      {laneLabel(task.lane)}
                    </span>
                    <span class="text-xs text-muted-foreground">{task.minutes} min</span>
                  </div>
                  <div>
                    <p
                      class={`text-base font-semibold ${task.done ? 'text-muted-foreground line-through' : 'text-foreground'}`}
                    >
                      {task.title}
                    </p>
                    <p class="mt-1 text-sm leading-6 text-muted-foreground">{task.note}</p>
                  </div>
                </div>

                <div
                  class={`mt-0.5 inline-flex size-6 shrink-0 items-center justify-center rounded-full border text-xs font-bold ${task.done ? 'border-chart-2/30 bg-chart-2/10 text-chart-2' : 'border-border bg-background text-muted-foreground group-hover:border-primary/30 group-hover:text-primary'}`}
                >
                  {task.done ? '✓' : ''}
                </div>
              </div>
            </button>
          {:else}
            <div class="rounded-3xl border border-border/70 bg-muted/35 p-4 text-sm text-muted-foreground">
              No tasks in this lane. Switch the filter to explore the other preview items.
            </div>
          {/each}
        </div>
      </Card.Content>
    </Card.Root>
  </section>

  <section class="grid gap-6 lg:grid-cols-[0.9fr_1.1fr]">
    <Card.Root class="rounded-4xl border-border/80 shadow-sm">
      <Card.Header class="px-6 py-6 sm:px-7">
        <Card.Description class="tracking-[0.28em] text-chart-5 uppercase">Selected task</Card.Description>
        <Card.Title class="text-2xl">Why this pattern works</Card.Title>
      </Card.Header>
      <Card.Content class="space-y-5 px-6 pb-6 sm:px-7">
        <div class="rounded-3xl border border-border/70 bg-muted/40 p-5">
          <p class="text-sm font-medium tracking-[0.24em] text-muted-foreground uppercase">Current item</p>
          <p class="mt-3 text-xl font-semibold">{selectedTask.title}</p>
          <p class="mt-2 text-sm leading-6 text-muted-foreground">{selectedTask.note}</p>
        </div>

        <div class="space-y-3">
          <div class="h-3 overflow-hidden rounded-full bg-muted">
            <div
              class="h-full bg-linear-to-r from-primary via-chart-1 to-chart-2 transition-[width] duration-300"
              style={`width:${completionRate}%`}
            ></div>
          </div>
          <p class="text-sm text-muted-foreground">
            {completedCount} of {demoTasks.length} tasks completed in the current queue.
          </p>
        </div>
      </Card.Content>
    </Card.Root>

    <Card.Root class="rounded-4xl border-border/80 shadow-sm">
      <Card.Header class="px-6 py-6 sm:px-7">
        <Card.Description class="tracking-[0.28em] text-chart-1 uppercase">Start clean</Card.Description>
        <Card.Title class="text-2xl">Move from preview into the real workspace</Card.Title>
      </Card.Header>
      <Card.Content class="space-y-5 px-6 pb-6 sm:px-7">
        <div class="rounded-3xl border border-border/70 bg-muted/40 p-5 sm:p-6">
          <p class="max-w-2xl text-sm leading-7 text-muted-foreground sm:text-base">
            Explore the interaction, get a feel for the workflow, and step into the app when you are ready. If you need
            details before signing up, the contact page is available directly from the header.
          </p>
        </div>

        <div class="flex flex-wrap gap-3">
          <Button href={resolve('/register')} class="rounded-full px-5">Create account</Button>
          <Button href="/contact" variant="outline" class="rounded-full px-5">Contact us</Button>
        </div>
      </Card.Content>
    </Card.Root>
  </section>
</div>
