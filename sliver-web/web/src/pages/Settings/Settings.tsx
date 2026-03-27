import React from 'react'
import { Card, Tabs, Form, Input, Button, Switch, message, Table } from 'antd'
import { useAppSelector } from '../../hooks/useStore'
import { userAPI, botAPI } from '../../services/api'

export const Settings: React.FC = () => {
  const { user } = useAppSelector((state) => state.auth)
  const isAdmin = user?.role === 'super_admin'

  return (
    <div>
      <h1>Settings</h1>
      <Tabs defaultActiveKey="profile">
        <Tabs.TabPane key="profile" tab="Profile">
          <Card title="Profile Settings">
            <p>Username: {user?.username}</p>
            <p>Role: {user?.role}</p>
          </Card>
        </Tabs.TabPane>
        {isAdmin && (
          <Tabs.TabPane key="users" tab="User Management">
            <Card title="Users">
              <Table
                columns={[
                  { title: 'ID', dataIndex: 'id', key: 'id' },
                  { title: 'Username', dataIndex: 'username', key: 'username' },
                  { title: 'Role', dataIndex: 'role', key: 'role' },
                ]}
                dataSource={[]}
              />
            </Card>
          </Tabs.TabPane>
        )}
        {isAdmin && (
          <Tabs.TabPane key="bot" tab="Telegram Bot">
            <Card title="Bot Configuration">
              <Form layout="vertical">
                <Form.Item label="Bot Token">
                  <Input.Password placeholder="Enter bot token" />
                </Form.Item>
                <Form.Item label="Enabled">
                  <Switch />
                </Form.Item>
                <Button type="primary">Save</Button>
              </Form>
            </Card>
          </Tabs.TabPane>
        )}
        <Tabs.TabPane key="about" tab="About">
          <Card title="About Sliver Web">
            <p>Sliver Web UI - v1.0.0</p>
            <p>Web-based management interface for Sliver C2 framework</p>
          </Card>
        </Tabs.TabPane>
      </Tabs>
    </div>
  )
}
