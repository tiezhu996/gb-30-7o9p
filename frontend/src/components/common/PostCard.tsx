import { Card, Tag } from 'antd'
import { useNavigate } from 'react-router-dom'
import type { CommunityPost } from '@/types/api'
import { formatDateTime } from '@/utils/dateFormat'

export default function PostCard({ post }: { post: CommunityPost }) {
  const navigate = useNavigate()
  return (
    <Card hoverable style={{ marginBottom: 12 }} onClick={() => navigate(`/community/${post.id}`)}>
      <Tag color={post.post_type === 'story' ? 'green' : 'orange'}>
        {post.post_type === 'story' ? '救助故事' : '寻主公告'}
      </Tag>
      <h3 style={{ margin: '8px 0' }}>{post.title}</h3>
      <p style={{ color: '#888' }}>{post.content?.slice(0, 80)}</p>
      <div style={{ color: '#aaa', fontSize: 12 }}>
        {formatDateTime(post.created_at)} · 👍 {post.like_count} · 💬 {post.comment_count}
      </div>
    </Card>
  )
}
