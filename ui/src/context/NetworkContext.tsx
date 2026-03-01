import { createContext, useContext, useEffect, useState, type ReactNode } from "react"
import { fetchNetworkInfo, type NetworkResponse } from "../api/network"
import type { ChatMessage } from "../hooks/useChat"

interface NetworkContextType {
  data: NetworkResponse | null
  loading: boolean
  error: string
  refresh: () => Promise<void>
  chatData: ChatMessage[]
  setChatData: React.Dispatch<React.SetStateAction<ChatMessage[]>>
}

const NetworkContext = createContext<NetworkContextType | undefined>(undefined)

export function NetworkProvider({ children }: { children: ReactNode }) {
  const [data, setData] = useState<NetworkResponse | null>(null)
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(true)
  const [chatData, setChatData] = useState<ChatMessage[]>([]);

  const fetchData = async () => {
    try {
      setLoading(true)
      const result = await fetchNetworkInfo()
      setData(result)
      setError("")
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  return (
    <NetworkContext.Provider
      value={{
        data,
        loading,
        error,
        refresh: fetchData,
        chatData,
        setChatData
      }}
    >
      {children}
    </NetworkContext.Provider>
  )
}

export function useNetwork() {
  const context = useContext(NetworkContext)
  if (!context) {
    throw new Error("useNetwork must be used inside NetworkProvider")
  }
  return context
}