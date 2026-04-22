import { browser } from '$app/environment'
import { urls } from '$lib/api/urls'
import { redirect } from '@sveltejs/kit'

export const ssr = false

const redirectablePaths = new Set(['/', '/login', '/register'])

export const load = async ({ fetch, url }) => {
  if (!browser) {
    return {}
  }

  const response = await fetch(urls.users.getMe, {
    credentials: 'include',
    headers: { accept: 'application/json' },
  })
  if (response.status !== 401 && redirectablePaths.has(url.pathname)) {
    throw redirect(307, '/app/dashboard')
  } else if (response.status === 401 && url.pathname.startsWith('/app')) {
    throw redirect(307, '/login')
  }
  return {}
}
