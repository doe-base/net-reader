import { useParams } from "react-router-dom"

export default function PeerFileSystem() {
  const { ip } = useParams<{ ip: string }>()

  return (
    <main className="container">
      <h2>Remote File System</h2>
      <p>Browsing files from peer: <strong>{ip}</strong></p>

      {/* File tree goes here */}
    </main>
  )
}