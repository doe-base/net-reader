import { useEffect, useState } from "react"
import { fetchNetworkInfo, type NetworkResponse, type NetworkDevice, type PeerInfo } from "./api/network"
import NetworkCard from "./components/NetworkCard"
import ChatPage from "./pages/chat-page"
import { Routes, Route } from 'react-router-dom';
import PeerFileSystem from "./pages/peer-files-system";

export default function App() {
  const [data, setData] = useState<NetworkResponse | null>(null)
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchNetworkInfo()
      .then(setData)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  // Helper to get subnet prefix (first 3 octets)
  const getSubnetPrefix = (ip?: string) => {
    if (!ip) return "" // return empty string if IP is undefined
    const parts = ip.split(".")
    return parts.slice(0, 3).join(".")
  }

  return (
    <main className="container">
      <header>
        <h1>Net Reader</h1>
        <p>Inspect your local network in real time</p>
      </header>

      {loading && <p className="status">Loading network data…</p>}
      {error && <p className="error">{error}</p>}

      {data &&
        data.devices.map((device: NetworkDevice) => {
          const subnet = getSubnetPrefix(device.localIp)

         const devicePeers: PeerInfo[] = subnet
          ? data.peers.filter((peer) => peer.ip.startsWith(subnet))
          : []
          
          return (
            <>
              <Routes>
                <Route path="/" element={<NetworkCard
                  device={device}
                  peers={devicePeers}
                />} />
                <Route path="/chat-room" element={<ChatPage  peers={devicePeers}/>} />
                <Route path="/peer/:ip" element={<PeerFileSystem />} />
              </Routes>
            </>
          )
        })}
    </main>
  )
}
