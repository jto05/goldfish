import { useEffect, useState } from 'react'
import Console from './Console'

const SERVER_URL = 'http://dev.homelab.internal:8080/api/server'

function App() {

  // running tracks whether Minecraft server is currently running. null means we haven't received
  // a status yet (initial load)
  const [running, setRunning] = useState<boolean | null>(null)

  // fetchStatus polls backkend for the current server status and updates states
  const fetchStatus = async () => {
    const res = await fetch(`${SERVER_URL}/status`)
    const data = await res.json()
    setRunning(data.status === 'running')
  }

  // Poll status on startup and every 5 seconds to keep button in sync
  useEffect( () => {
    fetchStatus()
    const interval = setInterval(fetchStatus, 5000)
    return () => clearInterval( interval )

  }, [])
  

  const startServer = async () => {
    await fetch(`${SERVER_URL}/start`, { method: 'POST' })
    fetchStatus()
  }

  const stopServer = async () => {
    await fetch(`${SERVER_URL}/stop`, { method: 'POST' })
    fetchStatus()
  }

  return (
    <div>
      <h1>GOLDFISH</h1>
      <div style={{ position: 'fixed', top: '16px', right: '16px', display: 'flex', gap: '8px' }}>
        {
          running === null ? null : running ? <button onClick={stopServer}>Stop</button>
          : <button onClick={startServer}>Start</button>
        }
      </div>
      <Console />
    </div>
  )
}

export default App
