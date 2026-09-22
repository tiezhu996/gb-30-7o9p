import { Card } from 'antd'
import { useParams } from 'react-router-dom'
import { useEffect, useState } from 'react'
import ApplicationForm from '@/components/common/ApplicationForm'
import { getPet } from '@/api/pet'
import type { Pet } from '@/constants/pet'

export default function Apply() {
  const { petId } = useParams()
  const [pet, setPet] = useState<Pet | null>(null)

  useEffect(() => {
    getPet(petId!).then(setPet)
  }, [petId])

  return (
    <div style={{ maxWidth: 560, margin: '0 auto' }}>
      <h1>领养申请 - {pet?.name || ''}</h1>
      <Card>
        <ApplicationForm petId={Number(petId)} />
      </Card>
    </div>
  )
}
