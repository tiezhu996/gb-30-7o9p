import { useEffect, useState } from 'react'
import { listMyReviews } from '@/api/review'
import type { VisitReview } from '@/types/api'

export function useReviewTasks() {
  const [reviews, setReviews] = useState<VisitReview[]>([])
  const [loading, setLoading] = useState(false)

  async function load() {
    setLoading(true)
    try {
      setReviews(await listMyReviews())
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  const overdue = reviews.filter((r) => r.status === 'overdue').length
  const pending = reviews.filter((r) => r.status === 'pending').length

  return { reviews, loading, overdue, pending, reload: load }
}
