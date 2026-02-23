import axios from 'axios'
import type { User, Match, Message, Notification, APIResponse } from '../types'

export function getApiOrigin(): string {
  const env = import.meta.env.VITE_API_ORIGIN
  if (!env || typeof env !== 'string') return ''
  return String(env).replace(/\/$/, '')
}

export function getPhotoUrl(photoUrl: string): string {
  if (!photoUrl) return photoUrl
  if (photoUrl.startsWith('http://') || photoUrl.startsWith('https://')) return photoUrl
  const origin = getApiOrigin()
  return origin ? origin + photoUrl : photoUrl
}

const api = axios.create({
  baseURL: getApiOrigin() ? getApiOrigin() + '/api' : '/api',
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
})

let refreshPromise: Promise<boolean> | null = null

function redirectToLogin(): void {
  window.location.href = '/login'
}

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    if (error.response?.status !== 401 || originalRequest._retry) {
      return Promise.reject(error)
    }

    const isRefreshRequest = originalRequest.url?.includes('/auth/refresh')
    if (isRefreshRequest) {
      redirectToLogin()
      return Promise.reject(error)
    }

    originalRequest._retry = true

    if (!refreshPromise) {
      refreshPromise = api
        .post('/auth/refresh')
        .then(() => true)
        .catch((err) => {
          if (err.response?.status === 401) {
            redirectToLogin()
          }
          return false
        })
        .finally(() => {
          refreshPromise = null
        })
    }

    const success = await refreshPromise
    if (!success) {
      return Promise.reject(error)
    }

    return api(originalRequest)
  }
)

export const authAPI = {
  register: async (email: string, password: string, name: string) => {
    const response = await api.post<APIResponse<{ user: User }>>('/auth/register', {
      email,
      password,
      name,
    })
    return response.data.data!
  },

  login: async (email: string, password: string) => {
    const response = await api.post<APIResponse<{ user: User }>>('/auth/login', {
      email,
      password,
    })
    return response.data.data!
  },

  getMe: async () => {
    const response = await api.get<APIResponse<User>>('/auth/me')
    return response.data.data!
  },

  refresh: async () => {
    await api.post('/auth/refresh')
  },

  logout: async () => {
    await api.post('/auth/logout')
  },
}

export const userAPI = {
  getUser: async (userId: string) => {
    const response = await api.get<APIResponse<User>>(`/users/${userId}`)
    return response.data.data!
  },

  updateProfile: async (data: Partial<User>) => {
    const response = await api.put<APIResponse<User>>('/users/me', data)
    return response.data.data!
  },

  uploadPhoto: async (file: File) => {
    const formData = new FormData()
    formData.append('photo', file)
    const response = await api.post<APIResponse<{ id: string; photo_url: string }>>('/users/me/photos', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data.data!
  },

  deletePhoto: async (photoId: string) => {
    await api.delete(`/users/me/photos/${photoId}`)
  },

  updateLocation: async (latitude: number, longitude: number) => {
    await api.put('/users/me/location', { latitude, longitude })
  },
}

export const discoveryAPI = {
  getPotentialMatches: async (params?: {
    limit?: number
    min_age?: number
    max_age?: number
    max_distance?: number
    gender_preference?: string
  }) => {
    const response = await api.get<APIResponse<User[]>>('/discover', { params })
    return response.data.data!
  },
}

export const swipeAPI = {
  swipe: async (targetId: string, direction: 'like' | 'dislike') => {
    const response = await api.post<APIResponse<{ is_match: boolean; match?: Match }>>('/swipes', {
      target_id: targetId,
      direction,
    })
    return response.data.data!
  },
}

export const matchAPI = {
  getMatches: async () => {
    const response = await api.get<APIResponse<Match[]>>('/matches')
    return response.data.data!
  },

  getMatch: async (matchId: string) => {
    const response = await api.get<APIResponse<Match>>(`/matches/${matchId}`)
    return response.data.data!
  },
}

export const messageAPI = {
  getMessages: async (matchId: string, limit = 50, offset = 0) => {
    const response = await api.get<APIResponse<Message[]>>(`/matches/${matchId}/messages`, {
      params: { limit, offset },
    })
    return response.data.data!
  },

  sendMessage: async (matchId: string, content: string) => {
    const response = await api.post<APIResponse<Message>>(`/matches/${matchId}/messages`, { content })
    return response.data.data!
  },
}

export const notificationAPI = {
  getNotifications: async (limit = 20, offset = 0) => {
    const response = await api.get<APIResponse<{ notifications: Notification[]; unread_count: number }>>(
      '/notifications',
      { params: { limit, offset } }
    )
    return response.data.data!
  },

  markAsRead: async (notificationId: string) => {
    await api.put(`/notifications/${notificationId}/read`)
  },
}

export default api
