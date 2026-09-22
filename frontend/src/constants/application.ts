export type ApplicationStatus =
  | 'submitted'
  | 'org_review'
  | 'communicating'
  | 'confirmed'
  | 'offline_interview'
  | 'approved'
  | 'rejected'
  | 'reserved'
  | 'waitlisted'
  | 'cancelled'
  | 'closed'

export const ApplicationStatusMap: Record<ApplicationStatus, { text: string; color: string }> = {
  submitted: { text: '已提交', color: 'blue' },
  org_review: { text: '机构审核中', color: 'gold' },
  communicating: { text: '沟通中', color: 'cyan' },
  confirmed: { text: '已确认', color: 'geekblue' },
  offline_interview: { text: '线下面签', color: 'purple' },
  approved: { text: '已通过', color: 'green' },
  rejected: { text: '已拒绝', color: 'red' },
  reserved: { text: '预留中', color: 'volcano' },
  waitlisted: { text: '候补中', color: 'orange' },
  cancelled: { text: '已取消', color: 'default' },
  closed: { text: '已结束', color: 'default' },
}

export const ApplicationStatusOptions = Object.entries(ApplicationStatusMap).map(([value, meta]) => ({
  value,
  label: meta.text,
}))

// Statuses that are still in the open application period and can be selected.
export const OPEN_APPLICATION_STATUSES = [
  'submitted',
  'org_review',
  'communicating',
  'confirmed',
  'offline_interview',
]

export const TERMINAL_APPLICATION_STATUSES = ['approved', 'rejected', 'cancelled', 'closed']

export function isTerminalApplicationStatus(status: string) {
  return TERMINAL_APPLICATION_STATUSES.includes(status)
}
