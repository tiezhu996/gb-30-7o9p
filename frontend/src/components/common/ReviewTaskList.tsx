import { Button, Card, Tag } from 'antd'
import type { VisitReview } from '@/types/api'
import { formatDate } from '@/utils/dateFormat'

const statusMap: Record<string, { text: string; color: string }> = {
  pending: { text: '待提交', color: 'gold' },
  submitted: { text: '已提交', color: 'green' },
  overdue: { text: '已逾期', color: 'red' },
}

export default function ReviewTaskList({
  reviews,
  onSubmit,
}: {
  reviews: VisitReview[]
  onSubmit: (review: VisitReview) => void
}) {
  return (
    <div>
      {reviews.length === 0 && <Card>暂无回访任务</Card>}
      {reviews.map((r) => {
        const s = statusMap[r.status]
        return (
          <Card key={r.id} size="small" style={{ marginBottom: 8 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <span>
                回访任务 #{r.id} · 截止 {formatDate(r.due_date)} <Tag color={s.color}>{s.text}</Tag>
              </span>
              {r.status !== 'submitted' && (
                <Button size="small" type="primary" onClick={() => onSubmit(r)}>
                  提交回访
                </Button>
              )}
            </div>
          </Card>
        )
      })}
    </div>
  )
}
