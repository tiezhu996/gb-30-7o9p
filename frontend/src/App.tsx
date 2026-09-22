import { Layout, Menu, Dropdown, Button, Avatar } from 'antd'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { UserOutlined } from '@ant-design/icons'
import { useEffect } from 'react'
import { useAuth } from '@/hooks/useAuth'

const { Header, Content } = Layout

const NAV_ITEMS = [
  { key: '/', label: '首页' },
  { key: '/pets', label: '待领养动物' },
  { key: '/orgs', label: '救助机构' },
  { key: '/community', label: '救助社区' },
  { key: '/donate', label: '公益捐赠' },
  { key: '/reviews', label: '领养回访' },
  { key: '/applications', label: '我的申请' },
]

export default function App() {
  const navigate = useNavigate()
  const location = useLocation()
  const { token, user, logout, fetchProfile } = useAuth()

  useEffect(() => {
    if (token && !user) {
      fetchProfile().catch(() => undefined)
    }
  }, [token])

  const selected = NAV_ITEMS.find((n) => location.pathname.startsWith(n.key))?.key || '/'

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center', gap: 24, background: '#fff' }}>
        <div style={{ fontSize: 20, fontWeight: 700, color: '#2f6f4e', cursor: 'pointer', whiteSpace: 'nowrap' }} onClick={() => navigate('/')}>
          🐾 宠物领养平台
        </div>
        <Menu mode="horizontal" selectedKeys={[selected]} items={NAV_ITEMS} onClick={(e) => navigate(e.key)} style={{ flex: 1 }} />
        {token ? (
          <Dropdown
            menu={{
              items: [
                { key: 'profile', label: '个人中心' },
                { key: 'logout', label: '退出登录' },
              ],
              onClick: ({ key }) => {
                if (key === 'profile') navigate('/profile')
                if (key === 'logout') {
                  logout()
                  navigate('/')
                }
              },
            }}
          >
            <span style={{ cursor: 'pointer' }}>
              <Avatar size="small" icon={<UserOutlined />} style={{ marginRight: 6 }} />
              {user?.nickname || user?.username}
            </span>
          </Dropdown>
        ) : (
          <Button type="primary" onClick={() => navigate('/login')}>
            登录/注册
          </Button>
        )}
      </Header>
      <Content style={{ padding: 24, maxWidth: 1200, width: '100%', margin: '0 auto' }}>
        <Outlet />
      </Content>
    </Layout>
  )
}
