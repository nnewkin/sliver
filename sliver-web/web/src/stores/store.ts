import { configureStore } from '@reduxjs/toolkit'
import authReducer from './authSlice'
import sessionReducer from './sessionSlice'
import beaconReducer from './beaconSlice'
import listenerReducer from './listenerSlice'

export const store = configureStore({
  reducer: {
    auth: authReducer,
    sessions: sessionReducer,
    beacons: beaconReducer,
    listeners: listenerReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: false,
    }),
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
