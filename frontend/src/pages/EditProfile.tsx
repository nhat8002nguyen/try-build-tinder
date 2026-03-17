import { useState, useRef } from 'react'
import { isAxiosError } from 'axios'
import { useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { ArrowLeft, Plus, X, Save } from 'lucide-react'
import toast from 'react-hot-toast'
import { userAPI, getPhotoUrl } from '../services/api'
import { useAuthStore } from '../store/authStore'
import type { User, UserPhoto } from '../types'

export default function EditProfile() {
  const navigate = useNavigate()
  const { user, updateUser } = useAuthStore()
  const fileInputRef = useRef<HTMLInputElement>(null)

  const [name, setName] = useState(user?.name || '')
  const [bio, setBio] = useState(user?.bio || '')
  const [gender, setGender] = useState<'male' | 'female' | 'other'>(user?.gender || 'other')
  const [birthdate, setBirthdate] = useState(user?.birthdate?.split('T')[0] || '')
  const [photos, setPhotos] = useState<UserPhoto[]>(user?.photos || [])
  const [isUploading, setIsUploading] = useState(false)
  const [genderPreference, setGenderPreference] = useState<('male' | 'female' | 'other')[]>(
    user?.preferences?.gender_preference || ['male']
  )
  const [minAge, setMinAge] = useState(user?.preferences?.min_age ?? 18)
  const [maxAge, setMaxAge] = useState(user?.preferences?.max_age ?? 50)
  const [maxDistance, setMaxDistance] = useState(user?.preferences?.max_distance ?? 100)

  const updateMutation = useMutation({
    mutationFn: (data: Partial<User>) => userAPI.updateProfile(data),
    onSuccess: (updatedUser) => {
      updateUser(updatedUser)
      toast.success('Profile updated!')
      navigate('/profile')
    },
    onError: () => {
      toast.error('Failed to update profile')
    },
  })

  const handlePhotoUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return

    if (photos.length >= 6) {
      toast.error('Maximum 6 photos allowed')
      return
    }

    setIsUploading(true)
    try {
      const result = await userAPI.uploadPhoto(file)
      setPhotos([
        ...photos,
        {
          id: result.id,
          photo_url: result.photo_url,
          user_id: user?.id || '',
          display_order: photos.length,
          is_approved: true,
          created_at: new Date().toISOString(),
        },
      ])
      toast.success('Photo uploaded!')
    } catch (error) {
      if (isAxiosError(error)) {
        const message =
          (error.response?.data as { error?: string } | undefined)?.error ||
          'Failed to upload photo'
        toast.error(message)
      } else {
        toast.error('Failed to upload photo')
      }
    } finally {
      setIsUploading(false)
    }
  }

  const handleDeletePhoto = async (photoId: string) => {
    try {
      await userAPI.deletePhoto(photoId)
      setPhotos(photos.filter(p => p.id !== photoId))
      toast.success('Photo deleted')
    } catch {
      toast.error('Failed to delete photo')
    }
  }

  const handleSave = () => {
    const data: Partial<User> = {
      name,
      bio,
      gender,
      preferences: {
        min_age: minAge,
        max_age: maxAge,
        max_distance: maxDistance,
        gender_preference: genderPreference,
      },
    }

    if (birthdate) {
      data.birthdate = new Date(birthdate).toISOString()
    }

    updateMutation.mutate(data)
  }

  return (
    <div className="min-h-screen pb-24">
      {/* Header */}
      <header className="glass sticky top-0 z-20 px-4 py-3 flex items-center justify-between">
        <button
          onClick={() => navigate(-1)}
          className="p-2 hover:bg-white/10 rounded-full transition-colors"
        >
          <ArrowLeft className="w-5 h-5 text-white" />
        </button>
        <h1 className="font-display text-lg font-bold text-white">Edit Profile</h1>
        <button
          onClick={handleSave}
          disabled={updateMutation.isPending}
          className="p-2 hover:bg-white/10 rounded-full transition-colors text-primary-400"
        >
          <Save className="w-5 h-5" />
        </button>
      </header>

      <div className="p-4 space-y-6">
        {/* Photos */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
        >
          <h2 className="text-white font-semibold mb-3">Photos</h2>
          <div className="grid grid-cols-3 gap-2">
            {photos.map((photo, index) => (
              <div key={photo.id} className="relative aspect-[3/4] rounded-xl overflow-hidden">
                <img
                  src={getPhotoUrl(photo.photo_url)}
                  alt={`Photo ${index + 1}`}
                  className="w-full h-full object-cover"
                />
                <button
                  onClick={() => handleDeletePhoto(photo.id)}
                  className="absolute top-2 right-2 w-6 h-6 bg-black/50 rounded-full flex items-center justify-center hover:bg-red-500 transition-colors"
                >
                  <X className="w-4 h-4 text-white" />
                </button>
              </div>
            ))}
            {photos.length < 6 && (
              <button
                onClick={() => fileInputRef.current?.click()}
                disabled={isUploading}
                className="aspect-[3/4] rounded-xl border-2 border-dashed border-white/20 flex flex-col items-center justify-center gap-2 hover:border-primary-500/50 hover:bg-primary-500/10 transition-colors disabled:opacity-50"
              >
                {isUploading ? (
                  <div className="w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full animate-spin" />
                ) : (
                  <>
                    <Plus className="w-8 h-8 text-white/40" />
                    <span className="text-white/40 text-sm">Add Photo</span>
                  </>
                )}
              </button>
            )}
          </div>
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            onChange={handlePhotoUpload}
            className="hidden"
          />
        </motion.div>

        {/* Basic Info */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="space-y-4"
        >
          <h2 className="text-white font-semibold">Basic Info</h2>
          
          <div>
            <label className="text-white/60 text-sm mb-1 block">Name</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="input-field"
              placeholder="Your name"
            />
          </div>

          <div>
            <label className="text-white/60 text-sm mb-1 block">Birthday</label>
            <input
              type="date"
              value={birthdate}
              onChange={(e) => setBirthdate(e.target.value)}
              className="input-field"
            />
          </div>

          <div>
            <label className="text-white/60 text-sm mb-1 block">Gender</label>
            <div className="grid grid-cols-3 gap-2">
              {(['male', 'female', 'other'] as const).map((g) => (
                <button
                  key={g}
                  onClick={() => setGender(g)}
                  className={`py-3 rounded-xl border transition-colors capitalize ${
                    gender === g
                      ? 'bg-primary-500/20 border-primary-500 text-primary-400'
                      : 'bg-white/5 border-white/10 text-white/60 hover:border-white/20'
                  }`}
                >
                  {g}
                </button>
              ))}
            </div>
          </div>
        </motion.div>

        {/* Bio */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
        >
          <h2 className="text-white font-semibold mb-3">About Me</h2>
          <textarea
            value={bio}
            onChange={(e) => setBio(e.target.value)}
            placeholder="Write something about yourself..."
            rows={4}
            maxLength={500}
            className="input-field resize-none"
          />
          <p className="text-white/40 text-sm mt-1 text-right">{bio.length}/500</p>
        </motion.div>

        {/* Preferences */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
        >
          <h2 className="text-white font-semibold mb-3">Dating Preferences</h2>
          <div className="card p-4 space-y-4">
            <div>
              <label className="text-white/60 text-sm mb-2 block">Age Range</label>
              <div className="flex items-center gap-4">
                <input
                  type="number"
                  min={18}
                  max={100}
                  value={minAge}
                  onChange={(e) => setMinAge(Number(e.target.value) || 18)}
                  className="input-field w-20 text-center"
                />
                <span className="text-white/40">to</span>
                <input
                  type="number"
                  min={18}
                  max={100}
                  value={maxAge}
                  onChange={(e) => setMaxAge(Number(e.target.value) || 50)}
                  className="input-field w-20 text-center"
                />
              </div>
            </div>

            <div>
              <label className="text-white/60 text-sm mb-2 block">Maximum Distance (km)</label>
              <input
                type="number"
                min={1}
                max={500}
                value={maxDistance}
                onChange={(e) => setMaxDistance(Number(e.target.value) || 1)}
                className="input-field w-24"
              />
            </div>

            <div>
              <label className="text-white/60 text-sm mb-2 block">Show Me</label>
              <div className="grid grid-cols-3 gap-2">
                {(['male', 'female', 'other'] as const).map((g) => {
                  const isSelected = genderPreference.includes(g)
                  return (
                    <button
                      key={g}
                      onClick={() => {
                        const next: ('male' | 'female' | 'other')[] =
                          g === 'other' ? ['male', 'female', 'other'] : [g]
                        setGenderPreference(next)
                      }}
                      className={`py-3 rounded-xl border transition-colors capitalize ${
                        isSelected
                          ? 'bg-primary-500/20 border-primary-500 text-primary-400'
                          : 'bg-white/5 border-white/10 text-white/60 hover:border-white/20'
                      }`}
                    >
                      {g === 'male' ? 'Men' : g === 'female' ? 'Women' : 'Everyone'}
                    </button>
                  )
                })}
              </div>
            </div>
          </div>
        </motion.div>

        {/* Save button */}
        <motion.button
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
          onClick={handleSave}
          disabled={updateMutation.isPending}
          className="btn-primary w-full"
        >
          {updateMutation.isPending ? (
            <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin mx-auto" />
          ) : (
            'Save Changes'
          )}
        </motion.button>
      </div>
    </div>
  )
}
