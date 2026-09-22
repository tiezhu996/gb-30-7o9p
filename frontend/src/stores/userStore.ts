import { create } from 'zustand'
import { getProfile, updateProfile } from '@/api/user'
import type { UserInfo } from '@/types/api'

interface UserState {
  profile: UserInfo | null
  load: () => Promise<void>
  save: (p: { nickname?: string; avatar?: string; phone?: string; bio?: string }) => Promise<void>
}

export const useUserStore = create<UserState>((set) => ({
  profile: null,
  load: async () => set({ profile: await getProfile() }),
  save: async (p) => set({ profile: await updateProfile(p) }),
}))
