import { create } from 'zustand'
import { listMyApplications } from '@/api/application'
import type { AdoptionApplication } from '@/types/api'

interface ApplicationState {
  applications: AdoptionApplication[]
  load: () => Promise<void>
}

export const useApplicationStore = create<ApplicationState>((set) => ({
  applications: [],
  load: async () => set({ applications: await listMyApplications() }),
}))
