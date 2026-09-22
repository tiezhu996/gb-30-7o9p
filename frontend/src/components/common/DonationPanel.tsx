import { Button, Card, Form, InputNumber, message } from 'antd'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { donate } from '@/api/donation'
import { useAuth } from '@/hooks/useAuth'

export default function DonationPanel({ orgId, orgName }: { orgId: number; orgName: string }) {
  const { isLoggedIn } = useAuth()
  const navigate = useNavigate()
  const [amount, setAmount] = useState<number>(50)
  const [loading, setLoading] = useState(false)

  async function submit() {
    if (!isLoggedIn) {
      message.warning('请先登录')
      navigate('/login')
      return
    }
    if (!amount || amount <= 0) {
      message.warning('请输入有效金额')
      return
    }
    setLoading(true)
    try {
      await donate({ org_id: orgId, amount })
      message.success('捐赠成功，感谢您的爱心！')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card title={`向 ${orgName} 捐款`}>
      <Form layout="vertical">
        <Form.Item label="捐款金额（元）">
          <InputNumber min={1} value={amount} onChange={(v) => setAmount(v ?? 0)} style={{ width: '100%' }} />
        </Form.Item>
        <Button type="primary" block loading={loading} onClick={submit}>
          立即捐赠
        </Button>
      </Form>
    </Card>
  )
}
