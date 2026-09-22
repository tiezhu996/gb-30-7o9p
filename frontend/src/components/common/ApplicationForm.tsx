import { Button, Form, Input, message } from 'antd'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { submitApplication } from '@/api/application'
import { useAuth } from '@/hooks/useAuth'

export default function ApplicationForm({ petId }: { petId: number }) {
  const { isLoggedIn } = useAuth()
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()

  async function onFinish(values: { hasYard: string; petExperience: string; family: string }) {
    if (!isLoggedIn) {
      message.warning('请先登录')
      navigate('/login')
      return
    }
    setLoading(true)
    try {
      const questionnaire = JSON.stringify({
        has_yard: values.hasYard === 'yes',
        pet_experience: values.petExperience,
        family: values.family,
      })
      await submitApplication({ pet_id: petId, questionnaire })
      message.success('领养申请已提交')
      navigate('/applications')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Form form={form} layout="vertical" onFinish={onFinish}>
      <Form.Item name="family" label="家庭情况" rules={[{ required: true, message: '请填写家庭情况' }]}>
        <Input placeholder="例如：三口之家，有稳定住所" />
      </Form.Item>
      <Form.Item name="hasYard" label="是否有庭院/阳台" rules={[{ required: true, message: '请选择' }]}>
        <Input placeholder="是 / 否" />
      </Form.Item>
      <Form.Item name="petExperience" label="养宠经验" rules={[{ required: true, message: '请填写养宠经验' }]}>
        <Input.TextArea rows={3} placeholder="请描述你的养宠经验" />
      </Form.Item>
      <Button type="primary" htmlType="submit" loading={loading} block>
        提交领养申请
      </Button>
    </Form>
  )
}
