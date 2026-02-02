import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '../types'
import { authAPI } from '../services/api'

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  setUser: (user: User | null) => void
  setTokens: (accessToken: string, refreshToken: string) => void
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string, name: string) => Promise<void>
  logout: () => void
  fetchUser: () => Promise<void>
  updateUser: (user: Partial<User>) => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      isAuthenticated: false,
      isLoading: true,

      setUser: (user) => set({ user, isAuthenticated: !!user, isLoading: false }),

      setTokens: (accessToken, refreshToken) => {
        localStorage.setItem('access_token', accessToken)
        localStorage.setItem('refresh_token', refreshToken)
      },

      login: async (email, password) => {
        const { user, tokens } = await authAPI.login(email, password)
        get().setTokens(tokens.access_token, tokens.refresh_token)
        set({ user, isAuthenticated: true, isLoading: false })
      },

      register: async (email, password, name) => {
        const { user, tokens } = await authAPI.register(email, password, name)
        get().setTokens(tokens.access_token, tokens.refresh_token)
        set({ user, isAuthenticated: true, isLoading: false })
      },

      logout: () => {
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        set({ user: null, isAuthenticated: false, isLoading: false })
      },

      fetchUser: async () => {
        const token = localStorage.getItem('access_token')
        if (!token) {
          set({ isLoading: false })
          return
        }

        try {
          const user = await authAPI.getMe()
          set({ user, isAuthenticated: true, isLoading: false })
        } catch {
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
          set({ user: null, isAuthenticated: false, isLoading: false })
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
