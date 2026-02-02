import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Flame, Heart, Sparkles, MessageCircle } from 'lucide-react'

export default function Landing() {
  return (
    <div className="min-h-screen relative overflow-hidden">
      {/* Animated background */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-primary-500/20 rounded-full blur-3xl animate-pulse-slow" />
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-secondary-500/20 rounded-full blur-3xl animate-pulse-slow" style={{ animationDelay: '1s' }} />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-primary-600/10 rounded-full blur-3xl" />
      </div>

      <div className="relative z-10 min-h-screen flex flex-col">
        {/* Header */}
        <header className="p-6">
          <div className="flex items-center gap-2">
            <div className="w-10 h-10 bg-gradient-to-br from-primary-500 to-primary-600 rounded-xl flex items-center justify-center">
              <Flame className="w-6 h-6 text-white" />
            </div>
            <span className="font-display text-2xl font-bold text-white">Spark</span>
          </div>
        </header>

        {/* Main content */}
        <main className="flex-1 flex flex-col items-center justify-center px-6 pb-20">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            className="text-center max-w-md"
          >
            {/* Floating hearts animation */}
            <div className="relative mb-8">
              <motion.div
                animate={{ y: [-10, 10, -10] }}
                transition={{ duration: 4, repeat: Infinity, ease: 'easeInOut' }}
                className="w-24 h-24 mx-auto bg-gradient-to-br from-primary-500 to-primary-600 rounded-3xl flex items-center justify-center shadow-2xl shadow-primary-500/30"
              >
                <Heart className="w-12 h-12 text-white" fill="white" />
              </motion.div>
              
              <motion.div
                animate={{ y: [0, -15, 0], x: [0, 10, 0] }}
                transition={{ duration: 3, repeat: Infinity, ease: 'easeInOut', delay: 0.5 }}
                className="absolute -top-4 -right-4 w-12 h-12 bg-secondary-500/80 rounded-xl flex items-center justify-center"
              >
                <Sparkles className="w-6 h-6 text-white" />
              </motion.div>

              <motion.div
                animate={{ y: [0, 15, 0], x: [0, -10, 0] }}
                transition={{ duration: 3.5, repeat: Infinity, ease: 'easeInOut', delay: 1 }}
                className="absolute -bottom-2 -left-4 w-10 h-10 bg-pink-400/80 rounded-lg flex items-center justify-center"
              >
                <MessageCircle className="w-5 h-5 text-white" />
              </motion.div>
            </div>

            <h1 className="font-display text-4xl md:text-5xl font-bold text-white mb-4">
              Find Your <span className="text-gradient">Perfect Match</span>
            </h1>
            
            <p className="text-white/70 text-lg mb-8">
              Discover meaningful connections with people who share your interests and values.
            </p>

            <div className="flex flex-col gap-4">
              <Link to="/register" className="btn-primary text-center">
                Create Account
              </Link>
              <Link to="/login" className="btn-secondary text-center">
                Sign In
              </Link>
            </div>
          </motion.div>
        </main>

        {/* Features */}
        <motion.div
          initial={{ opacity: 0, y: 40 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.3 }}
          className="px-6 pb-10"
        >
          <div className="max-w-md mx-auto grid grid-cols-3 gap-4">
            {[
              { icon: Heart, label: 'Match' },
              { icon: MessageCircle, label: 'Chat' },
              { icon: Sparkles, label: 'Connect' },
            ].map(({ icon: Icon, label }) => (
              <div
                key={label}
                className="card p-4 text-center"
              >
                <div className="w-10 h-10 mx-auto mb-2 bg-primary-500/20 rounded-xl flex items-center justify-center">
                  <Icon className="w-5 h-5 text-primary-400" />
                </div>
                <span className="text-white/80 text-sm font-medium">{label}</span>
              </div>
            ))}
          </div>
        </motion.div>
      </div>
    </div>
  )
}
