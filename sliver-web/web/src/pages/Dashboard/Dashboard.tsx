import React, { useEffect } from 'react'
import { Row, Col, Card, Statistic, Table } from 'antd'
import { DesktopOutlined, AimOutlined, BellOutlined, AlertOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useAppSelector, useAppDispatch } from '../../hooks/useStore'
import { setSessions } from '../../stores/sessionSlice'
import { setBeacons } from '../../stores/beaconSlice'
import { setListeners } from '../../stores/listenerSlice'
import { sessionAPI, beaconAPI, listenerAPI } from '../../services/api'
import type { Session } from '../../stores/sessionSlice'
import type { Beacon } from '../../stores/beaconSlice'
import type { Listener } from '../../stores/listenerSlice'

export const Dashboard: React.FC = () => {
  const dispatch = useAppDispatch()
  const sessions = useAppSelector((state) => state.sessions.sessions)
  const beacons = useAppSelector((state) => state.beacons.beacons)
  const listeners = useAppSelector((state) => state.listeners.listeners)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const [sessionsRes, beaconsRes, listenersRes] = await Promise.all([
        sessionAPI.list(),
        beaconAPI.list(),
        listenerAPI.list(),
      ])
      dispatch(setSessions(sessionsRes.data.data || []))
      dispatch(setBeacons(beaconsRes.data.data || []))
      dispatch(setListeners(listenersRes.data.data || []))
    } catch (error) {
      console.error('Failed to load data:', error)
    }
  }

  const sessionColumns: ColumnsType<Session> = [
    { title: 'Name', dataIndex: 'name', key: 'name' },
    { title: 'Hostname', dataIndex: 'hostname', key: 'hostname' },
    { title: 'Username', dataIndex: 'username', key: 'username' },
    { title: 'OS', dataIndex: 'os', key: 'os' },
    { title: 'Status', dataIndex: 'is_active', key: 'is_active', render: (active) => active ? 'Active' : 'Inactive' },
  ]

  const beaconColumns: ColumnsType<Beacon> = [
    { title: 'Name', dataIndex: 'name', key: 'name' },
    { title: 'Hostname', dataIndex: 'hostname', key: 'hostname' },
    { title: 'Username', dataIndex: 'username', key: 'username' },
    { title: 'Interval', dataIndex: 'interval', key: 'interval' },
    { title: 'Last Checkin', dataIndex: 'last_checkin', key: 'last_checkin' },
  ]

  return (
    <div>
      <h1>Dashboard</h1>
      <Row gutter={16}>
        <Col span={6}>
          <Card>
            <Statistic
              title="Active Sessions"
              value={sessions.filter((s) => s.is_active).length}
              prefix={<DesktopOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="Active Beacons"
              value={beacons.filter((b) => b.is_active).length}
              prefix={<AimOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="Listeners"
              value={listeners.length}
              prefix={<BellOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="Alerts"
              value={0}
              prefix={<AlertOutlined />}
            />
          </Card>
        </Col>
      </Row>
      <Row gutter={16} style={{ marginTop: 24 }}>
        <Col span={12}>
          <Card title="Recent Sessions" size="small">
            <Table columns={sessionColumns} dataSource={sessions.slice(0, 5)} rowKey="id" pagination={false} />
          </Card>
        </Col>
        <Col span={12}>
          <Card title="Recent Beacons" size="small">
            <Table columns={beaconColumns} dataSource={beacons.slice(0, 5)} rowKey="id" pagination={false} />
          </Card>
        </Col>
      </Row>
    </div>
  )
}
