import { create } from 'zustand'
import { listMyReviews } from '@/api/review'
import type { VisitReview } from '@/types/api'

interface ReviewState {
  reviews: VisitReview[]
  load: () => Promise<void>
}

export const useReviewStore = create<ReviewState>((set) => ({
  reviews: [],
  load: async () => set({ reviews: await listMyReviews() }),
}))
