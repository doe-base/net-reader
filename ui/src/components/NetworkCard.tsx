import type { NetworkDevice, PeerInfo } from "../api/network"
import { Link } from "react-router-dom"

interface Props {
  device: NetworkDevice
  peers?: PeerInfo[] // optional list of peers
}

export default function NetworkCard({ device, peers }: Props) {
  return (
    <div className="card">
      <h2>{device.vendor || "Unknown Device"}</h2>

      <div className="row">
        <span>IP Address</span>
        <span>{device.ip}</span>
      </div>

      <div className="row">
        <span>MAC Address</span>
        <span>{device.mac}</span>
      </div>

      <div className="row">
        <span>Device Name</span>
        <span>{device.deviceName || "Unknown"}</span>
      </div>

      <div className="row">
        <span>Interface</span>
        <span>{device.interface}</span>
      </div>

      <div className="row">
        <span>Local IP</span>
        <span>{device.localIp}</span>
      </div>

      {/* --- Peer list --- */}
      {peers && peers.length > 0 && (
        <div className="peers">
          <h3>Active Peers</h3>
          <ul>
            {peers.map((peer) => (
              <li key={peer.ip}>
                <Link
                  to={`/peer/${peer.ip}`}
                  className="peer-link"
                >
                  <span className="peer-name">{peer.name}</span>
                  <span className="peer-ip">{peer.ip}</span>
                </Link>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  )
}
