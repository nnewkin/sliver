import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface Beacon {
  id: string
  name: string
  hostname: string
  username: string
  uuid: string
  os: string
  arch: string
  remote_address: string
  pid: number
  last_checkin: string
  next_checkin: string
  interval: number
  jitter: number
  is_active: boolean
  version: string
}

export interface BeaconTask {
  id: string
  beacon_id: string
  type: string
  command: string[]
  created_at: string
  completed: boolean
  success: boolean
}

export interface BeaconState {
  beacons: Beacon[]
  selectedBeacon: Beacon | null
  tasks: BeaconTask[]
  loading: boolean
  error: string | null
}

const initialState: BeaconState = {
  beacons: [],
  selectedBeacon: null,
  tasks: [],
  loading: false,
  error: null,
}

const beaconSlice = createSlice({
  name: 'beacons',
  initialState,
  reducers: {
    setLoading(state, action: PayloadAction<boolean>) {
      state.loading = action.payload
    },
    setBeacons(state, action: PayloadAction<Beacon[]>) {
      state.beacons = action.payload
      state.loading = false
    },
    addBeacon(state, action: PayloadAction<Beacon>) {
      state.beacons.push(action.payload)
    },
    updateBeacon(state, action: PayloadAction<Beacon>) {
      const index = state.beacons.findIndex((b) => b.id === action.payload.id)
      if (index !== -1) {
        state.beacons[index] = action.payload
      }
    },
    removeBeacon(state, action: PayloadAction<string>) {
      state.beacons = state.beacons.filter((b) => b.id !== action.payload)
      if (state.selectedBeacon?.id === action.payload) {
        state.selectedBeacon = null
      }
    },
    selectBeacon(state, action: PayloadAction<Beacon | null>) {
      state.selectedBeacon = action.payload
    },
    setTasks(state, action: PayloadAction<BeaconTask[]>) {
      state.tasks = action.payload
    },
    addTask(state, action: PayloadAction<BeaconTask>) {
      state.tasks.push(action.payload)
    },
    updateTask(state, action: PayloadAction<BeaconTask>) {
      const index = state.tasks.findIndex((t) => t.id === action.payload.id)
      if (index !== -1) {
        state.tasks[index] = action.payload
      }
    },
    setError(state, action: PayloadAction<string | null>) {
      state.error = action.payload
      state.loading = false
    },
  },
})

export const {
  setLoading,
  setBeacons,
  addBeacon,
  updateBeacon,
  removeBeacon,
  selectBeacon,
  setTasks,
  addTask,
  updateTask,
  setError,
} = beaconSlice.actions

export default beaconSlice.reducer
