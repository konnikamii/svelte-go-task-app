<script lang="ts">
  import { taskCreateSchema, type TaskCreate } from '$lib/api/tasks/crud-task.schema'
  import { Button } from '$lib/components/ui/button'
  import Calendar from '$lib/components/ui/calendar/calendar.svelte'
  import { Input } from '$lib/components/ui/input'
  import { Label } from '$lib/components/ui/label'
  import * as Popover from '$lib/components/ui/popover'
  import { Textarea } from '$lib/components/ui/textarea'
  import { getLocalTimeZone, parseDate, type CalendarDate } from '@internationalized/date'
  import ChevronDownIcon from '@lucide/svelte/icons/chevron-down'

  type TaskFormValues = {
    title: string
    description: string | null
    dueDate: string | null
    completed: boolean
  }

  type TaskFormErrors = Partial<Record<keyof TaskFormValues, string>>

  let {
    initialValues,
    pending = false,
    submitLabel = 'Save task',
    onsubmit,
    oncancel,
  }: {
    initialValues?: Partial<TaskFormValues>
    pending?: boolean
    submitLabel?: string
    onsubmit: (value: TaskCreate) => void | Promise<void>
    oncancel?: () => void
  } = $props()

  function toCalendarDate(value: string | null | undefined) {
    if (!value) return undefined

    const date = new Date(value)
    if (Number.isNaN(date.getTime())) return undefined

    const year = String(date.getFullYear()).padStart(4, '0')
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return parseDate(`${year}-${month}-${day}`)
  }

  function toIsoFromCalendarDate(value: CalendarDate | undefined) {
    if (!value) return null

    const date = value.toDate(getLocalTimeZone())
    if (Number.isNaN(date.getTime())) return null
    return date.toISOString()
  }

  function getSeededValues(): TaskFormValues {
    return {
      title: initialValues?.title ?? '',
      description: initialValues?.description ?? null,
      dueDate: initialValues?.dueDate ?? null,
      completed: initialValues?.completed ?? false,
    }
  }

  const seededValues = getSeededValues()

  let formValues = $state<TaskFormValues>(seededValues)
  let dueDateOpen = $state(false)
  let dueDateValue = $state<CalendarDate | undefined>(toCalendarDate(seededValues.dueDate))

  let errors = $state<TaskFormErrors>({})

  const dueDateLabel = $derived.by(() => {
    if (!dueDateValue) return 'Select due date'

    return dueDateValue.toDate(getLocalTimeZone()).toLocaleDateString()
  })

  function handleDueDateChange(value: CalendarDate | undefined) {
    dueDateValue = value
    formValues.dueDate = toIsoFromCalendarDate(value)
    dueDateOpen = false
  }

  function clearDueDate() {
    dueDateValue = undefined
    formValues.dueDate = null
  }

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault()

    const payload = {
      title: formValues.title,
      description: formValues.description?.trim() ? formValues.description : null,
      dueDate: formValues.dueDate,
      completed: formValues.completed,
    }

    const parsed = taskCreateSchema.safeParse(payload)

    if (!parsed.success) {
      const nextErrors: TaskFormErrors = {}

      for (const issue of parsed.error.issues) {
        const field = issue.path[0] as keyof TaskFormValues | undefined
        if (field && !nextErrors[field]) {
          nextErrors[field] = issue.message
        }
      }

      errors = nextErrors
      return
    }

    errors = {}
    await onsubmit(parsed.data)
  }
</script>

<form class="space-y-6" onsubmit={handleSubmit}>
  <div class="space-y-5">
    <div class="space-y-2">
      <label class="text-sm font-medium text-foreground" for="task-title">Title</label>
      <Input id="task-title" bind:value={formValues.title} placeholder="Task title" />
      {#if errors.title}
        <p class="text-sm text-destructive">{errors.title}</p>
      {/if}
    </div>

    <div class="space-y-2">
      <label class="text-sm font-medium text-foreground" for="task-description">Description</label>
      <Textarea
        id="task-description"
        bind:value={formValues.description}
        placeholder="Add context, notes, or next steps..."
      />
      {#if errors.description}
        <p class="text-sm text-destructive">{errors.description}</p>
      {/if}
    </div>

    <div class="space-y-2">
      <Label class="text-sm font-medium text-foreground" for="task-due-date">Due date</Label>
      <div class="flex flex-wrap items-center gap-3">
        <Popover.Root bind:open={dueDateOpen}>
          <Popover.Trigger id="task-due-date">
            {#snippet child({ props })}
              <Button
                {...props}
                variant="outline"
                class="w-full justify-between rounded-2xl border-border/70 bg-background px-4 py-5 font-normal sm:w-64"
              >
                {dueDateLabel}
                <ChevronDownIcon class="size-4 opacity-70" />
              </Button>
            {/snippet}
          </Popover.Trigger>
          <Popover.Content class="w-auto overflow-hidden rounded-3xl p-0" align="start">
            <Calendar
              type="single"
              bind:value={dueDateValue}
              onValueChange={() => handleDueDateChange(dueDateValue)}
              captionLayout="dropdown"
            />
          </Popover.Content>
        </Popover.Root>

        {#if dueDateValue}
          <Button type="button" variant="ghost" class="rounded-full px-4" onclick={clearDueDate}>Clear</Button>
        {/if}
      </div>
      {#if errors.dueDate}
        <p class="text-sm text-destructive">{errors.dueDate}</p>
      {/if}
    </div>

    <label
      class="flex items-center gap-3 rounded-2xl border border-border/70 bg-muted/35 px-4 py-3 text-sm text-foreground"
    >
      <input bind:checked={formValues.completed} type="checkbox" class="size-4 rounded border-border" />
      <span>Mark as completed</span>
    </label>
  </div>

  <div class="flex flex-wrap justify-end gap-3">
    {#if oncancel}
      <Button type="button" variant="outline" class="rounded-full px-5" onclick={oncancel}>Cancel</Button>
    {/if}
    <Button type="submit" class="rounded-full px-5" disabled={pending}>
      {pending ? 'Saving...' : submitLabel}
    </Button>
  </div>
</form>
