import { Button, Form, Input } from 'antd'
import { useState } from 'react'

export default function SearchFilter({
  onSearch,
  children,
}: {
  onSearch: (keyword: string) => void
  children?: React.ReactNode
}) {
  const [keyword, setKeyword] = useState('')
  return (
    <Form layout="inline" style={{ marginBottom: 16 }}>
      {children}
      <Form.Item>
        <Input.Search
          placeholder="输入关键词搜索"
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          onSearch={onSearch}
          allowClear
        />
      </Form.Item>
      <Form.Item>
        <Button onClick={() => onSearch('')}>重置</Button>
      </Form.Item>
    </Form>
  )
}
