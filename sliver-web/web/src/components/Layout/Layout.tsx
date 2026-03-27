import React from 'react'
import { Outlet } from 'react-router-dom'
import { Layout as AntLayout, Menu, Avatar, Dropdown, Badge } from 'antd'
import {
  DashboardOutlined,
  DesktopOutlined,
  BellOutlined,
  SettingOutlined,
  LogoutOutlined,
} from '@ant-design/icons'
import { useNavigate, useLocation } from 'react-router-dom'
import { useAppSelector, useAppDispatch } from '../../hooks/useStore'
import { logout } from '../../stores/authSlice'

const { Header, Sider, Content } = AntLayout

export const Layout: React.FC = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const dispatch = useAppDispatch()
  const { user } = useAppSelector((state) => state.auth)

  const menuItems = [
    { key: '/', icon: <DashboardOutlined />, label: 'Dashboard' },
    { key: '/sessions', icon: <DesktopOutlined />, label: 'Sessions' },
    { key: '/beacons', icon: <DesktopOutlined />, label: 'Beacons' },
    { key: '/listeners', icon: <BellOutlined />, label: 'Listeners' },
    { key: '/settings', icon: <SettingOutlined />, label: 'Settings' },
  ]

  const handleLogout = () => {
    dispatch(logout())
    navigate('/login')
  }

  const userMenuItems = [
    {
      key: 'profile',
      label: user?.username || 'User',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: 'Logout',
      onClick: handleLogout,
    },
  ]

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      <Sider theme="dark" breakpoint="lg" collapsedWidth="0">
        <div style={{ height: 64, display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#fff', fontSize: 18, fontWeight: 'bold' }}>
          Sliver Web
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <AntLayout>
        <Header style={{ background: '#001529', padding: '0 24px', display: 'flex', alignItems: 'center', justifyContent: 'flex-end' }}>
          <Badge count={0} size="small">
            <BellOutlined style={{ fontSize: 18, color: '#fff', marginRight: 24, cursor: 'pointer' }} />
          </Badge>
          <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
            <Avatar style={{ cursor: 'pointer' }}>{user?.username?.[0]?.toUpperCase() || 'U'}</Avatar>
          </Dropdown>
        </Header>
        <Content style={{ margin: 24 }}>
          <div style={{ padding: 24, background: '#fff', minHeight: '85vh', borderRadius: 8 }}>
            <Outlet />
          </div>
        </Content>
      </AntLayout>
    </AntLayout>
  )
}
