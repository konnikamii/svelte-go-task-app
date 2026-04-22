<script lang="ts">
  import { useCreateContactRequest } from '$lib/api/contact/contact'
  import { contactCreateSchema } from '$lib/api/contact/contact.schema'
  import { Button } from '$lib/components/ui/button'
  import * as Card from '$lib/components/ui/card'
  import * as Form from '$lib/components/ui/form'
  import { Input } from '$lib/components/ui/input'
  import { Textarea } from '$lib/components/ui/textarea'
  import { toast } from 'svelte-sonner'
  import { defaults, superForm } from 'sveltekit-superforms'
  import { zod4 } from 'sveltekit-superforms/adapters'

  const contactMutation = useCreateContactRequest()

  const form = superForm(defaults(zod4(contactCreateSchema)), {
    validators: zod4(contactCreateSchema),
    SPA: true,
  })

  const { form: formData, enhance, reset } = form

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault()

    const parsed = contactCreateSchema.safeParse($formData)
    if (!parsed.success) {
      toast.error('Please fix the errors in the contact form.')
      return
    }

    try {
      toast.loading('Sending message...')
      await contactMutation.mutateAsync(parsed.data)
      toast.success('Your message has been sent.')
      reset()
    } catch (error) {
      console.error('Contact request failed:', error)
      toast.error('We could not send your message right now. Please try again.')
    }
  }
</script>

<svelte:head>
  <title>Taskify | Contact</title>
  <meta
    name="description"
    content="Send a question to Taskify. Use the public contact form to reach out with your email, a title, and a message."
  />
</svelte:head>

<div class="mx-auto flex w-full max-w-6xl flex-col gap-6 px-4 py-6 sm:px-6 lg:px-8">
  <section class="grid gap-6 lg:grid-cols-[0.95fr_1.05fr]">
    <Card.Root
      class="overflow-hidden rounded-4xl border-primary/15 bg-linear-to-br from-card via-card to-primary/6 shadow-xl shadow-primary/5"
    >
      <Card.Header class="gap-5 px-6 py-8 sm:px-8">
        <div
          class="inline-flex w-fit rounded-full border border-primary/15 bg-primary/10 px-3 py-1 text-xs font-semibold tracking-[0.28em] text-primary uppercase"
        >
          Contact
        </div>
        <div class="space-y-3">
          <Card.Title class="max-w-2xl text-3xl font-semibold tracking-tight sm:text-4xl">
            Ask a question without creating an account.
          </Card.Title>
          <Card.Description class="max-w-xl text-sm leading-6 sm:text-base">
            If someone wants to reach out before onboarding, they can send a public contact request here with an email,
            a title, and a clear message.
          </Card.Description>
        </div>
      </Card.Header>
      <Card.Content class="px-6 pb-8 sm:px-8">
        <div class="rounded-3xl border border-border/70 bg-muted/35 p-5 sm:p-6">
          <p class="max-w-2xl text-sm leading-7 text-muted-foreground sm:text-base">
            Use this page for product questions, partnership requests, onboarding help, or anything that needs a direct
            response from the team.
          </p>
        </div>
      </Card.Content>
    </Card.Root>

    <Card.Root class="rounded-4xl border-border/80 shadow-sm">
      <Card.Header class="px-6 py-6 sm:px-7">
        <Card.Description class="tracking-[0.28em] text-chart-2 uppercase">Contact form</Card.Description>
        <Card.Title class="text-2xl">Send a message</Card.Title>
      </Card.Header>

      <form method="POST" use:enhance onsubmit={handleSubmit}>
        <Card.Content class="space-y-6 px-6 pb-6 sm:px-7">
          <Form.Field {form} name="email">
            <Form.Control>
              {#snippet children({ props })}
                <Form.Label>Email</Form.Label>
                <Input {...props} bind:value={$formData.email} placeholder="you@example.com" />
              {/snippet}
            </Form.Control>
            <Form.FieldErrors />
          </Form.Field>

          <Form.Field {form} name="title">
            <Form.Control>
              {#snippet children({ props })}
                <Form.Label>Title</Form.Label>
                <Input {...props} bind:value={$formData.title} placeholder="What would you like to ask?" />
              {/snippet}
            </Form.Control>
            <Form.FieldErrors />
          </Form.Field>

          <Form.Field {form} name="message">
            <Form.Control>
              {#snippet children({ props })}
                <Form.Label>Message</Form.Label>
                <Textarea
                  {...props}
                  bind:value={$formData.message}
                  placeholder="Describe your question, issue, or request..."
                />
              {/snippet}
            </Form.Control>
            <Form.FieldErrors />
          </Form.Field>
        </Card.Content>

        <Card.Footer class="flex-col items-stretch gap-3 px-6 pb-6 sm:px-7">
          <Button type="submit" class="w-full" disabled={contactMutation.isPending}>
            {contactMutation.isPending ? 'Sending...' : 'Send message'}
          </Button>
        </Card.Footer>
      </form>
    </Card.Root>
  </section>
</div>
