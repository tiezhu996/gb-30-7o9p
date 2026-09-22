import request from '@/utils/request'
import type { VisitReview } from '@/types/api'

export function listMyReviews() {
  return request.get<never, VisitReview[]>('/reviews/me')
}

export function listOrgReviews() {
  return request.get<never, VisitReview[]>('/reviews/org')
}

export function createReview(payload: { application_id: number; scheduled_days?: number }) {
  return request.post<never, VisitReview>('/reviews', payload)
}

export function submitReview(id: number, payload: { photos?: string; note?: string }) {
  return request.put<never, VisitReview>(`/reviews/${id}/submit`, payload)
}
