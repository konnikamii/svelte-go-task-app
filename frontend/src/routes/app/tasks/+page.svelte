<script lang="ts">
  import type { TaskCreate, TasksPaginatedRequest } from '$lib/api/tasks/crud-task.schema'
  import { useCreateTask, useDeleteTask, useGetPaginatedTasks } from '$lib/api/tasks/tasks'
  import TaskForm from '$lib/components/custom/TaskForm.svelte'
  import { Button } from '$lib/components/ui/button'
  import * as Card from '$lib/components/ui/card'
  import * as Table from '$lib/components/ui/table'
  import { toast } from 'svelte-sonner'

  type TaskFilter = 'all' | 'open' | 'completed' | 'dueSoon'

  const taskQuery: TasksPaginatedRequest = {
    page: 1,
    pageSize: 100,
    sortBy: 'createdAt',
    sortType: 'desc',
  }

  const tasksQuery = useGetPaginatedTasks(taskQuery)
  const createTaskMutation = useCreateTask()
  const deleteTaskMutation = useDeleteTask()

  let activeFilter = $state<TaskFilter>('all')
  let isCreateModalOpen = $state(false)
  let deletingTaskId = $state<number | null>(null)
  let selectedTaskId = $state<number | null>(null)

  const tasks = $derived(tasksQuery.data?.entries ?? [])
  const now = $derived(new Date())
  const openTasks = $derived.by(() => tasks.filter((task) => !task.completed))
  const completedTasks = $derived.by(() => tasks.filter((task) => task.completed))
  const dueSoonTasks = $derived.by(() =>
    openTasks.filter((task) => {
      if (!task.dueDate) return false

      const dueTime = new Date(task.dueDate).getTime()
      const delta = dueTime - now.getTime()
      return delta >= 0 && delta <= 1000 * 60 * 60 * 24 * 3
    }),
  )
  const filteredTasks = $derived.by(() => {
    if (activeFilter === 'open') return openTasks
    if (activeFilter === 'completed') return completedTasks
    if (activeFilter === 'dueSoon') return dueSoonTasks
    return tasks
  })
  const selectedTask = $derived.by(
    () => tasks.find((task) => task.id === selectedTaskId) ?? filteredTasks[0] ?? tasks[0] ?? null,
  )
  const completionRate = $derived.by(() =>
    tasks.length === 0 ? 0 : Math.round((completedTasks.length / tasks.length) * 100),
  )

  function setFilter(nextFilter: TaskFilter) {
    activeFilter = nextFilter
    selectedTaskId = null
  }

  function selectTask(taskId: number) {
    selectedTaskId = taskId
  }

  async function handleCreateTask(value: TaskCreate) {
    try {
      await createTaskMutation.mutateAsync(value)
      isCreateModalOpen = false
      toast.success('Task created.')
    } catch (error) {
      console.error('Task creation failed:', error)
      toast.error('Task could not be created right now.')
    }
  }

  async function handleDeleteTask(taskId: number) {
    if (typeof window !== 'undefined' && !window.confirm('Delete this task?')) {
      return
    }

    deletingTaskId = taskId

    try {
      await deleteTaskMutation.mutateAsync(String(taskId))

      if (selectedTaskId === taskId) {
        selectedTaskId = null
      }

      toast.success('Task deleted.')
    } catch (error) {
      console.error('Task deletion failed:', error)
      toast.error('Task could not be deleted right now.')
    } finally {
      deletingTaskId = null
    }
  }

  function formatDate(value: string | null | undefined) {
    if (!value) return 'No due date'

    return new Intl.DateTimeFormat('en', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    }).format(new Date(value))
  }

  function dueTone(dueDate: string | null | undefined) {
    if (!dueDate) return 'text-muted-foreground'

    const dueTime = new Date(dueDate).getTime()
    if (dueTime < now.getTime()) return 'text-destructive'
    if (dueTime - now.getTime() <= 1000 * 60 * 60 * 24 * 3) return 'text-chart-4'
    return 'text-chart-2'
  }
</script>

<svelte:head>
  <title>Taskify | Tasks</title>
</svelte:head>

<div class="mx-auto flex w-full max-w-7xl flex-col gap-6 px-4 py-6 sm:px-6 lg:px-8">
  {#if isCreateModalOpen}
    <button
      type="button"
      class="fixed inset-0 z-40 bg-background/55 backdrop-blur-[1px]"
      aria-label="Close create task dialog"
      onclick={() => (isCreateModalOpen = false)}
    ></button>
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <Card.Root class="w-full max-w-2xl rounded-4xl border-border/80 shadow-2xl">
        <Card.Header class="px-6 py-6 sm:px-7">
          <Card.Description class="tracking-[0.28em] text-chart-2 uppercase">Create task</Card.Description>
          <Card.Title class="text-2xl">Add a new task</Card.Title>
        </Card.Header>
        <Card.Content class="px-6 pb-6 sm:px-7">
          <TaskForm
            onsubmit={handleCreateTask}
            oncancel={() => (isCreateModalOpen = false)}
            pending={createTaskMutation.isPending}
            submitLabel="Create task"
          />
        </Card.Content>
      </Card.Root>
    </div>
  {/if}

  <section class="grid gap-6 lg:grid-cols-[1.15fr_0.85fr]">
    <Card.Root
      class="overflow-hidden rounded-4xl border-primary/15 bg-linear-to-br from-card via-card to-primary/6 shadow-xl shadow-primary/5"
    >
      <Card.Header class="gap-5 px-6 py-8 sm:px-8">
        <div
          class="inline-flex w-fit rounded-full border border-primary/15 bg-primary/10 px-3 py-1 text-xs font-semibold tracking-[0.28em] text-primary uppercase"
        >
          Tasks
        </div>
        <div class="space-y-3">
          <Card.Title class="max-w-3xl text-3xl font-semibold tracking-tight sm:text-4xl">
            Review the queue and focus the next move.
          </Card.Title>
          <Card.Description class="max-w-2xl text-sm leading-6 sm:text-base">
            Keep the backlog visible, inspect a task quickly, and stay close to the work that needs attention next.
          </Card.Description>
        </div>
      </Card.Header>

      <Card.Content class="px-6 pb-8 sm:px-8">
        <div class="mb-5 flex flex-wrap gap-3">
          <Button class="rounded-full px-5" onclick={() => (isCreateModalOpen = true)}>New task</Button>
        </div>
        <div class="grid gap-4 sm:grid-cols-3">
          <Card.Root size="sm" class="rounded-3xl border-border/70 bg-muted/45 py-0 shadow-none">
            <Card.Header>
              <Card.Description>All tasks</Card.Description>
              <Card.Title class="text-2xl">{tasks.length}</Card.Title>
            </Card.Header>
          </Card.Root>

          <Card.Root size="sm" class="rounded-3xl border-border/70 bg-muted/45 py-0 shadow-none">
            <Card.Header>
              <Card.Description>Open / due soon</Card.Description>
              <Card.Title class="text-2xl">{openTasks.length} / {dueSoonTasks.length}</Card.Title>
            </Card.Header>
          </Card.Root>

          <Card.Root size="sm" class="rounded-3xl border-border/70 bg-muted/45 py-0 shadow-none">
            <Card.Header>
              <Card.Description>Completion rate</Card.Description>
              <Card.Title class="text-2xl">{completionRate}%</Card.Title>
            </Card.Header>
          </Card.Root>
        </div>
      </Card.Content>
    </Card.Root>

    <Card.Root class="rounded-4xl border-border/80 shadow-sm">
      <Card.Header class="px-6 py-6 sm:px-7">
        <Card.Description class="tracking-[0.28em] text-chart-2 uppercase">Filter the queue</Card.Description>
        <Card.Title class="text-2xl">Task lanes</Card.Title>
      </Card.Header>
      <Card.Content class="px-6 pb-6 sm:px-7">
        <div class="flex flex-wrap gap-2">
          <Button
            variant={activeFilter === 'all' ? 'default' : 'outline'}
            size="sm"
            class="rounded-full"
            onclick={() => setFilter('all')}>All</Button
          >
          <Button
            variant={activeFilter === 'open' ? 'default' : 'outline'}
            size="sm"
            class="rounded-full"
            onclick={() => setFilter('open')}>Open</Button
          >
          <Button
            variant={activeFilter === 'completed' ? 'default' : 'outline'}
            size="sm"
            class="rounded-full"
            onclick={() => setFilter('completed')}>Completed</Button
          >
          <Button
            variant={activeFilter === 'dueSoon' ? 'default' : 'outline'}
            size="sm"
            class="rounded-full"
            onclick={() => setFilter('dueSoon')}>Due soon</Button
          >
        </div>
      </Card.Content>
    </Card.Root>
  </section>

  <section class="grid gap-6 lg:grid-cols-[1.05fr_0.95fr]">
    <Card.Root class="rounded-4xl border-border/80 shadow-sm">
      <Card.Header class="px-6 py-6 sm:px-7">
        <Card.Description class="tracking-[0.28em] text-chart-5 uppercase">Task list</Card.Description>
        <Card.Title class="text-2xl">Browse and inspect</Card.Title>
      </Card.Header>
      <Card.Content class="space-y-3 px-6 pb-6 sm:px-7">
        {#if tasksQuery.isLoading}
          <div class="rounded-3xl bg-muted/45 p-4 text-sm text-muted-foreground">Loading tasks...</div>
        {:else if tasksQuery.isError}
          <div class="rounded-3xl border border-destructive/20 bg-destructive/10 p-4 text-sm text-destructive">
            Tasks could not be loaded right now.
          </div>
        {:else}
          {#each filteredTasks as task (task.id)}
            <button
              type="button"
              class={`group w-full rounded-3xl border p-4 text-left transition ${selectedTask?.id === task.id ? 'border-primary/35 bg-primary/8 shadow-lg shadow-primary/8' : 'border-border/70 bg-muted/35 hover:border-primary/20 hover:bg-muted/55'}`}
              onclick={() => selectTask(task.id)}
            >
              <div class="flex items-start justify-between gap-4">
                <div class="space-y-2">
                  <div class="flex flex-wrap items-center gap-2">
                    <span
                      class={`inline-flex rounded-full px-2.5 py-1 text-[11px] font-semibold tracking-[0.24em] uppercase ${task.completed ? 'bg-chart-2/10 text-chart-2' : 'bg-secondary text-secondary-foreground'}`}
                    >
                      {task.completed ? 'Completed' : 'Open'}
                    </span>
                    <span class={['text-xs font-medium', dueTone(task.dueDate)]}>{formatDate(task.dueDate)}</span>
                  </div>
                  <div>
                    <p class="text-base font-semibold text-foreground">{task.title}</p>
                    <p class="mt-1 text-sm leading-6 text-muted-foreground">
                      {task.description ?? 'No description yet. Add some detail to make handoffs easier.'}
                    </p>
                  </div>
                </div>
                <div class="rounded-full border border-border bg-background px-3 py-1 text-xs text-muted-foreground">
                  #{task.id}
                </div>
              </div>
            </button>
          {:else}
            <div class="rounded-3xl border border-border/70 bg-muted/35 p-4 text-sm text-muted-foreground">
              No tasks match this filter yet.
            </div>
          {/each}
        {/if}
      </Card.Content>
    </Card.Root>

    <Card.Root class="rounded-4xl border-border/80 shadow-sm">
      <Card.Header class="px-6 py-6 sm:px-7">
        <Card.Description class="tracking-[0.28em] text-chart-1 uppercase">Task inspector</Card.Description>
        <Card.Title class="text-2xl">Current selection</Card.Title>
      </Card.Header>
      <Card.Content class="space-y-5 px-6 pb-6 sm:px-7">
        {#if selectedTask}
          <div class="rounded-3xl border border-border/70 bg-muted/40 p-5">
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div>
                <p class="text-xl font-semibold">{selectedTask.title}</p>
                <p class={['mt-2 text-sm font-medium', dueTone(selectedTask.dueDate)]}>
                  {formatDate(selectedTask.dueDate)}
                </p>
              </div>
              <div
                class={`rounded-full px-3 py-1 text-xs font-medium ${selectedTask.completed ? 'bg-chart-2/10 text-chart-2' : 'bg-secondary text-secondary-foreground'}`}
              >
                {selectedTask.completed ? 'Completed' : 'Open'}
              </div>
            </div>
            <p class="mt-4 text-sm leading-6 text-muted-foreground">
              {selectedTask.description ?? 'This task does not have a description yet.'}
            </p>
            <div class="mt-5 flex flex-wrap gap-3">
              <Button href={`/app/tasks/${selectedTask.id}/edit`} variant="outline" class="rounded-full px-5">
                Edit task
              </Button>
            </div>
          </div>
        {:else}
          <div class="rounded-3xl border border-border/70 bg-muted/35 p-4 text-sm text-muted-foreground">
            Pick a task to inspect it here.
          </div>
        {/if}

        <div class="space-y-3">
          <div class="h-3 overflow-hidden rounded-full bg-muted">
            <div
              class="h-full bg-linear-to-r from-primary via-chart-1 to-chart-2 transition-[width] duration-300"
              style={`width:${completionRate}%`}
            ></div>
          </div>
          <p class="text-sm text-muted-foreground">
            {completedTasks.length} of {tasks.length} tasks completed across the loaded queue.
          </p>
        </div>
      </Card.Content>
    </Card.Root>
  </section>

  <Card.Root class="rounded-4xl border-border/80 shadow-sm">
    <Card.Header class="px-6 py-6 sm:px-7">
      <Card.Description class="tracking-[0.28em] text-chart-4 uppercase">All tasks</Card.Description>
      <Card.Title class="text-2xl">Your task table</Card.Title>
    </Card.Header>
    <Card.Content class="px-0 pb-2 sm:px-0">
      {#if tasksQuery.isLoading}
        <div class="px-6 pb-4 text-sm text-muted-foreground sm:px-7">Loading task table...</div>
      {:else if tasksQuery.isError}
        <div
          class="mx-6 mb-4 rounded-3xl border border-destructive/20 bg-destructive/10 p-4 text-sm text-destructive sm:mx-7"
        >
          The task table could not be loaded right now.
        </div>
      {:else}
        <Table.Root>
          <Table.Header>
            <Table.Row class="hover:bg-transparent">
              <Table.Head class="pl-6 sm:pl-7">Task</Table.Head>
              <Table.Head>Status</Table.Head>
              <Table.Head>Due date</Table.Head>
              <Table.Head>Actions</Table.Head>
              <Table.Head class="pr-6 text-right sm:pr-7">ID</Table.Head>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {#each tasks as task (task.id)}
              <Table.Row>
                <Table.Cell class="pl-6 sm:pl-7">
                  <div class="space-y-1">
                    <button
                      type="button"
                      class="text-left font-medium text-foreground hover:text-primary"
                      onclick={() => selectTask(task.id)}
                    >
                      {task.title}
                    </button>
                    <p class="max-w-xl truncate text-sm text-muted-foreground">
                      {task.description ?? 'No description'}
                    </p>
                  </div>
                </Table.Cell>
                <Table.Cell>
                  <span
                    class={`inline-flex rounded-full px-2.5 py-1 text-[11px] font-semibold tracking-[0.24em] uppercase ${task.completed ? 'bg-chart-2/10 text-chart-2' : 'bg-secondary text-secondary-foreground'}`}
                  >
                    {task.completed ? 'Completed' : 'Open'}
                  </span>
                </Table.Cell>
                <Table.Cell>
                  <span class={['font-medium', dueTone(task.dueDate)]}>{formatDate(task.dueDate)}</span>
                </Table.Cell>
                <Table.Cell>
                  <div class="flex flex-wrap gap-2">
                    <Button href={`/app/tasks/${task.id}/edit`} variant="ghost" size="sm" class="rounded-full px-3">
                      Edit
                    </Button>
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      class="rounded-full px-3 text-destructive hover:text-destructive"
                      disabled={deletingTaskId === task.id}
                      onclick={() => handleDeleteTask(task.id)}
                    >
                      {deletingTaskId === task.id ? 'Deleting...' : 'Delete'}
                    </Button>
                  </div>
                </Table.Cell>
                <Table.Cell class="pr-6 text-right text-muted-foreground sm:pr-7">#{task.id}</Table.Cell>
              </Table.Row>
            {:else}
              <Table.Row>
                <Table.Cell colspan={5} class="px-6 py-8 text-center text-muted-foreground sm:px-7">
                  No tasks yet.
                </Table.Cell>
              </Table.Row>
            {/each}
          </Table.Body>
        </Table.Root>
      {/if}
    </Card.Content>
  </Card.Root>
</div>
