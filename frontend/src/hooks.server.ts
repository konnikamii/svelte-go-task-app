import type { Handle } from '@sveltejs/kit'

export const handle: Handle = async ({ event, resolve }) => {
  event.locals.requestStartedAt = new Date().toISOString()

  const response = await resolve(event)
  response.headers.set('x-demo-app', 'sveltekit-task-example')

  return response
}
