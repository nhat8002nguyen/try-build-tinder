import { useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'

export default function AuthCallback() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const { setTokens, fetchUser } = useAuthStore()

  useEffect(() => {
    const accessToken = searchParams.get('access_token')
    const refreshToken = searchParams.get('refresh_token')

    if (accessToken && refreshToken) {
      setTokens(accessToken, refreshToken)
      fetchUser().then(() => {
        navigate('/discover')
      })
    } else {
      navigate('/login')
    }
  }, [searchParams, setTokens, fetchUser, navigate])

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <div className="w-12 h-12 border-4 border-primary-500 border-t-transparent rounded-full animate-spin mx-auto mb-4" />
        <p className="text-white/70">Completing authentication...</p>
      </div>
    </div>
  )
}
