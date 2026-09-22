import request from '@/utils/request'
import type { Pet } from '@/constants/pet'
import type { PageData } from '@/types/api'

export function listPets(params: { page?: number; page_size?: number; species?: string; status?: string; city?: string; keyword?: string }) {
  return request.get<never, PageData<Pet>>('/pets', { params })
}

export function getPet(id: number | string) {
  return request.get<never, Pet>(`/pets/${id}`)
}

export function publishPet(payload: Partial<Pet>) {
  return request.post<never, Pet>('/pets', payload)
}

export function updatePetStatus(id: number, status: string) {
  return request.put<never, Pet>(`/pets/${id}/status`, { status })
}
