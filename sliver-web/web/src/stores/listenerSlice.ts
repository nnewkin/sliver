import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface Listener {
  id: string
  name: string
  type: string
  status: string
  bind_address: string
  port: number
  domain: string
  connected: number
}

export interface ListenerState {
  listeners: Listener[]
  selectedListener: Listener | null
  loading: boolean
  error: string | null
}

const initialState: ListenerState = {
  listeners: [],
  selectedListener: null,
  loading: false,
  error: null,
}

const listenerSlice = createSlice({
  name: 'listeners',
  initialState,
  reducers: {
    setLoading(state, action: PayloadAction<boolean>) {
      state.loading = action.payload
    },
    setListeners(state, action: PayloadAction<Listener[]>) {
      state.listeners = action.payload
      state.loading = false
    },
    addListener(state, action: PayloadAction<Listener>) {
      state.listeners.push(action.payload)
    },
    updateListener(state, action: PayloadAction<Listener>) {
      const index = state.listeners.findIndex((l) => l.id === action.payload.id)
      if (index !== -1) {
        state.listeners[index] = action.payload
      }
    },
    removeListener(state, action: PayloadAction<string>) {
      state.listeners = state.listeners.filter((l) => l.id !== action.payload)
    },
    selectListener(state, action: PayloadAction<Listener | null>) {
      state.selectedListener = action.payload
    },
    setError(state, action: PayloadAction<string | null>) {
      state.error = action.payload
      state.loading = false
    },
  },
})

export const {
  setLoading,
  setListeners,
  addListener,
  updateListener,
  removeListener,
  selectListener,
  setError,
} = listenerSlice.actions

export default listenerSlice.reducer
