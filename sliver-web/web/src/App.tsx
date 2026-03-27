import { Routes, Route, Navigate } from 'react-router-dom'
import { Layout } from './components/Layout/Layout'
import { Login } from './pages/Login/Login'
import { Dashboard } from './pages/Dashboard/Dashboard'
import { Sessions } from './pages/Sessions/Sessions'
import { Beacons } from './pages/Beacons/Beacons'
import { Listeners } from './pages/Listeners/Listeners'
import { Generate } from './pages/Generate/Generate'
import { Loot } from './pages/Loot/Loot'
import { Settings } from './pages/Settings/Settings'
import { useAppSelector } from './hooks/useStore'

function App() {
  const { isAuthenticated } = useAppSelector((state) => state.auth)

  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      {isAuthenticated ? (
        <Route path="/" element={<Layout />}>
          <Route index element={<Dashboard />} />
          <Route path="sessions" element={<Sessions />} />
          <Route path="beacons" element={<Beacons />} />
          <Route path="listeners" element={<Listeners />} />
          <Route path="generate" element={<Generate />} />
          <Route path="loot" element={<Loot />} />
          <Route path="settings" element={<Settings />} />
        </Route>
      ) : (
        <Route path="*" element={<Navigate to="/login" replace />} />
      )}
    </Routes>
  )
}

export default App
