import { type LoginRequest } from '$lib/api/auth/auth.schema'
import { queryClient } from '$lib/api/client'
import { urls } from '$lib/api/urls'
import type { User } from '$lib/api/users/user.schema'
import api from '$lib/utils/fetch'
import { createMutation } from '@tanstack/svelte-query'

export const login = async (body: LoginRequest) => {
  const formData = new FormData()
  formData.set('email', body.email)
  formData.set('password', body.password)

  const response = await api.post<User>(urls.auth.login, formData)
  return response.data
}

export const useLogin = () =>
  createMutation(() => ({
    mutationFn: login,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user'] })
    },
  }))

export const logout = async () => {
  const response = await api.post(urls.auth.logout)
  return response.data
}

export const useLogout = () =>
  createMutation(() => ({
    mutationFn: logout,
    onSuccess: () => {
      queryClient.invalidateQueries()

      if (typeof window !== 'undefined') {
        window.location.href = '/login'
      }
    },
  }))
