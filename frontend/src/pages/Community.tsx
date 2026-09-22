import { Button, Card, Input, Modal, Radio, message } from 'antd'
import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import PostCard from '@/components/common/PostCard'
import SearchFilter from '@/components/common/SearchFilter'
import { listPosts, createPost } from '@/api/post'
import { useAuth } from '@/hooks/useAuth'
import type { CommunityPost } from '@/types/api'

export default function Community() {
  const navigate = useNavigate()
  const { isLoggedIn } = useAuth()
  const [posts, setPosts] = useState<CommunityPost[]>([])
  const [postType, setPostType] = useState('')
  const [open, setOpen] = useState(false)
  const [form, setForm] = useState({ title: '', content: '', post_type: 'story' })

  async function load(keyword = '') {
    const res = await listPosts({ page: 1, page_size: 20, post_type: postType, keyword })
    setPosts(res.list)
  }
  useEffect(() => {
    load()
  }, [postType])

  async function submit() {
    if (!isLoggedIn) {
      message.warning('请先登录')
      navigate('/login')
      return
    }
    await createPost(form)
    message.success('帖子发布成功')
    setOpen(false)
    setForm({ title: '', content: '', post_type: 'story' })
    await load()
  }

  return (
    <div>
      <h1>救助社区</h1>
      <SearchFilter onSearch={load}>
        <Radio.Group value={postType} onChange={(e) => setPostType(e.target.value)}>
          <Radio.Button value="">全部</Radio.Button>
          <Radio.Button value="story">救助故事</Radio.Button>
          <Radio.Button value="lost_notice">寻主公告</Radio.Button>
        </Radio.Group>
        <Button type="primary" onClick={() => setOpen(true)}>
          发布帖子
        </Button>
      </SearchFilter>
      {posts.map((p) => (
        <PostCard key={p.id} post={p} />
      ))}
      {!posts.length && <Card>暂无帖子</Card>}
      <Modal open={open} title="发布帖子" onOk={submit} onCancel={() => setOpen(false)} okText="发布">
        <Radio.Group value={form.post_type} onChange={(e) => setForm({ ...form, post_type: e.target.value })} style={{ marginBottom: 12 }}>
          <Radio.Button value="story">救助故事</Radio.Button>
          <Radio.Button value="lost_notice">寻主公告</Radio.Button>
        </Radio.Group>
        <Input placeholder="标题" value={form.title} onChange={(e) => setForm({ ...form, title: e.target.value })} style={{ marginBottom: 12 }} />
        <Input.TextArea rows={4} placeholder="内容" value={form.content} onChange={(e) => setForm({ ...form, content: e.target.value })} />
      </Modal>
    </div>
  )
}
