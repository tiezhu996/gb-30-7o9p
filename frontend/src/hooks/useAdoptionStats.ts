import { useEffect, useState } from 'react'
import { listMyApplications } from '@/api/application'
import type { AdoptionApplication } from '@/types/api'

export function useAdoptionStats() {
  const [apps, setApps] = useState<AdoptionApplication[]>([])
  const [loading, setLoading] = useState(false)

  async function load() {
    setLoading(true)
    try {
      setApps(await listMyApplications())
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  const byStatus = apps.reduce<Record<string, number>>((acc, a) => {
    acc[a.status] = (acc[a.status] || 0) + 1
    return acc
  }, {})
  const approved = apps.filter((a) => a.status === 'approved').length

  return { apps, loading, byStatus, approved, reload: load }
}
