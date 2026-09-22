import { Avatar, Button, Card, Col, Descriptions, Form, Input, Row, message } from 'antd'
import { useEffect } from 'react'
import { useUserStore } from '@/stores/userStore'

export default function Profile() {
  const { profile, load, save } = useUserStore()
  const [form] = Form.useForm()

  useEffect(() => {
    load().then(() => form.setFieldsValue(profile || {}))
  }, [])

  async function onFinish(values: { nickname: string; phone: string; bio: string }) {
    await save(values)
    message.success('资料已更新')
  }

  return (
    <Row gutter={16}>
      <Col xs={24} md={10}>
        <Card title="个人资料">
          <Avatar size={80} style={{ marginBottom: 12 }}>{profile?.nickname?.[0] || '宠'}</Avatar>
          <Descriptions column={1} size="small">
            <Descriptions.Item label="用户名">{profile?.username}</Descriptions.Item>
            <Descriptions.Item label="邮箱">{profile?.email}</Descriptions.Item>
            <Descriptions.Item label="角色">
              {profile?.role === 'org' ? '救助机构' : profile?.role === 'admin' ? '管理员' : '领养人'}
            </Descriptions.Item>
          </Descriptions>
          <Form form={form} layout="vertical" onFinish={onFinish} style={{ marginTop: 16 }}>
            <Form.Item name="nickname" label="昵称">
              <Input />
            </Form.Item>
            <Form.Item name="phone" label="电话">
              <Input />
            </Form.Item>
            <Form.Item name="bio" label="简介">
              <Input.TextArea rows={3} />
            </Form.Item>
            <Button type="primary" htmlType="submit">
              保存修改
            </Button>
          </Form>
        </Card>
      </Col>
    </Row>
  )
}
