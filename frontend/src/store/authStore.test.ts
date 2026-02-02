import { describe, it, expect, beforeEach } from 'vitest'
import { useAuthStore } from './authStore'

describe('Auth Store', () => {
  beforeEach(() => {
    // Clear store before each test
    useAuthStore.setState({
      user: null,
      isAuthenticated: false,
      isLoading: false,
    })
    localStorage.clear()
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
      birthdate: new Date('1990-01-01'),
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
      last_active_at: new Date(),
      created_at: new Date(),
      updated_at: new Date(),
    }

    useAuthStore.getState().setUser(testUser)
    const state = useAuthStore.getState()

    expect(state.user).toEqual(testUser)
    expect(state.isAuthenticated).toBe(true)
    expect(state.isLoading).toBe(false)
  })

  it('sets tokens correctly', () => {
    const accessToken = 'test-access-token'
    const refreshToken = 'test-refresh-token'

    useAuthStore.getState().setTokens(accessToken, refreshToken)

    expect(localStorage.getItem('access_token')).toBe(accessToken)
    expect(localStorage.getItem('refresh_token')).toBe(refreshToken)
  })

  it('logs out correctly', () => {
    // Set some data first
    localStorage.setItem('access_token', 'test-token')
    localStorage.setItem('refresh_token', 'test-refresh')
    
    useAuthStore.setState({
      user: {
        id: 'test-id',
        email: 'test@example.com',
        name: 'Test User',
      } as any,
      isAuthenticated: true,
    })

    useAuthStore.getState().logout()

    expect(useAuthStore.getState().user).toBeNull()
    expect(useAuthStore.getState().isAuthenticated).toBe(false)
    expect(localStorage.getItem('access_token')).toBeNull()
    expect(localStorage.getItem('refresh_token')).toBeNull()
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
