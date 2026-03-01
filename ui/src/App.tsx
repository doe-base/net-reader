import { Routes, Route } from 'react-router-dom';
import PeerFileSystem from "./pages/peer-files-system";
import NetworkInfoPage from "./pages/network-info";
import ChatPage from './pages/chat-page';

export default function App() {

  return (
    <main className="container">
        <>
          <Routes>
            <Route path="/" element={<NetworkInfoPage />} />
            <Route path='/chat-room' element={<ChatPage />}/>
            <Route path="/peer/:ip" element={<PeerFileSystem />} />
          </Routes>
        </>
    </main>
  )
}
