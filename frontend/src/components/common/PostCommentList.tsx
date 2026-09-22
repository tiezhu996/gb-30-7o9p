import { Button, Input, List, message } from 'antd'
import { useEffect, useState } from 'react'
import { listComments, createComment } from '@/api/post'
import type { PostComment } from '@/types/api'
import { formatDateTime } from '@/utils/dateFormat'
import { useAuth } from '@/hooks/useAuth'
import { useNavigate } from 'react-router-dom'

export default function PostCommentList({ postId }: { postId: number }) {
  const { isLoggedIn } = useAuth()
  const navigate = useNavigate()
  const [comments, setComments] = useState<PostComment[]>([])
  const [content, setContent] = useState('')

  async function load() {
    setComments(await listComments(postId))
  }
  useEffect(() => {
    load()
  }, [postId])

  async function submit() {
    if (!isLoggedIn) {
      message.warning('请先登录')
      navigate('/login')
      return
    }
    if (!content.trim()) return
    await createComment(postId, content)
    setContent('')
    await load()
  }

  return (
    <div>
      <List
        dataSource={comments}
        locale={{ emptyText: '暂无评论' }}
        renderItem={(c) => (
          <List.Item>
            <div>
              <div style={{ color: '#888', fontSize: 12 }}>用户 #{c.user_id} · {formatDateTime(c.created_at)}</div>
              <div>{c.content}</div>
            </div>
          </List.Item>
        )}
      />
      <Input.TextArea rows={2} value={content} onChange={(e) => setContent(e.target.value)} placeholder="写下你的评论…" style={{ marginTop: 8 }} />
      <Button type="primary" style={{ marginTop: 8 }} onClick={submit}>
        发表评论
      </Button>
    </div>
  )
}
