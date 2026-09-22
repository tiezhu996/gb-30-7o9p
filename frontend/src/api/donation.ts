import request from '@/utils/request'
import type { Donation, DonationUsage } from '@/types/api'

export function donate(payload: { org_id: number; amount: number }) {
  return request.post<never, Donation>('/donations', payload)
}

export function listMyDonations() {
  return request.get<never, Donation[]>('/donations/me')
}

export function createUsage(payload: { donation_id: number; amount: number; usage_desc: string; proof_url?: string }) {
  return request.post<never, DonationUsage>('/donations/usage', payload)
}

export function listUsages(orgId: number) {
  return request.get<never, DonationUsage[]>(`/orgs/${orgId}/usages`)
}

export function donationStats(orgId: number) {
  return request.get<never, { org_id: number; total_donated: number }>(`/orgs/${orgId}/donation-stats`)
}
