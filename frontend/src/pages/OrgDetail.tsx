import { Card, Col, Row, Table, Tag, Typography } from 'antd'
import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { getOrg } from '@/api/organization'
import { listPets } from '@/api/pet'
import PetCard from '@/components/common/PetCard'
import { listUsages, donationStats } from '@/api/donation'
import { OrganizationStatusMap, type Organization } from '@/constants/organization'
import type { Pet } from '@/constants/pet'
import type { DonationUsage } from '@/types/api'
import { formatDate } from '@/utils/dateFormat'

export default function OrgDetail() {
  const { id } = useParams()
  const [org, setOrg] = useState<Organization | null>(null)
  const [pets, setPets] = useState<Pet[]>([])
  const [usages, setUsages] = useState<DonationUsage[]>([])
  const [stats, setStats] = useState(0)

  useEffect(() => {
    getOrg(id!).then(async (o) => {
      setOrg(o)
      const res = await listPets({ page: 1, page_size: 20 })
      setPets(res.list.filter((p) => p.org_id === o.id))
      setUsages(await listUsages(o.id))
      setStats((await donationStats(o.id)).total_donated)
    })
  }, [id])

  if (!org) return <p>加载中…</p>
  const status = OrganizationStatusMap[org.status]

  return (
    <div>
      <h1>
        {org.name} <Tag color={status.color}>{status.text}</Tag>
      </h1>
      <Typography.Paragraph>{org.description}</Typography.Paragraph>
      <Typography.Text type="secondary">
        {org.city} · {org.contact}
      </Typography.Text>
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} md={16}>
          <Card title="待领养动物">
            <Row gutter={[12, 12]}>
              {pets.map((p) => (
                <Col xs={12} sm={8} key={p.id}>
                  <PetCard pet={p} />
                </Col>
              ))}
            </Row>
            {!pets.length && <p style={{ color: '#999' }}>暂无动物</p>}
          </Card>
        </Col>
        <Col xs={24} md={8}>
          <Card title="捐款使用明细">
            <p style={{ fontSize: 20, color: '#2f6f4e' }}>累计收到捐赠：¥{stats.toFixed(2)}</p>
            <Table
              size="small"
              rowKey="id"
              dataSource={usages}
              pagination={false}
              columns={[
                { title: '用途', dataIndex: 'usage_desc' },
                { title: '金额', dataIndex: 'amount', render: (v: number) => `¥${Number(v).toFixed(2)}` },
                { title: '时间', dataIndex: 'created_at', render: (v: string) => formatDate(v) },
              ]}
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}
