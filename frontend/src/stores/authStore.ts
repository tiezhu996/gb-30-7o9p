import { create } from 'zustand'
import { login as apiLogin, register as apiRegister, getProfile } from '@/api/user'
import type { UserInfo } from '@/types/api'

interface AuthState {
  token: string
  user: UserInfo | null
  setToken: (t: string) => void
  setUser: (u: UserInfo | null) => void
  login: (username: string, password: string) => Promise<void>
  register: (p: { username: string; email: string; password: string; nickname?: string; phone?: string }) => Promise<void>
  fetchProfile: () => Promise<void>
  logout: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem('gbadopt_token') || '',
  user: null,
  setToken: (t) => {
    localStorage.setItem('gbadopt_token', t)
    set({ token: t })
  },
  setUser: (u) => set({ user: u }),
  login: async (username, password) => {
    const res = await apiLogin({ username, password })
    localStorage.setItem('gbadopt_token', res.token)
    set({ token: res.token, user: res.user })
  },
  register: async (payload) => {
    const res = await apiRegister(payload)
    localStorage.setItem('gbadopt_token', res.token)
    set({ token: res.token, user: res.user })
  },
  fetchProfile: async () => {
    const u = await getProfile()
    set({ user: u })
  },
  logout: () => {
    localStorage.removeItem('gbadopt_token')
    set({ token: '', user: null })
  },
}))
