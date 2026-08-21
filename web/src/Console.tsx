import { useEffect, useRef, useState } from 'react'

// WebSocket and server command endpoint served by the Go backend
const WS_URL = 'ws://dev.homelab.internal:8080/ws/console'
const COMMAND_URL = 'http://dev.homelab.internal:8080/api/server/command'

function Console() {
  // lines holds the list of console output strings received from the server.
  // useState returns the current value and a setter — calling setLines causes
  // the component to re-render with the new value. <string[]> is the TypeScript
  // type annotation; [] is the initial value (empty array).
  const [lines, setLines] = useState<string[]>([])

  // input holds the current value of the command input field.
  // Controlled inputs in React keep their value in state — every keystroke
  // calls setInput, which re-renders the input with the new value.
  const [input, setInput] = useState('')

  // bottomRef is a pointer to the invisible div at the bottom of the list.
  // Unlike state, changing a ref does not trigger a re-render — it's just a
  // direct reference to the DOM node used for scrolling. Starts as null until
  // the component mounts and React attaches it to the actual element.
  const bottomRef = useRef<HTMLDivElement>(null)

  // Open the WebSocket connection once when the component mounts.
  // The empty [] dependency array means this effect runs only on mount, not
  // on every render. The returned cleanup function closes the socket when the
  // component unmounts (e.g. the user navigates away).
  useEffect(() => {
    const ws = new WebSocket(WS_URL)

    // ws.onmessage fires every time the server sends a line.
    // e.data is the raw string payload. We append it to lines by spreading
    // the previous array into a new one — React requires state to be replaced,
    // not mutated in place, so [...prev, e.data] creates a fresh array each time.
    ws.onmessage = (e) => {
      setLines((prev) => [...prev, e.data])
    }

    // Cleanup: close the WebSocket when the component unmounts
    return () => ws.close()
  }, [])

  // Scroll to the bottom whenever lines changes (i.e. a new line arrives).
  // bottomRef.current is the actual DOM element — ?. (optional chaining) means
  // "only call scrollIntoView if bottomRef.current is not null."
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [lines])

  // sendCommand POSTs the current input value to the backend as a server command.
  // async/await pauses execution until fetch resolves — similar to Go's blocking
  // calls but single-threaded. After sending, input is cleared regardless of
  // whether the command succeeded (the server's response will appear in the console).
  const sendCommand = async () => {
    if (!input.trim()) return
    await fetch(COMMAND_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ command: input }),
    })
    setInput('')
  }

  // handleKeyDown fires on every keypress in the input field.
  // We only act on Enter — all other keys are ignored and handled
  // natively by the input element.
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') sendCommand()
  }

  return (
    <div style={{ fontFamily: 'monospace' }}>
      <div style={{ background: '#111', color: '#eee', height: '400px', overflowY: 'scroll', padding: '8px' }}>
        {/* lines.map() transforms each string into a <div>.
            key={i} is required by React to efficiently track which elements
            changed between renders — the index is fine here since lines only
            ever grow (we never reorder or delete them). */}
        {lines.map((line, i) => (
          <div key={i}>{line}</div>
        ))}

        {/* Invisible anchor div at the end of the list. bottomRef points here
            so scrollIntoView always scrolls to the very bottom. */}
        <div ref={bottomRef} />
      </div>

      {/* Controlled input — value is always in sync with the input state.
          onChange updates state on every keystroke; onKeyDown sends on Enter. */}
      <input
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="Enter command..."
        style={{ width: '100%', fontFamily: 'monospace', background: '#222', color: '#eee', border: 'none', padding: '8px', boxSizing: 'border-box' }}
      />
    </div>
  )
}

export default Console
