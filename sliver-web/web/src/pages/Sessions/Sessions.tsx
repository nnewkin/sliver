import React, { useEffect } from 'react'
import { Table, Button, Space, message, Modal } from 'antd'
import { DeleteOutlined, EditOutlined, DesktopOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useAppSelector, useAppDispatch } from '../../hooks/useStore'
import { setSessions, removeSession } from '../../stores/sessionSlice'
import { sessionAPI } from '../../services/api'
import type { Session } from '../../stores/sessionSlice'

export const Sessions: React.FC = () => {
  const dispatch = useAppDispatch()
  const { sessions, loading } = useAppSelector((state) => state.sessions)

  useEffect(() => {
    loadSessions()
  }, [])

  const loadSessions = async () => {
    try {
      const response = await sessionAPI.list()
      dispatch(setSessions(response.data.data || []))
    } catch (error) {
      message.error('Failed to load sessions')
    }
  }

  const handleKill = async (id: string) => {
    Modal.confirm({
      title: 'Kill Session',
      content: 'Are you sure you want to kill this session?',
      onOk: async () => {
        try {
          await sessionAPI.kill(id)
          dispatch(removeSession(id))
          message.success('Session killed')
        } catch (error) {
          message.error('Failed to kill session')
        }
      },
    })
  }

  const columns: ColumnsType<Session> = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: 'Name', dataIndex: 'name', key: 'name' },
    { title: 'Hostname', dataIndex: 'hostname', key: 'hostname' },
    { title: 'Username', dataIndex: 'username', key: 'username' },
    { title: 'OS', dataIndex: 'os', key: 'os' },
    { title: 'Arch', dataIndex: 'arch', key: 'arch' },
    { title: 'PID', dataIndex: 'pid', key: 'pid' },
    { title: 'Status', dataIndex: 'is_active', key: 'is_active', render: (active) => active ? 'Active' : 'Inactive' },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button type="link" icon={<DesktopOutlined />} onClick={() => {}}>Interact</Button>
          <Button type="link" icon={<EditOutlined />} onClick={() => {}}>Rename</Button>
          <Button type="link" danger icon={<DeleteOutlined />} onClick={() => handleKill(record.id)}>Kill</Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <h1>Sessions</h1>
      <Table columns={columns} dataSource={sessions} rowKey="id" loading={loading} />
    </div>
  )
}
