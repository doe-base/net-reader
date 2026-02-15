import { useEffect, useRef, useState, useCallback } from "react"

export interface ChatMessage {
  from: string
  to: string
  type: string
  content: string
  timestamp: number
}

export function useChat() {
  const ws = useRef<WebSocket | null>(null)
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [isConnected, setIsConnected] = useState(false)

  useEffect(() => {
    // Prevent multiple connections in Dev mode
    if (ws.current?.readyState === WebSocket.OPEN) return;

    const socket = new WebSocket("ws://localhost:8080/ws/chat")

    socket.onopen = () => {
        setIsConnected(true)
        // IMPORTANT: Tell the backend who we are as soon as we connect
        // You'd get 'myHostname' from your backend discovery service
        socket.send(JSON.stringify({
        type: "identify",
        from: "Daniel-PC", 
        content: "192.168.0.161" 
        }))
    }
    socket.onclose = () => setIsConnected(false)
    
    socket.onmessage = (e) => {
      try {
        const msg: ChatMessage = JSON.parse(e.data)
        setMessages((prev) => [...prev, msg])
      } catch (err) {
        console.error("Failed to parse message:", err)
      }
    }

    ws.current = socket

    return () => {
      socket.close()
    }
  }, [])

  // Use useCallback so this function reference doesn't change on every render
  const sendMessage = useCallback((to: string, content: string) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(
        JSON.stringify({
          to,
          type: "chat",
          content,
        })
      )
    } else {
      console.warn("WebSocket is not connected.")
    }
  }, [])

  return { messages, sendMessage, isConnected }
}