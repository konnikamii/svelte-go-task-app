import { browser } from '$app/environment'
import { urls } from '$lib/api/urls'
import { redirect } from '@sveltejs/kit'

const redirectablePaths = new Set(['/', '/login', '/register'])

export const load = async ({ fetch, url }) => {
  if (!browser || !redirectablePaths.has(url.pathname)) {
    return {}
  }

  const response = await fetch(urls.users.getMe, {
    credentials: 'include',
    headers: { accept: 'application/json' },
  })
  console.log('first')

  if (!response.ok) {
    throw redirect(307, '/login')
  }
  return {}
}
