import { create } from 'zustand'
import { listPets } from '@/api/pet'
import type { Pet } from '@/constants/pet'

interface PetState {
  pets: Pet[]
  total: number
  load: (params?: { page?: number; page_size?: number; species?: string; status?: string; city?: string; keyword?: string }) => Promise<void>
}

export const usePetStore = create<PetState>((set) => ({
  pets: [],
  total: 0,
  load: async (params = {}) => {
    const res = await listPets(params)
    set({ pets: res.list, total: res.total })
  },
}))
