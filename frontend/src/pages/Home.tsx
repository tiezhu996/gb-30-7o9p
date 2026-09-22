import { Card, Col, Row, Typography, Input } from 'antd'
import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import PetCard from '@/components/common/PetCard'
import PostCard from '@/components/common/PostCard'
import OrgCard from '@/components/common/OrgCard'
import type { Pet } from '@/constants/pet'
import type { Organization } from '@/constants/organization'
import type { CommunityPost } from '@/types/api'

export default function Home() {
  const navigate = useNavigate()
  const [hotPets, setHotPets] = useState<Pet[]>([])
  const [posts, setPosts] = useState<CommunityPost[]>([])
  const [orgs, setOrgs] = useState<Organization[]>([])
  useEffect(() => {
    axios.get('/api/v1/home/overview').then((res) => {
      setHotPets(res.data.data.hot_pets || [])
      setPosts(res.data.data.latest_posts || [])
      setOrgs(res.data.data.orgs || [])
    })
  }, [])

  return (
    <div>
      <div style={{ textAlign: 'center', padding: '32px 0 20px' }}>
        <Typography.Title level={1}>让每个流浪生命找到新家</Typography.Title>
        <Input.Search
          size="large"
          placeholder="搜索待领养动物"
          style={{ maxWidth: 560 }}
          onSearch={(v) => navigate(`/pets?keyword=${encodeURIComponent(v)}`)}
        />
      </div>
      <Typography.Title level={3}>🐕 热门待领养动物</Typography.Title>
      <Row gutter={[16, 16]}>
        {hotPets.map((p) => (
          <Col xs={12} sm={8} md={6} key={p.id}>
            <PetCard pet={p} />
          </Col>
        ))}
      </Row>
      <Typography.Title level={3}>📖 最新救助故事</Typography.Title>
      <Row gutter={[16, 16]}>
        {posts.map((p) => (
          <Col xs={24} sm={12} md={8} key={p.id}>
            <PostCard post={p} />
          </Col>
        ))}
      </Row>
      <Typography.Title level={3}>🏠 推荐机构</Typography.Title>
      <Row gutter={[16, 16]}>
        {orgs.map((o) => (
          <Col xs={24} sm={12} md={6} key={o.id}>
            <OrgCard org={o} />
          </Col>
        ))}
      </Row>
      {!hotPets.length && !posts.length && <Card>暂无数据</Card>}
    </div>
  )
}
