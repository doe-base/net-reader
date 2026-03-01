import { useEffect, useRef, useState, useCallback, type ReactNode } from "react"
import { useNetwork } from "../context/NetworkContext"

/* ---------- Types ---------- */

export interface ChatMessage {
  content: ReactNode
  id: string
  from: string
  body: string
  timestamp: number
}

/* ---------- Helpers ---------- */

const generateId = () => crypto.randomUUID()

function getOrCreateDeviceId(): string {
  const existing = localStorage.getItem("device_id")
  if (existing) return existing

  const newId = generateId()
  localStorage.setItem("device_id", newId)
  return newId
}

/* ---------- Hook ---------- */

export function useChat() {
  const { setChatData } = useNetwork()
  const deviceId = useRef<string>(getOrCreateDeviceId())
  const messageCache = useRef<Set<string>>(new Set())

  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [isConnected, setIsConnected] = useState(true)

  /* ---------- Fetch Messages ---------- */

  const fetchMessages = useCallback(async () => {
    try {
      const res = await fetch("http://localhost:8080/api/messages")
      const data = await res.json()

      if (!Array.isArray(data)) {
        console.error("Expected array, got:", data)
        return
      }

      setChatData(data)

      setMessages((prev) => {
        const updated = [...prev]

        for (const msg of data) {
          if (!messageCache.current.has(msg.id)) {
            messageCache.current.add(msg.id)
            updated.push(msg)
          }
        }

        return updated.sort((a, b) => a.timestamp - b.timestamp)
      })
    } catch (err) {
      console.error("Failed to fetch messages:", err)
      setIsConnected(false)
    }
  }, [])

  /* ---------- Polling ---------- */

  useEffect(() => {
    fetchMessages()

    const interval = setInterval(fetchMessages, 1500)
    return () => clearInterval(interval)
  }, [fetchMessages])

  /* ---------- Send Message ---------- */

  const sendMessage = useCallback(async (content: string) => {
    const message: ChatMessage = {
      id: generateId(),
      from: deviceId.current,
      body: content,
      timestamp: Date.now(),
      content: undefined
    }

    // Optimistic UI
    setMessages((prev) => [...prev, message])
    messageCache.current.add(message.id)

    try {
      await fetch("http://localhost:8080/api/send", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(message),
      })
    } catch (err) {
      console.error("Send failed:", err)
    }
  }, [])

  return {
    messages,
    sendMessage,
    isConnected,
    deviceId: deviceId.current,
  }
}