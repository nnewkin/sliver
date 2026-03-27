import React, { useState } from 'react'
import { Card, Form, Input, Select, Button, message, Divider, Switch, InputNumber, Space } from 'antd'
import { DownloadOutlined } from '@ant-design/icons'
import { implantAPI } from '../../services/api'

const { Option } = Select

export const Generate: React.FC = () => {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)

  const onFinish = async (values: any) => {
    setLoading(true)
    try {
      const config = {
        format: values.format,
        goos: values.platform,
        goarch: values.arch,
        mtls_addresses: values.mtls_addresses?.split(',') || [],
        dns_addresses: values.dns_addresses?.split(',') || [],
        http_addresses: values.http_addresses?.split(',') || [],
        https_addresses: values.https_addresses?.split(',') || [],
        interval: values.interval || 60,
        jitter: values.jitter || 10,
        name: values.name,
        evasion: values.evasion || false,
      }
      const response = await implantAPI.generate(config)
      const blob = new Blob([response.data.data.binary], { type: 'application/octet-stream' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = values.name || 'implant.exe'
      a.click()
      URL.revokeObjectURL(url)
      message.success('Implant generated successfully')
    } catch (error) {
      message.error('Failed to generate implant')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div>
      <h1>Generate Implant</h1>
      <Card>
        <Form form={form} layout="vertical" onFinish={onFinish}>
          <Divider>Basic Settings</Divider>
          <Space style={{ display: 'flex' }}>
            <Form.Item name="name" label="Name" rules={[{ required: true }]}>
              <Input placeholder="my-implant" />
            </Form.Item>
            <Form.Item name="platform" label="Platform" rules={[{ required: true }]}>
              <Select style={{ width: 120 }}>
                <Option value="windows">Windows</Option>
                <Option value="linux">Linux</Option>
                <Option value="darwin">macOS</Option>
              </Select>
            </Form.Item>
            <Form.Item name="arch" label="Architecture" rules={[{ required: true }]}>
              <Select style={{ width: 120 }}>
                <Option value="amd64">x64</Option>
                <Option value="386">x86</Option>
                <Option value="arm64">ARM64</Option>
                <Option value="arm">ARM</Option>
              </Select>
            </Form.Item>
            <Form.Item name="format" label="Format" rules={[{ required: true }]}>
              <Select style={{ width: 120 }}>
                <Option value="EXE">EXE</Option>
                <Option value="DLL">DLL</Option>
                <Option value="SHARED_LIBRARY">Shared Lib</Option>
                <Option value="SHELLCODE">Shellcode</Option>
              </Select>
            </Form.Item>
          </Space>

          <Divider>C2 Addresses</Divider>
          <Space style={{ display: 'flex' }}>
            <Form.Item name="mtls_addresses" label="mTLS Addresses">
              <Input placeholder="host:port,host:port" style={{ width: 300 }} />
            </Form.Item>
            <Form.Item name="http_addresses" label="HTTP Addresses">
              <Input placeholder="https://domain.com,https://domain2.com" style={{ width: 300 }} />
            </Form.Item>
          </Space>

          <Divider>Beacon Settings</Divider>
          <Space style={{ display: 'flex' }}>
            <Form.Item name="interval" label="Interval (seconds)">
              <InputNumber min={1} defaultValue={60} />
            </Form.Item>
            <Form.Item name="jitter" label="Jitter (%)">
              <InputNumber min={0} max={100} defaultValue={10} />
            </Form.Item>
          </Space>

          <Divider>Evasion</Divider>
          <Form.Item name="evasion" label="Enable Evasion" valuePropName="checked">
            <Switch />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} icon={<DownloadOutlined />}>
              Generate
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}
