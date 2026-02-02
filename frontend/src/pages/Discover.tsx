import { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { motion, AnimatePresence, useMotionValue, useTransform } from 'framer-motion'
import { Heart, X, MapPin, RefreshCw } from 'lucide-react'
import toast from 'react-hot-toast'
import { discoveryAPI, swipeAPI, userAPI } from '../services/api'
import type { User } from '../types'

function SwipeCard({ 
  user, 
  onSwipe, 
  isTop 
}: { 
  user: User
  onSwipe: (direction: 'like' | 'dislike') => void
  isTop: boolean 
}) {
  const [currentPhotoIndex, setCurrentPhotoIndex] = useState(0)
  const x = useMotionValue(0)
  const rotate = useTransform(x, [-200, 200], [-25, 25])
  const likeOpacity = useTransform(x, [0, 100], [0, 1])
  const nopeOpacity = useTransform(x, [-100, 0], [1, 0])

  const handleDragEnd = (_: never, info: { offset: { x: number }; velocity: { x: number } }) => {
    if (info.offset.x > 100 || info.velocity.x > 500) {
      onSwipe('like')
    } else if (info.offset.x < -100 || info.velocity.x < -500) {
      onSwipe('dislike')
    }
  }

  const photos = user.photos?.length > 0 
    ? user.photos 
    : [{ id: 'placeholder', photo_url: `https://api.dicebear.com/7.x/avataaars/svg?seed=${user.id}`, display_order: 0, user_id: user.id, is_approved: true, created_at: '' }]

  const calculateAge = (birthdate: string | null) => {
    if (!birthdate) return null
    const today = new Date()
    const birth = new Date(birthdate)
    let age = today.getFullYear() - birth.getFullYear()
    const monthDiff = today.getMonth() - birth.getMonth()
    if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birth.getDate())) {
      age--
    }
    return age
  }

  const age = calculateAge(user.birthdate)

  return (
    <motion.div
      className={`absolute inset-0 ${isTop ? 'z-10' : 'z-0'}`}
      style={{ x, rotate }}
      drag={isTop ? 'x' : false}
      dragConstraints={{ left: 0, right: 0 }}
      onDragEnd={handleDragEnd}
      initial={{ scale: isTop ? 1 : 0.95, opacity: isTop ? 1 : 0.5 }}
      animate={{ scale: isTop ? 1 : 0.95, opacity: isTop ? 1 : 0.5 }}
      exit={{ 
        x: x.get() > 0 ? 300 : -300, 
        opacity: 0,
        transition: { duration: 0.3 }
      }}
    >
      <div className="relative w-full h-full rounded-3xl overflow-hidden shadow-2xl">
        {/* Photo */}
        <div className="absolute inset-0">
          <img
            src={photos[currentPhotoIndex].photo_url}
            alt={user.name}
            className="w-full h-full object-cover"
          />
          {/* Photo indicators */}
          {photos.length > 1 && (
            <div className="absolute top-4 left-4 right-4 flex gap-1">
              {photos.map((_, index) => (
                <div
                  key={index}
                  className={`flex-1 h-1 rounded-full transition-colors ${
                    index === currentPhotoIndex ? 'bg-white' : 'bg-white/40'
                  }`}
                />
              ))}
            </div>
          )}
          {/* Photo navigation */}
          {photos.length > 1 && (
            <>
              <button
                className="absolute left-0 top-0 w-1/3 h-2/3"
                onClick={() => setCurrentPhotoIndex(Math.max(0, currentPhotoIndex - 1))}
              />
              <button
                className="absolute right-0 top-0 w-1/3 h-2/3"
                onClick={() => setCurrentPhotoIndex(Math.min(photos.length - 1, currentPhotoIndex + 1))}
              />
            </>
          )}
        </div>

        {/* Gradient overlay */}
        <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent" />

        {/* Like/Nope indicators */}
        <motion.div
          className="absolute top-20 right-6 px-4 py-2 border-4 border-green-500 rounded-lg rotate-12"
          style={{ opacity: likeOpacity }}
        >
          <span className="text-green-500 font-bold text-3xl">LIKE</span>
        </motion.div>
        <motion.div
          className="absolute top-20 left-6 px-4 py-2 border-4 border-red-500 rounded-lg -rotate-12"
          style={{ opacity: nopeOpacity }}
        >
          <span className="text-red-500 font-bold text-3xl">NOPE</span>
        </motion.div>

        {/* User info */}
        <div className="absolute bottom-0 left-0 right-0 p-6">
          <h2 className="text-white text-3xl font-bold">
            {user.name}{age && <span className="font-normal">, {age}</span>}
          </h2>
          {user.location && (user.location.latitude !== 0 || user.location.longitude !== 0) && (
            <div className="flex items-center gap-1 text-white/70 mt-1">
              <MapPin className="w-4 h-4" />
              <span className="text-sm">Nearby</span>
            </div>
          )}
          {user.bio && (
            <p className="text-white/80 mt-2 line-clamp-2">{user.bio}</p>
          )}
        </div>
      </div>
    </motion.div>
  )
}

export default function Discover() {
  const queryClient = useQueryClient()

  useEffect(() => {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          userAPI.updateLocation(position.coords.latitude, position.coords.longitude)
        },
        (error) => {
          console.log('Location permission denied:', error)
        }
      )
    }
  }, [])

  const { data: users, isLoading, refetch } = useQuery({
    queryKey: ['discover'],
    queryFn: () => discoveryAPI.getPotentialMatches({ limit: 10 }),
  })

  const swipeMutation = useMutation({
    mutationFn: ({ targetId, direction }: { targetId: string; direction: 'like' | 'dislike' }) =>
      swipeAPI.swipe(targetId, direction),
    onSuccess: (data) => {
      if (data.is_match) {
        toast.success("It's a match! 🎉", {
          duration: 3000,
          icon: '❤️',
        })
        queryClient.invalidateQueries({ queryKey: ['matches'] })
      }
    },
  })

  const [currentIndex, setCurrentIndex] = useState(0)

  const handleSwipe = (direction: 'like' | 'dislike') => {
    if (!users || currentIndex >= users.length) return

    const targetUser = users[currentIndex]
    swipeMutation.mutate({ targetId: targetUser.id, direction })
    setCurrentIndex((prev) => prev + 1)
  }

  const visibleUsers = users?.slice(currentIndex, currentIndex + 2) || []

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="w-12 h-12 border-4 border-primary-500 border-t-transparent rounded-full animate-spin" />
      </div>
    )
  }

  if (!users || users.length === 0 || currentIndex >= users.length) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center px-6">
        <div className="text-center">
          <div className="w-20 h-20 mx-auto mb-6 bg-primary-500/20 rounded-full flex items-center justify-center">
            <Heart className="w-10 h-10 text-primary-400" />
          </div>
          <h2 className="font-display text-2xl font-bold text-white mb-2">No more profiles</h2>
          <p className="text-white/60 mb-6">Check back later for new matches!</p>
          <button
            onClick={() => {
              setCurrentIndex(0)
              refetch()
            }}
            className="btn-primary flex items-center gap-2 mx-auto"
          >
            <RefreshCw className="w-5 h-5" />
            Refresh
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen flex flex-col">
      {/* Header */}
      <header className="p-4 flex items-center justify-center">
        <h1 className="font-display text-xl font-bold text-white">Discover</h1>
      </header>

      {/* Card stack */}
      <div className="flex-1 px-4 pb-32">
        <div className="relative max-w-md mx-auto h-[calc(100vh-200px)] min-h-[500px]">
          <AnimatePresence mode="popLayout">
            {visibleUsers.map((user, index) => (
              <SwipeCard
                key={user.id}
                user={user}
                onSwipe={handleSwipe}
                isTop={index === 0}
              />
            ))}
          </AnimatePresence>
        </div>
      </div>

      {/* Action buttons */}
      <div className="fixed bottom-24 left-0 right-0">
        <div className="flex items-center justify-center gap-8">
          <button
            onClick={() => handleSwipe('dislike')}
            className="w-16 h-16 bg-white/10 backdrop-blur-xl border border-white/20 rounded-full flex items-center justify-center hover:bg-red-500/20 hover:border-red-500/50 transition-all active:scale-95"
          >
            <X className="w-8 h-8 text-red-500" />
          </button>
          <button
            onClick={() => handleSwipe('like')}
            className="w-20 h-20 bg-gradient-to-br from-primary-500 to-primary-600 rounded-full flex items-center justify-center shadow-lg shadow-primary-500/30 hover:shadow-primary-500/50 transition-all active:scale-95"
          >
            <Heart className="w-10 h-10 text-white" fill="white" />
          </button>
        </div>
      </div>
    </div>
  )
}
