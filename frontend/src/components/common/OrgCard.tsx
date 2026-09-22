import { Card, Tag } from 'antd'
import { useNavigate } from 'react-router-dom'
import { OrganizationStatusMap, type Organization } from '@/constants/organization'

export default function OrgCard({ org }: { org: Organization }) {
  const navigate = useNavigate()
  const status = OrganizationStatusMap[org.status]
  return (
    <Card hoverable onClick={() => navigate(`/orgs/${org.id}`)} style={{ marginBottom: 12 }}>
      <Card.Meta
        title={
          <span>
            {org.name} <Tag color={status.color}>{status.text}</Tag>
          </span>
        }
        description={`${org.city || '-'} · ${org.contact || '-'}`}
      />
    </Card>
  )
}
