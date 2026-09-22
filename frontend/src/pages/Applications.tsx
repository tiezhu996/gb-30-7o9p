import { Button, Card, Popconfirm, Select, Space, Table, Tag, message } from 'antd'
import { useEffect, useState } from 'react'
import {
  listMyApplications,
  listOrgApplications,
  releaseReservation,
  selectAdopter,
  updateApplicationStatus,
} from '@/api/application'
import ApplicationStatusBadge from '@/components/common/ApplicationStatusBadge'
import { ApplicationStatusOptions, isTerminalApplicationStatus } from '@/constants/application'
import { useAuth } from '@/hooks/useAuth'
import type { AdoptionApplication } from '@/types/api'
import { formatDate } from '@/utils/dateFormat'

// Open-period review steps share status names with the holder's pipeline;
// the holder is the row belonging to the current reservation holder.
function isHolder(r: AdoptionApplication) {
  return r.pet_status === 'reserved' && r.reserved_user_id === r.user_id
}

export default function Applications() {
  const { isOrg } = useAuth()
  const [apps, setApps] = useState<AdoptionApplication[]>([])
  const [status, setStatus] = useState('')

  async function load(s = status) {
    if (isOrg) {
      setApps(await listOrgApplications(s))
    } else {
      setApps(await listMyApplications())
    }
  }
  useEffect(() => {
    load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isOrg])

  async function selectHolder(id: number) {
    await selectAdopter(id)
    message.success('已选中领养人，宠物进入预留')
    await load()
  }

  async function cancelReserved(id: number) {
    await releaseReservation(id, 'cancel')
    message.success('已取消预留，首位候补将自动递补')
    await load()
  }

  async function abandon(id: number) {
    await releaseReservation(id, 'abandon')
    message.success('已放弃预留')
    await load()
  }

  async function changeStatus(id: number, next: string) {
    await updateApplicationStatus(id, next)
    message.success('状态已更新')
    await load()
  }

  const columns = [
    { title: 'ID', dataIndex: 'id' },
    { title: '宠物', dataIndex: 'pet_name', render: (v: string, r: AdoptionApplication) => v || `宠物 #${r.pet_id}` },
    { title: '申请时间', dataIndex: 'created_at', render: (v: string) => formatDate(v) },
    {
      title: '状态',
      dataIndex: 'status',
      render: (v: string) => <ApplicationStatusBadge status={v} />,
    },
    {
      title: '预留人',
      key: 'reserved',
      render: (_: unknown, r: AdoptionApplication) =>
        r.reserved_name ? (
          <Tag color={isHolder(r) ? 'volcano' : 'default'}>
            {isHolder(r) ? '本人预留' : `预留人：${r.reserved_name}`}
          </Tag>
        ) : (
          '-'
        ),
    },
    {
      title: '候补名次',
      key: 'rank',
      render: (_: unknown, r: AdoptionApplication) => {
        if (r.status === 'waitlisted') {
          return <Tag color="orange">候补第 {r.waitlist_rank} 名</Tag>
        }
        if (r.status === 'closed' && r.end_reason === 'adopted_by_other' && r.waitlist_rank > 0) {
          return <Tag>原候补第 {r.waitlist_rank} 名</Tag>
        }
        return '-'
      },
    },
    {
      title: '结束原因',
      dataIndex: 'end_reason_text',
      render: (v: string) => v || '-',
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: unknown, r: AdoptionApplication) => {
        if (isTerminalApplicationStatus(r.status)) return <span style={{ color: '#999' }}>已结束</span>
        if (isOrg) {
          const holder = isHolder(r)
          return (
            <Space wrap>
              {!holder && r.status === 'submitted' && (
                <>
                  <Popconfirm title="确认选中该领养人？其余申请将进入候补" onConfirm={() => selectHolder(r.id)}>
                    <Button size="small" type="primary">
                      选中预留
                    </Button>
                  </Popconfirm>
                  <Button size="small" onClick={() => changeStatus(r.id, 'org_review')}>
                    开始审核
                  </Button>
                  <Button size="small" danger onClick={() => changeStatus(r.id, 'rejected')}>
                    拒绝
                  </Button>
                </>
              )}
              {!holder && r.status === 'org_review' && (
                <>
                  <Popconfirm title="确认选中该领养人？其余申请将进入候补" onConfirm={() => selectHolder(r.id)}>
                    <Button size="small" type="primary">
                      选中预留
                    </Button>
                  </Popconfirm>
                  <Button size="small" type="primary" onClick={() => changeStatus(r.id, 'communicating')}>
                    推进到沟通
                  </Button>
                  <Button size="small" danger onClick={() => changeStatus(r.id, 'rejected')}>
                    拒绝
                  </Button>
                </>
              )}
              {!holder && r.status === 'communicating' && (
                <>
                  <Popconfirm title="确认选中该领养人？其余申请将进入候补" onConfirm={() => selectHolder(r.id)}>
                    <Button size="small" type="primary">
                      选中预留
                    </Button>
                  </Popconfirm>
                  <Button size="small" type="primary" onClick={() => changeStatus(r.id, 'confirmed')}>
                    确认意向
                  </Button>
                  <Button size="small" danger onClick={() => changeStatus(r.id, 'rejected')}>
                    拒绝
                  </Button>
                </>
              )}
              {!holder && r.status === 'confirmed' && (
                <>
                  <Popconfirm title="确认选中该领养人？其余申请将进入候补" onConfirm={() => selectHolder(r.id)}>
                    <Button size="small" type="primary">
                      选中预留
                    </Button>
                  </Popconfirm>
                  <Button size="small" type="primary" onClick={() => changeStatus(r.id, 'offline_interview')}>
                    安排面签
                  </Button>
                </>
              )}
              {!holder && r.status === 'offline_interview' && (
                <Popconfirm title="确认选中该领养人？其余申请将进入候补" onConfirm={() => selectHolder(r.id)}>
                  <Button size="small" type="primary">
                    选中预留
                  </Button>
                </Popconfirm>
              )}
              {holder && r.status === 'reserved' && (
                <>
                  <Button size="small" type="primary" onClick={() => changeStatus(r.id, 'org_review')}>
                    开始审核
                  </Button>
                  <Button size="small" danger onClick={() => changeStatus(r.id, 'rejected')}>
                    最终拒绝
                  </Button>
                  <Popconfirm title="确认取消预留？将由首位候补递补" onConfirm={() => cancelReserved(r.id)}>
                    <Button size="small">取消预留</Button>
                  </Popconfirm>
                </>
              )}
              {holder && r.status === 'org_review' && (
                <>
                  <Button size="small" type="primary" onClick={() => changeStatus(r.id, 'communicating')}>
                    推进到沟通
                  </Button>
                  <Button size="small" danger onClick={() => changeStatus(r.id, 'rejected')}>
                    最终拒绝
                  </Button>
                  <Popconfirm title="确认取消预留？将由首位候补递补" onConfirm={() => cancelReserved(r.id)}>
                    <Button size="small">取消预留</Button>
                  </Popconfirm>
                </>
              )}
              {holder && r.status === 'communicating' && (
                <>
                  <Button size="small" type="primary" onClick={() => changeStatus(r.id, 'confirmed')}>
                    确认意向
                  </Button>
                  <Button size="small" danger onClick={() => changeStatus(r.id, 'rejected')}>
                    最终拒绝
                  </Button>
                  <Popconfirm title="确认取消预留？将由首位候补递补" onConfirm={() => cancelReserved(r.id)}>
                    <Button size="small">取消预留</Button>
                  </Popconfirm>
                </>
              )}
              {holder && r.status === 'confirmed' && (
                <>
                  <Button size="small" type="primary" onClick={() => changeStatus(r.id, 'offline_interview')}>
                    安排面签
                  </Button>
                  <Button size="small" danger onClick={() => changeStatus(r.id, 'rejected')}>
                    最终拒绝
                  </Button>
                  <Popconfirm title="确认取消预留？将由首位候补递补" onConfirm={() => cancelReserved(r.id)}>
                    <Button size="small">取消预留</Button>
                  </Popconfirm>
                </>
              )}
              {holder && r.status === 'offline_interview' && (
                <>
                  <Button size="small" type="primary" onClick={() => changeStatus(r.id, 'approved')}>
                    确认最终领养
                  </Button>
                  <Button size="small" danger onClick={() => changeStatus(r.id, 'rejected')}>
                    最终拒绝
                  </Button>
                  <Popconfirm title="确认取消预留？将由首位候补递补" onConfirm={() => cancelReserved(r.id)}>
                    <Button size="small">取消预留</Button>
                  </Popconfirm>
                </>
              )}
              {r.status === 'waitlisted' && <Tag color="orange">候补第 {r.waitlist_rank} 名</Tag>}
            </Space>
          )
        }
        // Applicant side
        if (isHolder(r)) {
          return (
            <Space direction="vertical" size={4}>
              <Tag color="volcano">预留人</Tag>
              <Popconfirm title="确认放弃预留？名额将由首位候补递补" onConfirm={() => abandon(r.id)}>
                <Button size="small" danger>
                  放弃预留
                </Button>
              </Popconfirm>
            </Space>
          )
        }
        if (r.status === 'waitlisted') {
          return <span style={{ color: '#fa8c16' }}>候补第 {r.waitlist_rank} 名，等待递补</span>
        }
        return <span style={{ color: '#999' }}>审核中</span>
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
          options={ApplicationStatusOptions}
        />
      )}
      <Table rowKey="id" dataSource={apps} pagination={false} columns={columns} scroll={{ x: 1200 }} />
      {!apps.length && <Card style={{ marginTop: 12 }}>暂无申请记录</Card>}
    </div>
  )
}
