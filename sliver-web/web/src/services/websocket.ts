import { store } from '../stores/store'
import { setSessions, addSession, updateSession, removeSession } from '../stores/sessionSlice'
import { setBeacons, addBeacon, updateBeacon, removeBeacon } from '../stores/beaconSlice'

class WebSocketService {
  private ws: WebSocket | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 3000

  connect() {
    const token = localStorage.getItem('token')
    if (!token) {
      console.error('No token found')
      return
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws?token=${token}`

    this.ws = new WebSocket(wsUrl)

    this.ws.onopen = () => {
      console.log('WebSocket connected')
      this.reconnectAttempts = 0
      this.subscribe(['sessions', 'beacons', 'tasks', 'alerts'])
    }

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        this.handleMessage(message)
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
      }
    }

    this.ws.onclose = () => {
      console.log('WebSocket disconnected')
      this.reconnect()
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
  }

  private reconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      console.log(`Reconnecting... attempt ${this.reconnectAttempts}`)
      setTimeout(() => this.connect(), this.reconnectDelay)
    }
  }

  private subscribe(channels: string[]) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({
        action: 'subscribe',
        channels,
      }))
    }
  }

  private handleMessage(message: any) {
    switch (message.type) {
      case 'event':
        this.handleEvent(message)
        break
      case 'pong':
        break
      default:
        console.log('Unknown message type:', message.type)
    }
  }

  private handleEvent(message: any) {
    const { channel, data } = message
    const event = data?.event

    switch (channel) {
      case 'sessions':
        if (event === 'session:online') {
          store.dispatch(addSession(data.session))
        } else if (event === 'session:offline') {
          store.dispatch(removeSession(data.session.id))
        }
        break
      case 'beacons':
        if (event === 'beacon:checkin') {
          store.dispatch(updateBeacon(data.beacon))
        }
        break
      case 'tasks':
        if (event === 'beacon:task_complete') {
          console.log('Task completed:', data)
        }
        break
      case 'alerts':
        console.log('Alert:', data)
        break
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }

  send(data: any) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data))
    }
  }
}

export const wsService = new WebSocketService()
