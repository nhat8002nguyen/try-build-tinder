import { useState, useEffect, useRef, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { motion, AnimatePresence } from 'framer-motion'
import { ArrowLeft, Send, MoreVertical } from 'lucide-react'
import { matchAPI, messageAPI, getPhotoUrl } from '../services/api'
import { useAuthStore } from '../store/authStore'
import { useWebSocket } from '../contexts/WebSocketContext'
import type { Message } from '../types'

export default function Chat() {
  const { matchId } = useParams<{ matchId: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { user } = useAuthStore()
  const { addMessageHandler, removeMessageHandler, sendTyping } = useWebSocket()
  
  const [newMessage, setNewMessage] = useState('')
  const [isTyping, setIsTyping] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  const { data: match, isLoading: matchLoading } = useQuery({
    queryKey: ['match', matchId],
    queryFn: () => matchAPI.getMatch(matchId!),
    enabled: !!matchId,
  })

  const { data: messages, isLoading: messagesLoading } = useQuery({
    queryKey: ['messages', matchId],
    queryFn: () => messageAPI.getMessages(matchId!),
    enabled: !!matchId,
    refetchInterval: 5000,
  })

  const sendMessageMutation = useMutation({
    mutationFn: (content: string) => messageAPI.sendMessage(matchId!, content),
    onSuccess: (newMsg) => {
      queryClient.setQueryData(['messages', matchId], (old: Message[] | undefined) => 
        old ? [...old, newMsg] : [newMsg]
      )
      setNewMessage('')
    },
  })

  const handleWebSocketMessage = useCallback((wsMessage: { type: string; payload: Record<string, unknown> }) => {
    if (wsMessage.type === 'message' && wsMessage.payload.match_id === matchId) {
      queryClient.invalidateQueries({ queryKey: ['messages', matchId] })
    }
    if (wsMessage.type === 'typing' && wsMessage.payload.match_id === matchId) {
      setIsTyping(true)
      setTimeout(() => setIsTyping(false), 3000)
    }
  }, [matchId, queryClient])

  useEffect(() => {
    addMessageHandler(handleWebSocketMessage)
    return () => removeMessageHandler(handleWebSocketMessage)
  }, [addMessageHandler, removeMessageHandler, handleWebSocketMessage])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const handleSend = (e: React.FormEvent) => {
    e.preventDefault()
    if (!newMessage.trim()) return
    sendMessageMutation.mutate(newMessage.trim())
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNewMessage(e.target.value)
    if (matchId) {
      sendTyping(matchId)
    }
  }

  if (matchLoading || messagesLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="w-12 h-12 border-4 border-primary-500 border-t-transparent rounded-full animate-spin" />
      </div>
    )
  }

  if (!match) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-white/60">Match not found</p>
      </div>
    )
  }

  const otherUser = match.other_user
  const rawPhoto = otherUser.photos?.[0]?.photo_url
  const photo = rawPhoto ? getPhotoUrl(rawPhoto) : `https://api.dicebear.com/7.x/avataaars/svg?seed=${otherUser.id}`

  return (
    <div className="min-h-screen flex flex-col">
      {/* Header */}
      <header className="glass sticky top-0 z-20 px-4 py-3 flex items-center gap-4">
        <button
          onClick={() => navigate('/matches')}
          className="p-2 hover:bg-white/10 rounded-full transition-colors"
        >
          <ArrowLeft className="w-5 h-5 text-white" />
        </button>

        <div className="flex items-center gap-3 flex-1">
          <img
            src={photo}
            alt={otherUser.name}
            className="w-10 h-10 rounded-full object-cover"
          />
          <div>
            <h2 className="font-semibold text-white">{otherUser.name}</h2>
            <AnimatePresence>
              {isTyping && (
                <motion.p
                  initial={{ opacity: 0, y: -10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -10 }}
                  className="text-primary-400 text-xs"
                >
                  typing...
                </motion.p>
              )}
            </AnimatePresence>
          </div>
        </div>

        <button className="p-2 hover:bg-white/10 rounded-full transition-colors">
          <MoreVertical className="w-5 h-5 text-white" />
        </button>
      </header>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-3">
        {messages?.map((message, index) => {
          const isOwn = message.sender_id === user?.id
          const showAvatar = !isOwn && (
            index === 0 || messages[index - 1]?.sender_id !== message.sender_id
          )

          return (
            <motion.div
              key={message.id}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              className={`flex ${isOwn ? 'justify-end' : 'justify-start'}`}
            >
              <div className={`flex items-end gap-2 max-w-[80%] ${isOwn ? 'flex-row-reverse' : ''}`}>
                {!isOwn && (
                  <div className="w-8 h-8">
                    {showAvatar && (
                      <img
                        src={photo}
                        alt={otherUser.name}
                        className="w-8 h-8 rounded-full object-cover"
                      />
                    )}
                  </div>
                )}
                <div
                  className={`px-4 py-2 rounded-2xl ${
                    isOwn
                      ? 'bg-primary-500 text-white rounded-br-md'
                      : 'bg-white/10 text-white rounded-bl-md'
                  }`}
                >
                  <p className="text-sm">{message.content}</p>
                </div>
              </div>
            </motion.div>
          )
        })}
        <div ref={messagesEndRef} />
      </div>

      {/* Input */}
      <form onSubmit={handleSend} className="glass sticky bottom-0 p-4">
        <div className="flex items-center gap-3">
          <input
            ref={inputRef}
            type="text"
            value={newMessage}
            onChange={handleInputChange}
            placeholder="Type a message..."
            className="flex-1 px-4 py-3 bg-white/5 border border-white/10 rounded-full text-white placeholder:text-white/40 focus:outline-none focus:border-primary-500/50"
          />
          <button
            type="submit"
            disabled={!newMessage.trim() || sendMessageMutation.isPending}
            className="w-12 h-12 bg-gradient-to-br from-primary-500 to-primary-600 rounded-full flex items-center justify-center disabled:opacity-50 transition-opacity"
          >
            <Send className="w-5 h-5 text-white" />
          </button>
        </div>
      </form>
    </div>
  )
}
