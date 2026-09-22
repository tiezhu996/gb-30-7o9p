import { Result, Button } from 'antd'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'

export default function RoleGuard({ roles, children }: { roles: string[]; children: React.ReactNode }) {
  const { user, isLoggedIn } = useAuth()
  const navigate = useNavigate()
  if (!isLoggedIn) {
    return <Result status="warning" title="请先登录" extra={<Button type="primary" onClick={() => navigate('/login')}>去登录</Button>} />
  }
  if (user && !roles.includes(user.role)) {
    return <Result status="403" title="无权限访问" subTitle="当前角色无权查看该内容" extra={<Button onClick={() => navigate('/')}>返回首页</Button>} />
  }
  return <>{children}</>
}
