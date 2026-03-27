import React, { useState } from 'react'
import { Form, Input, Button, Card, message, Tabs } from 'antd'
import { UserOutlined, LockOutlined, SafetyOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { useAppDispatch, useAppSelector } from '../../hooks/useStore'
import { loginStart, loginSuccess, loginFailure, verify2FASuccess } from '../../stores/authSlice'
import { authAPI } from '../../services/api'

export const Login: React.FC = () => {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const { loading, error, need2FA } = useAppSelector((state) => state.auth)
  const [form] = Form.useForm()
  const [verifyForm] = Form.useForm()
  const [tempToken, setTempToken] = useState('')

  const onLoginFinish = async (values: { username: string; password: string }) => {
    dispatch(loginStart())
    try {
      const response = await authAPI.login(values)
      if (response.data.data.need_2fa) {
        setTempToken(response.data.data.token)
        message.info('Please enter your 2FA code')
      } else {
        dispatch(loginSuccess({
          token: response.data.data.token,
          user: response.data.data.user,
          need2FA: false,
        }))
        navigate('/')
      }
    } catch (err: any) {
      dispatch(loginFailure(err.response?.data?.message || 'Login failed'))
      message.error(err.response?.data?.message || 'Login failed')
    }
  }

  const onVerifyFinish = async (values: { code: string }) => {
    try {
      const response = await authAPI.verify2FA({
        token: tempToken,
        code: values.code,
      })
      dispatch(verify2FASuccess({
        token: response.data.data.token,
        user: response.data.data.user,
      }))
      navigate('/')
    } catch (err: any) {
      message.error(err.response?.data?.message || '2FA verification failed')
    }
  }

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh', background: '#f0f2f5' }}>
      <Card title="Sliver Web UI" style={{ width: 400 }}>
        <Tabs defaultActiveKey={need2FA ? '2fa' : 'login'}>
          <Tabs.TabPane key="login" tab="Login">
            <Form
              form={form}
              name="login"
              onFinish={onLoginFinish}
              layout="vertical"
            >
              <Form.Item
                name="username"
                rules={[{ required: true, message: 'Please enter your username' }]}
              >
                <Input prefix={<UserOutlined />} placeholder="Username" size="large" />
              </Form.Item>
              <Form.Item
                name="password"
                rules={[{ required: true, message: 'Please enter your password' }]}
              >
                <Input.Password prefix={<LockOutlined />} placeholder="Password" size="large" />
              </Form.Item>
              {error && <div style={{ color: 'red', marginBottom: 16 }}>{error}</div>}
              <Form.Item>
                <Button type="primary" htmlType="submit" loading={loading} block size="large">
                  Login
                </Button>
              </Form.Item>
            </Form>
          </Tabs.TabPane>
          <Tabs.TabPane key="2fa" tab="2FA Verification">
            <Form
              form={verifyForm}
              name="verify"
              onFinish={onVerifyFinish}
              layout="vertical"
            >
              <Form.Item
                name="code"
                rules={[{ required: true, message: 'Please enter your 2FA code' }]}
              >
                <Input prefix={<SafetyOutlined />} placeholder="Enter 2FA Code" size="large" />
              </Form.Item>
              <Form.Item>
                <Button type="primary" htmlType="submit" block size="large">
                  Verify
                </Button>
              </Form.Item>
            </Form>
          </Tabs.TabPane>
        </Tabs>
      </Card>
    </div>
  )
}
