<script lang="ts">
  import { page } from '$app/state'
  import type { TaskUpdate } from '$lib/api/tasks/crud-task.schema'
  import { getTaskById, useDeleteTask, useUpdateTask } from '$lib/api/tasks/tasks'
  import TaskForm from '$lib/components/custom/TaskForm.svelte'
  import { Button } from '$lib/components/ui/button'
  import * as Card from '$lib/components/ui/card'
  import { createQuery } from '@tanstack/svelte-query'
  import { toast } from 'svelte-sonner'

  const currentTaskId = $derived(page.params.id ?? '')
  const taskQuery = createQuery(() => ({
    queryKey: ['task', currentTaskId],
    queryFn: () => getTaskById(currentTaskId),
    enabled: currentTaskId.length > 0,
  }))
  const deleteTaskMutation = useDeleteTask()
  const updateTaskMutation = useUpdateTask()

  async function handleUpdateTask(value: {
    title: string
    description?: string | null
    dueDate?: string | null
    completed?: boolean
  }) {
    if (!currentTaskId) {
      toast.error('Task id is missing.')
      return
    }

    const payload: TaskUpdate = {
      title: value.title,
      description: value.description ?? null,
      dueDate: value.dueDate ?? null,
      completed: value.completed ?? false,
    }

    try {
      await updateTaskMutation.mutateAsync({ id: currentTaskId, body: payload })
      toast.success('Task updated.')

      if (typeof window !== 'undefined') {
        window.location.href = '/app/tasks'
      }
    } catch (error) {
      console.error('Task update failed:', error)
      toast.error('Task could not be updated right now.')
    }
  }

  async function handleDeleteTask() {
    if (!currentTaskId) {
      toast.error('Task id is missing.')
      return
    }

    if (typeof window !== 'undefined' && !window.confirm('Delete this task?')) {
      return
    }

    try {
      await deleteTaskMutation.mutateAsync(currentTaskId)
      toast.success('Task deleted.')

      if (typeof window !== 'undefined') {
        window.location.href = '/app/tasks'
      }
    } catch (error) {
      console.error('Task deletion failed:', error)
      toast.error('Task could not be deleted right now.')
    }
  }
</script>

<svelte:head>
  <title>Taskify | Edit Task</title>
</svelte:head>

<div class="mx-auto flex w-full max-w-5xl flex-col gap-6 px-4 py-6 sm:px-6 lg:px-8">
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
        <Card.Title class="max-w-3xl text-3xl font-semibold tracking-tight sm:text-4xl">Edit task</Card.Title>
        <Card.Description class="max-w-2xl text-sm leading-6 sm:text-base">
          Update the details, due date, or completion state for this task.
        </Card.Description>
      </div>
    </Card.Header>
  </Card.Root>

  <Card.Root class="rounded-4xl border-border/80 shadow-sm">
    <Card.Header class="px-6 py-6 sm:px-7">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <Card.Description class="tracking-[0.28em] text-chart-2 uppercase">Task editor</Card.Description>
          <Card.Title class="text-2xl">Update task details</Card.Title>
        </div>
        <div class="flex flex-wrap gap-3">
          <Button
            type="button"
            variant="outline"
            class="rounded-full px-5 text-destructive hover:text-destructive"
            disabled={deleteTaskMutation.isPending}
            onclick={handleDeleteTask}
          >
            {deleteTaskMutation.isPending ? 'Deleting...' : 'Delete task'}
          </Button>
          <Button href="/app/tasks" variant="outline" class="rounded-full px-5">Back to tasks</Button>
        </div>
      </div>
    </Card.Header>
    <Card.Content class="px-6 pb-6 sm:px-7">
      {#if taskQuery.isLoading}
        <div class="rounded-3xl bg-muted/45 p-4 text-sm text-muted-foreground">Loading task...</div>
      {:else if taskQuery.isError || !taskQuery.data}
        <div class="rounded-3xl border border-destructive/20 bg-destructive/10 p-4 text-sm text-destructive">
          This task could not be loaded right now.
        </div>
      {:else}
        <TaskForm
          initialValues={{
            title: taskQuery.data.title,
            description: taskQuery.data.description,
            dueDate: taskQuery.data.dueDate,
            completed: taskQuery.data.completed,
          }}
          onsubmit={handleUpdateTask}
          pending={updateTaskMutation.isPending}
          submitLabel="Save changes"
        />
      {/if}
    </Card.Content>
  </Card.Root>
</div>
