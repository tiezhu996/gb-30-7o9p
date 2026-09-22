import { Button, message } from 'antd'
import { StarFilled, StarOutlined } from '@ant-design/icons'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { addFavorite, removeFavorite } from '@/api/favorite'

export default function FavoriteButton({ targetType, targetId }: { targetType: string; targetId: number }) {
  const { isLoggedIn } = useAuth()
  const navigate = useNavigate()
  const [favorited, setFavorited] = useState(false)
  const [loading, setLoading] = useState(false)

  async function toggle() {
    if (!isLoggedIn) {
      message.warning('请先登录')
      navigate('/login')
      return
    }
    setLoading(true)
    try {
      if (favorited) {
        await removeFavorite(targetType, targetId)
        setFavorited(false)
        message.success('已取消收藏')
      } else {
        await addFavorite(targetType, targetId)
        setFavorited(true)
        message.success('收藏成功')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <Button icon={favorited ? <StarFilled /> : <StarOutlined />} loading={loading} onClick={toggle}>
      {favorited ? '已收藏' : '收藏'}
    </Button>
  )
}
