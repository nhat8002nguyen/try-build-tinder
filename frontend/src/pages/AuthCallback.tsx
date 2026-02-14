import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'

export default function AuthCallback() {
  const navigate = useNavigate()
  const { fetchUser } = useAuthStore()

  useEffect(() => {
    fetchUser().then((ok) => {
      navigate(ok ? '/discover' : '/login')
    })
  }, [fetchUser, navigate])

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <div className="w-12 h-12 border-4 border-primary-500 border-t-transparent rounded-full animate-spin mx-auto mb-4" />
        <p className="text-white/70">Completing authentication...</p>
      </div>
    </div>
  )
}
