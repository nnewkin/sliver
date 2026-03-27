import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface Session {
  id: string
  name: string
  hostname: string
  username: string
  uuid: string
  os: string
  arch: string
  transport: string
  remote_address: string
  pid: number
  last_checkin: string
  next_checkin: string
  is_active: boolean
  version: string
}

export interface SessionState {
  sessions: Session[]
  selectedSession: Session | null
  loading: boolean
  error: string | null
}

const initialState: SessionState = {
  sessions: [],
  selectedSession: null,
  loading: false,
  error: null,
}

const sessionSlice = createSlice({
  name: 'sessions',
  initialState,
  reducers: {
    setLoading(state, action: PayloadAction<boolean>) {
      state.loading = action.payload
    },
    setSessions(state, action: PayloadAction<Session[]>) {
      state.sessions = action.payload
      state.loading = false
    },
    addSession(state, action: PayloadAction<Session>) {
      state.sessions.push(action.payload)
    },
    updateSession(state, action: PayloadAction<Session>) {
      const index = state.sessions.findIndex((s) => s.id === action.payload.id)
      if (index !== -1) {
        state.sessions[index] = action.payload
      }
    },
    removeSession(state, action: PayloadAction<string>) {
      state.sessions = state.sessions.filter((s) => s.id !== action.payload)
      if (state.selectedSession?.id === action.payload) {
        state.selectedSession = null
      }
    },
    selectSession(state, action: PayloadAction<Session | null>) {
      state.selectedSession = action.payload
    },
    setError(state, action: PayloadAction<string | null>) {
      state.error = action.payload
      state.loading = false
    },
  },
})

export const {
  setLoading,
  setSessions,
  addSession,
  updateSession,
  removeSession,
  selectSession,
  setError,
} = sessionSlice.actions

export default sessionSlice.reducer
