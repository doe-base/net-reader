export interface NetworkDevice {
  index: number
  interface: string
  mac: string
  deviceName: string
  ip: string
  vendor: string
  localIp: string
}

export interface NetworkResponse {
  count: number
  devices: NetworkDevice[]
  defaultGateway: string
  peers: PeerInfo[]
}
export interface PeerInfo {
  name: string
  ip: string
}


export async function fetchNetworkInfo(): Promise<NetworkResponse> {
  const res = await fetch("http://localhost:8080/get-my-network-info")

  if (!res.ok) {
    throw new Error("Failed to fetch network info")
  }

  const data = await res.json()  // read body ONCE
  return data                     // return to caller
}