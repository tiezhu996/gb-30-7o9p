import { Button, Card, Popconfirm, Select, Space, Table, Tag, message } from 'antd'
import { useEffect, useState } from 'react'
import { listMyApplications, updateApplicationStatus } from '@/api/application'
import ApplicationStatusBadge from '@/components/common/ApplicationStatusBadge'
import {
  ApplicationEndReasonText,
  type ApplicationEndReason,
} from '@/constants/application'
import { useAuth } from '@/hooks/useAuth'
import type { AdoptionApplication } from '@/types/api'
import { formatDate } from '@/utils/dateFormat'

const STATUS_OPTIONS = [
  'submitted',
  'org_review',
  'communicating',
  'confirmed',
  'offline_interview',
  'reserved',
  'waitlisted',
  'approved',
  'rejected',
  'withdrawn',
  'cancelled',
  'closed',
]

const REVIEW_NEXT: Record<string, { status: string; label: string }> = {
  submitted: { status: 'org_review', label: '开始审核' },
  reserved: { status: 'org_review', label: '开始审核' },
  org_review: { status: 'communicating', label: '进入沟通' },
  communicating: { status: 'confirmed', label: '确认意向' },
  confirmed: { status: 'offline_interview', label: '安排面签' },
}

export default function Applications() {
  const { isOrg, user } = useAuth()
  const [apps, setApps] = useState<AdoptionApplication[]>([])
  const [status, setStatus] = useState('')

  async function load(s = status) {
    if (isOrg) {
      const res = await import('@/api/application').then((m) => m.listOrgApplications(s))
      setApps(res)
    } else {
      setApps(await listMyApplications())
    }
  }
  useEffect(() => {
    load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isOrg])

  async function changeStatus(id: number, next: string) {
    await updateApplicationStatus(id, next)
    message.success('状态已更新')
    await load()
  }

  function isReservationHolder(r: AdoptionApplication) {
    return (
      r.reserved_user_id === r.user_id ||
      ['org_review', 'communicating', 'confirmed', 'offline_interview'].includes(r.status)
    )
  }

  const columns = [
    { title: 'ID', dataIndex: 'id' },
    { title: '宠物 ID', dataIndex: 'pet_id' },
    { title: '申请时间', dataIndex: 'created_at', render: (v: string) => formatDate(v) },
    { title: '状态', dataIndex: 'status', render: (v: string) => <ApplicationStatusBadge status={v} /> },
    {
      title: '预留人 / 名次',
      render: (_: unknown, r: AdoptionApplication) => {
        if (r.status === 'approved') return <Tag color="green">最终领养人</Tag>
        if (isReservationHolder(r)) {
          const name = r.reserved_user_name || `用户 ${r.reserved_user_id || r.user_id}`
          return <Tag color="orange">预留人：{name}</Tag>
        }
        if (r.status === 'waitlisted') return <Tag color="volcano">候补第 {r.waitlist_position} 名</Tag>
        return '-'
      },
    },
    {
      title: '结束原因',
      dataIndex: 'end_reason',
      render: (v: string) =>
        v ? ApplicationEndReasonText[v as ApplicationEndReason] || v : '-',
    },
    {
      title: '操作',
      render: (_: unknown, r: AdoptionApplication) => {
        if (!user) return null
        const review = REVIEW_NEXT[r.status]
        if (isOrg) {
          const active = ['submitted', 'reserved', 'org_review', 'communicating', 'confirmed', 'offline_interview'].includes(r.status)
          if (!active) return '-'
          const holder = isReservationHolder(r)
          return (
            <Space wrap>
              {!holder && (
                <Button size="small" type="primary" onClick={() => changeStatus(r.id, 'reserved')}>
                  选中预留
                </Button>
              )}
              {holder && review && (
                <Button size="small" onClick={() => changeStatus(r.id, review.status)}>
                  {review.label}
                </Button>
              )}
              {holder && r.status === 'offline_interview' && (
                <Popconfirm title="确认最终领养？" onConfirm={() => changeStatus(r.id, 'approved')}>
                  <Button size="small" type="primary">
                    确认领养
                  </Button>
                </Popconfirm>
              )}
              {holder && (
                <Popconfirm title="取消预留并释放给首位候补？" onConfirm={() => changeStatus(r.id, 'cancelled')}>
                  <Button size="small">取消预留</Button>
                </Popconfirm>
              )}
              <Popconfirm title="拒绝该申请？" onConfirm={() => changeStatus(r.id, 'rejected')}>
                <Button size="small" danger>
                  拒绝
                </Button>
              </Popconfirm>
            </Space>
          )
        }
        const canGiveUp = ['reserved', 'org_review', 'communicating', 'confirmed', 'offline_interview'].includes(r.status)
        const canWithdraw = r.status === 'submitted' || r.status === 'waitlisted' || canGiveUp
        if (!canWithdraw) return '-'
        return (
          <Popconfirm
            title={canGiveUp ? '放弃预留？名额将释放给首位候补' : '撤回该申请？'}
            onConfirm={() => changeStatus(r.id, 'withdrawn')}
          >
            <Button size="small" danger>
              {canGiveUp ? '放弃领养' : '撤回申请'}
            </Button>
          </Popconfirm>
        )
      },
    },
  ]

  return (
    <div>
      <h1>领养申请</h1>
      {isOrg && (
        <Select
          style={{ width: 200, marginBottom: 12 }}
          placeholder="按状态筛选"
          allowClear
          value={status || undefined}
          onChange={(v) => {
            setStatus(v || '')
            load(v || '')
          }}
          options={STATUS_OPTIONS.map((s) => ({ value: s, label: s }))}
        />
      )}
      <Table rowKey="id" dataSource={apps} pagination={false} columns={columns} />
      {!apps.length && <Card style={{ marginTop: 12 }}>暂无申请记录</Card>}
    </div>
  )
}
