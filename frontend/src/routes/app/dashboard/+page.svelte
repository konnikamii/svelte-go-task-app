<script lang="ts">
  import type { TasksPaginatedRequest } from '$lib/api/tasks/crud-task.schema'
  import type { Task } from '$lib/api/tasks/task.schema'
  import { useGetPaginatedTasks } from '$lib/api/tasks/tasks'
  import { useGetMe } from '$lib/api/users/users'
  import Tour, { startTour } from '$lib/components/custom/Tour.svelte'
  import { Button } from '$lib/components/ui/button'
  import * as Card from '$lib/components/ui/card'

  const dashboardTaskQuery: TasksPaginatedRequest = {
    page: 1,
    pageSize: 100,
    sortBy: 'createdAt',
    sortType: 'desc',
    filters: {},
  }

  function launchDashboardTour() {
    startTour([
      {
        element: '#dashboard-tour-start',
        title: 'Start here',
        message: 'This button can reopen the dashboard tour any time you want a quick walkthrough.',
      },
      {
        element: '#dashboard-overview',
        title: 'Dashboard overview',
        message:
          'This intro card is a good place to explain what changed today or what the user should focus on first.',
      },
      {
        element: '#dashboard-focus',
        title: 'Primary focus',
        message:
          'Use this section for the most important work queue, deadline reminder, or personal productivity summary.',
      },
      {
        element: '#dashboard-stats',
        title: 'Stats grid',
        message: 'This panel summarizes your workload, urgency, and overall completion rate using live task data.',
      },
      {
        element: '#dashboard-priority',
        title: 'Priority lane',
        message:
          'The dashboard surfaces your most urgent open task first so you can move straight into the highest-risk item.',
      },
      {
        element: '#dashboard-upcoming',
        title: 'Upcoming work',
        message: 'This list keeps the next due tasks visible so deadlines do not disappear into a larger backlog.',
      },
      {
        title: 'You are ready',
        message:
          'The tour can also show a centered message when you want to finish with a general tip instead of pointing at a specific element.',
      },
    ])
  }

  const userQuery = useGetMe()
  const tasksQuery = useGetPaginatedTasks(dashboardTaskQuery)

  const user = $derived(userQuery.data)
  const tasks = $derived(tasksQuery.data?.entries ?? [])
  const totalEntries = $derived(tasksQuery.data?.totalEntries ?? 0)
  const now = $derived(new Date())

  const openTasks = $derived.by(() => tasks.filter((task) => !task.completed))
  const completedTasks = $derived.by(() => tasks.filter((task) => task.completed))
  const overdueTasks = $derived.by(() =>
    openTasks.filter((task) => {
      if (!task.dueDate) return false
      return new Date(task.dueDate).getTime() < now.getTime()
    }),
  )
  const dueSoonTasks = $derived.by(() =>
    openTasks
      .filter((task) => {
        if (!task.dueDate) return false
        const dueTime = new Date(task.dueDate).getTime()
        const delta = dueTime - now.getTime()
        return delta >= 0 && delta <= 1000 * 60 * 60 * 24 * 3
      })
      .sort((left, right) => new Date(left.dueDate!).getTime() - new Date(right.dueDate!).getTime()),
  )
  const unscheduledTasks = $derived.by(() => openTasks.filter((task) => !task.dueDate))
  const completionRate = $derived(tasks.length === 0 ? 0 : Math.round((completedTasks.length / tasks.length) * 100))
  const priorityTask = $derived.by(() => {
    const ranked = [...openTasks].sort((left, right) => {
      const leftDue = left.dueDate ? new Date(left.dueDate).getTime() : Number.POSITIVE_INFINITY
      const rightDue = right.dueDate ? new Date(right.dueDate).getTime() : Number.POSITIVE_INFINITY
      return leftDue - rightDue
    })

    return ranked[0] ?? null
  })
  const recentTasks = $derived.by(() => tasks.slice(0, 5))
  const workloadLabel = $derived(
    overdueTasks.length > 0
      ? 'Needs attention'
      : dueSoonTasks.length > 0
        ? 'On track'
        : openTasks.length > 0
          ? 'Stable'
          : 'Clear',
  )

  function formatDate(value: string | null | undefined) {
    if (!value) return 'No due date'

    return new Intl.DateTimeFormat('en', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    }).format(new Date(value))
  }

  function taskTone(task: Task) {
    if (!task.dueDate) return 'text-slate-500 dark:text-slate-400'

    const dueTime = new Date(task.dueDate).getTime()
    if (dueTime < now.getTime()) return 'text-rose-600 dark:text-rose-400'
    if (dueTime - now.getTime() <= 1000 * 60 * 60 * 24 * 3) return 'text-amber-600 dark:text-amber-400'
    return 'text-emerald-600 dark:text-emerald-400'
  }
</script>

<svelte:head>
  <title>Taskify | Dashboard</title>
</svelte:head>

<Tour />

<div class="space-y-8">
  <Card.Root
    id="dashboard-overview"
    class="overflow-hidden rounded-4xl border-primary/15 bg-linear-to-br from-card via-card to-primary/5 shadow-xl shadow-primary/5"
  >
    <Card.Header class="gap-6 px-6 py-8 sm:px-8">
      <div class="flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
        <div class="max-w-2xl space-y-4">
          <div
            class="inline-flex w-fit rounded-full border border-primary/15 bg-primary/10 px-3 py-1 text-xs font-semibold tracking-[0.28em] text-primary uppercase"
          >
            Dashboard
          </div>
          <div class="space-y-2">
            <Card.Title class="text-3xl font-semibold tracking-tight sm:text-4xl">
              Welcome back{user?.username ? `, ${user.username}` : ''}
            </Card.Title>
            <Card.Description class="max-w-xl text-sm leading-6 sm:text-base">
              Keep your work visible, track delivery risk, and focus the day around the tasks that need action first.
            </Card.Description>
          </div>
        </div>

        <Card.Action class="flex flex-wrap items-center gap-3">
          <div class="rounded-full border border-primary/15 bg-primary/10 px-4 py-2 text-sm font-medium text-primary">
            {workloadLabel}
          </div>

          <Button id="dashboard-tour-start" size="lg" onclick={launchDashboardTour}>Start tour</Button>
        </Card.Action>
      </div>
    </Card.Header>
  </Card.Root>

  <section class="grid gap-4 md:grid-cols-3" id="dashboard-stats">
    <Card.Root class="rounded-3xl border-border/80 shadow-sm">
      <Card.Header>
        <Card.Description>Open tasks</Card.Description>
        <Card.Title class="text-3xl">{openTasks.length}</Card.Title>
      </Card.Header>
      <Card.Content>
        <p class="text-sm text-muted-foreground">Tasks that still need active work.</p>
      </Card.Content>
    </Card.Root>

    <Card.Root class="rounded-3xl border-border/80 shadow-sm">
      <Card.Header>
        <Card.Description>Completion rate</Card.Description>
        <Card.Title class="text-3xl">{completionRate}%</Card.Title>
      </Card.Header>
      <Card.Content>
        <p class="text-sm text-chart-2">{completedTasks.length} of {tasks.length} tasks are done.</p>
      </Card.Content>
    </Card.Root>

    <Card.Root class="rounded-3xl border-border/80 shadow-sm">
      <Card.Header>
        <Card.Description>Due soon / overdue</Card.Description>
        <Card.Title class="text-3xl">{dueSoonTasks.length + overdueTasks.length}</Card.Title>
      </Card.Header>
      <Card.Content>
        <p class="text-sm text-chart-4">{overdueTasks.length} overdue, {dueSoonTasks.length} due in 3 days.</p>
      </Card.Content>
    </Card.Root>
  </section>

  <section class="grid gap-4 lg:grid-cols-[1.3fr_0.9fr]">
    <Card.Root id="dashboard-priority" class="rounded-3xl border-border/80 shadow-sm">
      <Card.Header class="pb-0">
        <Card.Description class="tracking-[0.28em] text-chart-2 uppercase">Priority lane</Card.Description>
        <Card.Title class="text-2xl">Move the most urgent work first</Card.Title>
        <Card.Action>
          <div class="rounded-full border border-chart-2/20 bg-chart-2/10 px-3 py-1 text-xs font-medium text-chart-2">
            Focus
          </div>
        </Card.Action>
      </Card.Header>
      <Card.Content>
        {#if tasksQuery.isLoading}
          <div class="rounded-2xl bg-muted/60 p-4 text-sm text-muted-foreground">Loading task priorities...</div>
        {:else if tasksQuery.isError}
          <div class="rounded-2xl border border-destructive/20 bg-destructive/10 p-4 text-sm text-destructive">
            The dashboard could not load your task summary right now.
          </div>
        {:else if priorityTask}
          <Card.Root id="dashboard-focus" class="rounded-3xl border-border/80 bg-muted/35 py-0 shadow-none">
            <Card.Header>
              <Card.Title class="text-xl">{priorityTask.title}</Card.Title>
              <Card.Description class={['text-sm font-medium', taskTone(priorityTask)]}>
                {formatDate(priorityTask.dueDate)}
              </Card.Description>
              <Card.Action>
                <div
                  class="rounded-full border border-border bg-background px-3 py-1 text-xs font-medium text-muted-foreground"
                >
                  {priorityTask.completed ? 'Completed' : 'Open'}
                </div>
              </Card.Action>
            </Card.Header>
            <Card.Content>
              <p class="text-sm leading-6 text-muted-foreground">
                {priorityTask.description ??
                  'This task has no description yet. Add a short brief so collaborators can move faster.'}
              </p>
            </Card.Content>
          </Card.Root>
        {:else}
          <div class="rounded-2xl border border-chart-1/20 bg-chart-1/10 p-4 text-sm text-chart-3">
            No urgent tasks right now. Your active queue is clear.
          </div>
        {/if}
      </Card.Content>
    </Card.Root>

    <Card.Root id="dashboard-upcoming" class="rounded-3xl border-border/80 shadow-sm">
      <Card.Header>
        <Card.Description class="tracking-[0.28em] text-chart-5 uppercase">Upcoming work</Card.Description>
        <Card.Title class="text-2xl">What needs attention next</Card.Title>
      </Card.Header>
      <Card.Content>
        {#if tasksQuery.isLoading}
          <p class="text-sm text-muted-foreground">Loading upcoming tasks...</p>
        {:else}
          <div class="space-y-3">
            {#each dueSoonTasks.slice(0, 4) as task (task.id)}
              <Card.Root size="sm" class="rounded-2xl border-border/70 bg-muted/45 py-0 shadow-none">
                <Card.Header>
                  <Card.Title class="text-sm">{task.title}</Card.Title>
                  <Card.Description class={['mt-1 text-xs font-medium', taskTone(task)]}
                    >{formatDate(task.dueDate)}</Card.Description
                  >
                  <Card.Action>
                    <div
                      class="rounded-full border border-border bg-background px-2 py-1 text-[11px] text-muted-foreground"
                    >
                      Due
                    </div>
                  </Card.Action>
                </Card.Header>
              </Card.Root>
            {:else}
              <div class="rounded-2xl bg-muted/45 p-4 text-sm text-muted-foreground">
                No deadlines in the next three days.
              </div>
            {/each}
          </div>
        {/if}
      </Card.Content>
    </Card.Root>
  </section>

  <section class="grid gap-4 xl:grid-cols-[1.1fr_0.9fr]">
    <Card.Root class="rounded-3xl border-border/80 shadow-sm">
      <Card.Header>
        <Card.Description class="tracking-[0.28em] text-chart-4 uppercase">Workload health</Card.Description>
        <Card.Title class="text-2xl">Snapshot of your current queue</Card.Title>
        <Card.Action>
          <div class="rounded-full border border-border px-3 py-1 text-xs text-muted-foreground">
            {Math.min(tasks.length, totalEntries)} of {totalEntries}
          </div>
        </Card.Action>
      </Card.Header>
      <Card.Content>
        <div class="grid gap-3 sm:grid-cols-3">
          <Card.Root size="sm" class="rounded-2xl border-border/70 bg-muted/50 py-0 shadow-none">
            <Card.Header>
              <Card.Description>Overdue</Card.Description>
              <Card.Title class="text-2xl">{overdueTasks.length}</Card.Title>
            </Card.Header>
          </Card.Root>
          <Card.Root size="sm" class="rounded-2xl border-border/70 bg-muted/50 py-0 shadow-none">
            <Card.Header>
              <Card.Description>Scheduled</Card.Description>
              <Card.Title class="text-2xl">{openTasks.length - unscheduledTasks.length}</Card.Title>
            </Card.Header>
          </Card.Root>
          <Card.Root size="sm" class="rounded-2xl border-border/70 bg-muted/50 py-0 shadow-none">
            <Card.Header>
              <Card.Description>Unscheduled</Card.Description>
              <Card.Title class="text-2xl">{unscheduledTasks.length}</Card.Title>
            </Card.Header>
          </Card.Root>
        </div>

        <div class="mt-6 h-3 overflow-hidden rounded-full bg-muted">
          <div class="h-full bg-linear-to-r from-primary to-chart-2" style={`width:${completionRate}%`}></div>
        </div>
        <p class="mt-3 text-sm text-muted-foreground">Completion rate across the loaded task set: {completionRate}%.</p>
      </Card.Content>
    </Card.Root>

    <Card.Root class="rounded-3xl border-border/80 shadow-sm">
      <Card.Header>
        <Card.Description class="tracking-[0.28em] text-chart-1 uppercase">Recent activity</Card.Description>
        <Card.Title class="text-2xl">Latest tasks in your workspace</Card.Title>
      </Card.Header>
      <Card.Content>
        <div class="space-y-3">
          {#each recentTasks as task (task.id)}
            <Card.Root size="sm" class="rounded-2xl border-border/70 bg-muted/45 py-0 shadow-none">
              <Card.Header>
                <Card.Title class="text-sm">{task.title}</Card.Title>
                <Card.Description>Created {formatDate(task.createdAt)}</Card.Description>
                <Card.Action>
                  <div
                    class={{
                      'rounded-full px-2 py-1 text-[11px] font-medium': true,
                      'bg-chart-2/10 text-chart-2': task.completed,
                      'bg-secondary text-secondary-foreground': !task.completed,
                    }}
                  >
                    {task.completed ? 'Done' : 'Open'}
                  </div>
                </Card.Action>
              </Card.Header>
            </Card.Root>
          {:else}
            <div class="rounded-2xl bg-muted/45 p-4 text-sm text-muted-foreground">
              No task activity yet. Create your first task to start tracking momentum.
            </div>
          {/each}
        </div>
      </Card.Content>
    </Card.Root>
  </section>
</div>
