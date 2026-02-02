import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { MessageCircle, Heart } from 'lucide-react'
import { matchAPI } from '../services/api'
import type { Match } from '../types'

function MatchCard({ match, index }: { match: Match; index: number }) {
  const otherUser = match.other_user
  const photo = otherUser.photos?.[0]?.photo_url || 
    `https://api.dicebear.com/7.x/avataaars/svg?seed=${otherUser.id}`

  const formatTime = (dateString: string | null) => {
    if (!dateString) return ''
    const date = new Date(dateString)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const days = Math.floor(diff / (1000 * 60 * 60 * 24))
    
    if (days === 0) {
      const hours = Math.floor(diff / (1000 * 60 * 60))
      if (hours === 0) {
        const minutes = Math.floor(diff / (1000 * 60))
        return minutes <= 1 ? 'Just now' : `${minutes}m ago`
      }
      return `${hours}h ago`
    }
    if (days === 1) return 'Yesterday'
    if (days < 7) return `${days}d ago`
    return date.toLocaleDateString()
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, delay: index * 0.05 }}
    >
      <Link
        to={`/matches/${match.id}`}
        className="flex items-center gap-4 p-4 card hover:bg-white/10 transition-colors"
      >
        <div className="relative">
          <img
            src={photo}
            alt={otherUser.name}
            className="w-16 h-16 rounded-full object-cover"
          />
          {/* Online indicator - could be dynamic */}
          <div className="absolute bottom-0 right-0 w-4 h-4 bg-green-500 rounded-full border-2 border-[#1a1a2e]" />
        </div>

        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between">
            <h3 className="font-semibold text-white truncate">{otherUser.name}</h3>
            {match.last_message_at && (
              <span className="text-white/40 text-sm">{formatTime(match.last_message_at)}</span>
            )}
          </div>
          <p className="text-white/60 text-sm truncate mt-1">
            {match.last_message_at ? 'Tap to continue chatting' : 'Say hi! 👋'}
          </p>
        </div>

        <MessageCircle className="w-5 h-5 text-primary-400" />
      </Link>
    </motion.div>
  )
}

export default function Matches() {
  const { data: matches, isLoading } = useQuery({
    queryKey: ['matches'],
    queryFn: matchAPI.getMatches,
  })

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="w-12 h-12 border-4 border-primary-500 border-t-transparent rounded-full animate-spin" />
      </div>
    )
  }

  return (
    <div className="min-h-screen">
      {/* Header */}
      <header className="p-4 flex items-center justify-center border-b border-white/10">
        <h1 className="font-display text-xl font-bold text-white">Messages</h1>
      </header>

      <div className="p-4">
        {/* New Matches section */}
        {matches && matches.filter(m => !m.last_message_at).length > 0 && (
          <div className="mb-6">
            <h2 className="text-white/60 text-sm font-medium mb-3">New Matches</h2>
            <div className="flex gap-4 overflow-x-auto pb-2 -mx-4 px-4">
              {matches
                .filter(m => !m.last_message_at)
                .map((match) => (
                  <Link
                    key={match.id}
                    to={`/matches/${match.id}`}
                    className="flex-shrink-0"
                  >
                    <motion.div
                      whileHover={{ scale: 1.05 }}
                      whileTap={{ scale: 0.95 }}
                      className="relative"
                    >
                      <img
                        src={match.other_user.photos?.[0]?.photo_url || 
                          `https://api.dicebear.com/7.x/avataaars/svg?seed=${match.other_user.id}`}
                        alt={match.other_user.name}
                        className="w-20 h-20 rounded-full object-cover ring-2 ring-primary-500"
                      />
                      <div className="absolute -bottom-1 -right-1 w-6 h-6 bg-primary-500 rounded-full flex items-center justify-center">
                        <Heart className="w-3 h-3 text-white" fill="white" />
                      </div>
                    </motion.div>
                    <p className="text-white text-center text-sm mt-2 max-w-[80px] truncate">
                      {match.other_user.name}
                    </p>
                  </Link>
                ))}
            </div>
          </div>
        )}

        {/* Messages section */}
        <div>
          <h2 className="text-white/60 text-sm font-medium mb-3">Messages</h2>
          
          {matches && matches.filter(m => m.last_message_at).length > 0 ? (
            <div className="space-y-2">
              {matches
                .filter(m => m.last_message_at)
                .map((match, index) => (
                  <MatchCard key={match.id} match={match} index={index} />
                ))}
            </div>
          ) : matches && matches.length > 0 ? (
            <div className="text-center py-12">
              <MessageCircle className="w-12 h-12 text-white/20 mx-auto mb-4" />
              <p className="text-white/40">No messages yet</p>
              <p className="text-white/40 text-sm mt-1">Start a conversation with your matches!</p>
            </div>
          ) : (
            <div className="text-center py-12">
              <Heart className="w-12 h-12 text-white/20 mx-auto mb-4" />
              <p className="text-white/40">No matches yet</p>
              <p className="text-white/40 text-sm mt-1">Keep swiping to find your match!</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
