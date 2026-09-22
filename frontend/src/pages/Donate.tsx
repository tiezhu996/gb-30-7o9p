import { Col, Row } from 'antd'
import { useEffect, useState } from 'react'
import DonationPanel from '@/components/common/DonationPanel'
import OrgCard from '@/components/common/OrgCard'
import { listOrgs } from '@/api/organization'
import type { Organization } from '@/constants/organization'

export default function Donate() {
  const [orgs, setOrgs] = useState<Organization[]>([])

  useEffect(() => {
    listOrgs({ page: 1, page_size: 20 }).then((res) => setOrgs(res.list))
  }, [])

  return (
    <div>
      <h1>公益捐赠</h1>
      <Row gutter={[16, 16]}>
        {orgs.map((o) => (
          <Col xs={24} md={12} key={o.id}>
            <Row gutter={[12, 12]}>
              <Col span={14}>
                <OrgCard org={o} />
              </Col>
              <Col span={10}>
                <DonationPanel orgId={o.id} orgName={o.name} />
              </Col>
            </Row>
          </Col>
        ))}
      </Row>
    </div>
  )
}
