import { useState } from 'react';

// Define the shape of our file data
interface PeerFile {
  name: string;
  size?: number; // Optional, if you decide to add more info later
}

const FileExplorer = () => {
  const [peerFiles, setPeerFiles] = useState<string[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const fetchPeerFiles = async (peerIP: string): Promise<void> => {
    setLoading(true);
    setError(null);

    try {
      // 1. Fetch from the specific Peer's Go server
      const response = await fetch(`http://${peerIP}:8080/api/list`, {
        method: 'GET',
        // Requesting JSON specifically
        headers: {
          'Accept': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`Failed to connect to peer: ${response.statusText}`);
      }

      // 2. Parse the string array from your Go backend
      const files: string[] = await response.json();
      setPeerFiles(files);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unknown error';
      setError(`Could not reach peer at ${peerIP}: ${message}`);
      console.error("Peer fetch error:", err);
    } finally {
      setLoading(false);
    }
  };

  return { fetchPeerFiles, peerFiles, loading, error };
};