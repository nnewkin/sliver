import React, { useEffect } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Select } from 'antd'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useAppSelector, useAppDispatch } from '../../hooks/useStore'
import { setListeners, addListener, removeListener } from '../../stores/listenerSlice'
import { listenerAPI } from '../../services/api'
import type { Listener } from '../../stores/listenerSlice'

export const Listeners: React.FC = () => {
  const dispatch = useAppDispatch()
  const { listeners, loading } = useAppSelector((state) => state.listeners)
  const [form] = Form.useForm()

  useEffect(() => {
    loadListeners()
  }, [])

  const loadListeners = async () => {
    try {
      const response = await listenerAPI.list()
      dispatch(setListeners(response.data.data || []))
    } catch (error) {
      message.error('Failed to load listeners')
    }
  }

  const handleStart = async (values: any) => {
    try {
      let response
      switch (values.type) {
        case 'mtls':
          response = await listenerAPI.startMTLS(values.name, values.bind_address)
          break
        case 'dns':
          response = await listenerAPI.startDNS(values.name, values.bind_address, values.domains?.split(',') || [])
          break
        case 'http':
          response = await listenerAPI.startHTTP(values.name, values.bind_address, values.domains?.split(',') || [])
          break
        case 'https':
          response = await listenerAPI.startHTTPS(values.name, values.bind_address, values.domains?.split(',') || [])
          break
        default:
          message.error('Unknown listener type')
          return
      }
      dispatch(addListener(response.data.data))
      message.success('Listener started')
    } catch (error) {
      message.error('Failed to start listener')
    }
  }

  const handleStop = async (id: string) => {
    Modal.confirm({
      title: 'Stop Listener',
      content: 'Are you sure you want to stop this listener?',
      onOk: async () => {
        try {
          await listenerAPI.stop(id)
          dispatch(removeListener(id))
          message.success('Listener stopped')
        } catch (error) {
          message.error('Failed to stop listener')
        }
      },
    })
  }

  const columns: ColumnsType<Listener> = [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    { title: 'Name', dataIndex: 'name', key: 'name' },
    { title: 'Type', dataIndex: 'type', key: 'type' },
    { title: 'Status', dataIndex: 'status', key: 'status' },
    { title: 'Bind Address', dataIndex: 'bind_address', key: 'bind_address' },
    { title: 'Port', dataIndex: 'port', key: 'port' },
    { title: 'Connected', dataIndex: 'connected', key: 'connected' },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button type="link" danger icon={<DeleteOutlined />} onClick={() => handleStop(record.id)}>Stop</Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <h1>Listeners</h1>
      <Space style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => {
          Modal.confirm({
            title: 'Start New Listener',
            content: (
              <Form form={form} layout="vertical">
                <Form.Item name="type" label="Type" rules={[{ required: true }]}>
                  <Select>
                    <Select.Option value="mtls">mTLS</Select.Option>
                    <Select.Option value="dns">DNS</Select.Option>
                    <Select.Option value="http">HTTP</Select.Option>
                    <Select.Option value="https">HTTPS</Select.Option>
                  </Select>
                </Form.Item>
                <Form.Item name="name" label="Name" rules={[{ required: true }]}>
                  <Input />
                </Form.Item>
                <Form.Item name="bind_address" label="Bind Address" rules={[{ required: true }]}>
                  <Input placeholder="0.0.0.0:8888" />
                </Form.Item>
                <Form.Item name="domains" label="Domains (for DNS/HTTP/HTTPS)">
                  <Input placeholder="domain1.com,domain2.com" />
                </Form.Item>
              </Form>
            ),
            onOk: () => form.validateFields().then(handleStart),
          })
        }}>
          Start Listener
        </Button>
      </Space>
      <Table columns={columns} dataSource={listeners} rowKey="id" loading={loading} />
    </div>
  )
}
