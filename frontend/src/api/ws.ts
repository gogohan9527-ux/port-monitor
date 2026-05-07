import { usePortsStore } from '@/stores/ports'

let socket: WebSocket | null = null
let retryDelay = 1000
let manualClose = false

export function connectWS(token: string) {
  manualClose = false
  open(token)
}

export function disconnectWS() {
  manualClose = true
  if (socket) {
    socket.close()
    socket = null
  }
}

function open(token: string) {
  const proto = location.protocol === 'https:' ? 'wss' : 'ws'
  const url = `${proto}://${location.host}/api/ws?token=${encodeURIComponent(token)}`
  socket = new WebSocket(url)

  socket.onopen = () => {
    retryDelay = 1000
  }

  socket.onmessage = (ev) => {
    try {
      const frame = JSON.parse(ev.data)
      if (frame.type === 'ports' && Array.isArray(frame.ports)) {
        usePortsStore().setList(frame.ports, frame.timestamp)
      }
    } catch {
      // ignore malformed
    }
  }

  socket.onclose = () => {
    socket = null
    if (manualClose) return
    setTimeout(() => open(token), retryDelay)
    retryDelay = Math.min(retryDelay * 2, 10000)
  }

  socket.onerror = () => {
    socket?.close()
  }
}
