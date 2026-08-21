import { useEffect, useRef, useState } from 'react'

// WebSocket endpoint for live console streaming
const WS_URL = 'ws://dev.homelab.internal:8080/ws/console'

// REST endpoint for sending commands to the server
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
  // async/await pauses execution until fetch resolves. After sending, input is
  // cleared — the server's response will appear in the console output above.
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
  // We only act on Enter — all other keys are handled natively by the input.
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') sendCommand()
  }

  return (
    // rounded-lg gives rounded corners; border and bg-zinc-950 style it as a
    // dark terminal panel regardless of the app's light/dark theme.
    <div className="rounded-lg border border-border overflow-hidden">

      {/* Scrollable output area — flex-col-reverse would auto-scroll but
          we use a ref-based approach instead for more control. font-mono
          ensures monospace rendering for console output. */}
      <div className="h-[500px] overflow-y-scroll bg-zinc-950 text-zinc-100 p-4 font-mono text-sm">
        {/* lines.map() transforms each string into a <div>.
            key={i} is required by React to efficiently track which elements
            changed between renders — index is fine here since lines only grow. */}
        {lines.map((line, i) => (
          <div key={i} className="leading-5 whitespace-pre-wrap break-all">{line}</div>
        ))}

        {/* Invisible anchor div at the bottom — bottomRef points here so
            scrollIntoView always scrolls to the very end of the output. */}
        <div ref={bottomRef} />
      </div>

      {/* Command input — sits flush at the bottom of the console panel.
          border-t separates it from the output area.
          value + onChange makes this a controlled input (value lives in state).
          onKeyDown sends the command when Enter is pressed. */}
      <div className="flex border-t border-border bg-zinc-900">
        <span className="pl-4 py-3 text-zinc-500 font-mono text-sm select-none">{'>'}</span>
        <input
          className="flex-1 bg-transparent px-2 py-3 font-mono text-sm text-zinc-100 placeholder:text-zinc-600 focus:outline-none"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Enter command..."
        />
      </div>
    </div>
  )
}

export default Console
