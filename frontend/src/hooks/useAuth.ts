import { useAuthStore } from '@/stores/authStore'

export function useAuth() {
  const token = useAuthStore((s) => s.token)
  const user = useAuthStore((s) => s.user)
  const logout = useAuthStore((s) => s.logout)
  const fetchProfile = useAuthStore((s) => s.fetchProfile)
  const isLoggedIn = !!token
  const isAdmin = user?.role === 'admin'
  const isOrg = user?.role === 'org'
  return { token, user, logout, fetchProfile, isLoggedIn, isAdmin, isOrg }
}
