export interface User {
  id: string
  email: string
  name: string
  gender: 'male' | 'female' | 'other'
  birthdate: string | null
  bio: string
  location: Location
  preferences: UserPreferences
  is_verified: boolean
  is_active: boolean
  last_active_at: string
  created_at: string
  updated_at: string
  photos: UserPhoto[]
}

export interface Location {
  latitude: number
  longitude: number
}

export interface UserPreferences {
  min_age: number
  max_age: number
  max_distance: number
  gender_preference: ('male' | 'female' | 'other')[]
}

export interface UserPhoto {
  id: string
  user_id: string
  photo_url: string
  display_order: number
  is_approved: boolean
  created_at: string
}

export interface Match {
  id: string
  other_user: User
  matched_at: string
  last_message_at: string | null
}

export interface Message {
  id: string
  match_id: string
  sender_id: string
  content: string
  is_read: boolean
  created_at: string
  sender?: User
}

export interface Notification {
  id: string
  user_id: string
  type: 'match' | 'message' | 'like'
  payload: Record<string, unknown>
  is_read: boolean
  created_at: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  expires_in: number
}

export interface APIResponse<T> {
  success: boolean
  data?: T
  error?: string
  message?: string
}
