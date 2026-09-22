import { Modal, Input, message } from 'antd'
import { useEffect, useState } from 'react'
import { listMyReviews, submitReview } from '@/api/review'
import ReviewTaskList from '@/components/common/ReviewTaskList'
import { useAuth } from '@/hooks/useAuth'
import type { VisitReview } from '@/types/api'

export default function Reviews() {
  const { isOrg } = useAuth()
  const [reviews, setReviews] = useState<VisitReview[]>([])
  const [current, setCurrent] = useState<VisitReview | null>(null)
  const [note, setNote] = useState('')

  async function load() {
    if (isOrg) {
      setReviews(await import('@/api/review').then((m) => m.listOrgReviews()))
    } else {
      setReviews(await listMyReviews())
    }
  }
  useEffect(() => {
    load()
  }, [isOrg])

  async function onSubmit(r: VisitReview) {
    setCurrent(r)
    setNote('')
  }

  async function confirmSubmit() {
    if (!current) return
    await submitReview(current.id, { note })
    message.success('回访已提交')
    setCurrent(null)
    await load()
  }

  return (
    <div>
      <h1>领养回访任务</h1>
      <ReviewTaskList reviews={reviews} onSubmit={onSubmit} />
      <Modal open={!!current} title="提交回访" onOk={confirmSubmit} onCancel={() => setCurrent(null)}>
        <Input.TextArea rows={3} value={note} onChange={(e) => setNote(e.target.value)} placeholder="回访备注（可选）" />
      </Modal>
    </div>
  )
}
