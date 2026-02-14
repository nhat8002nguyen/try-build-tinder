import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '../types'
import { authAPI } from '../services/api'
import { queryClient } from '../queryClient'

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  setUser: (user: User | null) => void
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string, name: string) => Promise<void>
  logout: () => Promise<void>
  fetchUser: () => Promise<boolean>
  updateUser: (user: Partial<User>) => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      isAuthenticated: false,
      isLoading: true,

      setUser: (user) => set({ user, isAuthenticated: !!user, isLoading: false }),

      login: async (email, password) => {
        const { user } = await authAPI.login(email, password)
        set({ user, isAuthenticated: true, isLoading: false })
      },

      register: async (email, password, name) => {
        const { user } = await authAPI.register(email, password, name)
        set({ user, isAuthenticated: true, isLoading: false })
      },

      logout: async () => {
        try {
          await authAPI.logout()
        } finally {
          queryClient.clear()
          set({ user: null, isAuthenticated: false, isLoading: false })
        }
      },

      fetchUser: async () => {
        try {
          const user = await authAPI.getMe()
          set({ user, isAuthenticated: true, isLoading: false })
          return true
        } catch {
          set({ user: null, isAuthenticated: false, isLoading: false })
          return false
        } finally {
          set({ isLoading: false })
        }
      },

      updateUser: (updates) => {
        const currentUser = get().user
        if (currentUser) {
          set({ user: { ...currentUser, ...updates } })
        }
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ isAuthenticated: state.isAuthenticated }),
      onRehydrateStorage: () => (state) => {
        if (state?.isAuthenticated) {
          state.fetchUser()
        } else {
          state?.setUser(null)
        }
      },
    }
  )
)
