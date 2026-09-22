import request from '@/utils/request'
import type { Organization } from '@/constants/organization'
import type { PageData } from '@/types/api'

export function listOrgs(params: { page?: number; page_size?: number; status?: string; keyword?: string }) {
  return request.get<never, PageData<Organization>>('/orgs', { params })
}

export function getOrg(id: number | string) {
  return request.get<never, Organization>(`/orgs/${id}`)
}

export function registerOrg(payload: { name: string; license_url?: string; cert_type?: string; contact?: string; city?: string; description?: string }) {
  return request.post<never, Organization>('/orgs', payload)
}

export function reviewOrg(id: number, status: string) {
  return request.put<never, Organization>(`/orgs/${id}/review`, { status })
}
