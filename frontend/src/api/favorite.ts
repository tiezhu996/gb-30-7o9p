import request from '@/utils/request'
import type { Favorite } from '@/types/api'

export function listFavorites() {
  return request.get<never, Favorite[]>('/favorites')
}

export function addFavorite(targetType: string, targetId: number) {
  return request.post<never, Favorite>('/favorites', { target_type: targetType, target_id: targetId })
}

export function removeFavorite(targetType: string, targetId: number) {
  return request.delete<never, { removed: boolean }>(`/favorites/${targetType}/${targetId}`)
}
