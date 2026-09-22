import { Button, Card, Col, Descriptions, Row, Space, Tag, message } from 'antd'
import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { getPet } from '@/api/pet'
import FavoriteButton from '@/components/common/FavoriteButton'
import ImageGallery from '@/components/common/ImageGallery'
import OrgCard from '@/components/common/OrgCard'
import { PetSpeciesMap, PetStatusMap, type Pet } from '@/constants/pet'
import { getOrg } from '@/api/organization'
import type { Organization } from '@/constants/organization'
import { useAuth } from '@/hooks/useAuth'

export default function PetDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { isLoggedIn } = useAuth()
  const [pet, setPet] = useState<Pet | null>(null)
  const [org, setOrg] = useState<Organization | null>(null)

  useEffect(() => {
    getPet(id!).then(async (p) => {
      setPet(p)
      try {
        setOrg(await getOrg(p.org_id))
      } catch {
        setOrg(null)
      }
    })
  }, [id])

  if (!pet) return <p>加载中…</p>

  const status = PetStatusMap[pet.status]

  return (
    <div>
      <h1>
        {pet.name} <Tag color="blue">{PetSpeciesMap[pet.species]}</Tag>
        <Tag color={status.color}>{status.text}</Tag>
      </h1>
      <Row gutter={16}>
        <Col xs={24} md={14}>
          <ImageGallery imageUrls={pet.image_urls} />
          <Card title="健康档案" style={{ marginTop: 16 }}>
            <Descriptions column={2} size="small">
              <Descriptions.Item label="品种">{pet.breed || '-'}</Descriptions.Item>
              <Descriptions.Item label="年龄">{pet.age} 岁</Descriptions.Item>
              <Descriptions.Item label="性别">{pet.gender === 'male' ? '公' : pet.gender === 'female' ? '母' : '-'}</Descriptions.Item>
              <Descriptions.Item label="体型">{pet.size || '-'}</Descriptions.Item>
              <Descriptions.Item label="绝育">{pet.neutered ? '已绝育' : '未绝育'}</Descriptions.Item>
              <Descriptions.Item label="疫苗">{pet.vaccinated ? '已接种' : '未接种'}</Descriptions.Item>
              <Descriptions.Item label="健康状况">{pet.health_status || '-'}</Descriptions.Item>
              <Descriptions.Item label="性格">{pet.personality || '-'}</Descriptions.Item>
            </Descriptions>
            <p style={{ marginTop: 12 }}>{pet.description}</p>
          </Card>
        </Col>
        <Col xs={24} md={10}>
          {org && <OrgCard org={org} />}
          <Card title="领养操作" style={{ marginTop: 16 }}>
            <Space direction="vertical" style={{ width: '100%' }}>
              <FavoriteButton targetType="pet" targetId={pet.id} />
              <Button
                type="primary"
                block
                disabled={pet.status !== 'available'}
                onClick={() => {
                  if (!isLoggedIn) {
                    message.warning('请先登录')
                    navigate('/login')
                    return
                  }
                  navigate(`/apply/${pet.id}`)
                }}
              >
                {pet.status === 'available' ? '发起领养申请' : '暂不可申请'}
              </Button>
            </Space>
          </Card>
        </Col>
      </Row>
    </div>
  )
}
