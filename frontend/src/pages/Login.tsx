import { Button, Card, Form, Input, Segmented, message } from 'antd'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'

export default function Login() {
  const navigate = useNavigate()
  const login = useAuthStore((s) => s.login)
  const register = useAuthStore((s) => s.register)
  const [mode, setMode] = useState<'login' | 'register'>('login')
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()

  async function onFinish(values: { username: string; email?: string; password: string; nickname?: string; phone?: string }) {
    setLoading(true)
    try {
      if (mode === 'login') {
        await login(values.username, values.password)
        message.success('登录成功')
      } else {
        await register({ username: values.username, email: values.email || '', password: values.password, nickname: values.nickname, phone: values.phone })
        message.success('注册成功')
      }
      navigate('/')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ maxWidth: 420, margin: '40px auto' }}>
      <Card>
        <h2 style={{ textAlign: 'center' }}>{mode === 'login' ? '登录' : '注册'}</h2>
        <Segmented
          block
          value={mode}
          onChange={(v) => setMode(v as 'login' | 'register')}
          options={[
            { label: '登录', value: 'login' },
            { label: '注册', value: 'register' },
          ]}
          style={{ marginBottom: 16 }}
        />
        <Form form={form} layout="vertical" onFinish={onFinish}>
          <Form.Item name="username" label="用户名" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          {mode === 'register' && (
            <>
              <Form.Item name="email" label="邮箱" rules={[{ required: true, type: 'email' }]}>
                <Input />
              </Form.Item>
              <Form.Item name="nickname" label="昵称">
                <Input />
              </Form.Item>
              <Form.Item name="phone" label="电话">
                <Input />
              </Form.Item>
            </>
          )}
          <Form.Item name="password" label="密码" rules={[{ required: true, min: 6 }]}>
            <Input.Password />
          </Form.Item>
          <Button type="primary" htmlType="submit" block loading={loading}>
            {mode === 'login' ? '登录' : '注册并登录'}
          </Button>
        </Form>
      </Card>
    </div>
  )
}
