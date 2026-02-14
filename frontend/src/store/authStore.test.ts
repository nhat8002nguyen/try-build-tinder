import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useAuthStore } from './authStore'
import { authAPI } from '../services/api'

vi.mock('../services/api', () => ({
  authAPI: {
    logout: vi.fn().mockResolvedValue(undefined),
  },
}))

describe('Auth Store', () => {
  beforeEach(() => {
    useAuthStore.setState({
      user: null,
      isAuthenticated: false,
      isLoading: false,
    })
    vi.mocked(authAPI.logout).mockClear()
  })

  it('initializes with default values', () => {
    const state = useAuthStore.getState()
    expect(state.user).toBeNull()
    expect(state.isAuthenticated).toBe(false)
    expect(state.isLoading).toBe(false)
  })

  it('sets user correctly', () => {
    const testUser = {
      id: 'test-id',
      email: 'test@example.com',
      name: 'Test User',
      gender: 'male' as const,
      birthdate: '1990-01-01',
      bio: '',
      photos: [],
      location: { latitude: 0, longitude: 0 },
      preferences: {
        min_age: 18,
        max_age: 50,
        max_distance: 100,
        gender_preference: ['female' as const],
      },
      is_verified: false,
      is_active: true,
      last_active_at: new Date().toISOString(),
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }

    useAuthStore.getState().setUser(testUser)
    const state = useAuthStore.getState()

    expect(state.user).toEqual(testUser)
    expect(state.isAuthenticated).toBe(true)
    expect(state.isLoading).toBe(false)
  })

  it('logs out correctly', async () => {
    useAuthStore.setState({
      user: {
        id: 'test-id',
        email: 'test@example.com',
        name: 'Test User',
      } as any,
      isAuthenticated: true,
    })

    await useAuthStore.getState().logout()

    expect(authAPI.logout).toHaveBeenCalled()
    expect(useAuthStore.getState().user).toBeNull()
    expect(useAuthStore.getState().isAuthenticated).toBe(false)
  })

  it('updates user correctly', () => {
    const initialUser = {
      id: 'test-id',
      email: 'test@example.com',
      name: 'Test User',
      bio: 'Old bio',
    } as any

    useAuthStore.setState({ user: initialUser })

    useAuthStore.getState().updateUser({ bio: 'New bio' })

    expect(useAuthStore.getState().user?.bio).toBe('New bio')
    expect(useAuthStore.getState().user?.name).toBe('Test User')
  })
})
