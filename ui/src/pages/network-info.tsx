import { type NetworkDevice, type PeerInfo } from "../api/network"
import NetworkCard from "../components/NetworkCard"
import { useNetwork } from "../context/NetworkContext"

export default function NetworkInfoPage() {
  const { data, loading, error } = useNetwork()

  const getSubnetPrefix = (ip?: string) => {
    if (!ip) return ""
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
            ? data.peers.filter((peer) =>
                peer.ip.startsWith(subnet)
              )
            : []

          return (
            <NetworkCard
              key={device.id}
              device={device}
              peers={devicePeers}
            />
          )
        })}
    </main>
  )
}