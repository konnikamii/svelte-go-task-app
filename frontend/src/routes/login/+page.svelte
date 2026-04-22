<script lang="ts">
  import { useLogin } from '$lib/api/auth/auth'
  import { loginSchema } from '$lib/api/auth/auth.schema'
  import { Button } from '$lib/components/ui/button/index.js'
  import * as Card from '$lib/components/ui/card/index.js'
  import * as Form from '$lib/components/ui/form/index.js'
  import { Input } from '$lib/components/ui/input/index.js'
  import { toast } from 'svelte-sonner'
  import { defaults, superForm } from 'sveltekit-superforms'
  import { zod4 } from 'sveltekit-superforms/adapters'

  const loginMutation = useLogin()

  const demoCredentials = [
    { username: 'demo-admin', email: 'admin@taskify.local', password: 'Taskify123' },
    { username: 'demo-olivia', email: 'olivia@taskify.local', password: 'Taskify123' },
    { username: 'demo-mateo', email: 'mateo@taskify.local', password: 'Taskify123' },
  ]

  const form = superForm(defaults(zod4(loginSchema)), {
    validators: zod4(loginSchema),
    SPA: true,
  })

  const { form: formData, enhance } = form

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault()

    const parsed = loginSchema.safeParse($formData)

    if (!parsed.success) {
      toast.error('Please fix the errors in the form.')
      return
    }

    try {
      toast.loading('Logging in...')
      await loginMutation.mutateAsync(parsed.data)
      toast.success('Login successful!')

      if (typeof window !== 'undefined') {
        window.location.href = '/app/dashboard'
      }
    } catch (error) {
      console.error('Login failed:', error)
      toast.error('Login failed. Please check your credentials and try again.')
    }
  }
</script>

<svelte:head>
  <title>Login</title>
</svelte:head>

<div class="mx-auto flex w-full max-w-md justify-center py-8">
  <Card.Root class="w-full">
    <Card.Header>
      <Card.Title>Login</Card.Title>
      <Card.Description>Enter your credentials below to login to your account</Card.Description>
      <Card.Action>
        <Button variant="link" href="/register">Sign Up</Button>
      </Card.Action>
    </Card.Header>
    <button onclick={() => formData.set(demoCredentials[Math.floor(Math.random() * demoCredentials.length)])}
      >Use demo</button
    >

    <form method="POST" use:enhance onsubmit={handleSubmit}>
      <Card.Content class="space-y-6 pb-6">
        <Form.Field {form} name="email">
          <Form.Control>
            {#snippet children({ props })}
              <Form.Label>Username / Email</Form.Label>
              <Input {...props} bind:value={$formData.email} placeholder="my@example.com" />
            {/snippet}
          </Form.Control>
          <Form.FieldErrors />
        </Form.Field>

        <Form.Field {form} name="password">
          <Form.Control>
            {#snippet children({ props })}
              <div class="flex items-center justify-between">
                <Form.Label>Password</Form.Label>
                <a href="##" class="text-sm underline-offset-4 hover:underline">Forgot your password?</a>
              </div>
              <Input {...props} bind:value={$formData.password} type="password" />
            {/snippet}
          </Form.Control>
          <Form.FieldErrors />
        </Form.Field>
      </Card.Content>

      <Card.Footer class=" flex-col gap-2">
        <Button type="submit" class="w-full" disabled={loginMutation.isPending}>
          {loginMutation.isPending ? 'Logging in...' : 'Login'}
        </Button>
      </Card.Footer>
    </form>
  </Card.Root>
</div>
