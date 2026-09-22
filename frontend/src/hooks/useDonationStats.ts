import { useEffect, useState } from 'react'
import { listMyDonations, listUsages } from '@/api/donation'
import type { Donation, DonationUsage } from '@/types/api'

export function useDonationStats(orgId?: number) {
  const [donations, setDonations] = useState<Donation[]>([])
  const [usages, setUsages] = useState<DonationUsage[]>([])

  async function load() {
    setDonations(await listMyDonations())
    if (orgId) {
      setUsages(await listUsages(orgId))
    }
  }

  useEffect(() => {
    load()
  }, [orgId])

  const total = donations.reduce((s, d) => s + Number(d.amount), 0)
  return { donations, usages, total, reload: load }
}
