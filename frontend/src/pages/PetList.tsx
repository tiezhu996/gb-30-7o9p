import { Col, Radio, Row, Select, Space } from 'antd'
import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import PetCard from '@/components/common/PetCard'
import SearchFilter from '@/components/common/SearchFilter'
import { usePetStore } from '@/stores/petStore'
import { PetSpeciesMap } from '@/constants/pet'

export default function PetList() {
  const [params] = useSearchParams()
  const { pets, total, load } = usePetStore()
  const [species, setSpecies] = useState<string>('')
  const [status, setStatus] = useState<string>('')
  const [keyword, setKeyword] = useState<string>(params.get('keyword') || '')

  useEffect(() => {
    load({ page: 1, page_size: 12, species, status, keyword })
  }, [species, status])

  function onSearch(kw: string) {
    setKeyword(kw)
    load({ page: 1, page_size: 12, species, status, keyword: kw })
  }

  return (
    <div>
      <h1>待领养动物</h1>
      <SearchFilter onSearch={onSearch}>
        <Space wrap>
          <Select
            placeholder="动物种类"
            allowClear
            style={{ width: 140 }}
            value={species || undefined}
            onChange={(v) => setSpecies(v || '')}
            options={Object.entries(PetSpeciesMap).map(([value, label]) => ({ value, label }))}
          />
          <Radio.Group value={status} onChange={(e) => setStatus(e.target.value)}>
            <Radio.Button value="">全部</Radio.Button>
            <Radio.Button value="available">可领养</Radio.Button>
            <Radio.Button value="adopted">已领养</Radio.Button>
          </Radio.Group>
        </Space>
      </SearchFilter>
      <Row gutter={[16, 16]}>
        {pets.map((p) => (
          <Col xs={12} sm={8} md={6} key={p.id}>
            <PetCard pet={p} />
          </Col>
        ))}
      </Row>
      {!pets.length && <p style={{ color: '#999' }}>暂无匹配动物（共 {total} 条）</p>}
    </div>
  )
}
