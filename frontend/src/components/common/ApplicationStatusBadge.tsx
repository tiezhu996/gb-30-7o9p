import { Tag } from 'antd'
import { ApplicationStatusMap, type ApplicationStatus } from '@/constants/application'

export default function ApplicationStatusBadge({ status }: { status: string }) {
  const meta = ApplicationStatusMap[status as ApplicationStatus]
  if (!meta) return <Tag>{status}</Tag>
  return <Tag color={meta.color}>{meta.text}</Tag>
}
