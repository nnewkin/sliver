import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  expires_at: string
  user: {
    id: number
    username: string
    role: string
  }
  need_2fa: boolean
}

export interface Verify2FARequest {
  token: string
  code: string
}

export const authAPI = {
  login: (data: LoginRequest) => api.post<{ code: number; data: LoginResponse }>('/auth/login', data),
  verify2FA: (data: Verify2FARequest) => api.post('/auth/2fa/verify', data),
  logout: () => api.post('/auth/logout'),
  getProfile: () => api.get('/auth/profile'),
}

export const sessionAPI = {
  list: () => api.get('/sessions'),
  get: (id: string) => api.get(`/sessions/${id}`),
  kill: (id: string) => api.delete(`/sessions/${id}`),
  rename: (id: string, newName: string) => api.post(`/sessions/${id}/rename`, { new_name: newName }),
  getProcesses: (id: string) => api.get(`/sessions/${id}/processes`),
  execute: (id: string, command: string) => api.post(`/sessions/${id}/exec`, { command }),
  listDirectory: (id: string, path: string) => api.get(`/sessions/${id}/ls?path=${path}`),
  download: (id: string, path: string) => api.get(`/sessions/${id}/download?path=${path}`, { responseType: 'blob' }),
  upload: (id: string, path: string, data: string) => api.post(`/sessions/${id}/upload`, { path, data }),
  screenshot: (id: string) => api.get(`/sessions/${id}/screenshot`, { responseType: 'blob' }),
  getPrivs: (id: string) => api.get(`/sessions/${id}/privs`),
  makeToken: (id: string, username: string, password: string) =>
    api.post(`/sessions/${id}/token`, { username, password }),
}

export const beaconAPI = {
  list: () => api.get('/beacons'),
  get: (id: string) => api.get(`/beacons/${id}`),
  getTasks: (id: string) => api.get(`/beacons/${id}/tasks`),
  executeTask: (id: string, command: string) => api.post(`/beacons/${id}/task`, { command }),
  cancelTask: (id: string, taskId: string) => api.delete(`/beacons/${id}/tasks/${taskId}`),
  getTaskOutput: (id: string, taskId: string) => api.get(`/beacons/${id}/tasks/${taskId}/output`),
}

export const listenerAPI = {
  list: () => api.get('/listeners'),
  startMTLS: (name: string, bindAddr: string) => api.post('/listeners/mtls', { name, bind_address: bindAddr }),
  startDNS: (name: string, bindAddr: string, domains: string[]) =>
    api.post('/listeners/dns', { name, bind_address: bindAddr, domains }),
  startHTTP: (name: string, bindAddr: string, domains: string[]) =>
    api.post('/listeners/http', { name, bind_address: bindAddr, domains }),
  startHTTPS: (name: string, bindAddr: string, domains: string[]) =>
    api.post('/listeners/https', { name, bind_address: bindAddr, domains }),
  startWG: (name: string, bindAddr: string, port: number) =>
    api.post('/listeners/wg', { name, bind_address: bindAddr, port }),
  stop: (id: string) => api.delete(`/listeners/${id}`),
}

export const implantAPI = {
  listProfiles: () => api.get('/profiles'),
  createProfile: (name: string, config: any) => api.post('/profiles', { name, config }),
  updateProfile: (name: string, config: any) => api.put(`/profiles/${name}`, { config }),
  deleteProfile: (name: string) => api.delete(`/profiles/${name}`),
  generate: (config: any) => api.post('/generate', config),
  download: (name: string) => api.get(`/generate/download/${name}`, { responseType: 'blob' }),
}

export const userAPI = {
  list: () => api.get('/users'),
  create: (data: any) => api.post('/users', data),
  get: (id: number) => api.get(`/users/${id}`),
  update: (id: number, data: any) => api.put(`/users/${id}`, data),
  delete: (id: number) => api.delete(`/users/${id}`),
  changePassword: (id: number, oldPassword: string, newPassword: string) =>
    api.post(`/users/${id}/password`, { old_password: oldPassword, new_password: newPassword }),
}

export const botAPI = {
  getConfig: () => api.get('/bot/config'),
  updateConfig: (data: any) => api.put('/bot/config', data),
  start: () => api.post('/bot/start'),
  stop: () => api.post('/bot/stop'),
  listWhitelist: () => api.get('/bot/whitelist'),
  addWhitelist: (data: any) => api.post('/bot/whitelist', data),
  removeWhitelist: (id: number) => api.delete(`/bot/whitelist/${id}`),
  listAlertRules: () => api.get('/bot/alerts'),
  createAlertRule: (data: any) => api.post('/bot/alerts', data),
  updateAlertRule: (id: number, data: any) => api.put(`/bot/alerts/${id}`, data),
  deleteAlertRule: (id: number) => api.delete(`/bot/alerts/${id}`),
}

export const auditAPI = {
  list: (params?: any) => api.get('/audit', { params }),
  get: (id: number) => api.get(`/audit/${id}`),
}

export default api
