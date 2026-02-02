import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Settings, Edit, LogOut, ChevronRight, Heart, Shield, Bell, HelpCircle } from 'lucide-react'
import { useAuthStore } from '../store/authStore'

export default function Profile() {
  const { user, logout } = useAuthStore()

  if (!user) return null

  const photo = user.photos?.[0]?.photo_url || 
    `https://api.dicebear.com/7.x/avataaars/svg?seed=${user.id}`

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

  const menuItems = [
    { icon: Settings, label: 'Settings', to: '/settings' },
    { icon: Bell, label: 'Notifications', to: '/notifications' },
    { icon: Shield, label: 'Privacy & Safety', to: '/privacy' },
    { icon: HelpCircle, label: 'Help & Support', to: '/help' },
  ]

  return (
    <div className="min-h-screen pb-24">
      {/* Header */}
      <header className="p-4 flex items-center justify-between">
        <h1 className="font-display text-xl font-bold text-white">Profile</h1>
        <Link to="/profile/edit" className="p-2 hover:bg-white/10 rounded-full transition-colors">
          <Edit className="w-5 h-5 text-white" />
        </Link>
      </header>

      <div className="px-4">
        {/* Profile card */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="card p-6 text-center mb-6"
        >
          <div className="relative w-32 h-32 mx-auto mb-4">
            <img
              src={photo}
              alt={user.name}
              className="w-full h-full rounded-full object-cover ring-4 ring-primary-500/30"
            />
            {user.is_verified && (
              <div className="absolute bottom-0 right-0 w-8 h-8 bg-primary-500 rounded-full flex items-center justify-center border-4 border-[#1a1a2e]">
                <Heart className="w-4 h-4 text-white" fill="white" />
              </div>
            )}
          </div>

          <h2 className="font-display text-2xl font-bold text-white">
            {user.name}{age && <span className="font-normal text-white/60">, {age}</span>}
          </h2>
          
          <p className="text-white/60 mt-1">{user.email}</p>

          {user.bio && (
            <p className="text-white/80 mt-4 max-w-xs mx-auto">{user.bio}</p>
          )}

          <Link
            to="/profile/edit"
            className="btn-primary mt-6 inline-flex items-center gap-2"
          >
            <Edit className="w-4 h-4" />
            Edit Profile
          </Link>
        </motion.div>

        {/* Stats */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="grid grid-cols-3 gap-4 mb-6"
        >
          {[
            { label: 'Photos', value: user.photos?.length || 0 },
            { label: 'Likes', value: '--' },
            { label: 'Matches', value: '--' },
          ].map(({ label, value }) => (
            <div key={label} className="card p-4 text-center">
              <p className="font-display text-2xl font-bold text-white">{value}</p>
              <p className="text-white/60 text-sm">{label}</p>
            </div>
          ))}
        </motion.div>

        {/* Menu */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="card overflow-hidden mb-6"
        >
          {menuItems.map(({ icon: Icon, label, to }, index) => (
            <Link
              key={label}
              to={to}
              className={`flex items-center justify-between p-4 hover:bg-white/5 transition-colors ${
                index !== menuItems.length - 1 ? 'border-b border-white/10' : ''
              }`}
            >
              <div className="flex items-center gap-3">
                <Icon className="w-5 h-5 text-white/60" />
                <span className="text-white">{label}</span>
              </div>
              <ChevronRight className="w-5 h-5 text-white/40" />
            </Link>
          ))}
        </motion.div>

        {/* Logout */}
        <motion.button
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
          onClick={logout}
          className="w-full card p-4 flex items-center justify-center gap-2 text-red-400 hover:bg-white/5 transition-colors"
        >
          <LogOut className="w-5 h-5" />
          <span>Log Out</span>
        </motion.button>
      </div>
    </div>
  )
}
