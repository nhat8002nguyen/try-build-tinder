import axios from 'axios'
import type { User, Match, Message, Notification, TokenPair, APIResponse } from '../types'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      const refreshToken = localStorage.getItem('refresh_token')
      if (refreshToken) {
        try {
          const response = await axios.post('/api/auth/refresh', {
            refresh_token: refreshToken,
          })

          const { access_token, refresh_token } = response.data.data
          localStorage.setItem('access_token', access_token)
          localStorage.setItem('refresh_token', refresh_token)

          originalRequest.headers.Authorization = `Bearer ${access_token}`
          return api(originalRequest)
        } catch {
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
          window.location.href = '/login'
        }
      }
    }

    return Promise.reject(error)
  }
)

export const authAPI = {
  register: async (email: string, password: string, name: string) => {
    const response = await api.post<APIResponse<{ user: User; tokens: TokenPair }>>('/auth/register', {
      email,
      password,
      name,
    })
    return response.data.data!
  },

  login: async (email: string, password: string) => {
    const response = await api.post<APIResponse<{ user: User; tokens: TokenPair }>>('/auth/login', {
      email,
      password,
    })
    return response.data.data!
  },

  getMe: async () => {
    const response = await api.get<APIResponse<User>>('/auth/me')
    return response.data.data!
  },

  refresh: async (refreshToken: string) => {
    const response = await api.post<APIResponse<TokenPair>>('/auth/refresh', {
      refresh_token: refreshToken,
    })
    return response.data.data!
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
