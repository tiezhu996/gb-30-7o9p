export type PetSpecies = 'dog' | 'cat' | 'rabbit' | 'other'
export type PetStatus = 'available' | 'pending' | 'adopted'

export const PetSpeciesMap: Record<PetSpecies, string> = {
  dog: '犬',
  cat: '猫',
  rabbit: '兔',
  other: '其他',
}

export const PetStatusMap: Record<PetStatus, { text: string; color: string }> = {
  available: { text: '可领养', color: 'green' },
  pending: { text: '申请中', color: 'gold' },
  adopted: { text: '已领养', color: 'default' },
}

export interface Pet {
  id: number
  org_id: number
  name: string
  species: PetSpecies
  breed: string
  age: number
  gender: string
  size: string
  city: string
  description: string
  personality: string
  health_status: string
  neutered: boolean
  vaccinated: boolean
  image_urls: string
  status: PetStatus
  created_at: string
}

export function parseImages(raw: string): string[] {
  try {
    const arr = JSON.parse(raw || '[]')
    return Array.isArray(arr) ? arr : []
  } catch {
    return []
  }
}
