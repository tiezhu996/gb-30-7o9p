import request from '@/utils/request'
import type { AdoptionApplication } from '@/types/api'

export function submitApplication(payload: { pet_id: number; questionnaire?: string }) {
  return request.post<never, AdoptionApplication>('/applications', payload)
}

export function listMyApplications() {
  return request.get<never, AdoptionApplication[]>('/applications/me')
}

export function listOrgApplications(status?: string) {
  return request.get<never, AdoptionApplication[]>('/applications/org', { params: { status } })
}

export function updateApplicationStatus(id: number, status: string) {
  return request.put<never, AdoptionApplication>(`/applications/${id}/status`, { status })
}
