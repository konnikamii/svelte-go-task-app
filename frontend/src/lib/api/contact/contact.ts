import type { ContactCreate, ContactRequestResponse } from '$lib/api/contact/contact.schema'
import { urls } from '$lib/api/urls'
import api from '$lib/utils/fetch'
import { createMutation } from '@tanstack/svelte-query'

export const createContactRequest = async (body: ContactCreate) => {
  const response = await api.post<ContactRequestResponse>(urls.contact.createContactRequest, body)
  return response.data
}

export const useCreateContactRequest = () =>
  createMutation(() => ({
    mutationFn: createContactRequest,
  }))
