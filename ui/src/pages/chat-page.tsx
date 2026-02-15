import { useState } from "react"
import { useChat } from "../hooks/useChat"
import type { PeerInfo } from "../api/network"
import "../styles/chat.css"

interface Props {
  peers: PeerInfo[]
}

export default function ChatPage({ peers }: Props) {
  const { messages, sendMessage } = useChat()
  const [text, setText] = useState<string>("")
  const [target, setTarget] = useState<string>("all")

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
        {messages.map((msg, idx) => (
          <div key={idx}>
            <strong>{msg.from}</strong>: {msg.content}
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
            sendMessage(target, text)
            setText("")
          }}
        >
          Send
        </button>
      </div>
    </div>
  )
}
