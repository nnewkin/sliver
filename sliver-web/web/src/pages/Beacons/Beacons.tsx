import React, { useEffect } from 'react'
import { Table, Button, Space, message } from 'antd'
import { AimOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useAppSelector, useAppDispatch } from '../../hooks/useStore'
import { setBeacons } from '../../stores/beaconSlice'
import { beaconAPI } from '../../services/api'
import type { Beacon } from '../../stores/beaconSlice'

export const Beacons: React.FC = () => {
  const dispatch = useAppDispatch()
  const { beacons, loading } = useAppSelector((state) => state.beacons)

  useEffect(() => {
    loadBeacons()
  }, [])

  const loadBeacons = async () => {
    try {
      const response = await beaconAPI.list()
      dispatch(setBeacons(response.data.data || []))
    } catch (error) {
      message.error('Failed to load beacons')
    }
  }

  const columns: ColumnsType<Beacon> = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: 'Name', dataIndex: 'name', key: 'name' },
    { title: 'Hostname', dataIndex: 'hostname', key: 'hostname' },
    { title: 'Username', dataIndex: 'username', key: 'username' },
    { title: 'OS', dataIndex: 'os', key: 'os' },
    { title: 'Interval', dataIndex: 'interval', key: 'interval' },
    { title: 'Jitter', dataIndex: 'jitter', key: 'jitter' },
    { title: 'Last Checkin', dataIndex: 'last_checkin', key: 'last_checkin' },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button type="link" icon={<AimOutlined />} onClick={() => {}}>Interact</Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <h1>Beacons</h1>
      <Table columns={columns} dataSource={beacons} rowKey="id" loading={loading} />
    </div>
  )
}
