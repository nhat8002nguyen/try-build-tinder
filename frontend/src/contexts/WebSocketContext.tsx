import { createContext, useContext, useEffect, useCallback, useState, type ReactNode } from 'react'
import { wsService, type WebSocketMessage } from '../services/websocket'

interface WebSocketContextType {
  isConnected: boolean
  sendMessage: (message: WebSocketMessage) => void
  addMessageHandler: (handler: (message: WebSocketMessage) => void) => void
  removeMessageHandler: (handler: (message: WebSocketMessage) => void) => void
  sendTyping: (matchId: string) => void
}

const WebSocketContext = createContext<WebSocketContextType | null>(null)

export function WebSocketProvider({ 
  children, 
  enabled 
}: { 
  children: ReactNode
  enabled: boolean 
}) {
  const [isConnected, setIsConnected] = useState(false)

  useEffect(() => {
    if (!enabled) {
      wsService.disconnect()
      setIsConnected(false)
      return
    }

    const token = localStorage.getItem('access_token')
    if (token) {
      wsService.connect(token)

      const connectionHandler = (message: WebSocketMessage) => {
        if (message.type === 'pong') {
          setIsConnected(true)
        }
      }

      wsService.addMessageHandler(connectionHandler)

      const interval = setInterval(() => {
        setIsConnected(wsService.isConnected)
      }, 1000)

      return () => {
        wsService.removeMessageHandler(connectionHandler)
        clearInterval(interval)
        wsService.disconnect()
      }
    }
  }, [enabled])

  const sendMessage = useCallback((message: WebSocketMessage) => {
    wsService.send(message)
  }, [])

  const addMessageHandler = useCallback((handler: (message: WebSocketMessage) => void) => {
    wsService.addMessageHandler(handler)
  }, [])

  const removeMessageHandler = useCallback((handler: (message: WebSocketMessage) => void) => {
    wsService.removeMessageHandler(handler)
  }, [])

  const sendTyping = useCallback((matchId: string) => {
    wsService.sendTyping(matchId)
  }, [])

  return (
    <WebSocketContext.Provider
      value={{
        isConnected,
        sendMessage,
        addMessageHandler,
        removeMessageHandler,
        sendTyping,
      }}
    >
      {children}
    </WebSocketContext.Provider>
  )
}

export function useWebSocket() {
  const context = useContext(WebSocketContext)
  if (!context) {
    throw new Error('useWebSocket must be used within a WebSocketProvider')
  }
  return context
}
