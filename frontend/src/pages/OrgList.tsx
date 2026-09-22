import { Col, Row, Select } from 'antd'
import { useEffect, useState } from 'react'
import OrgCard from '@/components/common/OrgCard'
import SearchFilter from '@/components/common/SearchFilter'
import { listOrgs } from '@/api/organization'
import type { Organization } from '@/constants/organization'

export default function OrgList() {
  const [orgs, setOrgs] = useState<Organization[]>([])
  const [city, setCity] = useState('')

  async function load(keyword: string) {
    const res = await listOrgs({ page: 1, page_size: 20, keyword })
    setOrgs(res.list)
  }
  useEffect(() => {
    load('')
  }, [city])

  return (
    <div>
      <h1>救助机构</h1>
      <SearchFilter onSearch={load}>
        <Select
          placeholder="城市"
          allowClear
          style={{ width: 140 }}
          value={city || undefined}
          onChange={(v) => setCity(v || '')}
          options={['上海', '北京', '广州', '深圳'].map((c) => ({ value: c, label: c }))}
        />
      </SearchFilter>
      <Row gutter={[16, 16]}>
        {orgs.map((o) => (
          <Col xs={24} sm={12} md={8} key={o.id}>
            <OrgCard org={o} />
          </Col>
        ))}
      </Row>
      {!orgs.length && <p style={{ color: '#999' }}>暂无机构</p>}
    </div>
  )
}
