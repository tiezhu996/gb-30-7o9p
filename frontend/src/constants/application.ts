export type ApplicationStatus =
  | 'submitted'
  | 'org_review'
  | 'communicating'
  | 'confirmed'
  | 'offline_interview'
  | 'reserved'
  | 'waitlisted'
  | 'approved'
  | 'rejected'
  | 'withdrawn'
  | 'cancelled'
  | 'closed'

export const ApplicationStatusMap: Record<ApplicationStatus, { text: string; color: string }> = {
  submitted: { text: '已提交', color: 'blue' },
  org_review: { text: '机构审核中', color: 'gold' },
  communicating: { text: '沟通中', color: 'cyan' },
  confirmed: { text: '已确认', color: 'geekblue' },
  offline_interview: { text: '线下面签', color: 'purple' },
  reserved: { text: '预留中', color: 'orange' },
  waitlisted: { text: '候补中', color: 'volcano' },
  approved: { text: '已领养', color: 'green' },
  rejected: { text: '已拒绝', color: 'red' },
  withdrawn: { text: '已放弃', color: 'default' },
  cancelled: { text: '预留已取消', color: 'default' },
  closed: { text: '已结束', color: 'default' },
}

export type ApplicationEndReason =
  | 'final_adoption'
  | 'adopter_gave_up'
  | 'reservation_cancelled'
  | 'final_rejection'
  | 'withdrawn'
  | 'superseded'

export const ApplicationEndReasonText: Record<ApplicationEndReason, string> = {
  final_adoption: '已确认最终领养',
  adopter_gave_up: '获选人放弃，名额已释放',
  reservation_cancelled: '机构取消预留，名额已释放',
  final_rejection: '最终审核拒绝',
  withdrawn: '申请人主动撤回',
  superseded: '已有其他申请人完成领养',
}
