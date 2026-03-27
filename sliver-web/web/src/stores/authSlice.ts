import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface User {
  id: number
  username: string
  role: string
}

export interface AuthState {
  isAuthenticated: boolean
  token: string | null
  user: User | null
  need2FA: boolean
  loading: boolean
  error: string | null
}

const initialState: AuthState = {
  isAuthenticated: false,
  token: localStorage.getItem('token'),
  user: null,
  need2FA: false,
  loading: false,
  error: null,
}

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    loginStart(state) {
      state.loading = true
      state.error = null
    },
    loginSuccess(state, action: PayloadAction<{ token: string; user: User; need2FA: boolean }>) {
      state.loading = false
      state.isAuthenticated = !action.payload.need2FA
      state.token = action.payload.token
      state.user = action.payload.user
      state.need2FA = action.payload.need2FA
      if (action.payload.token) {
        localStorage.setItem('token', action.payload.token)
      }
    },
    loginFailure(state, action: PayloadAction<string>) {
      state.loading = false
      state.error = action.payload
    },
    verify2FASuccess(state, action: PayloadAction<{ token: string; user: User }>) {
      state.loading = false
      state.isAuthenticated = true
      state.token = action.payload.token
      state.user = action.payload.user
      state.need2FA = false
      localStorage.setItem('token', action.payload.token)
    },
    logout(state) {
      state.isAuthenticated = false
      state.token = null
      state.user = null
      state.need2FA = false
      localStorage.removeItem('token')
    },
    clearError(state) {
      state.error = null
    },
  },
})

export const {
  loginStart,
  loginSuccess,
  loginFailure,
  verify2FASuccess,
  logout,
  clearError,
} = authSlice.actions

export default authSlice.reducer
