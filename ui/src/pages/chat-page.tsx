import { useState, useEffect } from "react"
import { useChat } from "../hooks/useChat"
import { useNetwork } from "../context/NetworkContext"
import "../styles/chat.css"

export default function ChatPage() {
  const { data, loading, error, chatData } = useNetwork()
  const { messages, sendMessage } = useChat()

  const [text, setText] = useState<string>("")
  const [target, setTarget] = useState<string>("all")

  const peers = data?.peers ?? []

  useEffect(() => {
    if (!peers.find(p => p.ip === target)) {
      setTarget("all")
    }
  }, [peers])


  if (loading) {
    return <div className="chat-page">Loading network...</div>
  }

  if (error) {
    return <div className="chat-page error">{error}</div>
  }

  return (
    <div className="chat-page">
      <h2>LAN Chat</h2>

      {/* Target selector */}
      <select
        value={target}
        onChange={(e) => setTarget(e.target.value)}
      >
        <option value="all">All peers</option>

        {peers.map((peer) => (
          <option key={peer.ip} value={peer.ip}>
            {peer.name} ({peer.ip})
          </option>
        ))}
      </select>

      {/* Messages */}
      <div className="messages">
        {chatData.map((msg, idx) => (
          <div key={idx}>
            <strong>{msg.from}</strong>: {msg.body} {/* <-- use body */}
          </div>
        ))}
      </div>

      {/* Input */}
      <div className="input-row">
        <input
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder="Type a message"
        />

        <button
          disabled={!text.trim()}
          onClick={() => {
            sendMessage(text)
            setText("")
          }}
        >
          Send
        </button>
      </div>
    </div>
  )
}