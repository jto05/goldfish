import { useEffect, useState } from 'react'
import { Moon, Sun } from 'lucide-react'
import { Button } from '@/components/ui/button'
import Console from './Console'

// Base URL for all server control REST endpoints
const SERVER_URL = 'http://dev.homelab.internal:8080/api/server'

function App() {
  // running tracks whether the Minecraft server is currently running.
  // null = status not yet fetched (initial load), true = running, false = stopped.
  const [running, setRunning] = useState<boolean | null>(null)

  // dark tracks whether dark mode is active. Defaults to true since this is
  // a terminal-style dashboard that looks better dark.
  const [dark, setDark] = useState(true)

  // Toggle the "dark" class on <html> whenever dark changes.
  // shadcn's CSS variables switch automatically when this class is present.
  useEffect(() => {
    document.documentElement.classList.toggle('dark', dark)
  }, [dark])

  // fetchStatus calls GET /api/server/status and updates the running state.
  const fetchStatus = async () => {
    const res = await fetch(`${SERVER_URL}/status`)
    const data = await res.json()
    setRunning(data.status === 'running')
  }

  // Poll status once on mount and then every 5 seconds to keep the UI in sync.
  useEffect(() => {
    fetchStatus()
    const interval = setInterval(fetchStatus, 5000)
    return () => clearInterval(interval)
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
    <div className="min-h-screen bg-background text-foreground p-6">
      <div className="max-w-4xl mx-auto">

        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold tracking-tight">Goldfish</h1>

          <div className="flex items-center gap-3">
            {/* Status indicator */}
            <span className="flex items-center gap-2 text-sm text-muted-foreground">
              <span className={`w-2 h-2 rounded-full ${running ? 'bg-green-500' : 'bg-red-500'}`} />
              {running === null ? 'Loading...' : running ? 'Running' : 'Stopped'}
            </span>

            {/* Start/stop button */}
            {running === null ? null : running
              ? <Button variant="destructive" size="sm" onClick={stopServer}>Stop</Button>
              : <Button size="sm" onClick={startServer}>Start</Button>
            }

            {/* Dark mode toggle — shows Sun in dark mode, Moon in light mode.
                variant="ghost" makes it a borderless icon button.
                size="icon" makes it square to fit just the icon. */}
            <Button variant="ghost" size="icon" onClick={() => setDark(d => !d)}>
              {dark ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
            </Button>
          </div>
        </div>

        <Console />
      </div>
    </div>
  )
}

export default App
