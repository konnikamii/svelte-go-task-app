<script lang="ts">
  import { userCreateSchema } from '$lib/api/users/crud-user.schema'
  import { useCreateUser } from '$lib/api/users/users'
  import { Button } from '$lib/components/ui/button/index.js'
  import * as Card from '$lib/components/ui/card/index.js'
  import * as Form from '$lib/components/ui/form/index.js'
  import { Input } from '$lib/components/ui/input/index.js'
  import { toast } from 'svelte-sonner'
  import { defaults, superForm } from 'sveltekit-superforms'
  import { zod4 } from 'sveltekit-superforms/adapters'

  const form = superForm(defaults(zod4(userCreateSchema)), {
    validators: zod4(userCreateSchema),
    SPA: true,
  })

  const { form: formData, enhance } = form

  const createUserMutation = useCreateUser()

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault()

    const parsed = userCreateSchema.safeParse($formData)

    if (!parsed.success) {
      toast.error('Please fix the errors in the form.')
      return
    }

    try {
      toast.loading('Registering in...')
      await createUserMutation.mutateAsync(parsed.data)
      toast.success('Registration successful!')

      if (typeof window !== 'undefined') {
        window.location.href = '/app/dashboard'
      }
    } catch (error) {
      console.error('Registration failed:', error)
      toast.error('Registration failed. Please try again.')
    }
  }
</script>

<svelte:head>
  <title>Register</title>
</svelte:head>

<div class="mx-auto flex w-full max-w-md justify-center py-8">
  <Card.Root class="w-full">
    <Card.Header>
      <Card.Title>Register</Card.Title>
      <Card.Description>Enter your credentials below to register a new account</Card.Description>
      <Card.Action>
        <Button variant="link" href="/login">Login</Button>
      </Card.Action>
    </Card.Header>

    <form method="POST" use:enhance onsubmit={handleSubmit}>
      <Card.Content class="space-y-6 pb-6">
        <Form.Field {form} name="username">
          <Form.Control>
            {#snippet children({ props })}
              <Form.Label>Username</Form.Label>
              <Input {...props} bind:value={$formData.username} placeholder="Username..." />
            {/snippet}
          </Form.Control>
          <Form.FieldErrors />
        </Form.Field>
        <Form.Field {form} name="email">
          <Form.Control>
            {#snippet children({ props })}
              <Form.Label>Email</Form.Label>
              <Input {...props} bind:value={$formData.email} placeholder="my@example.com" />
            {/snippet}
          </Form.Control>
          <Form.FieldErrors />
        </Form.Field>

        <Form.Field {form} name="password">
          <Form.Control>
            {#snippet children({ props })}
              <Form.Label>Password</Form.Label>
              <Input {...props} bind:value={$formData.password} type="password" placeholder="Password..." />
            {/snippet}
          </Form.Control>
          <Form.FieldErrors />
        </Form.Field>
      </Card.Content>

      <Card.Footer class=" flex-col gap-2">
        <Button type="submit" class="w-full" disabled={createUserMutation.isPending}>
          {createUserMutation.isPending ? 'Loading...' : 'Register'}
        </Button>
      </Card.Footer>
    </form>
  </Card.Root>
</div>
