export type OrganizationStatus = 'pending' | 'approved' | 'rejected'

export const OrganizationStatusMap: Record<OrganizationStatus, { text: string; color: string }> = {
  pending: { text: '待审核', color: 'gold' },
  approved: { text: '已认证', color: 'green' },
  rejected: { text: '已驳回', color: 'red' },
}

export interface Organization {
  id: number
  user_id: number
  name: string
  license_url: string
  cert_type: string
  status: OrganizationStatus
  contact: string
  city: string
  description: string
  created_at: string
}
