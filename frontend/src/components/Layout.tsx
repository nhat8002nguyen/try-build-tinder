import { Outlet, NavLink, useLocation } from 'react-router-dom'
import { Flame, Heart, User } from 'lucide-react'
import { motion } from 'framer-motion'

const navItems = [
  { to: '/discover', icon: Flame, label: 'Discover' },
  { to: '/matches', icon: Heart, label: 'Matches' },
  { to: '/profile', icon: User, label: 'Profile' },
]

export default function Layout() {
  const location = useLocation()
  const isChatPage = location.pathname.includes('/matches/') && location.pathname !== '/matches'

  return (
    <div className="min-h-screen flex flex-col">
      <main className="flex-1 pb-20">
        <Outlet />
      </main>

      {!isChatPage && (
        <nav className="fixed bottom-0 left-0 right-0 z-50">
          <div className="mx-auto max-w-lg px-4 pb-4">
            <div className="glass rounded-2xl px-6 py-3">
              <div className="flex items-center justify-around">
                {navItems.map(({ to, icon: Icon, label }) => {
                  const isActive = location.pathname === to || 
                    (to === '/matches' && location.pathname.startsWith('/matches'))
                  
                  return (
                    <NavLink
                      key={to}
                      to={to}
                      className="relative flex flex-col items-center gap-1 px-4 py-2"
                    >
                      {isActive && (
                        <motion.div
                          layoutId="activeTab"
                          className="absolute inset-0 bg-primary-500/20 rounded-xl"
                          transition={{ type: 'spring', bounce: 0.2, duration: 0.6 }}
                        />
                      )}
                      <Icon
                        className={`w-6 h-6 relative z-10 transition-colors ${
                          isActive ? 'text-primary-500' : 'text-white/60'
                        }`}
                      />
                      <span
                        className={`text-xs font-medium relative z-10 transition-colors ${
                          isActive ? 'text-primary-500' : 'text-white/60'
                        }`}
                      >
                        {label}
                      </span>
                    </NavLink>
                  )
                })}
              </div>
            </div>
          </div>
        </nav>
      )}
    </div>
  )
}
